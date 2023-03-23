package main

import (
	"context"
	"flag"
	"log"
	"time"

	proto "IPBlockerService/apis/grpc/proto"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var (
	addr = flag.String("addr", "127.0.0.1:64075", "the address to connect to")
)

// This client program is made available to test the server, nothing more.
func main() {
	flag.Parse()
	conn, err := grpc.Dial(*addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	c := proto.NewIPBlockerServiceClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	healthCheckRes, err := c.HealthCheck(ctx, &proto.HealthCheckRequest{})
	if err != nil {
		log.Fatalf("could not HealthCheck: %v", err)
	}
	log.Printf("Server health: %+v", healthCheckRes)

	authorizeIPRes, err := c.AuthorizeIP(ctx, &proto.AuthroizeIPRequest{
		IPAddress: "162.226.203.50",
		ValidCountries: []string{
			"United States",
			"United Kingdom",
			"Canada",
		},
	})
	if err != nil {
		log.Fatalf("could not Authorize IP: %v", err)
	}
	log.Printf("Authorize IP: %+v", authorizeIPRes)
}
