package service

import (
	"gitlab.com/distributed_lab/kit/copus/types"
	"gitlab.com/distributed_lab/logan/v3"
	"helper/internal/config"
)

type service struct {
	cfg   config.Config
	log   *logan.Entry
	copus types.Copus
}

func (s *service) run(cfg config.Config) error {
	err := Start(cfg)
	return err
}

func newService(cfg config.Config) *service {
	return &service{
		cfg:   cfg,
		log:   cfg.Log(),
		copus: cfg.Copus(),
	}
}

func Run(cfg config.Config) {
	if err := newService(cfg).run(cfg); err != nil {
		panic(err)
	}

}
