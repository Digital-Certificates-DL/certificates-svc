package config

import (
	"github.com/pkg/errors"
	"gitlab.com/distributed_lab/figure/v3"
	"gitlab.com/distributed_lab/kit/comfig"
	"gitlab.com/distributed_lab/kit/kv"
)

type QRCoder interface {
	QRCode() *QRCode
}

type QRCode struct {
	QRPath   string `fig:"qr_path"`
	Template string `fig:"template"`
}

func NewQRCoder(getter kv.Getter) QRCoder {
	return &qrcoder{
		getter: getter,
	}
}

type qrcoder struct {
	getter kv.Getter
	once   comfig.Once
}

func (c *qrcoder) QRCode() *QRCode {
	return c.once.Do(func() interface{} {
		raw := kv.MustGetStringMap(c.getter, "qr_code")
		config := QRCode{}
		err := figure.Out(&config).From(raw).Please()
		if err != nil {
			panic(errors.Wrap(err, "failed to figure out"))
		}
		return &config
	}).(*QRCode)
}
