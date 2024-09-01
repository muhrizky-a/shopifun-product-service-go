package ports

import (
	"codebase-app/internal/module/product/entity"
	"context"
)

type ProductRepository interface {
	CreateProduct(ctx context.Context, req *entity.CreateProductRequest) (*entity.CreateProductResponse, error)
	GetProduct(ctx context.Context, req *entity.GetProductRequest) (*entity.GetProductResponse, error)
	DeleteProduct(ctx context.Context, req *entity.DeleteProductRequest) error
	UpdateProduct(ctx context.Context, req *entity.UpdateProductRequest) (*entity.UpdateProductResponse, error)
	GetProducts(ctx context.Context, req *entity.ProductsRequest) (*entity.ProductsResponse, error)
}

type ProductService interface {
	CreateProduct(ctx context.Context, req *entity.CreateProductRequest) (*entity.CreateProductResponse, error)
	GetProduct(ctx context.Context, req *entity.GetProductRequest) (*entity.GetProductResponse, error)
	DeleteProduct(ctx context.Context, req *entity.DeleteProductRequest) error
	UpdateProduct(ctx context.Context, req *entity.UpdateProductRequest) (*entity.UpdateProductResponse, error)
	GetProducts(ctx context.Context, req *entity.ProductsRequest) (*entity.ProductsResponse, error)
}
