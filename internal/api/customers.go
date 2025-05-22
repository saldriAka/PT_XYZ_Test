package api

import (
	"context"
	"time"

	"saldri/test_pt_xyz/domain"
	"saldri/test_pt_xyz/dto"
	"saldri/test_pt_xyz/internal/config"
	"saldri/test_pt_xyz/internal/util"

	"github.com/gofiber/fiber/v2"
)

type customersApi struct {
	conf             *config.Config
	customersService domain.CustomersService
}

func NewCustomers(app *fiber.App, conf *config.Config, customersService domain.CustomersService, authMid fiber.Handler) {
	ca := &customersApi{
		conf:             conf,
		customersService: customersService,
	}

	customers := app.Group("/api", authMid)
	customers.Get("/customers", ca.Index)
	customers.Post("/customers", ca.Create)
	customers.Get("/customers/:id", ca.Show)
	customers.Put("/:id", ca.Update)
	customers.Put("/customers/assets/:id", ca.UpdateAssets)
	customers.Delete("/customers/:id", ca.Delete)
	app.Static("/media", conf.Storage.BasePath)
}

func (ca customersApi) Index(ctx *fiber.Ctx) error {
	c, cancel := context.WithTimeout(ctx.Context(), 10*time.Second)
	defer cancel()

	page, limit, fails := util.SafePaginationParams(ctx)
	if len(fails) > 0 {
		return ctx.Status(fiber.StatusBadRequest).
			JSON(dto.CreateResponseErrorData("validasi query gagal", fails))
	}
	res, total, err := ca.customersService.Index(c, page, limit)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).
			JSON(dto.CreateResponseError(err.Error()))
	}

	return dto.PaginateAndRespond(ctx, res, int(total))
}

func (ca customersApi) Create(ctx *fiber.Ctx) error {
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

func (ca customersApi) Show(ctx *fiber.Ctx) error {
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

func (ca customersApi) Update(ctx *fiber.Ctx) error {
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

func (ca customersApi) Delete(ctx *fiber.Ctx) error {
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

func (ca customersApi) UpdateAssets(ctx *fiber.Ctx) error {
	c, cancel := context.WithTimeout(ctx.Context(), 10*time.Second)
	defer cancel()

	id := ctx.Params("id")
	if id == "" {
		return ctx.Status(fiber.StatusBadRequest).
			JSON(dto.CreateResponseError("id customer tidak boleh kosong"))
	}

	ktpURL, err := util.ProcessAndSaveImageFile(ctx, util.ImageSaveOptions{
		FieldName: "ktp_photo",
		BasePath:  ca.conf.Storage.BasePath,
		PublicURL: ca.conf.Server.Assets,
		MaxSizeMB: 2,
	})
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).
			JSON(dto.CreateResponseError("gagal upload foto KTP: " + err.Error()))
	}

	selfieURL, err := util.ProcessAndSaveImageFile(ctx, util.ImageSaveOptions{
		FieldName: "selfie_photo",
		BasePath:  ca.conf.Storage.BasePath,
		PublicURL: ca.conf.Server.Assets,
		MaxSizeMB: 2,
	})
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).
			JSON(dto.CreateResponseError("gagal upload foto selfie: " + err.Error()))
	}

	// Simpan ke DB
	req := dto.UpdateAssetsCustomersRequest{
		KTPPhotoURL:    ktpURL,
		SelfiePhotoURL: selfieURL,
	}
	if err := ca.customersService.UpdateAssets(c, id, req); err != nil {
		return ctx.Status(fiber.StatusInternalServerError).
			JSON(dto.CreateResponseError("gagal mengupdate foto: " + err.Error()))
	}

	return ctx.Status(fiber.StatusOK).
		JSON(dto.CreateResponseSuccess("berhasil mengupdate foto customer"))
}
