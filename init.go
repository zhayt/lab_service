package main

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"
	"gopkg.in/yaml.v2"
	"log"
	"net/http"
	"os"
)

var dsn string

func init() {
	godotenv.Load(".env")
	yamlFile := os.Getenv("CONFIG_FILE_PATH")
	file, err := os.OpenFile(yamlFile, os.O_RDONLY, 0644)
	if err != nil {
		log.Fatal()
	}
	defer file.Close()
	var config struct {
		App struct {
			Host string `yaml:"host"`
			Port string `yaml:"port"`
			Mode string `yaml:"mode"`
		} `yaml:"app"`
		DB struct {
			Driver   string `yaml:"driver"`
			Protocol string `yaml:"protocol"`
			URL      string `yaml:"url"`
			Port     string `yaml:"port"`
			User     string `yaml:"user"`
			Password string `yaml:"password"`
			Database string `yaml:"database"`
		} `yaml:"db"`
	}
	data, err := os.ReadFile(yamlFile)
	if err != nil {
		log.Fatal(err)
	}
	cfg := config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		log.Fatal(err)
	}
	dsn = cfg.DB.User + ":" + cfg.DB.Password + "@" + cfg.DB.Protocol + "(" + cfg.DB.URL + ":" + cfg.DB.Port + ")/" + cfg.DB.Database
	db, err = sql.Open(cfg.DB.Driver, dsn)
	if err != nil {
		log.Fatal(err)
	}
	if cfg.App.Mode == "prod" {
		logger = log.New(os.Stdout, "INFO", log.Ldate|log.Ltime)
	} else {
		logger = log.New(os.Stdout, "INFO DEV", log.Ldate|log.Ltime|log.Lshortfile)
	}
	server = http.Server{
		Addr: cfg.App.Host + ":" + cfg.App.Port,
	}
}
