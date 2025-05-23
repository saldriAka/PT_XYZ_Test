package api

import (
	"context"
	"time"

	"saldri/test_pt_xyz/domain"
	"saldri/test_pt_xyz/dto"
	"saldri/test_pt_xyz/internal/util"

	"github.com/gofiber/fiber/v2"
)

type transactionsApi struct {
	transactionsService domain.TransactionsService
}

func NewTransactions(app *fiber.App, transactionsService domain.TransactionsService, authMid fiber.Handler) {
	tc := &transactionsApi{
		transactionsService: transactionsService,
	}

	transactions := app.Group("/api", authMid)
	transactions.Get("/transactions", tc.Index)
	transactions.Post("/transactions", tc.Create)
	transactions.Get("/transactions/:id", tc.Show)
	transactions.Get("/transactions/customer/:customer_id", tc.CustomerShow)
	transactions.Put("/transactions/:id", tc.Update)
	transactions.Delete("/transactions/:id", tc.Delete)
}

func (tc transactionsApi) Index(ctx *fiber.Ctx) error {
	c, cancel := context.WithTimeout(ctx.Context(), 10*time.Second)
	defer cancel()

	page, limit, fails := util.SafePaginationParams(ctx)
	if len(fails) > 0 {
		return ctx.Status(fiber.StatusBadRequest).
			JSON(dto.CreateResponseErrorData("validasi query gagal", fails))
	}

	res, total, err := tc.transactionsService.Index(c, page, limit)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).
			JSON(dto.CreateResponseError(err.Error()))
	}

	return dto.PaginateAndRespond(ctx, res, int(total))
}

func (tc transactionsApi) Create(ctx *fiber.Ctx) error {
	c, cancel := context.WithTimeout(ctx.Context(), 10*time.Second)
	defer cancel()

	var req dto.CreateTransactionsRequest
	if err := ctx.BodyParser(&req); err != nil {
		return ctx.SendStatus(fiber.StatusUnprocessableEntity)
	}

	fails := util.Validate(req)
	if len(fails) > 0 {
		return ctx.Status(fiber.StatusBadRequest).
			JSON(dto.CreateResponseErrorData("validasi gagal", fails))
	}

	err := tc.transactionsService.Create(c, req)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).
			JSON(dto.CreateResponseError(err.Error()))
	}

	return ctx.Status(fiber.StatusCreated).
		JSON(dto.CreateResponseSuccess("data transaksi berhasil dibuat"))
}

func (tc transactionsApi) Show(ctx *fiber.Ctx) error {
	c, cancel := context.WithTimeout(ctx.Context(), 10*time.Second)
	defer cancel()

	id := ctx.Params("id")
	res, err := tc.transactionsService.Show(c, id)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).
			JSON(dto.CreateResponseError(err.Error()))
	}

	return ctx.JSON(dto.CreateResponseSuccess(res))
}

func (tc transactionsApi) CustomerShow(ctx *fiber.Ctx) error {
	c, cancel := context.WithTimeout(ctx.Context(), 10*time.Second)
	defer cancel()

	id := ctx.Params("customer_id")
	res, err := tc.transactionsService.CustomerShow(c, id)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).
			JSON(dto.CreateResponseError(err.Error()))
	}

	return ctx.JSON(dto.CreateResponseSuccess(res))
}

func (tc transactionsApi) Update(ctx *fiber.Ctx) error {
	c, cancel := context.WithTimeout(ctx.Context(), 10*time.Second)
	defer cancel()

	var req dto.UpdateTransactionsRequest
	if err := ctx.BodyParser(&req); err != nil {
		return ctx.SendStatus(fiber.StatusUnprocessableEntity)
	}

	req.ID = ctx.Params("id")

	fails := util.Validate(req)

	if len(fails) > 0 {
		return ctx.Status(fiber.StatusBadRequest).
			JSON(dto.CreateResponseError("validasi gagal"))
	}

	err := tc.transactionsService.Update(c, req)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).
			JSON(dto.CreateResponseError(err.Error()))
	}

	return ctx.JSON(dto.CreateResponseSuccess("data transaksi berhasil diperbarui"))
}

func (tc transactionsApi) Delete(ctx *fiber.Ctx) error {
	c, cancel := context.WithTimeout(ctx.Context(), 10*time.Second)
	defer cancel()

	id := ctx.Params("id")
	err := tc.transactionsService.Delete(c, id)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).
			JSON(dto.CreateResponseError(err.Error()))
	}

	return ctx.JSON(dto.CreateResponseSuccess("data transaksi berhasil dihapus"))
}
