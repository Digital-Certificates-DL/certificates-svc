package config

import (
	"gitlab.com/distributed_lab/figure"
	"gitlab.com/distributed_lab/kit/comfig"
	"gitlab.com/distributed_lab/kit/kv"
	"gitlab.com/distributed_lab/logan/v3/errors"
)

type TableConfiger interface {
	Table() *Table
}

type Table struct {
	Input  string `fig:"input" `
	Result string `fig:"result" `
}

func NewTableConfiger(getter kv.Getter) TableConfiger {
	return &tableConfiger{
		getter: getter,
	}
}

type tableConfiger struct {
	getter kv.Getter
	once   comfig.Once
}

func (c *tableConfiger) Table() *Table {
	return c.once.Do(func() interface{} {
		raw := kv.MustGetStringMap(c.getter, "tables")
		config := Table{}
		err := figure.Out(&config).From(raw).Please()
		if err != nil {
			panic(errors.Wrap(err, "failed to figure out"))
		}
		return &config
	}).(*Table)
}
