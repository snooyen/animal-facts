package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/go-kit/kit/log"
	"github.com/go-redis/redis/v8"
	flag "github.com/spf13/pflag"
	"github.com/twilio/twilio-go"

	"github.com/snooyen/animal-facts/factadmin/pkg/admin"
	"github.com/snooyen/animal-facts/factadmin/pkg/version"
)

var (
	// commandline flags
	versionInfo   = flag.Bool("version", false, "prints the version information")
	host          = flag.String("host", "localhost", "service hostname")
	port          = flag.String("port", "3002", "Port to service requests on")
	redisHost     = flag.String("redisHost", "localhost", "Hostname/address of redis")
	redisPort     = flag.String("redisPort", "6379", "Port with which to connect to redis")
	redisPassword = flag.String("redisPassword", "password123!", "Password to authenticate to redis")
	redisDB       = flag.Int("redisDB", 0, "Redis DB id")
)

func main() {
	// Parse commandline flags
	flag.Parse()
	listen := fmt.Sprintf("%s:%s", *host, *port)

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
	redisAddr := fmt.Sprintf("%s:%s", *redisHost, *redisPort)
	logger.Log("redisAddr", redisAddr, "redisDB", *redisDB)
	rdb := redis.NewClient(&redis.Options{
		Addr:     redisAddr,
		Password: *redisPassword,
		DB:       *redisDB,
	})

	// Initialize Twilio Client
	twilioClient := twilio.NewRestClient()

	ctx, cancel := context.WithCancel(context.Background())
	// signalChan will catch SIGINT and SIGTERM and allow the parser to cleanup before exiting the program
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM)
	defer func() {
		signal.Stop(signalChan)
		cancel()
	}()

	// Create admin Service
	s := admin.New(rdb, twilioClient, logger)
	s = admin.LoggingMiddleware(logger)(s)

	httpServer := http.Server{
		Addr:    listen,
		Handler: admin.NewHTTPHandler(s),
	}

	err = s.ProcessApprovalRequests(ctx) // Background process handles redis pubsub messages
	if err != nil {
		logger.Log("err", fmt.Sprintf("failed to start background processes, err: %+v", err))
		os.Exit(1)
	}

	// Background goroutine waits for ctx to be canceled or SIGINT/SIGTERM to be caught
	go func() {
		select {
		case <-signalChan: //First SIGINT/SIGTERM; cancel the context & allow parser to cleanup
			httpServer.Shutdown(ctx)
			cancel()
		case <-ctx.Done():
		}
		<-signalChan // Second SIGINT/SIGTERM; Hard exit (cleanup skipped).
		os.Exit(2)
	}()

	logger.Log("msg", "HTTP", "addr", listen)
	logger.Log("err", httpServer.ListenAndServe())
}
