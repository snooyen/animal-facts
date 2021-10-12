package publisher

import (
	"fmt"
	"google.golang.org/grpc"
	"os"
	"time"

	//	"github.com/go-kit/kit/log"
	//	grpctransport "github.com/go-kit/kit/transport/grpc"

	pb "github.com/snooyen/animal-facts/facts/pb"
)

func NewFactsClient(grpcAddr string) pb.FactsClient {
	conn, err := grpc.Dial(grpcAddr, grpc.WithInsecure(), grpc.WithTimeout(time.Second))
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v", err)
		os.Exit(1)
	}

	c := pb.NewFactsClient(conn)

	return c
}
