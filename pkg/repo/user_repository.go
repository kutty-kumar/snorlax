package repo

import (
	"context"

	charminder "github.com/kutty-kumar/charminder/pkg"
	"github.com/kutty-kumar/user_service/pkg/domain/entity"
)

type UserRepo interface {
	charminder.BaseRepository
	GetUserByEmailPassword(ctx context.Context, email string, password string) (error, entity.User)
}
