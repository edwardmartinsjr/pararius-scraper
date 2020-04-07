package main

import (
	"log"
	"os"

	"github.com/edwardmartinsjr/pararius-scraper/config"
	"github.com/edwardmartinsjr/pararius-scraper/worker/scraper"
	utilsLogger "github.com/jexia-com/jexia-go-common/logger/utils"
	"github.com/jexia-com/yggi-go-sdk/logger/masker"
	"github.com/sirupsen/logrus"
)

func main() {
	// get config values
	config, err := config.GetConfig()
	if err != nil {
		errExit(err, nil)
	}

	// instantiate logger for service
	logger, err := utilsLogger.NewLogger(config.Service.Name, &config.Logger)
	if err != nil {
		errExit(err, nil)
	}

	// print the configs on startup if debug config equals true
	masker.DebugfWithMaskedSecrets(logger, "pararius-scraper configuration: %s", config)

	errorChan := make(chan error)
	go func() {
		err = scraper.InitAndStart(logger, &config.Scraper, &config.Authenticate, &config.Project)
		if err != nil {
			errorChan <- err
			return
		}
	}()

	err = <-errorChan
	if err != nil {
		errExit(err, logger)
	}
}

// errExit logs and handles errors in the service startup
func errExit(err error, l logrus.FieldLogger) {
	if l != nil {
		l.Fatal(err)
	} else {
		log.Fatal(err)
	}

	os.Exit(1)
}
