package service

import (
	"codebase-app/internal/module/category/entity"
	"codebase-app/internal/module/category/ports"
	"context"
)

var _ ports.CategoryService = &categoryService{}

type categoryService struct {
	repo ports.CategoryRepository
}

func NewCategoryService(repo ports.CategoryRepository) *categoryService {
	return &categoryService{
		repo: repo,
	}
}

func (s *categoryService) GetCategories(ctx context.Context, req *entity.CategoriesRequest) (*entity.CategoriesResponse, error) {
	return s.repo.GetCategories(ctx, req)
}
