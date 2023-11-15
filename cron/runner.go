package cron

import (
	"log"
	"time"

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

		log.Printf("[%s] GET %s Status: UP\n\tTCP Connection: %dms TLS Handshake: %dms DNS Lookup: %dms Time To First Byte: %dms", service.Name, service.Url, status.Statistics.TcpDone.Milliseconds(), time.Duration(status.Statistics.TlsHandshakeDone-status.Statistics.TlsHandshakeStart).Milliseconds(), status.Statistics.DnsDone.Milliseconds(), status.Statistics.GotFirstByte.Milliseconds())
	}
}

func makeJobRunner(services []*config.Service) JobRunnerFunc {
	return func() {
		for _, service := range services {
			go runJob(service)
		}
	}
}
