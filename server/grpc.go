package server

import (
	"fmt"
	"net"

	"github.com/bradenrayhorn/ledger-protos/session"
	"github.com/go-redis/redis/v8"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"google.golang.org/grpc"
)

type GRPCServer struct {
	rdb *redis.Client
}

func NewGRPCServer(client *redis.Client) GRPCServer {
	return GRPCServer{
		rdb: client,
	}
}

func (s GRPCServer) Start() {
	requestedPort := viper.GetString("grpc_port")
	zap.S().Infof("starting grpc server on port %s", requestedPort)

	lis, err := net.Listen("tcp", fmt.Sprintf(":%s", requestedPort))
	if err != nil {
		zap.S().Fatalf("failed to bind grpc port %s: %v", requestedPort, err)
	}

	grpcServer := grpc.NewServer()

	session.RegisterSessionAuthenticatorServer(grpcServer, NewSessionAuthenticatorServer(s.rdb))

	if err := grpcServer.Serve(lis); err != nil {
		zap.S().Fatalf("failed to start grpc server: %s", err)
	}

}
