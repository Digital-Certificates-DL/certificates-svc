package handlers

import (
	"context"
	"fmt"
	"github.com/pkg/errors"
	"gitlab.com/distributed_lab/logan/v3"
	"gitlab.com/distributed_lab/running"
	"helper/internal/config"
	"helper/internal/data"
	"helper/internal/service/google"
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

//function func(ctx context.Context, log *logan.Entry, name string, runner func(ctx2 context.Context)) (bool, error)

func (h *Handler) StartDriveRunner() { //todo do more useful
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

//
//func (h *Handler) StartSheetRunner() { //todo do more useful
//	for i := 0; i < h.running; i++ {
//		h.log.Debug("start ", i)
//		go func(name string) {
//			defer h.decrement()
//			defer h.log.Debug("quit ", i)
//			for path := range h.chInput {
//				running.UntilSuccess(context.Background(), h.log, h.name, func(ctx context.Context) (bool, error) {
//					link, err := h.googleClient.Update(path.FilesBytes)
//					if err != nil {
//						h.log.Error(h.name, "--->", "error: ", err)
//						return false, err
//					}
//					path.FilesBytes = link
//					h.chOutput <- path
//					h.log.Debug("send ", name)
//					h.count++
//					return true, err
//				}, time.Millisecond*150, time.Millisecond*180)
//			}
//		}(fmt.Sprintf("%s-%d", h.name, i))
//	}
//}

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

func (h *Handler) Read(users []*data.User) []*data.User {
	for {
		select {
		case path := <-h.chOutput:
			h.log.Debug("read")
			users[path.ID].DigitalCertificate = path.Link
		case <-h.ctx.Done():
			h.log.Debug("out")
			return users
		}
	}
}

func (h *Handler) insertData(paths []FilesBytes) {
	for _, path := range paths {
		h.chInput <- path
	}
	close(h.chInput)
}

func Drive(googleClient *google.Google, cfg config.Config, log *logan.Entry, paths []FilesBytes, users []*data.User) ([]*data.User, error) {
	var err error
	input := make(chan FilesBytes)
	output := make(chan FilesBytes)
	sendToDrive := cfg.Google().Enable

	ctx := context.Background()
	if sendToDrive {
		err = googleClient.CreateFolder(cfg.Google().QRPath)
		if err != nil {
			return users, errors.Wrap(err, "failed to create folder")
		}
		handler := NewHandler(input, output, log, googleClient, 10, ctx)
		handler.StartDriveRunner()
		go handler.insertData(paths)
		users = handler.Read(users)
		log.Info("sent to drive: ", handler.count)
	}

	return users, err
}
