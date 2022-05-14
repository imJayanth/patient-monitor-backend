package routers

import (
	"patient-monitor-backend/internal/controllers"
)

func (r *Router) SetAuthRoutes() {
	userController := controllers.NewUserController()

	api := r.app.Group("/api/v1/patient")
	api.Post("/login", userController.Login)
	api.Get("/:ID/contacts", userController.GetContactsById)
	api.Post("/contact", userController.AddContact)
	api.Get("/:ID/sensor-data", userController.GetSensorDataById)
}
