package google

import (
	"context"
	"fmt"
	"github.com/pkg/errors"
	"gitlab.com/distributed_lab/logan/v3"
	"gitlab.com/distributed_lab/running"
	"gitlab.com/tokend/course-certificates/ccp/internal/service/core/helpers"
	"strings"
	"time"
)

type FilesBytes struct {
	File []byte
	ID   int
	Name string
	Type string
	Link string
}

type Handler struct {
	running      int
	chInput      chan FilesBytes
	chOutput     chan FilesBytes
	minTime      time.Duration
	maxTime      time.Duration
	name         string
	log          *logan.Entry
	ctx          context.Context
	cancel       context.CancelFunc
	googleClient *Google
	count        int32
}

func NewHandler(input chan FilesBytes, output chan FilesBytes, log *logan.Entry, google *Google, ran int, ctx context.Context) Handler {
	ctxInner, cancel := context.WithCancel(ctx)
	return Handler{
		running:      ran,
		chInput:      input,
		chOutput:     output,
		log:          log,
		googleClient: google,
		name:         "handler",
		cancel:       cancel,
		ctx:          ctxInner,
	}
}

func (h *Handler) StartDriveRunner() {
	for i := 0; i < h.running; i++ {
		h.log.Debug("start worker ", i)
		go func(name string) {
			defer h.decrement()
			defer h.log.Debug("quit from worker: ", i)
			for path := range h.chInput {
				running.UntilSuccess(context.Background(), h.log, h.name, func(ctx context.Context) (bool, error) {
					link, err := h.googleClient.Update(path.Name, path.File, path.Type)
					if err != nil {
						h.log.Error(h.name, "--->", "error: ", err)
						return false, err
					}
					path.Link = link
					h.chOutput <- path
					h.log.Debug("send file to cloud and contain link on chan  ", name)
					h.count++
					return true, err
				}, time.Millisecond*150, time.Millisecond*180)
			}
		}(fmt.Sprintf("%s-%d", h.name, i))
	}
}

func (h *Handler) decrement() {
	h.running--
	if h.running == 0 {
		close(h.chOutput)
		h.cancel()
		h.log.Debug("worker context  is done")
	}
}

func (h *Handler) Read(users []*helpers.Certificate, flag string) []*helpers.Certificate {
	for {
		select {
		case path := <-h.chOutput:
			h.log.Debug("read from  worker chan")
			for id, u := range users {
				if u.ID == path.ID {
					users[id] = h.setLink(users[id], path, flag)
					h.log.Debug("set external link: ", flag)
					break
				}
			}
		case <-h.ctx.Done():
			h.log.Debug("exit from read func")
			return users
		}
	}
}

func (h *Handler) setLink(user *helpers.Certificate, path FilesBytes, flag string) *helpers.Certificate {
	switch strings.ToLower(flag) {
	case "qr":
		user.DigitalCertificate = strings.ReplaceAll(path.Link, "https://", "")
		return user
	case "certificate":
		user.Certificate = strings.ReplaceAll(path.Link, "https://", "")
		return user
	}
	return user
}

func (h *Handler) insertData(files []FilesBytes) {
	for _, path := range files {
		h.chInput <- path
	}
	close(h.chInput)
}

func Drive(googleClient *Google, log *logan.Entry, files []FilesBytes, users []*helpers.Certificate, flag string, folderName string) ([]*helpers.Certificate, error) {
	var err error
	input := make(chan FilesBytes)
	output := make(chan FilesBytes)

	if err := googleClient.CreateFolder(folderName); err != nil {
		return users, errors.Wrap(err, "failed to create folder")
	}

	handler := NewHandler(input, output, log, googleClient, 10, context.Background())
	handler.StartDriveRunner()
	go handler.insertData(files)
	users = handler.Read(users, flag)
	log.Info("sSent to drive: ", handler.count)

	return users, err
}
