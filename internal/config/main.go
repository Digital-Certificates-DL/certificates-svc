package config

import (
	"gitlab.com/distributed_lab/kit/comfig"
	"gitlab.com/distributed_lab/kit/copus"
	"gitlab.com/distributed_lab/kit/copus/types"
	"gitlab.com/distributed_lab/kit/kv"
)

type Config interface {
	TableConfiger
	QRCoder
	Signer
	GoogleConfiger
	comfig.Logger
	TemplatesConfiger
	comfig.Listenerer
	types.Copuser
}

type config struct {
	TableConfiger
	Signer
	GoogleConfiger
	TemplatesConfiger
	comfig.Logger
	QRCoder
	types.Copuser
	getter kv.Getter
	comfig.Listenerer
}

func New(getter kv.Getter) Config {
	return &config{
		getter:            getter,
		Listenerer:        comfig.NewListenerer(getter),
		Copuser:           copus.NewCopuser(getter),
		TemplatesConfiger: NewTemplatesConfiger(getter),
		QRCoder:           NewQRCoder(getter),
		GoogleConfiger:    NewGoogler(getter),
		Signer:            NewKeyer(getter),
		TableConfiger:     NewTableConfiger(getter),
		Logger:            comfig.NewLogger(getter, comfig.LoggerOpts{}),
	}
}
