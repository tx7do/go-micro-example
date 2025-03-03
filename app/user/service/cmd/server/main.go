package main

import (
	"context"
	"flag"
	"log"

	"go-micro-example/app/user/service/internal/data"
	"go-micro-example/app/user/service/internal/server"

	"go-micro-example/pkg/app"
	"go-micro-example/pkg/service"
)

var version = flag.String("version", "1.0.0", "service version")
var confPath = flag.String("conf", "./configs/server.yaml", "service config path")

// go build -ldflags "-X main.version=x.y.z"

func init() {
	flag.Parse()
}

func main() {
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	aApp := app.New()

	// 初始化服务
	if err := aApp.Start(ctx, *confPath, service.UserService, *version); err != nil {
		panic(err)
	}
	defer aApp.Stop()

	db := data.NewGormClient(aApp.Config().Data, aApp.Logger())

	userRepo := data.NewUserRepo(db)

	// 初始化rpc服务
	server.InitMicroServer(aApp.Service(), userRepo)

	// 启动服务
	if err := aApp.Run(); err != nil {
		log.Fatalf("Failed to run service: %v", err)
	}
}
