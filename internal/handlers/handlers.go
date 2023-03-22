package handlers

import (
	"context"
	"fmt"
	"github.com/pkg/errors"
	"gitlab.com/distributed_lab/running"
	"gitlab.com/tokend/course-certificates/ccp/internal/config"
	"gitlab.com/tokend/course-certificates/ccp/internal/data"
	"gitlab.com/tokend/course-certificates/ccp/internal/service/google"
	"strings"
	"sync"
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
	mu           sync.Mutex
	running      int
	chInput      chan FilesBytes
	chOutput     chan FilesBytes
	minTime      time.Duration
	maxTime      time.Duration
	name         string
	log          *logan.Entry
	ctx          context.Context
	cancel       context.CancelFunc
	googleClient *google.Google
	count        int32
}

func NewHandler(input chan FilesBytes, output chan FilesBytes, log *logan.Entry, google *google.Google, ran int, ctx context.Context) Handler {
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
		h.log.Debug("start ", i)
		go func(name string) {
			defer h.decrement()
			defer h.log.Debug("quit ", i)
			for path := range h.chInput {
				running.UntilSuccess(context.Background(), h.log, h.name, func(ctx context.Context) (bool, error) {
					link, err := h.googleClient.Update(path.Name, path.File, path.Type)
					if err != nil {
						h.log.Error(h.name, "--->", "error: ", err)
						return false, err
					}
					path.Link = link
					h.chOutput <- path
					h.log.Debug("send ", name)
					h.count++
					return true, err
				}, time.Millisecond*150, time.Millisecond*180)
			}
		}(fmt.Sprintf("%s-%d", h.name, i))
	}
}

func (h *Handler) decrement() {
	h.mu.Lock()
	h.running--
	h.mu.Unlock()
	if h.running == 0 {
		close(h.chOutput)
		h.cancel()
		h.log.Debug("ctx done")
	}
}

func (h *Handler) Read(users []*data.User, flag string) []*data.User {
	for {
		select {
		case path := <-h.chOutput:
			h.log.Debug("read")
			users[path.ID] = h.setLink(users[path.ID], path, flag)

		case <-h.ctx.Done():
			h.log.Debug("out")
			return users
		}
	}
}

func (h *Handler) setLink(user *data.User, path FilesBytes, flag string) *data.User {
	switch strings.ToLower(flag) {
	case "qr":
		user.DigitalCertificate = path.Link
		return user
	case "certificate":
		user.Certificate = path.Link
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

func Drive(googleClient *google.Google, cfg config.Config, log *logan.Entry, files []FilesBytes, users []*data.User, flag string) ([]*data.User, error) {
	var err error
	input := make(chan FilesBytes)
	output := make(chan FilesBytes)

	ctx := context.Background()

	err = googleClient.CreateFolder(cfg.Google().QRPath)
	if err != nil {
		return users, errors.Wrap(err, "failed to create folder")
	}
	handler := NewHandler(input, output, log, googleClient, 10, ctx)
	handler.StartDriveRunner()
	go handler.insertData(files)
	users = handler.Read(users, flag)
	log.Info("sent to drive: ", handler.count)

	return users, err
}
