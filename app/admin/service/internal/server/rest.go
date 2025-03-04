package server

import (
	"context"

	"github.com/gin-gonic/gin"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"

	"go-micro-example/app/admin/service/internal/data"
	"go-micro-example/app/admin/service/internal/service"

	adminV1 "go-micro-example/api/gen/go/admin/service/v1"
	userV1 "go-micro-example/api/gen/go/user/service/v1"

	"go-micro-example/pkg/app"
)

func InitServer(ctx context.Context, app *app.App) error {
	// 创建内部服务的客户端
	userServiceClient := data.NewUserServiceMicroClient(app.MicroClient())

	// 注册路由
	if app.Config().Server != nil && app.Config().Server.Rest != nil {
		if app.Config().Server.Rest.GetEnableGrpcGateway() {
			registerGrpcGatewayRouter(ctx, app.GrpcGatewayServeMux(), userServiceClient)
		} else {
			registerGinRouter(ctx, app.GinRouter(), userServiceClient)
		}
	}

	return nil
}

// 注册gRPC-gateway路由
func registerGrpcGatewayRouter(
	ctx context.Context,
	mux *runtime.ServeMux,
	userServiceClient userV1.UserService,
) {
	_ = adminV1.RegisterUserServiceGWServer(ctx, mux, service.NewUserService(userServiceClient))
}

// registerGinRouter 注册GIN路由
func registerGinRouter(
	_ context.Context,
	router *gin.Engine,
	userServiceClient userV1.UserService,
) {
	g := router.Group("/admin/v1")

	{
		userService := service.NewUserGinService(userServiceClient)
		g.GET("/users/:id", userService.GetUser)
		g.POST("/users", userService.CreateUser)
		g.PUT("/users/:id", userService.UpdateUser)
		g.DELETE("/users/:id", userService.DeleteUser)
	}
}
