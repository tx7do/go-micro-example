package server

import (
	"context"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"

	"go-micro-example/app/admin/service/internal/service"

	adminV1 "go-micro-example/api/gen/go/admin/service/v1"
	userV1 "go-micro-example/api/gen/go/user/service/v1"
)

func InitRestServer(
	ctx context.Context,
	mux *runtime.ServeMux,
	userServiceClient userV1.UserService,
) {
	_ = adminV1.RegisterUserServiceGWServer(ctx, mux, service.NewUserService(userServiceClient))
}
