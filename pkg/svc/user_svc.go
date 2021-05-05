package svc

import (
	"context"

	"github.com/infobloxopen/atlas-app-toolkit/errors"
	charminder "github.com/kutty-kumar/charminder/pkg"
	"github.com/kutty-kumar/ho_oh/user_service_v1"
	"github.com/kutty-kumar/user_service/pkg/domain/entity"
	"github.com/kutty-kumar/user_service/pkg/repo"
	"google.golang.org/grpc/codes"
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

func (us *UserSvc) CreateUser(ctx context.Context, in *user_service_v1.CreateUserRequest) (*user_service_v1.UserOperationResponse, error) {
	var resp user_service_v1.UserOperationResponse
	var user entity.User
	user.FillProperties(*in.Payload)
	err, createdUser := us.Create(ctx, &user)
	if err != nil {
		return nil, errors.New(ctx, codes.Internal, "error in creating user")
	}
	responseDto := createdUser.ToDto().(user_service_v1.UserDto)
	resp.Response = &responseDto
	return &resp, nil
}

func (us *UserSvc) UpdateUser(ctx context.Context, in *user_service_v1.UpdateUserRequest) (*user_service_v1.UserOperationResponse, error) {
	var resp user_service_v1.UserOperationResponse
	var user entity.User
	user.FillProperties(*in.Payload)
	err, updatedUser := us.Update(ctx, in.UserId, &user)
	if err != nil {
		return nil, errors.New(ctx, codes.Internal, "error in updating user")
	}
	responseDto := updatedUser.ToDto().(user_service_v1.UserDto)
	resp.Response = &responseDto
	return &resp, nil
}

func (us *UserSvc) GetUserByExternalId(ctx context.Context, in *user_service_v1.GetUserByExternalIdRequest) (*user_service_v1.UserOperationResponse, error) {
	var resp user_service_v1.UserOperationResponse
	err, user := us.FindByExternalId(ctx, in.UserId)
	if err != nil {
		return nil, errors.New(ctx, codes.NotFound, "user not found")
	}
	responseDto := user.ToDto().(user_service_v1.UserDto)
	resp.Response = &responseDto
	return &resp, nil
}

func (us *UserSvc) GetUserByEmailAndPassword(ctx context.Context, in *user_service_v1.GetUserByEmailAndPasswordRequest) (*user_service_v1.UserOperationResponse, error) {
	var resp user_service_v1.UserOperationResponse
	err, user := us.userRepo.GetUserByEmailPassword(ctx, in.Email, in.Password)
	if err != nil {
		return nil, errors.New(ctx, codes.NotFound, "user not found")
	}
	responseDto := user.ToDto().(user_service_v1.UserDto)
	resp.Response = &responseDto
	return &resp, nil
}
