package server

import (
	"go-micro.dev/v5"

	"go-micro-example/app/user/service/internal/data"
	"go-micro-example/app/user/service/internal/service"

	userV1 "go-micro-example/api/gen/go/user/service/v1"
)

func InitMicroServer(
	srv micro.Service,
	userRepo *data.UserRepo,
) {
	_ = userV1.RegisterUserServiceHandler(srv.Server(), service.NewUserMicroService(userRepo))
}
