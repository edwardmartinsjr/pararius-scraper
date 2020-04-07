package scraper

import (
	"fmt"
	"strings"

	"github.com/edwardmartinsjr/pararius-scraper/config"
	"github.com/edwardmartinsjr/pararius-scraper/restapi"
	"github.com/edwardmartinsjr/pararius-scraper/worker"
	"github.com/gocolly/colly"
	"github.com/sirupsen/logrus"
)

type scraperWorker struct {
	scraperConfig      *config.Scraper
	authenticateConfig *config.Authenticate
	projectConfig      *config.Project
	workers            []*worker.RepeatedTask
	logger             logrus.FieldLogger
	restapi            *restapi.ClientResp
}

// InitAndStart - initialize scraper workers
func InitAndStart(logger logrus.FieldLogger, scraperConfig *config.Scraper, authenticateConfig *config.Authenticate, projectConfig *config.Project) error {
	var auth restapi.Auth
	a, err := auth.NewAuth(logger, authenticateConfig, projectConfig)
	restapi, err := a.UserAuthenticate()
	if err != nil {
		logger.Errorln("error on user authenticate up:", err)
	}

	sw := &scraperWorker{
		scraperConfig:      scraperConfig,
		authenticateConfig: authenticateConfig,
		projectConfig:      projectConfig,
		logger:             logger.WithField("component", "scraper-workers"),
		restapi:            restapi,
	}

	sw.workers = append(sw.workers, sw.newScraperWorker())

	err = sw.Start()
	if err != nil {
		logger.Errorln("workers failed to start", err)
	}

	return err
}

// Start workers execution
func (sw *scraperWorker) Start() error {
	sw.logger.Infoln("workers started...")
	for _, worker := range sw.workers {
		worker.Start()
	}
	return nil
}

// Stop workers execution
func (sw *scraperWorker) Stop() error {
	for _, worker := range sw.workers {
		worker.Stop()
	}
	return nil
}
func (sw *scraperWorker) newScraperWorker() *worker.RepeatedTask {
	duration, err := sw.scraperConfig.ScraperWorkerInterval.ParseDuration()

	if err != nil {
		sw.logger.Errorln("error parsing scraper interval:", err)
	}

	return worker.NewRepeatedTask(func() error {
		err := sw.Scrap()
		if err != nil {
			sw.logger.Errorln("error scraping up:", err)
			return err
		}

		return nil
	}, duration, sw.logger)
}

func (sw *scraperWorker) Scrap() error {
	sw.logger.Infoln("scraping...")

	// TODO: handle pagination

	c := colly.NewCollector(
		colly.AllowedDomains(sw.scraperConfig.AllowedDomains),
	)

	c.OnRequest(func(r *colly.Request) {
		sw.logger.Infoln("Visiting", r.URL.String())
	})

	// Run ETL process
	// Extract - is the process of reading data from a source.
	// In this stage, the data is collected, often from multiple and different types of sources.
	c.OnHTML(".property-list-item-container", func(e *colly.HTMLElement) {

		area := e.ChildText(".surface")
		bedroom := e.ChildText(".bedrooms")
		price := e.ChildText(".price")

		// Transform - is the process of converting the extracted data from its previous form into the form it needs to be in so that it can be placed into another database.
		// Transformation occurs by using rules or lookup tables or by combining the data with other data.
		bedroom = strings.Replace(bedroom, "bedrooms", "", -1)
		bedroom = strings.Replace(bedroom, "bedroom", "", -1)
		bedroom = strings.TrimSpace(bedroom)

		area = strings.Replace(area, "m²", "", -1)
		area = strings.TrimSpace(area)

		price = strings.Replace(price, "\n", "", -1)
		price = strings.Replace(price, "/month", "", -1)
		price = strings.Replace(price, "(ex.)", "", -1)
		price = strings.Replace(price, "€", "", -1)
		price = strings.Replace(price, ",", "", -1)
		price = strings.Replace(price, "(incl.)", "", -1)
		price = strings.TrimSpace(price)

		fmt.Println("bedroom:", bedroom)
		fmt.Println("area:", area)
		fmt.Println("price:", price)

		// Load - is the process of writing the data into the target database.
		// fmt.Println("refreshtoken:", sw.restapi.RefreshToken)
		features := restapi.Features{
			Bedroom: bedroom,
			Area:    area,
			Price:   price,
		}
		var auth restapi.Auth
		a, _ := auth.NewAuth(sw.logger, sw.authenticateConfig, sw.projectConfig)
		a.Insert(features, sw.restapi.RefreshToken)

	})

	c.OnResponse(func(r *colly.Response) {
		sw.logger.Infoln("Visited", r.Request.URL)
	})

	c.OnScraped(func(r *colly.Response) {
		sw.logger.Infoln("Finished", r.Request.URL)
	})

	c.OnError(func(_ *colly.Response, err error) {
		sw.logger.Infoln("Something went wrong:", err)
	})

	c.Visit(sw.scraperConfig.URLToVisit)

	return nil
}
