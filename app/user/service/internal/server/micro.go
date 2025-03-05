package server

import (
	"context"

	"go-micro-example/app/user/service/internal/data"
	"go-micro-example/app/user/service/internal/service"

	userV1 "go-micro-example/api/gen/go/user/service/v1"

	"go-micro-example/pkg/app"
)

func InitMicroServer(_ context.Context, app *app.App) {
	userRepo := data.NewUserRepo(app.Gorm())

	_ = userV1.RegisterUserServiceHandler(app.MicroService().Server(), service.NewUserMicroService(userRepo))
}
