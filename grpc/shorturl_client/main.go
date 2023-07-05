package main

import (
	"context"
	"flag"
	"log"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	pb "short-url/grpc/protob"
)


var addr = flag.String("addr", "localhost:50051", "the address to connect to")

func main() {
	flag.Parse()
	// Set up a connection to the server.
	conn, err := grpc.Dial(*addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	c := pb.NewShortUrlerClient(conn)

	// Contact the server and print out its response.
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	r, err := c.GetShortUrl(ctx, &pb.GetShortUrlRequest{OrigUrl: "github.com"})
	if err != nil {
		log.Fatalf("could not greet: %v", err)
	}
	log.Println(r.GetShortUrl())
}