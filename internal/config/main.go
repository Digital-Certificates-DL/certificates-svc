package config

import (
	"gitlab.com/distributed_lab/kit/kv"
)

type Config interface {
	Tabler
	Googler
	Keyer
}

type config struct {
	Tabler
	Keyer
	Googler
	getter kv.Getter
}

func New(getter kv.Getter) Config {
	return &config{
		getter:  getter,
		Googler: NewGoogler(getter),
		Keyer:   NewKeyer(getter),
		Tabler:  NewTabler(getter),
	}
}
