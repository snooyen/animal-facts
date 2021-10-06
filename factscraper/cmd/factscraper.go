package main

import (
	"fmt"
	"net/http"
	"os"

	flag "github.com/spf13/pflag"
	"github.com/go-redis/redis/v8"

	"github.com/go-kit/kit/log"
	httptransport "github.com/go-kit/kit/transport/http"

	"github.com/snooyen/elephant-seal-facts/factscraper/pkg/scraper"
	"github.com/snooyen/elephant-seal-facts/factscraper/pkg/version"
)

var (
	animals = map[string]string {
		"elephant-seal": "https://elephantseal.org/about-the-seals/",
	}

	// commandline flags
	versionInfo    = flag.Bool("version", false, "prints the version information")
	port = flag.String("port", "3000", "Port to service requests on")
	redisHost = flag.String("redisHost", "localhost", "Hostname/address of redis")
	redisPort = flag.String("redisPort", "6379", "Port with which to connect to redis")
	redisPassword = flag.String("redisPassword", "password123!", "Password to authenticate to redis")
	redisDB = flag.Int("redisDB", 0, "Redis DB id")
)

func main() {
	// Parse commandline flags
	flag.Parse()
	listen := fmt.Sprintf(":%s", *port)

	var logger log.Logger
	logger = log.NewLogfmtLogger(os.Stderr)
	logger = log.With(logger, "listen", listen, "caller", log.DefaultCaller)

	// --version: print version info
	var err error
	if *versionInfo {
		if err = version.PrintVersion(); err != nil {
			logger.Log("err", fmt.Sprintf("failed to print version, err: %+v", err))
		}
		os.Exit(0)
	}

	rdb := redis.NewClient(&redis.Options{
		Addr:	  fmt.Sprintf("%s:%s", *redisHost, *redisPort),
		Password: *redisPassword,
		DB:		  *redisDB,
	})

	// Create Scraper Service
	s := scraper.New(animals, rdb)
	s = scraper.LoggingMiddleware(logger)(s)

	// Register Scraper Service Handlers
	scrapeHandler := httptransport.NewServer(
		scraper.MakeScrapeEndpoint(s),
		scraper.DecodeScrapeRequest,
		scraper.EncodeResponse,
	)

	http.Handle("/scrape", scrapeHandler)
	logger.Log("msg", "HTTP", "addr", listen)
	logger.Log("err", http.ListenAndServe(listen, nil))
}
