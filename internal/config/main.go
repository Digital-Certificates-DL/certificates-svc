package config

import (
	"gitlab.com/distributed_lab/kit/kv"
)

type Config interface {
	Tabler
	QRCoder
	Keyer
	Googler
}

type config struct {
	Tabler
	Keyer
	Googler
	QRCoder
	getter kv.Getter
}

func New(getter kv.Getter) Config {
	return &config{
		getter:  getter,
		QRCoder: NewQRCoder(getter),
		Googler: NewGoogler(getter),
		Keyer:   NewKeyer(getter),
		Tabler:  NewTabler(getter),
	}
}
