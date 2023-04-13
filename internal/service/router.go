package service

import (
	"github.com/go-chi/chi"
	"github.com/go-chi/cors"
	"gitlab.com/distributed_lab/ape"
	"gitlab.com/tokend/course-certificates/ccp/internal/config"
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
			helpers.CtxLog(s.log),
			helpers.CtxConfig(cfg),
		),
	)

	r.Route("/integrations/ccp", func(r chi.Router) {
		r.Post("/", handlers.GetUsers)
		r.Post("/empty", handlers.GetUsersEmpty)
		r.Post("/certificate", handlers.PrepareCertificate)
		r.Get("/template/", handlers.CreateTemplate)
		r.Post("/ipfs/", handlers.UploadFileToIpfs)
		r.Get("/test", handlers.Test)

	})
	return r
}
