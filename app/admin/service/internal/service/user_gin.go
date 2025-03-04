package service

import (
	"fmt"
	"strconv"

	"github.com/gin-gonic/gin"

	userV1 "go-micro-example/api/gen/go/user/service/v1"
)

type UserGinService struct {
	userServiceClient userV1.UserService
}

func NewUserGinService(userServiceClient userV1.UserService) *UserGinService {
	return &UserGinService{
		userServiceClient: userServiceClient,
	}
}

// GetUser 获取用户数据
func (s *UserGinService) GetUser(ctx *gin.Context) {
	var req userV1.GetUserRequest

	strId := ctx.Param("id")
	if id, err := strconv.Atoi(strId); err != nil {
		ctx.JSON(400, gin.H{"error": err.Error()})
		return
	} else {
		req.Id = uint32(id)
	}

	resp, err := s.userServiceClient.GetUser(ctx, &req)
	if err != nil {
		ctx.JSON(500, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(200, resp)

	fmt.Printf("[GIN] GET USER [%v] [%v] [%v] \n", &req, resp, err)
}

// CreateUser 创建用户
func (s *UserGinService) CreateUser(ctx *gin.Context) {
	var req userV1.CreateUserRequest
	if err := ctx.Bind(&req); err != nil {
		ctx.JSON(400, gin.H{"error": err.Error()})
		return
	}

	resp, err := s.userServiceClient.CreateUser(ctx, &req)
	if err != nil {
		ctx.JSON(500, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(200, resp)
}

// UpdateUser 更新用户
func (s *UserGinService) UpdateUser(ctx *gin.Context) {
	var req userV1.UpdateUserRequest
	if err := ctx.Bind(&req); err != nil {
		ctx.JSON(400, gin.H{"error": err.Error()})
		return
	}

	strId := ctx.Param("id")
	if id, err := strconv.Atoi(strId); err != nil {
		ctx.JSON(400, gin.H{"error": err.Error()})
		return
	} else {
		req.Id = uint32(id)
	}

	resp, err := s.userServiceClient.UpdateUser(ctx, &req)
	if err != nil {
		ctx.JSON(500, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(200, resp)
}

// DeleteUser 删除用户
func (s *UserGinService) DeleteUser(ctx *gin.Context) {
	var req userV1.DeleteUserRequest
	if err := ctx.Bind(&req); err != nil {
		ctx.JSON(400, gin.H{"error": err.Error()})
		return
	}

	strId := ctx.Param("id")
	if id, err := strconv.Atoi(strId); err != nil {
		ctx.JSON(400, gin.H{"error": err.Error()})
		return
	} else {
		req.Id = uint32(id)
	}

	resp, err := s.userServiceClient.DeleteUser(ctx, &req)
	if err != nil {
		ctx.JSON(500, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(200, resp)
}
