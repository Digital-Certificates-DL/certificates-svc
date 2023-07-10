package config

import (
	"gitlab.com/distributed_lab/figure/v3"
	"gitlab.com/distributed_lab/kit/comfig"
	"gitlab.com/distributed_lab/kit/kv"
	"gitlab.com/distributed_lab/logan/v3/errors"
)

type SbtConfiger interface {
	SbtConfig() *SbtConfig
}

type SbtConfig struct {
	ExternalURL string `fig:"external_url" `
}

func NewSbtConfiger(getter kv.Getter) SbtConfiger {
	return &sbtConfiger{
		getter: getter,
	}
}

type sbtConfiger struct {
	getter kv.Getter
	once   comfig.Once
}

func (c *sbtConfiger) SbtConfig() *SbtConfig {
	return c.once.Do(func() interface{} {
		raw := kv.MustGetStringMap(c.getter, "sbt")
		config := SbtConfig{}
		err := figure.Out(&config).From(raw).Please()
		if err != nil {
			panic(errors.Wrap(err, "failed to figure out"))
		}
		return &config
	}).(*SbtConfig)
}
