package config

import (
	figure "gitlab.com/distributed_lab/figure/v3"
	"gitlab.com/distributed_lab/kit/comfig"
	"gitlab.com/distributed_lab/kit/kv"
	"gitlab.com/distributed_lab/logan/v3/errors"
)

type TemplatesConfiger interface {
	TemplatesConfig() map[string]string
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

func (c *templatesConfig) TemplatesConfig() map[string]string {
	return c.once.Do(func() interface{} {
		raw := kv.MustGetStringMap(c.getter, "templates")
		config := struct {
			List map[string]string `fig:"list"`
		}{}
		err := figure.Out(&config).From(raw).Please()
		if err != nil {
			panic(errors.Wrap(err, "failed to figure out"))
		}
		return prepareMap(config.List)
	}).(map[string]string)
}

func prepareMap(in map[string]string) map[string]string {

	var templates = make(map[string]string)
	for key, temp := range in {
		templates[temp] = key
	}
	return templates
}
