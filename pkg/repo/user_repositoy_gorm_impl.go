package repo

import (
	"context"

	charminder "github.com/kutty-kumar/charminder/pkg"
	"github.com/kutty-kumar/snorlax/pkg/domain/entity"
)

type UserRepoGormImpl struct {
	charminder.BaseDao
}

func NewUserRepoGormImpl(dao charminder.BaseDao) UserRepoGormImpl {
	return UserRepoGormImpl{
		dao,
	}
}

func (ugr *UserRepoGormImpl) GetUserByEmailPassword(ctx context.Context, email string, password string) (error, entity.User) {
	var user entity.User
	err := ugr.GetDb().Model(user).Where("email = ?", email).Where("password = ?", password).First(&user).Error
	if err != nil {
		return err, entity.User{}
	}
	return nil, user
}
