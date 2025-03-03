package service

import (
	"context"
	"fmt"

	"google.golang.org/protobuf/types/known/emptypb"

	adminV1 "go-micro-example/api/gen/go/admin/service/v1"
	userV1 "go-micro-example/api/gen/go/user/service/v1"
)

type UserService struct {
	adminV1.UnimplementedUserServiceServer

	userServiceClient userV1.UserService
}

func NewUserService(userServiceClient userV1.UserService) *UserService {
	return &UserService{
		userServiceClient: userServiceClient,
	}
}

// GetUser 获取用户数据
func (s *UserService) GetUser(ctx context.Context, req *userV1.GetUserRequest) (*userV1.User, error) {
	resp, err := s.userServiceClient.GetUser(ctx, req)
	fmt.Printf("GET USER 1 [%v] [%v] [%v] \n", req, resp, err)
	return resp, err
}

// CreateUser 创建用户
func (s *UserService) CreateUser(ctx context.Context, req *userV1.CreateUserRequest) (*emptypb.Empty, error) {
	return s.userServiceClient.CreateUser(ctx, req)
}

// UpdateUser 更新用户
func (s *UserService) UpdateUser(ctx context.Context, req *userV1.UpdateUserRequest) (*emptypb.Empty, error) {
	return s.userServiceClient.UpdateUser(ctx, req)
}

// DeleteUser 删除用户
func (s *UserService) DeleteUser(ctx context.Context, req *userV1.DeleteUserRequest) (*emptypb.Empty, error) {
	return s.userServiceClient.DeleteUser(ctx, req)
}
