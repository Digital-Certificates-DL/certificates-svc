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
	"time"
)

type Path struct {
	Path string
	ID   int
}

type Handler struct {
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

func (h *Handler) StartRunner() {
	for i := 0; i < h.running; i++ {
		h.log.Debug("start ", i)
		go func(name string) {
			defer h.decrement()
			defer h.log.Debug("quit ", i)
			for path := range h.chInput {
				running.UntilSuccess(context.Background(), h.log, h.name, func(ctx context.Context) (bool, error) {
					link, success, err := h.googleClient.Update(path.Path)
					if err != nil {
						h.log.Error(h.name, "--->", "error: ", err)
						return false, err
					}

					path.Path = link
					h.chOutput <- path

					h.log.Debug("send ", name)
					return success, err
				}, time.Millisecond*150, time.Millisecond*180) //todo move config file
			}
		}(fmt.Sprintf("%s-%d", h.name, i))

	}
}

func (h *Handler) decrement() {
	h.running--
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
			h.log.Debug("read ")

			users[path.ID].DigitalCertificate = path.Path
		case <-h.ctx.Done():

			h.log.Info("out ")
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

		err = googleClient.Connect(cfg.Google().SecretPath, cfg.Google().Code)
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
			cfg.Log().Debug(err)
			return users, err
		}
		handler := NewHandler(input, output, log, googleClient, 10, ctx)
		handler.StartRunner()
		go handler.insertData(paths)

		users = handler.Read(users)
	}

	log.Debug("move out go")

	return users, err
}
