package server

import (
	"google.golang.org/grpc"

	"go-micro-example/app/user/service/internal/data"
	"go-micro-example/app/user/service/internal/service"

	userV1 "go-micro-example/api/gen/go/user/service/v1"
)

func InitGrpcServer(
	srv *grpc.Server,
	userRepo *data.UserRepo,
) {
	userV1.RegisterUserServiceServer(srv, service.NewUserGrpcService(userRepo))
}
