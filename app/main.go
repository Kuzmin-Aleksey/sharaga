package main

import (
	"app/config"
	"app/service"
	"app/storage"
	"app/ui"
	"log"
)

const TokenFile = "token.txt"
const ConfigFile = "config/config.yaml"

func main() {
	cfg, err := config.ReadConfig(ConfigFile)
	if err != nil {
		log.Fatal(err)
	}

	tokenStorage := &storage.TokenStorage{
		Path: TokenFile,
	}

	s, err := service.NewService(cfg.Service, tokenStorage)
	if err != nil {
		log.Fatal(err)
	}

	a := ui.NewApplication(s)
	a.Run()
}
