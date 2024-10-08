package handler

import (
	"codebase-app/internal/adapter"
	"codebase-app/internal/middleware"
	"codebase-app/internal/module/product/entity"
	"codebase-app/internal/module/product/ports"
	"codebase-app/internal/module/product/repository"
	"codebase-app/internal/module/product/service"
	"codebase-app/pkg/errmsg"
	"codebase-app/pkg/response"

	"github.com/gofiber/fiber/v2"
	"github.com/rs/zerolog/log"
)

type productHandler struct {
	service ports.ProductService
}

func NewProductHandler() *productHandler {
	var (
		handler = new(productHandler)
		repo    = repository.NewProductRepository(adapter.Adapters.ShopeefunPostgres)
		service = service.NewProductService(repo)
	)
	handler.service = service

	return handler
}

func (h *productHandler) Register(router fiber.Router) {
	router.Get("/products", h.GetProducts)
	router.Get("/shops/:shop_id/products", h.GetProductsByShopId)
	router.Post("/products", middleware.UserIdHeader, h.CreateProduct)
	router.Get("/products/:id", h.GetProduct)
	router.Delete("/products/:id", middleware.UserIdHeader, h.DeleteProduct)
	router.Patch("/products/:id", middleware.UserIdHeader, h.UpdateProduct)
}

func (h *productHandler) CreateProduct(c *fiber.Ctx) error {
	var (
		req = new(entity.CreateProductRequest)
		ctx = c.Context()
		v   = adapter.Adapters.Validator
	)

	if err := c.BodyParser(req); err != nil {
		log.Warn().Err(err).Msg("handler::CreateProduct - Parse request body")
		return c.Status(fiber.StatusBadRequest).JSON(response.Error(err))
	}

	if err := v.Validate(req); err != nil {
		log.Warn().Err(err).Any("payload", req).Msg("handler::CreateProduct - Validate request body")
		code, errs := errmsg.Errors(err, req)
		return c.Status(code).JSON(response.Error(errs))
	}

	resp, err := h.service.CreateProduct(ctx, req)
	if err != nil {
		code, errs := errmsg.Errors[error](err)
		return c.Status(code).JSON(response.Error(errs))
	}

	return c.Status(fiber.StatusCreated).JSON(response.Success(resp, ""))

}

func (h *productHandler) GetProduct(c *fiber.Ctx) error {
	var (
		req = new(entity.GetProductRequest)
		ctx = c.Context()
		v   = adapter.Adapters.Validator
	)

	req.Id = c.Params("id")

	if err := v.Validate(req); err != nil {
		log.Warn().Err(err).Any("payload", req).Msg("handler::GetProduct - Validate request body")
		code, errs := errmsg.Errors(err, req)
		return c.Status(code).JSON(response.Error(errs))
	}

	resp, err := h.service.GetProduct(ctx, req)
	if err != nil {
		code, errs := errmsg.Errors[error](err)
		return c.Status(code).JSON(response.Error(errs))
	}

	return c.Status(fiber.StatusOK).JSON(response.Success(resp, ""))
}

func (h *productHandler) DeleteProduct(c *fiber.Ctx) error {
	var (
		req           = new(entity.DeleteProductRequest)
		reqGetProduct = new(entity.GetProductRequest)
		ctx           = c.Context()
		v             = adapter.Adapters.Validator
		l             = middleware.GetLocals(c)
	)

	req.Id = c.Params("id")
	OwnerId := l.UserId

	if err := v.Validate(req); err != nil {
		log.Warn().Err(err).Any("payload", req).Msg("handler::DeleteProduct - Validate request body")
		code, errs := errmsg.Errors(err, req)
		return c.Status(code).JSON(response.Error(errs))
	}

	reqGetProduct.Id = req.Id
	respExistingProduct, err := h.service.VerifyProductExists(ctx, reqGetProduct)
	if err != nil {
		code, errs := errmsg.Errors[error](err)
		return c.Status(code).JSON(response.Error(errs))
	}

	if respExistingProduct.UserId != OwnerId {
		log.Warn().Err(err).Msg("handler::DeleteProduct - Unauthorized")
		return c.Status(403).JSON(response.Error(
			"Terlarang: anda tidak diizinkan untuk mengakses resource ini",
		))
	}

	err = h.service.DeleteProduct(ctx, req)
	if err != nil {
		code, errs := errmsg.Errors[error](err)
		return c.Status(code).JSON(response.Error(errs))
	}

	return c.Status(fiber.StatusOK).JSON(response.Success(nil, ""))
}

func (h *productHandler) UpdateProduct(c *fiber.Ctx) error {
	var (
		req           = new(entity.UpdateProductRequest)
		reqGetProduct = new(entity.GetProductRequest)
		ctx           = c.Context()
		v             = adapter.Adapters.Validator
		l             = middleware.GetLocals(c)
	)

	if err := c.BodyParser(req); err != nil {
		log.Warn().Err(err).Msg("handler::UpdateProduct - Parse request body")
		return c.Status(fiber.StatusBadRequest).JSON(response.Error(err))
	}

	req.Id = c.Params("id")
	OwnerId := l.UserId

	if err := v.Validate(req); err != nil {
		log.Warn().Err(err).Any("payload", req).Msg("handler::UpdateProduct - Validate request body")
		code, errs := errmsg.Errors(err, req)
		return c.Status(code).JSON(response.Error(errs))
	}

	reqGetProduct.Id = req.Id
	respExistingProduct, err := h.service.VerifyProductExists(ctx, reqGetProduct)
	if err != nil {
		code, errs := errmsg.Errors[error](err)
		return c.Status(code).JSON(response.Error(errs))
	}

	if respExistingProduct.UserId != OwnerId {
		log.Warn().Err(err).Msg("handler::UpdateProduct - Unauthorized")
		return c.Status(403).JSON(response.Error(
			"Terlarang: anda tidak diizinkan untuk mengakses resource ini",
		))
	}

	resp, err := h.service.UpdateProduct(ctx, req)
	if err != nil {
		code, errs := errmsg.Errors[error](err)
		return c.Status(code).JSON(response.Error(errs))
	}

	return c.Status(fiber.StatusOK).JSON(response.Success(resp, ""))
}

func (h *productHandler) GetProducts(c *fiber.Ctx) error {
	var (
		req = new(entity.ProductsRequest)
		ctx = c.Context()
		v   = adapter.Adapters.Validator
	)

	if err := c.QueryParser(req); err != nil {
		log.Warn().Err(err).Msg("handler::GetProducts - Parse request query")
		return c.Status(fiber.StatusBadRequest).JSON(response.Error(err))
	}

	req.SetDefault()

	if err := v.Validate(req); err != nil {
		log.Warn().Err(err).Any("payload", req).Msg("handler::GetProducts - Validate request body")
		code, errs := errmsg.Errors(err, req)
		return c.Status(code).JSON(response.Error(errs))
	}

	resp, err := h.service.GetProducts(ctx, req)
	if err != nil {
		code, errs := errmsg.Errors[error](err)
		return c.Status(code).JSON(response.Error(errs))
	}

	return c.Status(fiber.StatusOK).JSON(response.Success(resp, ""))
}

func (h *productHandler) GetProductsByShopId(c *fiber.Ctx) error {
	var (
		req = new(entity.ProductsByShopIdRequest)
		ctx = c.Context()
		v   = adapter.Adapters.Validator
	)

	if err := c.QueryParser(req); err != nil {
		log.Warn().Err(err).Msg("handler::GetProductsByShopId - Parse request query")
		return c.Status(fiber.StatusBadRequest).JSON(response.Error(err))
	}

	req.ShopId = c.Params("shop_id")
	req.SetDefault()

	if err := v.Validate(req); err != nil {
		log.Warn().Err(err).Any("payload", req).Msg("handler::GetProductsByShopId - Validate request body")
		code, errs := errmsg.Errors(err, req)
		return c.Status(code).JSON(response.Error(errs))
	}

	resp, err := h.service.GetProductsByShopId(ctx, req)
	if err != nil {
		code, errs := errmsg.Errors[error](err)
		return c.Status(code).JSON(response.Error(errs))
	}

	return c.Status(fiber.StatusOK).JSON(response.Success(resp, ""))
}
