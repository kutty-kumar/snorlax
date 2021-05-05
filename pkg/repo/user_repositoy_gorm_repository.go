package repo

import (
	"context"

	charminder "github.com/kutty-kumar/charminder/pkg"
	"github.com/kutty-kumar/user_service/pkg/domain/entity"
)

type UserGORMRepo struct {
	charminder.BaseDao
}

func NewUserGORMRepo(dao charminder.BaseDao) UserGORMRepo {
	return UserGORMRepo{
		dao,
	}
}

func (ugr *UserGORMRepo) GetUserByEmailPassword(ctx context.Context, email string, password string) (error, entity.User) {
	var user entity.User
	err := ugr.GetDb().First(&user).Where("email = ?", email).Where("password = ?", password).Error
	if err != nil {
		return err, entity.User{}
	}
	return nil, user
}
