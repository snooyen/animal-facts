package main

import (
	"context"
	"fmt"
	"net/http"
	"os"

	"github.com/go-redis/redis/v8"
	flag "github.com/spf13/pflag"

	"github.com/go-kit/kit/log"
	httptransport "github.com/go-kit/kit/transport/http"

	"github.com/snooyen/animal-facts/factpublisher/pkg/publisher"
	"github.com/snooyen/animal-facts/factpublisher/pkg/version"
)

var (
	// commandline flags
	versionInfo   = flag.Bool("version", false, "prints the version information")
	factsApiAddr  = flag.String("factsApiAddr", "facts-api:3081", "Address of facts-api grpc server")
	cronSchedule  = flag.String("schedule", "15 9 * * *", "cron schedule for publish jobs")
	port          = flag.String("port", "3001", "Port to service requests on")
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
	{
		logger = log.NewLogfmtLogger(os.Stderr)
		logger = log.With(logger, "ts", log.DefaultTimestampUTC)
		logger = log.With(logger, "caller", log.DefaultCaller)
	}

	// --version: print version info
	var err error
	if *versionInfo {
		if err = version.PrintVersion(); err != nil {
			logger.Log("err", fmt.Sprintf("failed to print version, err: %+v", err))
		}
		os.Exit(0)
	}

	redisAddr := fmt.Sprintf("%s:%s", *redisHost, *redisPort)
	logger.Log("redisAddr", redisAddr, "redisDB", *redisDB)
	rdb := redis.NewClient(&redis.Options{
		Addr:     redisAddr,
		Password: *redisPassword,
		DB:       *redisDB,
	})

	// Create Publisher Service
	s, err := publisher.New(context.Background(), rdb, logger, *factsApiAddr, *cronSchedule)
	if err != nil {
		logger.Log("err", err)
		os.Exit(1)
	}
	s = publisher.LoggingMiddleware(logger)(s)

	// Register Publisher Service Handlers
	publishHandler := httptransport.NewServer(
		publisher.MakePublishEndpoint(s),
		publisher.DecodePublishRequest,
		publisher.EncodeResponse,
	)

	http.Handle("/publish", publishHandler)
	logger.Log("msg", "HTTP", "addr", listen)
	logger.Log("err", http.ListenAndServe(listen, nil))
}
