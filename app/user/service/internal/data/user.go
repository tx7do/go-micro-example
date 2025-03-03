package data

import (
	"context"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"

	"go-micro-example/app/user/service/internal/data/models"

	userV1 "go-micro-example/api/gen/go/user/service/v1"
)

type UserRepo struct {
	db *gorm.DB
}

func NewUserRepo(db *gorm.DB) *UserRepo {
	return &UserRepo{
		db: db,
	}
}

func (r *UserRepo) ConvertModelToProto(in *models.User) *userV1.User {
	if in == nil {
		return nil
	}

	var out userV1.User
	r.CopyModelToProto(in, &out)

	return &out
}

func (r *UserRepo) CopyModelToProto(in *models.User, out *userV1.User) {
	if in == nil || out == nil {
		return
	}

	userId := (uint32)(in.ID)

	out.Id = &userId
	out.Username = &in.UserName
	out.Nickname = &in.NickName
}

func (r *UserRepo) GetUser(_ context.Context, req *userV1.GetUserRequest) (*models.User, error) {
	res := &models.User{}
	err := r.db.First(res, "id = ?", req.GetId()).Error
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (r *UserRepo) CreateUser(_ context.Context, req *userV1.CreateUserRequest) (*userV1.User, error) {
	res := &models.User{
		UserName: req.Data.GetUsername(),
		NickName: req.Data.GetNickname(),
	}

	result := r.db.Create(res)
	if result.Error != nil {
		return nil, result.Error
	}

	return r.ConvertModelToProto(res), nil
}

func (r *UserRepo) UpdateUser(_ context.Context, req *userV1.UpdateUserRequest) (*userV1.User, error) {
	res := &models.User{
		UserName: req.Data.GetUsername(),
		NickName: req.Data.GetNickname(),
	}

	result := r.db.Model(res).Updates(res)
	if result.Error != nil {
		return nil, result.Error
	}

	return r.ConvertModelToProto(res), nil
}

func (r *UserRepo) UpsertUser(_ context.Context, req *userV1.UpdateUserRequest) (*userV1.User, error) {
	res := &models.User{
		UserName: req.Data.GetUsername(),
		NickName: req.Data.GetNickname(),
	}

	result := r.db.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "id"}},
		UpdateAll: true,
	}).Create(res)
	if result.Error != nil {
		return nil, result.Error
	}

	return r.ConvertModelToProto(res), nil
}

func (r *UserRepo) DeleteUser(_ context.Context, req *userV1.DeleteUserRequest) (bool, error) {
	result := r.db.Delete(&models.User{}, req.GetId())
	if result.Error != nil {
		return false, result.Error
	}
	return true, nil
}
