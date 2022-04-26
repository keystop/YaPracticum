package defoptions

import (
	"flag"
	"fmt"
	"os"

	"github.com/caarlos0/env/v6"
	"github.com/keystop/YaPracticum.git/internal/models"
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

type Config struct {
	ServAddr     string `env:"SERVER_ADDRESS"`
	BaseURL      string `env:"BASE_URL"`
	RepoFileName string `env:"FILE_STORAGE_PATH"`
	DBConnString string `env:"DATABASE_DSN"`
}

func (d *defOptions) fillFromConf(c *Config) {
	if len(c.ServAddr) != 0 && len(d.servAddr) == 0 {
		d.servAddr = c.ServAddr
	}
	if len(c.BaseURL) != 0 && len(d.baseURL) == 0 {
		d.baseURL = c.BaseURL
	}
	if len(c.RepoFileName) != 0 && len(d.repoFileName) == 0 {
		d.repoFileName = c.RepoFileName
	}
	if len(c.DBConnString) != 0 && len(d.dbConnString) == 0 {
		d.dbConnString = c.DBConnString
	}
}

func (d *defOptions) setDefault(appDir string) {

	config := &Config{
		"localhost:8080",
		"http://localhost:8080",
		appDir + `/local.gob`,
		"user=kseikseich password=112233 dbname=yap sslmode=disable",
	}
	d.fillFromConf(config)
}

//checkEnv for get options from env to default application options.
func (d *defOptions) checkEnv() {

	e := &Config{}
	err := env.Parse(e)
	if err != nil {
		fmt.Println(err.Error())
	}
	d.fillFromConf(e)

}

//setFlags for get options from console to default application options.
func (d *defOptions) setFlags() {

	flag.StringVar(&d.servAddr, "a", d.servAddr, "a server address string")
	flag.StringVar(&d.baseURL, "b", d.baseURL, "a response address string")
	flag.StringVar(&d.repoFileName, "f", d.repoFileName, "a file storage path string")
	flag.StringVar(&d.dbConnString, "d", d.dbConnString, "a db connection string")

	flag.Parse()

}

// NewDefOptions return obj like Options interfasc.
func NewDefOptions() models.Options {
	appDir, err := os.Getwd()
	if err != nil {
		fmt.Println("Не удалось найти каталог программы!")
	}

	opt := &defOptions{}

	opt.setFlags()
	opt.checkEnv()

	opt.setDefault(appDir)

	return opt
}
