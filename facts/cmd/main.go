package main

import (
	"context"
	"fmt"
	"os"

	"github.com/go-redis/redis/v8"
	flag "github.com/spf13/pflag"

	"github.com/go-kit/kit/log"

	"github.com/snooyen/animal-facts/facts/pkg/api"
	"github.com/snooyen/animal-facts/facts/pkg/version"
)

var (
	// commandline flags
	versionInfo   = flag.Bool("version", false, "prints the version information")
	port          = flag.String("port", "3000", "Port to service requests on")
	redisHost     = flag.String("redisHost", "localhost", "Hostname/address of redis")
	redisPort     = flag.String("redisPort", "6379", "Port with which to connect to redis")
	redisPassword = flag.String("redisPassword", "password123!", "Password to authenticate to redis")
	redisDB       = flag.Int("redisDB", 0, "Redis DB id")
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

	// Initialize Redis Client
	rdb := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%s", *redisHost, *redisPort),
		Password: *redisPassword,
		DB:       *redisDB,
	})

	// Create Fact API Service
	s := api.New(rdb, logger)
	animals, _ := s.GetAnimals(context.Background())
	logger.Log("msg", fmt.Sprintf("%v", animals))
}
