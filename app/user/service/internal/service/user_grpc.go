package service

import (
	"context"

	"google.golang.org/protobuf/types/known/emptypb"

	"go-micro-example/app/user/service/internal/data"

	userV1 "go-micro-example/api/gen/go/user/service/v1"
)

type UserGrpcService struct {
	userV1.UnimplementedUserServiceServer

	userRepo *data.UserRepo
}

func NewUserGrpcService(userRepo *data.UserRepo) *UserGrpcService {
	return &UserGrpcService{
		userRepo: userRepo,
	}
}

// GetUser 获取用户数据
func (s *UserGrpcService) GetUser(ctx context.Context, req *userV1.GetUserRequest) (*userV1.User, error) {
	user, err := s.userRepo.GetUser(ctx, req)
	if err != nil {
		return nil, err
	}

	return s.userRepo.ConvertModelToProto(user), nil
}

// CreateUser 创建用户
func (s *UserGrpcService) CreateUser(ctx context.Context, req *userV1.CreateUserRequest) (*emptypb.Empty, error) {
	_, err := s.userRepo.CreateUser(ctx, req)

	return &emptypb.Empty{}, err
}

// UpdateUser 更新用户
func (s *UserGrpcService) UpdateUser(ctx context.Context, req *userV1.UpdateUserRequest) (*emptypb.Empty, error) {
	_, err := s.userRepo.UpdateUser(ctx, req)

	return &emptypb.Empty{}, err
}

// DeleteUser 删除用户
func (s *UserGrpcService) DeleteUser(ctx context.Context, req *userV1.DeleteUserRequest) (*emptypb.Empty, error) {
	_, err := s.userRepo.DeleteUser(ctx, req)

	return &emptypb.Empty{}, err
}
