package api

import (
	"context"
	"gitlab.com/distributed_lab/kit/copus/types"
	"gitlab.com/distributed_lab/logan/v3"
	"gitlab.com/distributed_lab/logan/v3/errors"
	"gitlab.com/tokend/course-certificates/ccp/internal/config"
	"gitlab.com/tokend/course-certificates/ccp/internal/service/core/pdf"
	"net"
	"net/http"
)

type service struct {
	cfg        config.Config
	log        *logan.Entry
	listener   net.Listener
	copus      types.Copus
	pdfCreator *pdf.CreatorPDFType
}

func (s *service) run(cfg config.Config) error {
	r := s.router(cfg)
	if err := s.copus.RegisterChi(r); err != nil {
		return errors.Wrap(err, "cop failed")
	}
	return http.Serve(s.listener, r)

}

func newService(cfg config.Config) *service {

	pdfCreator := pdf.NewPdfCreator(cfg.Log())
	go pdfCreator.Run(context.Background())
	return &service{
		cfg:        cfg,
		log:        cfg.Log(),
		copus:      cfg.Copus(),
		listener:   cfg.Listener(),
		pdfCreator: pdfCreator,
	}
}

func Run(cfg config.Config) {
	if err := newService(cfg).run(cfg); err != nil {
		panic(err)
	}

}
