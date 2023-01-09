package config

import (
	"gitlab.com/distributed_lab/kit/comfig"
	"gitlab.com/distributed_lab/kit/kv"
)

type Config interface {
	TableConfiger
	QRCoder
	Signer
	GoogleConfiger
	comfig.Logger
	TemplatesConfiger
}

type config struct {
	TableConfiger
	Signer
	GoogleConfiger
	TemplatesConfiger
	comfig.Logger
	QRCoder
	getter kv.Getter
}

func New(getter kv.Getter) Config {
	return &config{
		getter:            getter,
		TemplatesConfiger: NewTemplatesConfiger(getter),
		QRCoder:           NewQRCoder(getter),
		GoogleConfiger:    NewGoogler(getter),
		Signer:            NewKeyer(getter),
		TableConfiger:     NewTableConfiger(getter),
		Logger:            comfig.NewLogger(getter, comfig.LoggerOpts{}),
	}
}
