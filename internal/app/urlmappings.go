package app

import (
	"patient-monitor-backend/internal/routers"

	"github.com/gofiber/fiber/v2"
)

type Router struct {
	App *fiber.App
}

func NewRouter(app *fiber.App) *Router {
	return &Router{App: app}
}

func mapUrls(app *fiber.App) {
	app.Get("/ping", func(c *fiber.Ctx) error {
		return c.SendString("pong")
	})

	router := routers.NewRouter(app)
	router.SetAuthRoutes()
	// router.SetTestRoutes()
}
