package unit_tests

import (
	"context"
	"fmt"
	"log"
	"net"
	"testing"

	grpcapi "IPBlockerService/apis/grpc"
	proto "IPBlockerService/apis/grpc/proto"

	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/test/bufconn"
)

func testable_grpc_server(ctx context.Context) (proto.IPBlockerServiceClient, func()) {
	buffer := 101024 * 1024
	lis := bufconn.Listen(buffer)

	baseServer := grpc.NewServer()
	proto.RegisterIPBlockerServiceServer(baseServer, &grpcapi.IPBlockerServiceServer{})
	go func() {
		if err := baseServer.Serve(lis); err != nil {
			log.Printf("error serving server: %v", err)
		}
	}()

	conn, err := grpc.DialContext(ctx, "",
		grpc.WithContextDialer(func(context.Context, string) (net.Conn, error) {
			return lis.Dial()
		}), grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Printf("error connecting to server: %v", err)
	}

	closer := func() {
		err := lis.Close()
		if err != nil {
			log.Printf("error closing listener: %v", err)
		}
		baseServer.Stop()
	}

	client := proto.NewIPBlockerServiceClient(conn)

	return client, closer
}

func TestHealthCheckRoute_grpc(t *testing.T) {
	Setup()
	defer Teardown()
	ctx := context.Background()
	client, closer := testable_grpc_server(ctx)
	defer closer()

	out, err := client.HealthCheck(ctx, &proto.HealthCheckRequest{})

	println(fmt.Sprintf("%+v", out))
	assert.Nil(t, err)
	assert.Equal(t, out.ServiceOnline, true)
	assert.Equal(t, out.DatabaseAccessable, true)
}

func TestHealthCheckRoute_grpc_db_inaccessable(t *testing.T) {
	// Don't perform setup/teardown so the db file cannot be found
	ctx := context.Background()
	client, closer := testable_grpc_server(ctx)
	defer closer()

	out, err := client.HealthCheck(ctx, &proto.HealthCheckRequest{})

	println(fmt.Sprintf("%+v", out))
	assert.Nil(t, err)
	assert.Equal(t, out.ServiceOnline, true)
	assert.Equal(t, out.DatabaseAccessable, false)
}

func TestAuthroizeIPRoute_grpc(t *testing.T) {
	Setup()
	defer Teardown()
	ctx := context.Background()
	client, closer := testable_grpc_server(ctx)
	defer closer()

	requestData := proto.AuthroizeIPRequest{
		IPAddress: "162.226.203.50",
		ValidCountries: []string{
			"United States",
			"United Kingdom",
			"Canada",
		},
	}
	out, err := client.AuthorizeIP(ctx, &requestData)

	assert.Nil(t, err)
	assert.Equal(t, true, out.Authorized)
}


func TestAuthroizeIPRoute_grpc_invalid_ip(t *testing.T) {
	Setup()
	defer Teardown()
	ctx := context.Background()
	client, closer := testable_grpc_server(ctx)
	defer closer()

	requestData := proto.AuthroizeIPRequest{
		IPAddress: "103.136.43.2",
		ValidCountries: []string{
			"United States",
			"United Kingdom",
			"Canada",
		},
	}
	out, err := client.AuthorizeIP(ctx, &requestData)

	assert.Nil(t, err)
	assert.Equal(t, false, out.Authorized)
}