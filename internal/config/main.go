package config

import (
	"gitlab.com/distributed_lab/kit/comfig"
	"gitlab.com/distributed_lab/kit/copus"
	"gitlab.com/distributed_lab/kit/copus/types"
	"gitlab.com/distributed_lab/kit/kv"
	"gitlab.com/distributed_lab/kit/pgdb"
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
	TitlesConfiger
	NetworksConfiger
	ExamsConfiger
	pgdb.Databaser
}

type config struct {
	TableConfiger
	Signer
	pgdb.Databaser
	TitlesConfiger
	GoogleConfiger
	TemplatesConfiger
	ExamsConfiger
	comfig.Logger
	QRCoder
	types.Copuser
	NetworksConfiger
	getter kv.Getter
	comfig.Listenerer
}

func New(getter kv.Getter) Config {
	return &config{
		getter:            getter,
		ExamsConfiger:     NewExamsConfiger(getter),
		TitlesConfiger:    NewTitlesConfiger(getter),
		Listenerer:        comfig.NewListenerer(getter),
		Copuser:           copus.NewCopuser(getter),
		TemplatesConfiger: NewTemplatesConfiger(getter),
		QRCoder:           NewQRCoder(getter),
		GoogleConfiger:    NewGoogler(getter),
		Signer:            NewKeyer(getter),
		TableConfiger:     NewTableConfiger(getter),
		NetworksConfiger:  NewEthRPCConfiger(getter),
		Databaser:         pgdb.NewDatabaser(getter),
		Logger:            comfig.NewLogger(getter, comfig.LoggerOpts{}),
	}
}
