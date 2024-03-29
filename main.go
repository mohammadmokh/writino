package main

import (
	"log"

	"github.com/mohammadmokh/writino/app"
	"github.com/mohammadmokh/writino/config"
	v1 "github.com/mohammadmokh/writino/delivery/http/v1"
)

func main() {

	cfg, err := config.LoadCfg("config/config.yml")
	if err != nil {
		log.Fatal(err)
	}

	app, err := app.New(cfg)
	if err != nil {
		log.Fatal(err)
	}

	log.Fatal(v1.New(app, cfg).Run())
}
