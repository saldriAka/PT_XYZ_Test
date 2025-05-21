package controller

import (
	"context"
	"time"

	"saldri/test_pt_xyz/domain"
	"saldri/test_pt_xyz/dto"
	"saldri/test_pt_xyz/internal/util"

	"github.com/gofiber/fiber/v2"
)

type customersController struct {
	customersService domain.CustomersService
}

func NewCustomers(app *fiber.App, customersService domain.CustomersService) {
	ca := &customersController{
		customersService: customersService,
	}

	customers := app.Group("/customers")
	customers.Get("/", ca.Index)
	customers.Post("/", ca.Create)
	customers.Get("/:id", ca.Show)
	customers.Put("/:id", ca.Update)
	customers.Delete("/:id", ca.Delete)
}

func (ca customersController) Index(ctx *fiber.Ctx) error {
	c, cancel := context.WithTimeout(ctx.Context(), 10*time.Second)
	defer cancel()

	res, err := ca.customersService.Index(c)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).
			JSON(dto.CreateResponseError(err.Error()))
	}

	return ctx.JSON(dto.CreateResponseSuccess(res))
}

func (ca customersController) Create(ctx *fiber.Ctx) error {
	c, cancel := context.WithTimeout(ctx.Context(), 10*time.Second)
	defer cancel()

	var req dto.CreateCustomersRequest
	if err := ctx.BodyParser(&req); err != nil {
		return ctx.SendStatus(fiber.StatusUnprocessableEntity)
	}

	fails := util.Validate(req)
	if len(fails) > 0 {
		return ctx.Status(fiber.StatusBadRequest).
			JSON(dto.CreateResponseErrorData("validasi gagal", fails))
	}

	err := ca.customersService.Create(c, req)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).
			JSON(dto.CreateResponseError(err.Error()))
	}

	return ctx.Status(fiber.StatusCreated).
		JSON(dto.CreateResponseSuccess("data customer berhasil dibuat"))
}

func (ca customersController) Show(ctx *fiber.Ctx) error {
	c, cancel := context.WithTimeout(ctx.Context(), 10*time.Second)
	defer cancel()

	id := ctx.Params("id")
	res, err := ca.customersService.Show(c, id)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).
			JSON(dto.CreateResponseError(err.Error()))
	}

	return ctx.JSON(dto.CreateResponseSuccess(res))
}

func (ca customersController) Update(ctx *fiber.Ctx) error {
	c, cancel := context.WithTimeout(ctx.Context(), 10*time.Second)
	defer cancel()

	var req dto.UpdateCustomersRequest
	if err := ctx.BodyParser(&req); err != nil {
		return ctx.SendStatus(fiber.StatusUnprocessableEntity)
	}

	req.ID = ctx.Params("id")

	fails := util.Validate(req)
	if len(fails) > 0 {
		return ctx.Status(fiber.StatusBadRequest).
			JSON(dto.CreateResponseError("validasi gagal"))
	}

	err := ca.customersService.Update(c, req)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).
			JSON(dto.CreateResponseError(err.Error()))
	}

	return ctx.JSON(dto.CreateResponseSuccess("data customer berhasil diperbarui"))
}

func (ca customersController) Delete(ctx *fiber.Ctx) error {
	c, cancel := context.WithTimeout(ctx.Context(), 10*time.Second)
	defer cancel()

	id := ctx.Params("id")
	err := ca.customersService.Delete(c, id)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).
			JSON(dto.CreateResponseError(err.Error()))
	}

	return ctx.JSON(dto.CreateResponseSuccess("data customer berhasil dihapus"))
}
