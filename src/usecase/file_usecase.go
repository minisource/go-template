package usecase

import (
	"context"

	"github.com/minisource/template_go/config"
	"github.com/minisource/template_go/domain/model"
	"github.com/minisource/template_go/domain/repository"
	"github.com/minisource/template_go/usecase/dto"
	"github.com/minisource/go-common/filter"
)

type FileUsecase struct {
	base *BaseUsecase[model.File, dto.CreateFile, dto.UpdateFile, dto.File]
}

func NewFileUsecase(cfg *config.Config, repository repository.FileRepository) *FileUsecase {
	return &FileUsecase{
		base: NewBaseUsecase[model.File, dto.CreateFile, dto.UpdateFile, dto.File](cfg, repository),
	}
}

// Create
func (u *FileUsecase) Create(ctx context.Context, req dto.CreateFile) (dto.File, error) {
	return u.base.Create(ctx, req)
}

// Update
func (u *FileUsecase) Update(ctx context.Context, id int, req dto.UpdateFile) (dto.File, error) {
	return u.base.Update(ctx, id, req)
}

// Delete
func (u *FileUsecase) Delete(ctx context.Context, id int) error {
	return u.base.Delete(ctx, id)
}

// Get By Id
func (u *FileUsecase) GetById(ctx context.Context, id int) (dto.File, error) {
	return u.base.GetById(ctx, id)
}

// Get By Filter
func (u *FileUsecase) GetByFilter(ctx context.Context, req filter.PaginationInputWithFilter) (*filter.PagedList[dto.File], error) {
	return u.base.GetByFilter(ctx, req)
}
