package ports

import (
	"codebase-app/internal/module/category/entity"
	"context"
)

type CategoryRepository interface {
	GetCategories(ctx context.Context, req *entity.CategoriesRequest) (*entity.CategoriesResponse, error)
}

type CategoryService interface {
	GetCategories(ctx context.Context, req *entity.CategoriesRequest) (*entity.CategoriesResponse, error)
}
