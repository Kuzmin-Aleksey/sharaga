package main

import (
	"log"
	"sharaga/internal/app"
	"sharaga/internal/config"
)

const configPath = "config/config.go"

func main() {
	cfg, err := config.ReadConfig(configPath, ".env")
	if err != nil {
		log.Fatal(err)
	}

	if err := app.New(cfg).Run(); err != nil {
		log.Fatal(err)
	}
}
