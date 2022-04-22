package global

//Options interface for program options.
type Options interface {
	ServAddr() string
	RespBaseURL() string
	RepoFileName() string
}

//Repository interface repo urls.
type Repository interface {
	GetURL(string) (string, error)
	SaveURL([]byte) string
	Get() map[string]string
	ToSet() *map[string]string
}
