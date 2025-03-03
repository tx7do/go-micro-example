package data

import (
	"go-micro.dev/v5/client"

	userV1 "go-micro-example/api/gen/go/user/service/v1"

	"go-micro-example/pkg/service"
)

func NewUserServiceMicroClient(cli client.Client) userV1.UserService {
	return userV1.NewUserService(service.UserService, cli)
}
