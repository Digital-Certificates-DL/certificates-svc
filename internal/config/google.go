package config

import (
	"github.com/pkg/errors"
	"gitlab.com/distributed_lab/figure/v3"
	"gitlab.com/distributed_lab/kit/comfig"
	"gitlab.com/distributed_lab/kit/kv"
)

type Googler interface {
	Google() *Google
}

type Google struct {
	Code       string `fig:"code"`
	ApiKey     string `fig:"api_key"`
	SecretPath string `fig:"secret_path"`
}

func NewGoogler(getter kv.Getter) Googler {
	return &googler{
		getter: getter,
	}
}

type googler struct {
	getter kv.Getter
	once   comfig.Once
}

func (c *googler) Google() *Google {
	return c.once.Do(func() interface{} {
		raw := kv.MustGetStringMap(c.getter, "google")
		config := Google{}
		err := figure.Out(&config).From(raw).Please()
		if err != nil {
			panic(errors.Wrap(err, "failed to figure out"))
		}
		return &config
	}).(*Google)
}
