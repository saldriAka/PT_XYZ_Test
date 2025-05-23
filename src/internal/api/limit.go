package api

import (
	"context"
	"time"

	"saldri/test_pt_xyz/domain"
	"saldri/test_pt_xyz/dto"
	"saldri/test_pt_xyz/internal/util"

	"github.com/gofiber/fiber/v2"
)

type limitApi struct {
	limitService domain.LimitService
}

func NewLimit(app *fiber.App, limitService domain.LimitService, authMid fiber.Handler) {
	lmt := &limitApi{
		limitService: limitService,
	}

	limit := app.Group("/api", authMid)
	limit.Get("/limit", lmt.Index)
	limit.Get("/limit/:id", lmt.Show)
	limit.Post("/limit", lmt.Create)
	limit.Put("/limit/:id", lmt.Update)
	limit.Delete("/limit/:id", lmt.Delete)
}

func (lmt *limitApi) Index(ctx *fiber.Ctx) error {
	c, cancel := context.WithTimeout(ctx.Context(), 10*time.Second)
	defer cancel()

	page, limit, fails := util.SafePaginationParams(ctx)
	if len(fails) > 0 {
		return ctx.Status(fiber.StatusBadRequest).
			JSON(dto.CreateResponseErrorData("validasi query gagal", fails))
	}
	res, total, err := lmt.limitService.Index(c, page, limit)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).
			JSON(dto.CreateResponseError(err.Error()))
	}

	return dto.PaginateAndRespond(ctx, res, int(total))
}

func (lmt *limitApi) Show(ctx *fiber.Ctx) error {
	c, cancel := context.WithTimeout(ctx.Context(), 10*time.Second)
	defer cancel()

	id := ctx.Params("id")
	res, err := lmt.limitService.Show(c, id)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).
			JSON(dto.CreateResponseError(err.Error()))
	}

	return ctx.JSON(dto.CreateResponseSuccess(res))
}

func (lmt *limitApi) Create(ctx *fiber.Ctx) error {
	c, cancel := context.WithTimeout(ctx.Context(), 10*time.Second)
	defer cancel()

	var req dto.CreateLimitRequest
	if err := ctx.BodyParser(&req); err != nil {
		return ctx.SendStatus(fiber.StatusUnprocessableEntity)
	}

	fails := util.Validate(req)
	if len(fails) > 0 {
		return ctx.Status(fiber.StatusBadRequest).
			JSON(dto.CreateResponseErrorData("validasi gagal", fails))
	}

	err := lmt.limitService.Create(c, req)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).
			JSON(dto.CreateResponseError(err.Error()))
	}

	return ctx.Status(fiber.StatusCreated).
		JSON(dto.CreateResponseSuccess("data limit berhasil dibuat"))
}

func (lmt *limitApi) Update(ctx *fiber.Ctx) error {
	c, cancel := context.WithTimeout(ctx.Context(), 10*time.Second)
	defer cancel()

	var req dto.UpdateLimitRequest
	if err := ctx.BodyParser(&req); err != nil {
		return ctx.SendStatus(fiber.StatusUnprocessableEntity)
	}

	req.ID = ctx.Params("id")

	fails := util.Validate(req)
	if len(fails) > 0 {
		return ctx.Status(fiber.StatusBadRequest).
			JSON(dto.CreateResponseErrorData("validasi gagal", fails))
	}

	err := lmt.limitService.Update(c, req)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).
			JSON(dto.CreateResponseError(err.Error()))
	}

	return ctx.JSON(dto.CreateResponseSuccess("data limit berhasil diperbarui"))
}

func (lmt *limitApi) Delete(ctx *fiber.Ctx) error {
	c, cancel := context.WithTimeout(ctx.Context(), 10*time.Second)
	defer cancel()

	id := ctx.Params("id")
	err := lmt.limitService.Delete(c, id)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).
			JSON(dto.CreateResponseError(err.Error()))
	}

	return ctx.JSON(dto.CreateResponseSuccess("data limit berhasil dihapus"))
}
