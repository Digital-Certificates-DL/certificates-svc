package service

import (
	"github.com/go-chi/chi"
	"github.com/go-chi/cors"
	"gitlab.com/distributed_lab/ape"
	"gitlab.com/tokend/course-certificates/ccp/internal/config"
	"gitlab.com/tokend/course-certificates/ccp/internal/data/pg"
	"gitlab.com/tokend/course-certificates/ccp/internal/service/handlers"
	"gitlab.com/tokend/course-certificates/ccp/internal/service/helpers"
)

func (s *service) router(cfg config.Config) chi.Router {
	r := chi.NewRouter()

	r.Use(
		cors.Handler(cors.Options{}),
		ape.RecoverMiddleware(s.log),
		ape.LoganMiddleware(s.log),
		ape.CtxMiddleware(
			helpers.CtxClientQ(pg.NewClientQ(s.cfg.DB())),
			helpers.CtxTemplateQ(pg.NewTemplateQ(s.cfg.DB())),
			helpers.CtxLog(s.log),
			helpers.CtxConfig(cfg),
		),
	)

	r.Route("/integrations/ccp/", func(r chi.Router) {
		r.Route("/users", func(r chi.Router) {
			r.Post("/", handlers.GetUsers)
			r.Post("/empty", handlers.GetUsersEmpty)
			r.Post("/settings", handlers.SetSettings)
			r.Post("/token", handlers.UpdateToken)

		})
		r.Route("/certificate", func(r chi.Router) {
			r.Post("/", handlers.PrepareCertificate)
			r.Post("/bitcoin", handlers.UpdateCertificate)
			r.Post("/ipfs", handlers.UploadFileToIpfs)

			r.Route("/template", func(r chi.Router) {
				r.Post("/", handlers.CreateTemplate)
				r.Get("/{user}", handlers.GetTemplates)
				r.Get("/{user}/{name}", handlers.GetTemplateByName)
			})
		})
	})
	return r
}
