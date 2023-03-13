package service

import (
	"github.com/go-chi/chi"
	"gitlab.com/distributed_lab/ape"
	"helper/internal/config"
	"helper/internal/service/handlers"
	"helper/internal/service/helpers"
)

func (s *service) router(cfg config.Config) chi.Router {
	r := chi.NewRouter()

	r.Use(
		ape.RecoverMiddleware(s.log),
		ape.LoganMiddleware(s.log),
		ape.CtxMiddleware(
			helpers.CtxLog(s.log),
			helpers.CtxConfig(cfg),
		),
	)
	r.Route("/integrations/css", func(r chi.Router) {
		r.Post("/generate", handlers.GenerateTable)
		r.Get("/template", handlers.CreateTemplate)

	})
	return r
}
