package grpc

import (
	proto "IPBlockerService/apis/grpc/proto"
	"IPBlockerService/base"
	"IPBlockerService/errors"
	"IPBlockerService/handlers"
	"context"
	"net"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	guuid "github.com/google/uuid"
)

type IPBlockerServiceServer struct {
	proto.UnimplementedIPBlockerServiceServer
}

func NewIPBlockerServiceServer() proto.IPBlockerServiceServer {
	return &IPBlockerServiceServer{}
}

func addMsgId(ctx context.Context) context.Context {
	return context.WithValue(ctx, base.CONTEXT_MESSAGE_ID, guuid.New())
}

func (s *IPBlockerServiceServer) HealthCheck(ctx context.Context, in *proto.HealthCheckRequest) (*proto.HealthCheckResponse, error) {
	ctx = addMsgId(ctx)
	logWhenFinished, functionName := base.HandleNewMessage(ctx)
	defer logWhenFinished(ctx, functionName)

	res := proto.HealthCheckResponse{
		ServiceOnline:      true,
		DatabaseAccessable: handlers.IsDbAvailable(ctx),
	}

	base.LogResponse(ctx, &res)

	return &res, nil
}

func (s *IPBlockerServiceServer) AuthorizeIP(ctx context.Context, in *proto.AuthroizeIPRequest) (*proto.AuthroizeIPResponse, error) {
	ctx = addMsgId(ctx)
	logWhenFinished, functionName := base.HandleNewMessage(ctx)
	defer logWhenFinished(ctx, functionName)

	ip := net.ParseIP(in.IPAddress)
	if ip == nil {
		err := errors.InvalidIPAddressError{}
		return nil, status.Errorf(codes.InvalidArgument, err.Error())
	}

	authorized, err := handlers.AuthorizeIPAddress(ctx, ip, in.ValidCountries)

	if err != nil {
		return nil, status.Errorf(codes.Internal, "Internal Server Error")
	}

	res := proto.AuthroizeIPResponse{
		Authorized: authorized,
	}
	base.LogResponse(ctx, &res)

	return &res, nil
}
