package main

import (
	"fmt"
	"net"
	"net/http"
	"os"

	"github.com/go-redis/redis/v8"
	flag "github.com/spf13/pflag"
	"google.golang.org/grpc"

	"github.com/go-kit/kit/log"

	"github.com/snooyen/animal-facts/users/pb"
	"github.com/snooyen/animal-facts/users/pkg/api"
	"github.com/snooyen/animal-facts/users/pkg/version"
)

var (
	// commandline flags
	versionInfo    = flag.Bool("version", false, "prints the version information")
	httpListenPort = flag.String("httpPort", "3080", "Port to service HTTP requests on")
	grpcListenPort = flag.String("grpcPort", "3081", "Port to service GRPC requests on")
	redisHost      = flag.String("redisHost", "localhost", "Hostname/address of redis")
	redisPort      = flag.String("redisPort", "6379", "Port with which to connect to redis")
	redisPassword  = flag.String("redisPassword", "password123!", "Password to authenticate to redis")
	redisDB        = flag.Int("redisDB", 0, "Redis DB id")
)

func main() {
	// Parse commandline flags
	flag.Parse()
	httpListen := fmt.Sprintf(":%s", *httpListenPort)
	grpcListen := fmt.Sprintf(":%s", *grpcListenPort)

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

	// Initialize Redis Client
	rdb := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%s", *redisHost, *redisPort),
		Password: *redisPassword,
		DB:       *redisDB,
	})

	var (
		service     = api.New(rdb, logger)
		endpoints   = api.MakeServerEndpoints(service, logger)
		httpHandler = api.NewHTTPHandler(endpoints, logger)
		grpcServer  = api.NewGRPCServer(endpoints, logger)
	)

	grpcErr := make(chan error)
	httpErr := make(chan error)
	listenGRPC(grpcListen, grpcServer, logger, grpcErr)
	listenHTTP(httpListen, httpHandler, logger, httpErr)

	select {
	case <-httpErr:
		logger.Log("exit", <-httpErr)
	case <-grpcErr:
		logger.Log("exit", <-grpcErr)
	}

}

func listenGRPC(grpcAddr string, grpcServer pb.UsersServer, logger log.Logger, errChan chan error) net.Listener {
	logger = log.With(logger, "grpcListen", grpcAddr, "caller", log.DefaultCaller)
	grpcListener, err := net.Listen("tcp", grpcAddr)
	if err != nil {
		logger.Log("transport", "gRPC", "during", "Listen", "err", err)
		os.Exit(1)
	}
	logger.Log("transport", "gRPC", "addr", grpcAddr)
	baseGRPCServer := grpc.NewServer()
	pb.RegisterUsersServer(baseGRPCServer, grpcServer)

	go func() {
		errChan <- baseGRPCServer.Serve(grpcListener)
	}()

	return grpcListener
}

func listenHTTP(httpAddr string, handler http.Handler, logger log.Logger, errChan chan error) net.Listener {
	logger = log.With(logger, "httpListen", httpAddr, "caller", log.DefaultCaller)
	httpListener, err := net.Listen("tcp", httpAddr)
	if err != nil {
		logger.Log("transport", "HTTP", "during", "Listen", "err", err)
		os.Exit(1)
	}
	logger.Log("transport", "HTTP", "addr", httpAddr)

	go func() {
		errChan <- http.Serve(httpListener, handler)
	}()

	return httpListener
}
