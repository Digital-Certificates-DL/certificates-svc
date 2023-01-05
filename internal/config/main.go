package config

import (
	"gitlab.com/distributed_lab/kit/kv"
)

type Config interface {
	Tabler
	QRCoder
	Signer
	Googler
	TemplatesConfiger
}

type config struct {
	Tabler
	Signer
	Googler
	TemplatesConfiger
	QRCoder
	getter kv.Getter
}

func New(getter kv.Getter) Config {
	return &config{
		getter:            getter,
		TemplatesConfiger: NewTemplatesConfiger(getter),
		QRCoder:           NewQRCoder(getter),
		Googler:           NewGoogler(getter),
		Signer:            NewKeyer(getter),
		Tabler:            NewTabler(getter),
	}
}
