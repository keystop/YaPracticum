package main

import (
	"fmt"

	"github.com/keystop/YaPracticum.git/internal/defoptions"
	"github.com/keystop/YaPracticum.git/internal/repository"
	"github.com/keystop/YaPracticum.git/internal/server"
)

// Main.
func main() {

	opt := defoptions.NewDefOptions()
	sr, err := repository.NewServerRepo(opt.DBConnString())
	if err != nil {
		fmt.Println("Ошибка при подключении к БД: ", err)
		return
	}
	// serverRepo := repository.NewRepo(opt.RepoFileName())
	s := new(server.Server)
	s.Start(sr, opt)
	defer sr.Close()
}
