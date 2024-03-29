package pdf

import (
	"context"
	"gitlab.com/distributed_lab/logan/v3"
	"gitlab.com/tokend/course-certificates/ccp/internal/config"
	"gitlab.com/tokend/course-certificates/ccp/internal/data"
	"gitlab.com/tokend/course-certificates/ccp/internal/service/core/google"
	"gitlab.com/tokend/course-certificates/ccp/internal/service/core/helpers"
)

const Update = "Update"
const Generate = "generates"

type CreatorPDF interface {
	NewContainer(users []*helpers.Certificate, google *google.Google, address, sheetUrl string, owner *data.Client, masterQ data.MasterQ, process string) int
	CheckContainerState(containerID int) *Container
	Run(ctx context.Context)
	removeIndex(index int)
}

type CreatorPDFType struct {
	handlerChan     chan *Container
	lastContainerID int
	log             *logan.Entry
	config          config.Config
	readyContainers []*Container
}

func NewPdfCreator(log *logan.Entry, config config.Config) *CreatorPDFType {
	return &CreatorPDFType{
		handlerChan:     make(chan *Container),
		log:             log,
		config:          config,
		readyContainers: make([]*Container, 0),
	}
}

func (p *CreatorPDFType) NewContainer(certificates []*helpers.Certificate, google *google.Google, address, sheetUrl string, owner *data.Client, masterQ data.MasterQ, process string) int {
	p.lastContainerID++
	p.handlerChan <- &Container{
		Certificates: certificates,
		ID:           p.lastContainerID,
		Status:       isProcessingStatus,
		log:          p.log,
		config:       p.config,
		masterQ:      masterQ,
		googleClient: google,
		address:      address,
		sheetUrl:     sheetUrl,
		owner:        owner,
		process:      process,
	}
	return p.lastContainerID
}

func (p *CreatorPDFType) Run(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			return

		case container, ok := <-p.handlerChan:
			if !ok {
				ctx.Done()
				return
			}

			switch container.process {
			case Generate:
				err := container.Generate()
				if err != nil {
					p.log.Error(err, "failed to run container")
				}
				container.Status = readyStatus
				p.log.Debug("Success: ", container)

				p.readyContainers = append(p.readyContainers, container)
				break
			case Update:
				//todo make better
				err := container.Update()
				if err != nil {
					p.log.Error(err, "failed to run container")
				}
				container.Status = readyStatus
				p.readyContainers = append(p.readyContainers, container)
			}
			break
		}
	}
}

func (p *CreatorPDFType) CheckContainerState(containerID int) *Container {
	p.log.Debug("readyContainers: ", p.readyContainers)

	for _, container := range p.readyContainers {
		if container.ID == containerID {
			p.log.Debug("container: ", container)
			return container
		}
	}

	return nil
}

func (p *CreatorPDFType) removeIndex(index int) {
	ret := make([]*Container, 0)
	ret = append(ret, p.readyContainers[:index]...)
	p.readyContainers = append(ret, p.readyContainers[index+1:]...)
}
