package config

import (
	"github.com/pkg/errors"
	"gitlab.com/distributed_lab/figure/v3"
	"gitlab.com/distributed_lab/kit/comfig"
	"gitlab.com/distributed_lab/kit/kv"
)

type StaticConfiger interface {
	StaticConfig() *StaticConfig
}

type StaticConfig struct {
	Location string `fig:"location"`
}

func NewStaticConfiger(getter kv.Getter) StaticConfiger {
	return &staticConfig{
		getter: getter,
	}
}

type staticConfig struct {
	getter kv.Getter
	once   comfig.Once
}

func (c *staticConfig) StaticConfig() *StaticConfig {
	return c.once.Do(func() interface{} {
		raw := kv.MustGetStringMap(c.getter, "static")
		config := StaticConfig{}
		err := figure.Out(&config).From(raw).Please()
		if err != nil {
			panic(errors.Wrap(err, "failed to figure out"))
		}
		return &config
	}).(*StaticConfig)
}
