package handler

import (
	"codebase-app/internal/adapter"
	"codebase-app/internal/module/category/entity"
	"codebase-app/internal/module/category/ports"
	"codebase-app/internal/module/category/repository"
	"codebase-app/internal/module/category/service"
	"codebase-app/pkg/errmsg"
	"codebase-app/pkg/response"

	"github.com/gofiber/fiber/v2"
	"github.com/rs/zerolog/log"
)

type categoryHandler struct {
	service ports.CategoryService
}

func NewCategoryHandler() *categoryHandler {
	var (
		handler = new(categoryHandler)
		repo    = repository.NewCategoryRepository(adapter.Adapters.ShopeefunPostgres)
		service = service.NewCategoryService(repo)
	)
	handler.service = service

	return handler
}

func (h *categoryHandler) Register(router fiber.Router) {
	router.Get("/categories", h.GetCategories)
}

func (h *categoryHandler) GetCategories(c *fiber.Ctx) error {
	var (
		req = new(entity.CategoriesRequest)
		ctx = c.Context()
		v   = adapter.Adapters.Validator
	)

	if err := c.QueryParser(req); err != nil {
		log.Warn().Err(err).Msg("handler::GetCategories - Parse request query")
		return c.Status(fiber.StatusBadRequest).JSON(response.Error(err))
	}

	req.SetDefault()

	if err := v.Validate(req); err != nil {
		log.Warn().Err(err).Any("payload", req).Msg("handler::GetCategories - Validate request body")
		code, errs := errmsg.Errors(err, req)
		return c.Status(code).JSON(response.Error(errs))
	}

	resp, err := h.service.GetCategories(ctx, req)
	if err != nil {
		code, errs := errmsg.Errors[error](err)
		return c.Status(code).JSON(response.Error(errs))
	}

	return c.Status(fiber.StatusOK).JSON(response.Success(resp, ""))

}
