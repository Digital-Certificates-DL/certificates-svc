package config

import (
	figure "gitlab.com/distributed_lab/figure/v3"
	"gitlab.com/distributed_lab/kit/comfig"
	"gitlab.com/distributed_lab/kit/kv"
	"gitlab.com/distributed_lab/logan/v3/errors"
)

type TemplatesConfiger interface {
	TemplatesConfig() *RawFiltersConfig
}

type RawFiltersConfig struct {
	Filters []rawFilterConfig `fig:"list"`
}

type rawFilterConfig struct {
	Params map[string][]string `fig:"params"`
}

func NewTemplatesConfiger(getter kv.Getter) TemplatesConfiger {
	return &templatesConfig{
		getter: getter,
	}
}

type templatesConfig struct {
	getter kv.Getter
	once   comfig.Once
}

func (c *templatesConfig) TemplatesConfig() *RawFiltersConfig {
	return c.once.Do(func() interface{} {
		raw := kv.MustGetStringMap(c.getter, "templates")
		config := RawFiltersConfig{}
		err := figure.Out(&config).From(raw).Please()
		if err != nil {
			panic(errors.Wrap(err, "failed to figure out"))
		}
		return &config
	}).(*RawFiltersConfig)
}
