package repository

import (
	"codebase-app/internal/module/product/entity"
	"codebase-app/internal/module/product/ports"
	"codebase-app/pkg/errmsg"
	"context"
	"database/sql"
	"strings"

	"github.com/jmoiron/sqlx"
	"github.com/rs/zerolog/log"
)

var _ ports.ProductRepository = &productRepository{}

type productRepository struct {
	db *sqlx.DB
}

func NewProductRepository(db *sqlx.DB) *productRepository {
	return &productRepository{
		db: db,
	}
}

func (r *productRepository) CreateProduct(ctx context.Context, req *entity.CreateProductRequest) (*entity.CreateProductResponse, error) {
	var resp = new(entity.CreateProductResponse)
	// Your code here
	query := `
		INSERT INTO products (shop_id, category_id, name, description, price, stock)
		VALUES (?, ?, ?, ?, ?, ?) RETURNING id
	`

	err := r.db.QueryRowContext(ctx, r.db.Rebind(query),
		req.ShopId,
		req.CategoryId,
		req.Name,
		req.Description,
		req.Price,
		req.Stock,
	).Scan(&resp.Id)
	if err != nil {
		log.Error().Err(err).Any("payload", req).Msg("repository::CreateProduct - Failed to create product")
		return nil, err
	}

	return resp, nil
}

func (r *productRepository) GetProduct(ctx context.Context, req *entity.GetProductRequest) (*entity.GetProductResponse, error) {
	var (
		item = new(entity.GetProductItem)
		resp = new(entity.GetProductResponse)
	)

	// Your code here
	query := `
		SELECT
			p.name,
			p.description,
			p.price,
			p.stock,
			p.category_id,
			c.name AS category_name
		FROM products p
		LEFT JOIN
			categories c ON p.category_id = c.id
		WHERE
			p.deleted_at IS NULL
			AND p.id = ?
	`

	err := r.db.QueryRowxContext(ctx, r.db.Rebind(query), req.Id).StructScan(item)
	if err != nil {
		if err == sql.ErrNoRows {
			log.Error().Err(err).Any("payload", req).Msg("repository::GetProduct - Product not found")
			return nil, errmsg.NewCustomErrors(404, errmsg.WithMessage("Produk tidak ditemukan"))
		} else {
			log.Error().Err(err).Any("payload", req).Msg("repository::GetProduct - Failed to get product")
			return nil, err
		}
	}

	resp.Name = item.Name
	resp.Description = item.Description
	resp.Price = item.Price
	resp.Stock = item.Stock
	resp.Category.CategoryId = item.CategoryId
	resp.Category.CategoryName = item.CategoryName

	return resp, nil
}

func (r *productRepository) VerifyProductExists(ctx context.Context, req *entity.GetProductRequest) (*entity.GetExistingProductResponse, error) {
	var resp = new(entity.GetExistingProductResponse)

	// Your code here
	query := `
		SELECT
			p.id,
			p.shop_id,
			s.user_id
		FROM products p
		LEFT JOIN
			shops s ON p.shop_id = s.id
		WHERE
			p.deleted_at IS NULL
			AND s.deleted_at IS NULL
			AND p.id = ?
	`

	err := r.db.QueryRowxContext(ctx, r.db.Rebind(query), req.Id).StructScan(resp)
	if err != nil {
		if err == sql.ErrNoRows {
			log.Error().Err(err).Any("payload", req).Msg("repository::VerifyProductExists - Product not found")
			return nil, errmsg.NewCustomErrors(404, errmsg.WithMessage("Produk tidak ditemukan"))
		} else {
			log.Error().Err(err).Any("payload", req).Msg("repository::VerifyProductExists - Failed to get product")
			return nil, err
		}
	}

	return resp, nil
}

func (r *productRepository) DeleteProduct(ctx context.Context, req *entity.DeleteProductRequest) error {
	query := `
		UPDATE products
		SET deleted_at = NOW()
		WHERE id = ?
	`

	_, err := r.db.ExecContext(ctx, r.db.Rebind(query), req.Id)
	if err != nil {
		log.Error().Err(err).Any("payload", req).Msg("repository::DeleteProduct - Failed to delete product")
		return err
	}

	return nil
}

func (r *productRepository) UpdateProduct(ctx context.Context, req *entity.UpdateProductRequest) (*entity.UpdateProductResponse, error) {
	var resp = new(entity.UpdateProductResponse)

	query := `
		UPDATE products
		SET
			name = ?,
			description = ?,
			price = ?,
			stock = ?,
			category_id = ?,
			updated_at = NOW()
		WHERE
			deleted_at IS NULL
			AND id = ?
		RETURNING id
	`

	err := r.db.QueryRowxContext(ctx, r.db.Rebind(query),
		req.Name,
		req.Description,
		req.Price,
		req.Stock,
		req.CategoryId,
		req.Id).Scan(&resp.Id)
	if err != nil {
		if err == sql.ErrNoRows {
			log.Error().Err(err).Any("payload", req).Msg("repository::UpdateProduct - Product not found")
			return nil, errmsg.NewCustomErrors(404, errmsg.WithMessage("Produk tidak ditemukan"))
		} else {
			log.Error().Err(err).Any("payload", req).Msg("repository::UpdateProduct - Failed to update product")
			return nil, err
		}
	}

	return resp, nil
}

func (r *productRepository) GetProducts(ctx context.Context, req *entity.ProductsRequest) (*entity.ProductsResponse, error) {
	type dao struct {
		TotalData int `db:"total_data"`
		entity.ProductItem
	}

	var (
		resp = new(entity.ProductsResponse)
		data = make([]dao, 0, req.Paginate)
	)
	resp.Items = make([]entity.ProductItem, 0, req.Paginate)

	query := `
		SELECT
			COUNT(id) OVER() as total_data,
			id,
			name,
			price,
			stock
		FROM products
		WHERE
			deleted_at IS NULL
	`

	// Search and filter query
	queries := []interface{}{}

	/// Filter by Category Ids
	/// Example: category_ids=08362b22-f51d-40b1-a16b-49af90d561d9,3b4da768-e480-4cbb-b7fe-8b229123b50a
	if len(req.CategoryIds) > 0 {
		CategoryIdList := strings.Split(req.CategoryIds, ",")
		query += ` AND category_id IN (`

		for i, categoryId := range CategoryIdList {
			query += `?`
			if i < len(CategoryIdList)-1 {
				query += ` , `
			}

			queries = append(
				queries,
				categoryId,
			)
		}
		query += `)`
	}

	/// Filter by Price Range
	if req.MinPrice > 0 {
		query += ` AND price >= ?`
		queries = append(
			queries,
			req.MinPrice,
		)
	}
	if req.MaxPrice > req.MinPrice {
		query += ` AND price BETWEEN ? AND ?`
		queries = append(
			queries,
			req.MinPrice, req.MaxPrice,
		)
	}

	/// Filter by Keyword on either name or description
	if len(req.Keyword) > 0 {
		query += ` AND (
			name ILIKE ?
			OR
		 	description ILIKE ?
		)`
		queries = append(
			queries,
			"%"+req.Keyword+"%",
			"%"+req.Keyword+"%",
		)
	}

	// Pagination query
	query += ` LIMIT ? OFFSET ?`
	queries = append(
		queries,
		req.Paginate, req.Paginate*(req.Page-1),
	)

	err := r.db.SelectContext(ctx, &data, r.db.Rebind(query), queries...)
	if err != nil {
		log.Error().Err(err).Any("payload", req).Msg("repository::GetProducts - Failed to get products")
		return nil, err
	}

	if len(data) > 0 {
		resp.Meta.TotalData = data[0].TotalData
	}

	for _, d := range data {
		resp.Items = append(resp.Items, d.ProductItem)
	}

	resp.Meta.CountTotalPage(req.Page, req.Paginate, resp.Meta.TotalData)

	return resp, nil
}

func (r *productRepository) GetProductsByShopId(ctx context.Context, req *entity.ProductsByShopIdRequest) (*entity.ProductsResponse, error) {
	type dao struct {
		TotalData int `db:"total_data"`
		entity.ProductItem
	}

	var (
		resp = new(entity.ProductsResponse)
		data = make([]dao, 0, req.Paginate)
	)
	resp.Items = make([]entity.ProductItem, 0, req.Paginate)

	query := `
		SELECT
			COUNT(id) OVER() as total_data,
			id,
			name,
			price,
			stock
		FROM products
		WHERE
			deleted_at IS NULL
			AND shop_id = ?
	`

	// Search and filter query
	queries := []interface{}{req.ShopId}

	/// Filter by Category Ids
	/// Example: category_ids=08362b22-f51d-40b1-a16b-49af90d561d9,3b4da768-e480-4cbb-b7fe-8b229123b50a
	if len(req.CategoryIds) > 0 {
		CategoryIdList := strings.Split(req.CategoryIds, ",")
		query += ` AND category_id IN (`

		for i, categoryId := range CategoryIdList {
			query += `?`
			if i < len(CategoryIdList)-1 {
				query += ` , `
			}

			queries = append(
				queries,
				categoryId,
			)
		}
		query += `)`
	}

	/// Filter by Price Range
	if req.MinPrice > 0 {
		query += ` AND price >= ?`
		queries = append(
			queries,
			req.MinPrice,
		)
	}
	if req.MaxPrice > req.MinPrice {
		query += ` AND price BETWEEN ? AND ?`
		queries = append(
			queries,
			req.MinPrice, req.MaxPrice,
		)
	}

	/// Filter by Keyword on either name or description
	if len(req.Keyword) > 0 {
		query += ` AND (
			name ILIKE ?
			OR
		 	description ILIKE ?
		)`
		queries = append(
			queries,
			"%"+req.Keyword+"%",
			"%"+req.Keyword+"%",
		)
	}

	// Pagination query
	query += ` LIMIT ? OFFSET ?`
	queries = append(
		queries,
		req.Paginate, req.Paginate*(req.Page-1),
	)
	log.Debug().Any("req", req).Msg("repository::GetProductsByShopId - req")
	log.Debug().Any("query", query).Msg("repository::GetProductsByShopId - Failed to get products")
	log.Debug().Any("queries", queries).Msg("repository::GetProductsByShopId - Failed to get products")
	err := r.db.SelectContext(ctx, &data, r.db.Rebind(query), queries...)
	if err != nil {
		log.Error().Err(err).Any("payload", req).Msg("repository::GetProductsByShopId - Failed to get products")
		return nil, err
	}

	if len(data) > 0 {
		resp.Meta.TotalData = data[0].TotalData
	}

	for _, d := range data {
		resp.Items = append(resp.Items, d.ProductItem)
	}

	resp.Meta.CountTotalPage(req.Page, req.Paginate, resp.Meta.TotalData)

	return resp, nil
}
