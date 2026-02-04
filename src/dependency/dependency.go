package dependency

import (
	"github.com/minisource/template_go/config"
	"github.com/minisource/template_go/domain/model"
	contractRepository "github.com/minisource/template_go/domain/repository"
	infrarepository "github.com/minisource/template_go/infra/persistence/repository"
	gormdb "github.com/minisource/go-common/db/gorm"
)

// func GetUserRepository(cfg *config.Config) contractRepository.UserRepository {
// 	return infraRepository.NewUserRepository(cfg)
// }

func GetFileRepository(cfg *config.Config) contractRepository.FileRepository {
	var preloads []gormdb.PreloadEntity = []gormdb.PreloadEntity{}
	return infrarepository.NewBaseRepository[model.File](cfg, preloads)
}

func GetUserRepository(cfg *config.Config) contractRepository.UserRepository {
	return infrarepository.NewUserRepository(cfg)
}
