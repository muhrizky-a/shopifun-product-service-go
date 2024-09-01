package repository

import (
	"codebase-app/internal/module/category/entity"
	"codebase-app/internal/module/category/ports"
	"context"

	"github.com/jmoiron/sqlx"
	"github.com/rs/zerolog/log"
)

var _ ports.CategoryRepository = &categoryRepository{}

type categoryRepository struct {
	db *sqlx.DB
}

func NewCategoryRepository(db *sqlx.DB) *categoryRepository {
	return &categoryRepository{
		db: db,
	}
}

func (r *categoryRepository) GetCategories(ctx context.Context, req *entity.CategoriesRequest) (*entity.CategoriesResponse, error) {
	type dao struct {
		TotalData int `db:"total_data"`
		entity.CategoryItem
	}

	var (
		resp = new(entity.CategoriesResponse)
		data = make([]dao, 0, req.Paginate)
	)
	resp.Items = make([]entity.CategoryItem, 0, req.Paginate)

	query := `
		SELECT
			COUNT(id) OVER() as total_data,
			id,
			name
		FROM categories
		WHERE
			deleted_at IS NULL
	`

	err := r.db.SelectContext(ctx, &data, r.db.Rebind(query))
	if err != nil {
		log.Error().Err(err).Any("payload", req).Msg("repository::GetCategories - Failed to get categories")
		return nil, err
	}

	if len(data) > 0 {
		resp.Meta.TotalData = data[0].TotalData
	}

	for _, d := range data {
		resp.Items = append(resp.Items, d.CategoryItem)
	}

	resp.Meta.CountTotalPage(req.Page, req.Paginate, resp.Meta.TotalData)

	return resp, nil
}
