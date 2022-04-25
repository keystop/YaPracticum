package server

import (
	"net/http"

	"github.com/keystop/YaPracticum.git/internal/handlers"
	"github.com/keystop/YaPracticum.git/internal/middlewares"
	"github.com/keystop/YaPracticum.git/internal/models"
	"github.com/go-chi/chi/v5"
)

type Server struct {
	http.Server
}

//Start server with router.
func (s *Server) Start(repo models.Repository, opt models.Options) {
	r := chi.NewRouter()
	handlers.NewHandlers(repo, opt)
	middlewares.NewCookie(repo)
	r.Use(middlewares.SetCookieUser, middlewares.ZipHandlerRead, middlewares.ZipHandlerWrite)
	//r.Use(middlewares.ZipHandlerRead, middlewares.ZipHandlerWrite)

	r.Get("/user/urls", handlers.HandlerUserPostURLs)
	r.Get("/ping", handlers.HandlerCheckDBConnect)
	r.Route("/{id}", func(r chi.Router) {
		r.Use(middlewares.URLCtx)
		r.Get("/", handlers.HandlerURLGet)
	})
	r.Post("/", handlers.HandlerURLPost)
	r.Post("/api/shorten", handlers.HandlerAPIURLPost)
	r.Post("/api/shorten/batch", handlers.HandlerAPIURLsPost)

	s.Addr = opt.ServAddr()
	s.Handler = r
	s.ListenAndServe()
}
