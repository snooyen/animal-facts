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
	"github.com/twilio/twilio-go"
	"google.golang.org/grpc"
)

var (
	versionInfo     = flag.Bool("version", false, "prints the version information")
	redisHost       = flag.String("redisHost", "localhost", "Hostname/address of redis")
	redisPort       = flag.String("redisPort", "6379", "Port with which to connect to redis")
	redisPassword   = flag.String("redisPassword", "password123!", "Password to authenticate to redis")
	redisDB         = flag.Int("redisDB", 0, "Redis DB id")
	twilioNumber    = flag.String("twilioNumber", "(555) 555-5555", "Twilio number to send messages from")
	usersClientAddr = flag.String("usersClientAddr", "users-api:3081", "Address of users grpc service")
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

	twilioClient := twilio.NewRestClient()

	ctx, cancel := context.WithCancel(context.Background())
	conn, err := grpc.Dial(*usersClientAddr, grpc.WithInsecure())
	if err != nil {
		fmt.Errorf("failed to connect to users grpc service, err: %+v", err)
		os.Exit(1)
	}
	defer func() {
		cancel()
		conn.Close()
	}()

	sub := subscriptions.New("elephant-seal", logger, rdb, conn, twilioClient, *twilioNumber)
	sub.Start(ctx)
}
