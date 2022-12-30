package config

import (
	"gitlab.com/distributed_lab/figure"
	"gitlab.com/distributed_lab/kit/comfig"
	"gitlab.com/distributed_lab/kit/kv"
	"gitlab.com/distributed_lab/logan/v3/errors"
)

type Tabler interface {
	Table() *Table
}

type Table struct {
	Input  string `fig:"input" `
	Result string `fig:"result" `
}

func NewTabler(getter kv.Getter) Tabler {
	return &tabler{
		getter: getter,
	}
}

type tabler struct {
	getter kv.Getter
	once   comfig.Once
}

func (c *tabler) Table() *Table {
	return c.once.Do(func() interface{} {
		raw := kv.MustGetStringMap(c.getter, "aws")
		config := Table{}
		err := figure.Out(&config).From(raw).Please()
		if err != nil {
			panic(errors.Wrap(err, "failed to figure out"))
		}
		return &config
	}).(*Table)
}
