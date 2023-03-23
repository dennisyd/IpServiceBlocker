package main

import (
	"IPBlockerService/apis/grpc/proto"
	apishttp "IPBlockerService/apis/http"
	"IPBlockerService/base"
	"net"

	grpcapi "IPBlockerService/apis/grpc"

	zlog "github.com/rs/zerolog/log"

	"google.golang.org/grpc"
)

func startgRPCServer(config base.Configuration) {
	listener, err := net.Listen("tcp", config.IPBLOCKERSERVICE_SERVER_IP_ADDRESS+":"+config.IPBLOCKERSERVICE_SERVER_GRPC_PORT)
	if err != nil {
		zlog.Printf("failed to listen: %v", err)
		return
	}

	baseServer := grpc.NewServer()
	proto.RegisterIPBlockerServiceServer(baseServer, &grpcapi.IPBlockerServiceServer{})

	go func() {
		zlog.Printf("grpc server listening at %v", listener.Addr())
		if err := baseServer.Serve(listener); err != nil {
			zlog.Printf("error serving server: %v", err)
			return
		}
	}()
}

func startHTTPServer(config base.Configuration) {
	api := apishttp.SetupHTTPAPI()
	api.Run(config.IPBLOCKERSERVICE_SERVER_IP_ADDRESS + ":" + config.IPBLOCKERSERVICE_SERVER_HTTP_PORT)
}

func main() {
	config := base.GetConfig()
	startgRPCServer(config)
	startHTTPServer(config)
}
