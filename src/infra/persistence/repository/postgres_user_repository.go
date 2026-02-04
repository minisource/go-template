package repository

import (
	"context"

	"github.com/minisource/template_go/config"
	"github.com/minisource/template_go/domain/model"
	gormdb "github.com/minisource/go-common/db/gorm"
	"github.com/minisource/go-common/logging"
)

const userIdFilterExp string = "user_id = ?"
const countFilterExp string = "count(*) > 0"

type PostgresUserRepository struct {
	*BaseRepository[model.User]
}

func NewUserRepository(cfg *config.Config) *PostgresUserRepository {
	var preloads []gormdb.PreloadEntity = []gormdb.PreloadEntity{}
	return &PostgresUserRepository{BaseRepository: NewBaseRepository[model.User](cfg, preloads)}
}

func (r *PostgresUserRepository) CreateUser(ctx context.Context, u model.User) (model.User, error) {
	tx := r.database.WithContext(ctx).Begin()
	err := tx.Create(&u).Error
	if err != nil {
		tx.Rollback()
		r.logger.Error(logging.Postgres, logging.Rollback, err.Error(), nil)
		return u, err
	}
	tx.Commit()
	return u, nil
}


func (r *PostgresUserRepository) ExistsUserId(ctx context.Context, userId string) (bool, error) {
	var exists bool
	if err := r.database.WithContext(ctx).Model(&model.User{}).
		Select(countFilterExp).
		Where(userIdFilterExp, userId).
		Find(&exists).
		Error; err != nil {
		r.logger.Error(logging.Postgres, logging.Select, err.Error(), nil)
		return false, err
	}
	return exists, nil
}
