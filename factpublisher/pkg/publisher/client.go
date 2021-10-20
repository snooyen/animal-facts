package publisher

import (
	"context"
	"fmt"
	"os"
	"time"

	"google.golang.org/grpc"

	pb "github.com/snooyen/animal-facts/facts/pb"
)

func NewFactsClient(ctx context.Context, grpcAddr string) pb.FactsClient {
	ctx, cancel := context.WithTimeout(ctx, 100*time.Millisecond)
	defer cancel()
	conn, err := grpc.DialContext(ctx, grpcAddr, grpc.WithInsecure())
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v", err)
		os.Exit(1)
	}

	c := pb.NewFactsClient(conn)

	return c
}
