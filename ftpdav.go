package main

import (
	"encoding/json"
	"flag"
	"log"
	"os"

	"github.com/evgeni/ftpdav/driver/webdav"
	"goftp.io/server/v2"
)

var (
	cfgPath string
)

type Configuration struct {
	URL      string
	User     string
	Password string
}

func main() {
	flag.StringVar(&cfgPath, "config", "config.json", "config file path")
	flag.Parse()

	file, _ := os.Open(cfgPath)
	defer file.Close()

	decoder := json.NewDecoder(file)
	configuration := Configuration{}
	err := decoder.Decode(&configuration)
	if err != nil {
		log.Fatal(err)
	}

	driver, err := webdav.NewDriver(configuration.URL, configuration.User, configuration.Password, true)
	if err != nil {
		log.Fatal(err)
	}

	s, err := server.NewServer(&server.Options{
		Driver: driver,
		Auth: &server.SimpleAuth{
			Name:     "admin",
			Password: "admin",
		},
		Perm:      server.NewSimplePerm("root", "root"),
		RateLimit: 1000000, // 1MB/s limit
	})
	if err != nil {
		log.Fatal(err)
	}

	if err := s.ListenAndServe(); err != nil {
		log.Fatal(err)
	}
}
