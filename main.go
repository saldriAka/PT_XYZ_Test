package main

import (
	"saldri/test_pt_xyz/dto"
	"saldri/test_pt_xyz/internal/config"
	"saldri/test_pt_xyz/internal/connection"
	"saldri/test_pt_xyz/internal/controller"
	"saldri/test_pt_xyz/internal/repository"
	"saldri/test_pt_xyz/internal/service"

	jwtMid "github.com/gofiber/contrib/jwt"
	"github.com/gofiber/fiber/v2"
)

func main() {
	cnf := config.Get()
	dbConnection := connection.GetDatabase(cnf.Database)

	app := fiber.New()

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

	controller.NewCustomers(app, cnf, customerService, jwtMidd)
	controller.NewLimit(app, limitService, jwtMidd)
	controller.NewTransactions(app, transactionsService, jwtMidd)
	controller.NewAuth(app, authService)

	_ = app.Listen(cnf.Server.Host + ":" + cnf.Server.Port)
}
