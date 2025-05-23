package web

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"saldri/test_pt_xyz/domain"
	"saldri/test_pt_xyz/dto"
	"saldri/test_pt_xyz/internal/config"
	"saldri/test_pt_xyz/internal/util"
	"strconv"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/session"
)

type xyzWeb struct {
	conf                *config.Config
	customersService    domain.CustomersService
	limitService        domain.LimitService
	transactionsService domain.TransactionsService
	authService         domain.AuthService
	store               *session.Store
}

func NewWeb(app *fiber.App, customersService domain.CustomersService, limitService domain.LimitService, transactionsService domain.TransactionsService, authService domain.AuthService, conf *config.Config, store *session.Store) {
	web := &xyzWeb{
		conf:                conf,
		customersService:    customersService,
		limitService:        limitService,
		transactionsService: transactionsService,
		authService:         authService,
		store:               store,
	}

	app.Static("/assets", "./assets")
	app.Static("/storage", "./storage")

	app.Get("/login", web.Login)
	app.Post("/login", web.HandleLogin)
	app.Get("/logout", web.LogoutWeb)

	webGroup := app.Group("/", util.AuthRequired(store))
	webGroup.Get("/", web.Index)
	webGroup.Get("/profile/:id", web.EditProfile)
	webGroup.Post("/profile", web.UpdateProfile)
	webGroup.Get("/credits", web.Credits)
	webGroup.Get("/transaction", web.Transactions)
	webGroup.Post("/transaction", web.CreateTransactions)

}

func (xy xyzWeb) Index(ctx *fiber.Ctx) error {
	c, cancel := context.WithTimeout(ctx.Context(), 10*time.Second)
	defer cancel()

	sess, err := xy.store.Get(ctx)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).SendString("Gagal mengambil session")
	}

	id := sess.Get("user_id")
	if id == nil {
		return ctx.Redirect("/login")
	}

	userID, ok := id.(string)
	if !ok {
		return ctx.Status(fiber.StatusInternalServerError).SendString("user_id tidak valid")
	}

	res, err := xy.customersService.Show(c, userID)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).
			JSON(dto.CreateResponseError(err.Error()))
	}

	formattedSalary := util.FormatRupiah(res.Salary)
	formattedDOB, err := util.FormatTanggalIndo(res.DateOfBirth)
	if err != nil {
		formattedDOB = res.DateOfBirth
	}

	return ctx.Render("profile", fiber.Map{
		"Customer": fiber.Map{
			"ID":             userID,
			"NIK":            res.NIK,
			"FullName":       res.FullName,
			"LegalName":      res.LegalName,
			"PlaceOfBirth":   res.PlaceOfBirth,
			"DateOfBirth":    formattedDOB,
			"Salary":         formattedSalary,
			"KTPPhotoURL":    res.KTPPhotoURL,
			"SelfiePhotoURL": res.SelfiePhotoURL,
		},
	}, "layout/main")

}

func (xy xyzWeb) EditProfile(ctx *fiber.Ctx) error {
	c, cancel := context.WithTimeout(ctx.Context(), 10*time.Second)
	defer cancel()

	id := ctx.Params("id")
	res, err := xy.customersService.Show(c, id)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).
			JSON(dto.CreateResponseError(err.Error()))
	}

	formattedSalary := util.FormatRupiah(res.Salary)

	return ctx.Render("profile_edit", fiber.Map{
		"Customer": fiber.Map{
			"ID":             res.ID,
			"NIK":            res.NIK,
			"FullName":       res.FullName,
			"LegalName":      res.LegalName,
			"PlaceOfBirth":   res.PlaceOfBirth,
			"DateOfBirth":    res.DateOfBirth,
			"Salary":         formattedSalary,
			"KTPPhotoURL":    res.KTPPhotoURL,
			"SelfiePhotoURL": res.SelfiePhotoURL,
		},
	}, "layout/main")

}

func (xy xyzWeb) UpdateProfile(ctx *fiber.Ctx) error {
	c, cancel := context.WithTimeout(ctx.Context(), 10*time.Second)
	defer cancel()

	salaryStr := ctx.FormValue("salary")
	salary, err := util.ParseRupiahToFloat64(salaryStr)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Format salary tidak valid",
		})
	}

	req := dto.UpdateCustomersRequest{
		ID:           ctx.FormValue("id"),
		NIK:          ctx.FormValue("nik"),
		FullName:     ctx.FormValue("full_name"),
		LegalName:    ctx.FormValue("legal_name"),
		PlaceOfBirth: ctx.FormValue("place_of_birth"),
		DateOfBirth:  ctx.FormValue("date_of_birth"),
		Salary:       salary,
	}

	// Validasi
	fails := util.Validate(req)
	if len(fails) > 0 {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"code":    "99",
			"message": "validasi gagal",
			"errors":  fails,
		})
	}

	// Upload file (jika ada)
	ktpURL, _ := util.ProcessAndSaveImageFile(ctx, util.ImageSaveOptions{
		FieldName: "ktp_photo",
		BasePath:  xy.conf.Storage.BasePath,
		PublicURL: xy.conf.Server.Assets,
		MaxSizeMB: 2,
	})
	if ktpURL != "" {
		req.KTPPhotoURL = ktpURL
	}

	selfieURL, _ := util.ProcessAndSaveImageFile(ctx, util.ImageSaveOptions{
		FieldName: "selfie_photo",
		BasePath:  xy.conf.Storage.BasePath,
		PublicURL: xy.conf.Server.Assets,
		MaxSizeMB: 2,
	})
	if selfieURL != "" {
		req.SelfiePhotoURL = selfieURL
	}

	// Simpan
	if err := xy.customersService.Update(c, req); err != nil {
		return ctx.Status(fiber.StatusInternalServerError).
			JSON(dto.CreateResponseError(err.Error()))
	}

	return ctx.Redirect("/")
}

func (xy xyzWeb) Credits(ctx *fiber.Ctx) error {
	c, cancel := context.WithTimeout(ctx.Context(), 10*time.Second)
	defer cancel()

	sess, err := xy.store.Get(ctx)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).SendString("Gagal mengambil session")
	}

	id := sess.Get("user_id")
	if id == nil {
		return ctx.Redirect("/login")
	}

	customerID, ok := id.(string)
	if !ok {
		return ctx.Status(fiber.StatusInternalServerError).SendString("user_id tidak valid")
	}

	limitResp, err := xy.limitService.Show(c, customerID)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).
			JSON(dto.CreateResponseError("Gagal mengambil data limit: " + err.Error()))
	}

	transactionResp, err := xy.transactionsService.CustomerShow(c, customerID)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).
			JSON(dto.CreateResponseError("Gagal mengambil data transaksi: " + err.Error()))
	}

	formattedLimits := make([]fiber.Map, 0, len(limitResp.Limit))
	for i, limit := range limitResp.Limit {
		formattedLimits = append(formattedLimits, fiber.Map{
			"No":          i + 1,
			"TenorMonths": limit.TenorMonths,
			"LimitAmount": util.FormatRupiah(limit.LimitAmount),
			"Status":      limit.Status,
		})
	}

	formattedTransactions := make([]fiber.Map, 0, len(transactionResp))
	for _, t := range transactionResp {
		formattedTransactions = append(formattedTransactions, fiber.Map{
			"ContractNumber": t.ContractNumber,
			"AssetName":      t.AssetName,
			"Channel":        t.Channel,
			"OTRAmount":      util.FormatRupiah(t.OTRAmount),
			"AdminFee":       util.FormatRupiah(t.AdminFee),
			"Installment":    util.FormatRupiah(t.Installment),
			"Interest":       util.FormatRupiah(t.Interest),
			"TenorMonths":    t.TenorMonths,
		})
	}

	return ctx.Render("credits", fiber.Map{
		"Limits":      formattedLimits,
		"Transaction": formattedTransactions,
	}, "layout/main")
}

func (xy xyzWeb) Transactions(ctx *fiber.Ctx) error {
	_, cancel := context.WithTimeout(ctx.Context(), 10*time.Second)
	defer cancel()

	data, err := os.ReadFile("assets/product.json")
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).
			SendString("Gagal membaca file JSON")
	}

	var products []dto.Product
	if err := json.Unmarshal(data, &products); err != nil {
		return ctx.Status(fiber.StatusInternalServerError).
			SendString("Gagal memparsing JSON")
	}

	return ctx.Render("transactions", fiber.Map{
		"Products": products,
	}, "layout/main")
}

func (xy xyzWeb) CreateTransactions(ctx *fiber.Ctx) error {
	c, cancel := context.WithTimeout(ctx.Context(), 10*time.Second)
	defer cancel()

	// Ambil form data
	otrAmount, _ := strconv.ParseFloat(ctx.FormValue("otr_amount"), 64)
	adminFee, _ := strconv.ParseFloat(ctx.FormValue("admin_fee"), 64)
	installment, _ := strconv.ParseFloat(ctx.FormValue("installment"), 64)
	interest, _ := strconv.ParseFloat(ctx.FormValue("interest"), 64)
	tenorMonths, _ := strconv.Atoi(ctx.FormValue("tenor_months"))

	// Ambil session
	sess, err := xy.store.Get(ctx)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).SendString("Gagal mengambil session")
	}

	id := sess.Get("user_id")
	if id == nil {
		return ctx.Redirect("/login")
	}

	customerID, ok := id.(string)
	if !ok {
		return ctx.Status(fiber.StatusInternalServerError).SendString("user_id tidak valid")
	}

	// Ambil limit user
	limitResp, err := xy.limitService.Show(c, customerID)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).
			JSON(dto.CreateResponseError("Gagal mengambil data limit: " + err.Error()))
	}

	var matchedLimit *dto.LimitDetail
	for _, limit := range limitResp.Limit {
		if limit.TenorMonths == tenorMonths {
			if strings.ToLower(limit.Status) == "booked" {
				return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
					"code":    "98",
					"message": fmt.Sprintf("Tenor %d bulan sudah digunakan", tenorMonths),
				})
			}

			if limit.LimitAmount < otrAmount {
				return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
					"code":    "97",
					"message": "Harga dan tenor harus lebih besar dari harga OTR",
				})
			}

			matchedLimit = &limit
			break
		}
	}

	if matchedLimit == nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"code":    "97",
			"message": fmt.Sprintf("Tenor %d bulan tidak tersedia untuk user ini", tenorMonths),
		})
	}

	updateReq := dto.UpdateLimitRequest{
		ID:          matchedLimit.LimitID,
		CustomerId:  customerID,
		TenorMonths: matchedLimit.TenorMonths,
		LimitAmount: matchedLimit.LimitAmount,
		Status:      "booked",
	}

	if err := xy.limitService.Update(c, updateReq); err != nil {
		return ctx.Status(fiber.StatusInternalServerError).
			JSON(dto.CreateResponseError("Gagal update status limit: " + err.Error()))
	}

	single := dto.SingleTransactionRequest{
		CustomerID:     customerID,
		ContractNumber: ctx.FormValue("contract_number"),
		Channel:        ctx.FormValue("channel"),
		OTRAmount:      otrAmount,
		AdminFee:       adminFee,
		Installment:    installment,
		Interest:       interest,
		AssetName:      ctx.FormValue("asset_name"),
		TenorMonths:    tenorMonths,
	}

	fails := util.Validate(single)
	if len(fails) > 0 {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"code":    "99",
			"message": "validasi gagal",
			"errors":  fails,
		})
	}

	req := dto.CreateTransactionsRequest{
		Transactions: []dto.SingleTransactionRequest{single},
	}

	if err := xy.transactionsService.Create(c, req); err != nil {
		return ctx.Status(fiber.StatusInternalServerError).
			JSON(dto.CreateResponseError(err.Error()))
	}

	return ctx.Redirect("/credits", fiber.StatusSeeOther)
}

func (xy xyzWeb) Login(ctx *fiber.Ctx) error {
	_, cancel := context.WithTimeout(ctx.Context(), 10*time.Second)
	defer cancel()

	return ctx.Render("login", fiber.Map{})

}

func (xy xyzWeb) HandleLogin(ctx *fiber.Ctx) error {
	c, cancel := context.WithTimeout(ctx.Context(), 10*time.Second)
	defer cancel()

	email := ctx.FormValue("email")
	password := ctx.FormValue("password")

	req := dto.AuthRequest{
		Email:    email,
		Password: password,
	}

	resp, err := xy.authService.LoginWeb(c, req)
	if err != nil {
		return ctx.Render("login", fiber.Map{
			"Error": "Email atau password salah",
		})
	}

	sess, err := xy.store.Get(ctx)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).SendString("Gagal membuat session")
	}

	sess.Set("user_id", resp.ID)
	sess.Set("email", email)
	if err := sess.Save(); err != nil {
		return ctx.Status(fiber.StatusInternalServerError).SendString("Gagal menyimpan session")
	}

	return ctx.Redirect("/")
}

func (xy xyzWeb) LogoutWeb(ctx *fiber.Ctx) error {
	sess, err := xy.store.Get(ctx)
	if err == nil {
		sess.Destroy()
	}
	return ctx.Redirect("/login", fiber.StatusSeeOther)
}
