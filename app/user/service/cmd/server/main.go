package main

import (
	"context"
	"flag"
	"log"

	"go-micro-example/app/user/service/internal/data/models"
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

	a := app.New(
		app.WithConfigPath(*confPath),
		app.WithServiceName(service.UserService),
		app.WithVersion(*version),
		app.WithGormMigrators(models.GetMigrates()),
	)

	// 初始化服务
	if err := a.Start(ctx); err != nil {
		panic(err)
	}
	defer a.Stop()

	// 初始化rpc服务
	server.InitMicroServer(ctx, a)

	// 启动服务
	if err := a.Run(); err != nil {
		log.Fatalf("Failed to run service: %v", err)
	}
}
