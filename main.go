package main

import (
	"fmt"
	"log"

	"github.com/team-tritan/ping-pong/api"
	"github.com/team-tritan/ping-pong/config"
	"github.com/team-tritan/ping-pong/cron"
)

func main() {
	cfg, err := config.ReadConfig("./testconfig.yml")
	if err != nil {
		panic(err)
	}

	log.Println("Loaded config")

	scheduler, err := cron.NewCronScheduler(cfg)
	if err != nil {
		panic(err)
	}
	scheduler.Run()

	err = api.NewServer(fmt.Sprintf("%s:%d", cfg.Http.Host, cfg.Http.Port))
	if err != nil {
		panic(err)
	}
}
