package cron

import (
	"log"

	"github.com/team-tritan/ping-pong/config"
	"github.com/team-tritan/ping-pong/internal/status"
)

type JobRunnerFunc = func()

func runJob(service *config.Service) {
	switch service.ServiceType {
	case "http":
		httpStatusChecker := status.NewHTTPStatus()

		status, err := httpStatusChecker.CheckStatus("GET", service.Url, service.AllowRetries)
		if err != nil {
			panic(err) // Ok to panic here, we catch it.
		}

		log.Printf("%v", status.Statistics)
	}
}

func makeJobRunner(services []*config.Service) JobRunnerFunc {
	return func() {
		for _, service := range services {
			go runJob(service)
		}
	}
}
