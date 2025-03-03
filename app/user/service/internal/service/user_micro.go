package service

import (
	"context"

	"google.golang.org/protobuf/types/known/emptypb"

	"go-micro-example/app/user/service/internal/data"

	userV1 "go-micro-example/api/gen/go/user/service/v1"
)

type UserMicroService struct {
	userRepo *data.UserRepo
}

func NewUserMicroService(userRepo *data.UserRepo) *UserMicroService {
	return &UserMicroService{
		userRepo: userRepo,
	}
}

// GetUser 获取用户数据
func (s *UserMicroService) GetUser(ctx context.Context, req *userV1.GetUserRequest, resp *userV1.User) error {
	user, err := s.userRepo.GetUser(ctx, req)
	if err != nil {
		return err
	}

	s.userRepo.CopyModelToProto(user, resp)

	return nil
}

// CreateUser 创建用户
func (s *UserMicroService) CreateUser(ctx context.Context, req *userV1.CreateUserRequest, resp *emptypb.Empty) error {
	_, err := s.userRepo.CreateUser(ctx, req)

	resp = &emptypb.Empty{}

	return err
}

// UpdateUser 更新用户
func (s *UserMicroService) UpdateUser(ctx context.Context, req *userV1.UpdateUserRequest, resp *emptypb.Empty) error {
	_, err := s.userRepo.UpdateUser(ctx, req)

	resp = &emptypb.Empty{}

	return err
}

// DeleteUser 删除用户
func (s *UserMicroService) DeleteUser(ctx context.Context, req *userV1.DeleteUserRequest, resp *emptypb.Empty) error {
	_, err := s.userRepo.DeleteUser(ctx, req)

	resp = &emptypb.Empty{}

	return err
}
