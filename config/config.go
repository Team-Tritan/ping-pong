package config

import (
	"errors"
	"os"

	"gopkg.in/yaml.v3"
)

var (
	ErrInvalidServiceType = errors.New("invalid service type")
	ErrNoUrlSpecified     = errors.New("no URL specified")
	ErrNoHostSpecified    = errors.New("no host specified")
)

type ServiceType string

const (
	ServiceType_HTTP ServiceType = "http"
	ServiceType_ICMP ServiceType = "icmp"
)

type Service struct {
	ServiceType   ServiceType `yaml:"type"`
	Host          string      `yaml:"host"`
	Port          int         `yaml:"port"`
	Url           string      `yaml:"url"`
	AllowRetries  bool        `yaml:"retries"`
	MaxRetries    int         `yaml:"max_retries"`
	RetryInterval int         `yaml:"retry_interval"`
	CronSchedule  string      `yaml:"cron_schedule"`
	Every         string      `yaml:"every"`
	Name          string      `yaml:"name"`
}

type HttpServerConfig struct {
	Host string `yaml:"host" default:"127.0.0.1"`
	Port int    `yaml:"port" default:"8000"`
}

type Config struct {
	Services      map[string]*Service `yaml:"services"`
	AllowRetries  bool                `yaml:"retries"`
	MaxRetries    int                 `yaml:"max_retries" default:"10"`
	RetryInterval bool                `yaml:"retry_interval" default:"500"`
	LogLevel      string              `yaml:"log_level" default:"info"`
	Http          HttpServerConfig    `yaml:"http"`
	CronSchedule  string              `yaml:"cron_schedule" default:"*/5 * * * *"`
}

func (c *Config) ValidateConfig() error {
	for name, service := range c.Services {
		service.Name = name

		if service.CronSchedule == "" {
			service.CronSchedule = c.CronSchedule
		}

		switch service.ServiceType {
		case ServiceType_HTTP:
			if service.Url == "" {
				return ErrNoUrlSpecified
			}
			break
		case ServiceType_ICMP:
			if service.Host == "" {
				return ErrNoHostSpecified
			}
			break
		default:
			return ErrInvalidServiceType
		}

	}

	return nil
}

func ReadConfig(file string) (*Config, error) {
	data, err := os.ReadFile(file)
	if err != nil {
		return nil, err
	}

	cfg := &Config{}
	if err = yaml.Unmarshal(data, cfg); err != nil {
		return nil, err
	}

	cfg.ValidateConfig()

	return cfg, nil
}
