package config

import (
	"gitlab.com/distributed_lab/kit/kv"
)

type Config interface {
	Tabler
	QRCoder
	Signer
	Googler
}

type config struct {
	Tabler
	Signer
	Googler
	QRCoder
	getter kv.Getter
}

func New(getter kv.Getter) Config {
	return &config{
		getter:  getter,
		QRCoder: NewQRCoder(getter),
		Googler: NewGoogler(getter),
		Signer:  NewKeyer(getter),
		Tabler:  NewTabler(getter),
	}
}
