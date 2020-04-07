# Pararius web site scraper
Scraping pararius.com with go-colly.

## Disclaimers

The script scrapes the pararius.com web site to collect data about the shape of the company's business. No guarantees are made about the quality of data obtained using this script, statistically or about an individual page. So please check your results.

Maybe the pararius.com site refuses repeated requests. Try to run the script using a number of proxy IP addresses to avoid being turned away.

## Status and recent changes

### May 2019 (0.0.1)

Beta version.

## Prerequisites

- Golang go1.11.5
- Go Colly - https://github.com/gocolly/colly

## Using the script

You must be comfortable messing about with API consume and Golang to use this. 

To run this scraper you will need to use go1.11.5 or later and install the Go Colly properly.

This scraper stores the survey results in Jexia dataset, so, Jexia account is required. This kinda setup permits the distributed use of the scraped data.

### Installing Go Colly Pararius

```git clone git@github.com:edwardmartinsjr/pararius-scraper.git```

### Running a survey 

Just set env vars:
- *SCRAPER_WORKER_INTERVAL* (eg.: 1m)
- *CONFIG* (eg.: /Users/.../pararius-scraper/resources/pararius-scraper.toml)
- *EMAIL* (Authorized Jexia UMS e-mail)
- *PASSWORD*
- *PROJECT_ID* (any project that you created for this purposes...)
- *TRAINING_DATA_SET* (any dataset that you created for this purposes...)

and run: make run-service

## Configuration

There are some env variables you can set to modify pararius-scraper configs:

| NAME                               | DESCRIPTION
-------------------------------------|------------------------------------------------------------------------
| SERVICE_NAME                       | Modify the pararius-scraper service name
| SCRAPER_WORKER_INTERVAL            | Interval of time which the worker for scrap pararius web site will be ran. If the value is an integer it is processed as seconds, otherwise it is parsed as a duration: https://golang.org/pkg/time/#ParseDuration (Default: 1h)
| CONFIG                             | Change the config file path
| ALLOWED_DOMAINS                    | Sets the domain whitelist used by the Collector
| URL_TO_VISIT                       | Sets the URl to be visited by scraper
| EMAIL                              | User e-mail to authenticate REST API
| PASSWORD                           | User password to authenticate REST API
| PROJECT_ID                         | User Project ID
| TRAINING_DATA_SET                  | User Training Data Set
