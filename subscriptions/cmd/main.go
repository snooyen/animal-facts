package main

import (
	"context"
	"fmt"
	"os"

	"github.com/go-kit/kit/log"
	"github.com/go-redis/redis/v8"
	"github.com/snooyen/animal-facts/subscriptions/pkg/subscriptions"
	"github.com/snooyen/animal-facts/subscriptions/pkg/version"
	flag "github.com/spf13/pflag"
)

var (
	versionInfo   = flag.Bool("version", false, "prints the version information")
	redisHost     = flag.String("redisHost", "localhost", "Hostname/address of redis")
	redisPort     = flag.String("redisPort", "6379", "Port with which to connect to redis")
	redisPassword = flag.String("redisPassword", "password123!", "Password to authenticate to redis")
	redisDB       = flag.Int("redisDB", 0, "Redis DB id")
)

func main() {
	flag.Parse()

	// Create logger to pass to components
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

	// Create redis client
	redisAddr := fmt.Sprintf("%s:%s", *redisHost, *redisPort)
	rdb := redis.NewClient(&redis.Options{
		Addr:     redisAddr,
		Password: *redisPassword,
		DB:       *redisDB,
	})

	sub := subscriptions.New("elephant-seal", logger, rdb)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	sub.Start(ctx)
}
