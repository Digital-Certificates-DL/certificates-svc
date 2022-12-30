package service

import (
	"helper/internal/config"
)

type service struct {
	cfg config.Config
}

func (s *service) run(cfg config.Config) error {
	err := Start(cfg)
	return err
}

func newService(cfg config.Config) *service {
	return &service{
		cfg: cfg,
	}
}

func Run(cfg config.Config) {
	if err := newService(cfg).run(cfg); err != nil {
		panic(err)
	}
}
