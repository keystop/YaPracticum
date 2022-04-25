package models

type URLs struct {
	ShortURL    string `json:"short_url"`
	OriginalURL string `json:"original_url"`
}

//Repository interface repo urls.
type Repository interface {
	GetURL(string) (string, error)
	SaveURL(string, string, string) (string, error)
	SaveURLs(map[string]string, string, string) (map[string]string, error)
	FindUser(string) bool
	CreateUser() (string, error)
	GetUserURLs(string) ([]URLs, error)
	CheckDBConnection() error
}
