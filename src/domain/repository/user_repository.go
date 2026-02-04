package repository

import (
	"context"

	"github.com/minisource/template_go/domain/model"
)

type UserRepository interface {
	ExistsUserId(ctx context.Context, userId string) (bool, error)
	CreateUser(ctx context.Context, u model.User) (model.User, error)
}
