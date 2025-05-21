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

	customerService := service.NewCustomers(cnf, customersRepository)

	controller.NewCustomers(app, customerService)

	_ = app.Listen(cnf.Server.Host + ":" + cnf.Server.Port)
}
