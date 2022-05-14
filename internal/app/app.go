package app

import (
	"log"
	"patient-monitor-backend/internal/config"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/compress"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
)

var app = fiber.New()

func StartApplication() {
	appConfig := config.AppConfig

	app.Use(cors.New())

	if appConfig.Environment == "development" {
		app.Use(logger.New(logger.ConfigDefault))
	}

	app.Use(compress.New())
	app.Use(recover.New())

	mapUrls(app)

	log.Fatal(app.Listen(":" + strconv.Itoa(appConfig.ServerConfig.APIPORT)))
}
