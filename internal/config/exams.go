package config

import (
	figure "gitlab.com/distributed_lab/figure/v3"
	"gitlab.com/distributed_lab/kit/comfig"
	"gitlab.com/distributed_lab/kit/kv"
	"gitlab.com/distributed_lab/logan/v3/errors"
)

type ExamsConfiger interface {
	ExamsConfig() map[string]string
}

func NewExamsConfiger(getter kv.Getter) ExamsConfiger {
	return &examsConfiger{
		getter: getter,
	}
}

type examsConfiger struct {
	getter kv.Getter
	once   comfig.Once
}

func (c *examsConfiger) ExamsConfig() map[string]string {
	return c.once.Do(func() interface{} {
		raw := kv.MustGetStringMap(c.getter, "exams")
		config := struct {
			List map[string]string `fig:"list"`
		}{}
		err := figure.Out(&config).From(raw).Please()
		if err != nil {
			panic(errors.Wrap(err, "failed to figure out"))
		}
		return config.List
	}).(map[string]string)
}
