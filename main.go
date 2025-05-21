package main

import (
	"saldri/test_pt_xyz/internal/config"
	"saldri/test_pt_xyz/internal/connection"

	"github.com/gofiber/fiber/v2"
)

func main() {
	cnf := config.Get()
	dbConnection := connection.GetDatabase(cnf.Database)

	app := fiber.New()

	_ = app.Listen(cnf.Server.Host + ":" + cnf.Server.Port)
}
