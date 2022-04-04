package server

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/keystop/YaPracticum.git/internal/handlers"
	"github.com/keystop/YaPracticum.git/internal/middlewares"
	"github.com/keystop/YaPracticum.git/internal/repository"
)

type Server struct {
	http.Server
}

//Start server with router.
func (s *Server) Start(addr string, repo repository.Repository) {
	r := chi.NewRouter()
	r.Post("/", handlers.HandlerURLPost(repo))
	r.Route("/{id}", func(r chi.Router) {
		r.Use(middlewares.URLCtx)
		r.Get("/", handlers.HandlerURLGet(repo))
	})
	s.Addr = addr
	s.Handler = r
	s.ListenAndServe()
}
