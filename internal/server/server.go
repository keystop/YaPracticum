package server

import (
	"net/http"

	"github.com/keystop/YaPracticum.git/internal/global"
	"github.com/keystop/YaPracticum.git/internal/handlers"
	"github.com/keystop/YaPracticum.git/internal/middlewares"
	"github.com/go-chi/chi/v5"
)

type Server struct {
	http.Server
}

//Start server with router.
func (s *Server) Start(repo global.Repository, opt global.Options) {
	r := chi.NewRouter()
	baseURL := opt.RespBaseURL()
	r.Post("/", handlers.ZipHandlerRead(handlers.ZipHandlerWrite(handlers.HandlerURLPost(repo, baseURL))))
	r.Route("/{id}", func(r chi.Router) {
		r.Use(middlewares.URLCtx)
		r.Get("/", handlers.HandlerURLGet(repo))
	})
	r.Post("/api/shorten", handlers.ZipHandlerRead(handlers.ZipHandlerWrite(handlers.HandlerAPIURLPost(repo, baseURL))))

	s.Addr = opt.ServAddr()
	s.Handler = r
	s.ListenAndServe()
}
