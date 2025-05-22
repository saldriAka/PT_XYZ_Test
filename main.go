package main

import (
	"saldri/test_pt_xyz/dto"
	"saldri/test_pt_xyz/internal/api"
	"saldri/test_pt_xyz/internal/config"
	"saldri/test_pt_xyz/internal/connection"
	"saldri/test_pt_xyz/internal/repository"
	"saldri/test_pt_xyz/internal/service"
	"saldri/test_pt_xyz/internal/web"

	jwtMid "github.com/gofiber/contrib/jwt"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/session"
	"github.com/gofiber/template/html/v2"
)

var store = session.New()

func main() {
	cnf := config.Get()
	dbConnection := connection.GetDatabase(cnf.Database)

	var app *fiber.App

	if cnf.ServiceMode.ServiceMode == "web" {
		engine := html.New("./views", ".html")
		app = fiber.New(fiber.Config{
			Views: engine,
		})
	} else {
		app = fiber.New()
	}

	jwtMidd := jwtMid.New(jwtMid.Config{
		SigningKey: jwtMid.SigningKey{Key: []byte(cnf.Jwt.Key)},
		ErrorHandler: func(ctx *fiber.Ctx, err error) error {
			return ctx.Status(fiber.StatusUnauthorized).JSON(
				dto.CreateResponseError("endpoint perlu token, silakan login"),
			)
		},
	})

	customersRepository := repository.NewCustomers(dbConnection)
	limitRepository := repository.NewLimit(dbConnection)
	transactionsRepository := repository.NewTransactions(dbConnection)
	userRepository := repository.NewUser(dbConnection)

	customerService := service.NewCustomers(cnf, customersRepository)
	limitService := service.NewLimit(cnf, limitRepository)
	transactionsService := service.NewTransactions(cnf, transactionsRepository)
	authService := service.NewAuth(cnf, userRepository)

	api.NewCustomers(app, cnf, customerService, jwtMidd)
	api.NewLimit(app, limitService, jwtMidd)
	api.NewTransactions(app, transactionsService, jwtMidd)
	api.NewAuth(app, authService)

	web.NewWeb(app, customerService, limitService, transactionsService, authService, cnf, store)

	_ = app.Listen(cnf.Server.Host + ":" + cnf.Server.Port)
}
