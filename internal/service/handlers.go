package service

import (
	"context"
	"gitlab.com/distributed_lab/logan/v3"
	"gitlab.com/distributed_lab/running"
	"helper/internal/service/google"
	"sync"
	"time"
)

type Path struct {
	Path string
	ID   int
}

type Handler struct {
	data         []Path
	running      int
	chInput      chan Path
	chOutput     chan Path
	minTime      time.Duration
	maxTime      time.Duration
	name         string
	log          *logan.Entry
	googleClient *google.Google
}

func NewHandler(input chan Path, output chan Path, log *logan.Entry, name string, google *google.Google) Handler {
	return Handler{
		data:         make([]Path, 0),
		chInput:      input,
		chOutput:     output,
		log:          log,
		name:         name,
		googleClient: google,
	}
}

func (h *Handler) SetData(path Path) {
	h.data = append(h.data, path)
}
func (h *Handler) StartRunner(wg *sync.WaitGroup) {
	defer wg.Done()
	go func() {
		//defer h.cancel()
		h.data = append(h.data, <-h.chInput)
		for {
			h.data = append(h.data, <-h.chInput)
			if len(h.data) == 0 || h.data[0].Path == "" {
				return
			}

			running.UntilSuccess(context.Background(), h.log, h.name, func(ctx context.Context) (bool, error) {

				h.log.Debug(h.name, "---> processing: len of data", len(h.data))
				link, success, err := h.googleClient.Update(h.data[0].Path)
				if err != nil {
					h.log.Debug(h.name, "---> ", "error: ", err)
					//success = false
					return false, err

				}
				path := h.data[0]
				path.Path = link

				h.chOutput <- path
				h.pop()
				h.log.Debug(h.name, "---> delete el of data", len(h.data))
				return success, err

			}, time.Millisecond*150, time.Millisecond*180) //todo move config file
			if len(h.data) == 0 {
				h.log.Debug(h.name, "---> skip: len of data = 0")
				break
			}
		}
	}()

	//h.result = h.check(h.chLink, h.cancel, h.ctx, h.name)
}

func (h *Handler) pop() {
	h.data = h.data[1:]
}
