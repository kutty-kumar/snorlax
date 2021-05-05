package svc

import (
	"context"

	"errors"

	charminder "github.com/kutty-kumar/charminder/pkg"
	"github.com/kutty-kumar/ho_oh/snorlax_v1"
	"github.com/kutty-kumar/snorlax/pkg/domain/entity"
	"github.com/kutty-kumar/snorlax/pkg/repo"
)

type UserSvc struct {
	charminder.BaseSvc
	userRepo repo.UserRepo
}

func NewUserSvc(userRepo repo.UserRepo) UserSvc {
	userSvc := UserSvc{userRepo: userRepo}
	userSvc.BaseSvc.Init(userRepo)
	return userSvc
}

func (us *UserSvc) CreateUser(ctx context.Context, in *snorlax_v1.CreateUserRequest) (*snorlax_v1.UserOperationResponse, error) {
	var resp snorlax_v1.UserOperationResponse
	var user entity.User
	user.FillProperties(*in.Payload)
	err, createdUser := us.Create(ctx, &user)
	if err != nil {
		return nil, errors.New("error in creating user")
	}
	responseDto := createdUser.ToDto().(snorlax_v1.UserDto)
	resp.Response = &responseDto
	return &resp, nil
}

func (us *UserSvc) UpdateUser(ctx context.Context, in *snorlax_v1.UpdateUserRequest) (*snorlax_v1.UserOperationResponse, error) {
	var resp snorlax_v1.UserOperationResponse
	var user entity.User
	user.FillProperties(*in.Payload)
	err, updatedUser := us.Update(ctx, in.UserId, &user)
	if err != nil {
		return nil, errors.New("error in updating user")
	}
	responseDto := updatedUser.ToDto().(snorlax_v1.UserDto)
	resp.Response = &responseDto
	return &resp, nil
}

func (us *UserSvc) GetUserByExternalId(ctx context.Context, in *snorlax_v1.GetUserByExternalIdRequest) (*snorlax_v1.UserOperationResponse, error) {
	var resp snorlax_v1.UserOperationResponse
	err, user := us.FindByExternalId(ctx, in.UserId)
	if err != nil {
		return nil, errors.New("user not found")
	}
	responseDto := user.ToDto().(snorlax_v1.UserDto)
	resp.Response = &responseDto
	return &resp, nil
}

func (us *UserSvc) GetUserByEmailAndPassword(ctx context.Context, in *snorlax_v1.GetUserByEmailAndPasswordRequest) (*snorlax_v1.UserOperationResponse, error) {
	var resp snorlax_v1.UserOperationResponse
	err, user := us.userRepo.GetUserByEmailPassword(ctx, in.Email, in.Password)
	if err != nil {
		return nil, errors.New("user not found")
	}
	responseDto := user.ToDto().(snorlax_v1.UserDto)
	resp.Response = &responseDto
	return &resp, nil
}
