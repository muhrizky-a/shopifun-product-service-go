package entity

import "codebase-app/pkg/types"

type CreateProductRequest struct {
	ShopId     string `json:"shop_id" validate:"uuid" db:"shop_id"`
	CategoryId string `json:"category_id" validate:"uuid" db:"category_id"`

	Name        string `json:"name" validate:"required" db:"name"`
	Description string `json:"description" validate:"required,max=255" db:"description"`
	Price       int    `json:"price" validate:"required,gte=0" db:"price"`
	Stock       int    `json:"stock" validate:"required,gte=0" db:"stock"`
}

type CreateProductResponse struct {
	Id string `json:"id" db:"id"`
}

type GetProductRequest struct {
	Id string `validate:"uuid" db:"id"`
}

type GetExistingProductResponse struct {
	Id     string `json:"id" db:"id"`
	UserId string `json:"user_id" db:"user_id"`
	ShopId string `json:"shop_id" db:"shop_id"`
}

type GetProductItem struct {
	Name         string `json:"name" db:"name"`
	Description  string `json:"description" db:"description"`
	Price        int    `json:"price" validate:"required" db:"price"`
	Stock        int    `json:"stock" validate:"required" db:"stock"`
	CategoryId   string `json:"category_id" validate:"required" db:"category_id"`
	CategoryName string `json:"category_name" validate:"required" db:"category_name"`
}

type CategoryItem struct {
	CategoryId   string `json:"id" validate:"required" db:"category_id"`
	CategoryName string `json:"name" validate:"required" db:"category_name"`
}

type GetProductResponse struct {
	Name        string       `json:"name" db:"name"`
	Description string       `json:"description" db:"description"`
	Price       int          `json:"price" validate:"required" db:"price"`
	Stock       int          `json:"stock" validate:"required" db:"stock"`
	Category    CategoryItem `json:"category"`
}

type DeleteProductRequest struct {
	Id string `validate:"uuid" db:"id"`
}

type UpdateProductRequest struct {
	CategoryId string `json:"category_id" validate:"uuid" db:"category_id"`

	Id          string `params:"id" validate:"uuid" db:"id"`
	Name        string `json:"name" validate:"required" db:"name"`
	Description string `json:"description" validate:"required" db:"description"`
	Price       int    `json:"price" validate:"required" db:"price"`
	Stock       int    `json:"stock" validate:"required" db:"stock"`
}

type UpdateProductResponse struct {
	Id string `json:"id" db:"id"`
}

type ProductsRequest struct {
	ShopId   string `validate:"uuid" db:"shop_id"`
	Page     int    `query:"page" validate:"required"`
	Paginate int    `query:"paginate" validate:"required"`
}

func (r *ProductsRequest) SetDefault() {
	if r.Page < 1 {
		r.Page = 1
	}

	if r.Paginate < 1 {
		r.Paginate = 10
	}
}

type ProductItem struct {
	Id    string `json:"id" db:"id"`
	Name  string `json:"name" db:"name"`
	Price int    `json:"price" validate:"required" db:"price"`
	Stock int    `json:"stock" validate:"required" db:"stock"`
}

type ProductsResponse struct {
	Items []ProductItem `json:"items"`
	Meta  types.Meta    `json:"meta"`
}
