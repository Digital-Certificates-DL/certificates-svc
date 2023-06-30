package pdf

import (
	"context"
	"gitlab.com/distributed_lab/logan/v3"
	"gitlab.com/tokend/course-certificates/ccp/internal/service/helpers"
)

type PDFCreator interface {
	NewContainer(users []*helpers.User)
}

type PDFCreatorType struct {
	handlerChan   chan Container
	lastIndex     int
	log           *logan.Entry
	conainerState []Container
}

func NewPdfCreator(log *logan.Entry) *PDFCreatorType {
	return &PDFCreatorType{
		handlerChan: make(chan Container),
		log:         log,
	}
}

func (p *PDFCreatorType) NewContainer(users []*helpers.User) {
	p.handlerChan <- Container{
		users:  users,
		number: p.lastIndex + 1,
	}
}

func (p *PDFCreatorType) Run(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			return
		case container := <-p.handlerChan:
			err := container.run()
			if err != nil {
				p.log.Error(err, "failed to run container")
			}
			break
		}
	}
}
