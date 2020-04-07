package config

import (
	"strconv"
	"time"

	"github.com/jexia-com/jexia-config/reader/file"
	"github.com/jexia-com/jexia-config/toml"
	"github.com/jexia-com/jexia-go-common/logger/utils"
)

// Config contains the complete pararius-scraper configuration
type Config struct {
	Service      Service
	Logger       utils.LoggerConfig
	Scraper      Scraper
	Authenticate Authenticate
	Project      Project
}

// Service contains the service configuration
type Service struct {
	Version string `toml:"version"`
	Name    string `toml:"name" config:"SERVICE_NAME, pararius-scraper"`
}

// Scraper contains the scraper configuration
type Scraper struct {
	// ScraperWorkerInterval describes timeout interval
	// for the worker which will scrap site
	ScraperWorkerInterval Timeout `toml:"scraperworkerinterval" config:"SCRAPER_WORKER_INTERVAL, 1h"`
	// AllowedDomains sets the domain whitelist used by the Collector.
	AllowedDomains string `toml:"alloweddomains" config:"ALLOWED_DOMAINS, www.pararius.com"`
	// URLToVisit sets the URl to be visited by scraper
	URLToVisit string `toml:"urltovisit" config:"URL_TO_VISIT, https://www.pararius.com/apartments/amsterdam"`
}

// Authenticate contains the user authentication configuration
type Authenticate struct {
	Email    string `toml:"email" config:"EMAIL"`
	Password string `toml:"password" config:"PASSWORD"`
}

// Project contains the project configuration
type Project struct {
	ProjectID       string `config:"PROJECT_ID"`
	TrainingDataSet string `config:"TRAINING_DATA_SET"`
}

// Timeout is a string representing the maximum duration of a flow/call/TTL
// If the values is an integer it is processed as seconds, otherwise it is parsed as a duration: https://golang.org/pkg/time/#ParseDuration
type Timeout string

// ParseDuration converts the Timeout value into a Duration object
func (t Timeout) ParseDuration() (time.Duration, error) {
	if t == "" {
		return time.Duration(0), nil
	}
	if seconds, err := strconv.Atoi(string(t)); err == nil {
		return time.Duration(seconds) * time.Second, nil
	}

	duration, err := time.ParseDuration(string(t))
	if err != nil {
		return 0, ErrParsingTimeout{Value: string(t), Err: err}
	}

	return duration, nil
}

// GetConfig - get configuration values
func GetConfig() (*Config, error) {
	// get configuration file
	configPath, err := file.JexiaConfigLocation("pararius-scraper")
	if err != nil && err != file.ErrNoGoPath {
		return nil, err
	}

	configReader := file.Reader{
		ReaderDecoder: toml.NewReaderDecoder("", true),
		AppLocation:   configPath + "pararius-scraper.toml",
		EnvVar:        "CONFIG",
	}

	// get config values
	config := Config{}
	err = configReader.ReadConfig(&config)
	if err != nil {
		return nil, err
	}

	return &config, nil
}
