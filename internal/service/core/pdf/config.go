package pdf

import "gitlab.com/tokend/course-certificates/ccp/internal/config"

type PDFConfig struct {
	titles    map[string]string
	exams     map[string]string
	templates map[string]string
}

func NewPDFConfig(cfg config.Config) *PDFConfig {
	return &PDFConfig{
		titles:    cfg.TitlesConfig(),
		exams:     cfg.ExamsConfig(),
		templates: cfg.TemplatesConfig(),
	}
}
