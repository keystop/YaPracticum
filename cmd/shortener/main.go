package main

import (
	"github.com/keystop/YaPracticum.git/internal/defoptions"
	"github.com/keystop/YaPracticum.git/internal/repository"
	"github.com/keystop/YaPracticum.git/internal/serialize"
	"github.com/keystop/YaPracticum.git/internal/server"

)

func main() {
	opt := defoptions.NewDefOptions()

	urlRepo := repository.NewURLRepo()

	serialize.NewSerialize(opt.RepoFileName())
	serialize.ReadURLSFromFile(urlRepo)
	repository.SerializeURLRepo = serialize.SaveURLSToFile

	s := new(server.Server)
	s.Start(urlRepo, opt)
}
