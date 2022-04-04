package main

import (
	"github.com/keystop/YaPracticum.git/internal/repository"
	"github.com/keystop/YaPracticum.git/internal/server"
)

func main() {
	s := new(server.Server)
	urlRepo := make(repository.URLRepo)
	s.Start("localhost:8080", &urlRepo)
}
