package logic

import (
	"broabroad/internal/app/database"
	"context"
)

type StorageController struct {
	repo database.Repository
}

func NewStorageController(repository database.Repository) *StorageController {
	return &StorageController{
		repo: repository,
	}
}

func (sc *StorageController) CreateSeekRequest(ctx context.Context, sk database.SeekRequest) error {
	return sc.repo.CreateSeekRequest(ctx, sk)
}
