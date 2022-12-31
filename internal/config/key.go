package config

import (
	"gitlab.com/distributed_lab/figure/v3"
	"gitlab.com/distributed_lab/kit/comfig"
	"gitlab.com/distributed_lab/kit/kv"
	"gitlab.com/distributed_lab/logan/v3/errors"
)

type Keyer interface {
	Key() *Key
}

type Key struct {
	Private string `fig:"key" `
}

func NewKeyer(getter kv.Getter) Keyer {
	return &keyer{
		getter: getter,
	}
}

type keyer struct {
	getter kv.Getter
	once   comfig.Once
}

func (c *keyer) Key() *Key {
	return c.once.Do(func() interface{} {
		raw := kv.MustGetStringMap(c.getter, "signature")
		config := Key{}
		err := figure.Out(&config).From(raw).Please()
		if err != nil {
			panic(errors.Wrap(err, "failed to figure out"))
		}
		return &config
	}).(*Key)
}
