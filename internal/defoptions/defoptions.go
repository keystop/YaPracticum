package defoptions

import (
	"flag"
	"fmt"
	"os"

	"github.com/keystop/YaPracticum.git/internal/models"
	"github.com/caarlos0/env/v6"
)

type defOptions struct {
	servAddr     string
	baseURL      string
	repoFileName string
	dbConnString string
}

func (d defOptions) ServAddr() string {
	return d.servAddr
}

func (d defOptions) RespBaseURL() string {
	return d.baseURL
}

func (d defOptions) RepoFileName() string {
	return d.repoFileName
}

func (d defOptions) DBConnString() string {
	return d.dbConnString
}

type EnvOptions struct {
	ServAddr     string `env:"SERVER_ADDRESS"`
	BaseURL      string `env:"BASE_URL"`
	RepoFileName string `env:"FILE_STORAGE_PATH"`
	DBConnString string `env:"DATABASE_DSN"`
}

//checkEnv for get options from env to default application options.
func (d *defOptions) checkEnv() {

	e := &EnvOptions{}
	err := env.Parse(e)
	if err != nil {
		fmt.Println(err.Error())
	}
	if len(e.ServAddr) != 0 {
		d.servAddr = e.ServAddr
	}
	if len(e.BaseURL) != 0 {
		d.baseURL = e.BaseURL
	}
	if len(e.RepoFileName) != 0 {
		d.repoFileName = e.RepoFileName
	}
	if len(e.DBConnString) != 0 {
		d.dbConnString = e.DBConnString
	}
}

//setFlags for get options from console to default application options.
func (d *defOptions) setFlags() {
	appDir, err := os.Getwd()
	if err != nil {
		fmt.Println("Не удалось найти каталог программы!")
	}
	flag.StringVar(&d.servAddr, "a", "localhost:8080", "a server address string")
	flag.StringVar(&d.baseURL, "b", "http://localhost:8080", "a response address string")
	flag.StringVar(&d.repoFileName, "f", appDir+`/local.gob`, "a file storage path string")
	flag.StringVar(&d.dbConnString, "d", "user=kseikseich dbname=yap sslmode=disable", "a db connection string")
	flag.Parse()
}

// NewDefOptions return obj like Options interfase.
func NewDefOptions() models.Options {
	opt := new(defOptions)
	opt.setFlags()
	opt.checkEnv()
	return opt
}
