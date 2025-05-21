package main

import (
	"saldri/test_pt_xyz/internal/config"
	"saldri/test_pt_xyz/internal/connection"
	"saldri/test_pt_xyz/internal/controller"
	"saldri/test_pt_xyz/internal/repository"
	"saldri/test_pt_xyz/internal/service"

	"github.com/gofiber/fiber/v2"
)

func main() {
	cnf := config.Get()
	dbConnection := connection.GetDatabase(cnf.Database)

	app := fiber.New()

	customersRepository := repository.NewCustomers(dbConnection)
	limitRepository := repository.NewLimit(dbConnection)
	transactionsRepository := repository.NewTransactions(dbConnection)

	customerService := service.NewCustomers(cnf, customersRepository)
	limitService := service.NewLimit(cnf, limitRepository)
	transactionsService := service.NewTransactions(cnf, transactionsRepository)

	controller.NewCustomers(app, cnf, customerService)
	controller.NewLimit(app, limitService)
	controller.NewTransactions(app, transactionsService)

	_ = app.Listen(cnf.Server.Host + ":" + cnf.Server.Port)
}
