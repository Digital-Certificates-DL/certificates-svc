package service

import (
	"context"
	"fmt"
	"gitlab.com/distributed_lab/logan/v3"
	"gitlab.com/distributed_lab/running"
	"helper/internal/service/google"
	"sync"
	"time"
)

type Paths struct {
	Path string
	ID   int
}

type Handler struct {
	data         []Paths
	result       []Paths
	running      int
	ctx          context.Context
	cancel       context.CancelFunc
	chLink       chan Paths
	minTime      time.Duration
	maxTime      time.Duration
	name         string
	log          *logan.Entry
	googleClient *google.Google
}

func NewHandler(ctx context.Context, ch chan Paths, log *logan.Entry, name string, google *google.Google) Handler {
	ctxChild, cancel := context.WithCancel(ctx)

	return Handler{
		data:         make([]Paths, 0),
		chLink:       ch,
		ctx:          ctxChild,
		log:          log,
		name:         name,
		googleClient: google,
		cancel:       cancel,
	}
}

func (h *Handler) SetData(path Paths) {
	h.data = append(h.data, path)
}
func (h *Handler) StartRunner(wg *sync.WaitGroup) {

	defer wg.Done()
	go func() {
		//defer h.cancel()
		for {

			running.UntilSuccess(h.ctx, h.log, h.name, func(ctx context.Context) (bool, error) {
				h.log.Debug(h.name, "---> processing: len of data", len(h.data))
				link, success, err := h.googleClient.Update(h.data[len(h.data)-1].Path)
				if err != nil {
					h.log.Debug(h.name, "---> ", "error: ", err)
					success = false
				}
				if len(h.data) == 0 {
					h.log.Debug(h.name, " ---> skip: len of data = 0")
					return true, nil
				}
				path := h.data[len(h.data)-1]
				path.Path = link
				if success {
					h.chLink <- path
					h.pop()
					h.log.Debug(h.name, "---> delete el of data", len(h.data))
					return success, err
				}
				return false, err

			}, time.Millisecond*150, time.Millisecond*180) //todo move config file
			if len(h.data) == 0 {
				break
			}
		}
	}()

	h.result = h.check(h.chLink, h.cancel, h.ctx, h.name)
}

func (h *Handler) check(ch <-chan Paths, cancel context.CancelFunc, ctx context.Context, name string) []Paths {
	paths := make([]Paths, 0)
	defer close(h.chLink)
	for {
		select {
		case <-ch:
			paths = append(paths, <-ch)
			if len(h.data) == 0 {
				fmt.Println("process   ", name)
				cancel()
			}
		case <-ctx.Done():
			fmt.Println("end   ", name) //todo will use logs
			return paths
		}
	}

}

func (h *Handler) pop() {
	h.data = h.data[1:]
}

func (h *Handler) Result() ([]Paths, bool) {
	if h.result != nil {
		buf := h.result
		h.result = nil
		return buf, true
	}

	return nil, false

}
