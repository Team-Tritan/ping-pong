package cron

import (
	"log"
	"time"

	"github.com/go-co-op/gocron"
	"github.com/team-tritan/ping-pong/config"
)

type CronScheduler struct {
	s *gocron.Scheduler
}

func NewCronScheduler(cfg *config.Config) (*CronScheduler, error) {
	s := &CronScheduler{
		s: gocron.NewScheduler(time.UTC),
	}

	if err := s.scheduleJobs(cfg); err != nil {
		return nil, err
	}

	return s, nil
}

func (cs *CronScheduler) Run() {
	cs.s.StartAsync()
}

// Groups jobs together and schedules them
func (cs *CronScheduler) scheduleJobs(cfg *config.Config) error {
	jobs := map[string][]*config.Service{}

	for name, service := range cfg.Services {
		log.Printf("Loaded service %s", name)
		jobs[service.CronSchedule] = append(jobs[service.CronSchedule], service)
	}

	for schedule, group := range jobs {
		_, err := cs.s.Cron(schedule).Do(makeJobRunner(group))
		if err != nil {
			return err
		}
	}

	return nil
}
