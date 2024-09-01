package entity

import "codebase-app/pkg/types"

type CategoryItem struct {
	Id   string `json:"id" db:"id"`
	Name string `json:"name" db:"name"`
}

type CategoriesResponse struct {
	Items []CategoryItem `json:"items"`
	Meta  types.Meta     `json:"meta"`
}

type CategoriesRequest struct {
	Page     int `query:"page" validate:"required"`
	Paginate int `query:"paginate" validate:"required"`
}

func (r *CategoriesRequest) SetDefault() {
	if r.Page < 1 {
		r.Page = 1
	}

	if r.Paginate < 1 {
		r.Paginate = 10
	}
}
