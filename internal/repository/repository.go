package repository

import (
	"errors"

	"github.com/keystop/YaPracticum.git/internal/shorter"
)

type Key string

//Repository interface repo urls.
type Repository interface {
	GetURL(string) (string, error)
	SaveURL([]byte) string
}

//UrlsData repository of urls. Realize Repository interface.
type URLRepo map[string]string

func (u *URLRepo) SaveURL(url []byte) string {
	r := shorter.MakeShortner(url)
	(*u)[r] = string(url)
	return r
}

func (u *URLRepo) GetURL(id string) (string, error) {
	if r, ok := (*u)[id]; ok {
		return string(r), nil
	}
	return "", errors.New("not found")
}
