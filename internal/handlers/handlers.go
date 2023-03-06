package handlers

import (
	"bufio"
	"context"
	"fmt"
	"gitlab.com/distributed_lab/logan/v3"
	"gitlab.com/distributed_lab/running"
	"helper/internal/config"
	"helper/internal/data"
	"helper/internal/service/google"
	"os"
	"strings"
	"sync"
	"time"
)

type Path struct {
	Path string
	ID   int
}

type Handler struct {
	mu           sync.Mutex
	running      int
	chInput      chan Path
	chOutput     chan Path
	minTime      time.Duration
	maxTime      time.Duration
	name         string
	log          *logan.Entry
	ctx          context.Context
	cancel       context.CancelFunc
	googleClient *google.Google
	count        int32
}

func NewHandler(input chan Path, output chan Path, log *logan.Entry, google *google.Google, ran int, ctx context.Context) Handler {
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
					link, err := h.googleClient.Update(path.Path)
					if err != nil {
						h.log.Error(h.name, "--->", "error: ", err)
						return false, err
					}
					path.Path = link
					h.chOutput <- path
					h.log.Debug("send ", name)
					h.count++
					return true, err
				}, time.Millisecond*150, time.Millisecond*180)
			}
		}(fmt.Sprintf("%s-%d", h.name, i))
	}
}

func (h *Handler) StartSheetRunner() { //todo do more useful
	for i := 0; i < h.running; i++ {
		h.log.Debug("start ", i)
		go func(name string) {
			defer h.decrement()
			defer h.log.Debug("quit ", i)
			for path := range h.chInput {
				running.UntilSuccess(context.Background(), h.log, h.name, func(ctx context.Context) (bool, error) {
					link, err := h.googleClient.Update(path.Path)
					if err != nil {
						h.log.Error(h.name, "--->", "error: ", err)
						return false, err
					}
					path.Path = link
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

func (h *Handler) Read(users []*data.User) []*data.User {
	for {
		select {
		case path := <-h.chOutput:
			h.log.Debug("read")
			users[path.ID].DigitalCertificate = path.Path
		case <-h.ctx.Done():
			h.log.Debug("out")
			return users
		}
	}
}

func (h *Handler) insertData(paths []Path) {
	for _, path := range paths {
		h.chInput <- path
	}
	close(h.chInput)
}

func Drive(cfg config.Config, log *logan.Entry, paths []Path, users []*data.User) ([]*data.User, error) {
	var err error
	input := make(chan Path)
	output := make(chan Path)
	sendToDrive := cfg.Google().Enable

	var googleClient *google.Google
	if sendToDrive {
		googleClient = google.NewGoogleClient(cfg)

		err = googleClient.ConnectToDrive(cfg.Google().SecretPath, cfg.Google().Code)
		if err != nil {
			log.Info("Could you continue to work without google drive? (y)")
			reader := bufio.NewReader(os.Stdin)
			text, _ := reader.ReadString('\n')
			if strings.ToLower(text) != "y\n" {
				sendToDrive = true
			}
		}
	}
	ctx := context.Background()

	if sendToDrive {
		err = googleClient.CreateFolder(cfg.Google().QRPath)
		if err != nil {
			cfg.Log().Error(err)
			return users, err
		}
		handler := NewHandler(input, output, log, googleClient, 10, ctx)
		handler.StartDriveRunner()
		go handler.insertData(paths)
		users = handler.Read(users)
		log.Info("sent to drive: ", handler.count)
	}

	return users, err
}
