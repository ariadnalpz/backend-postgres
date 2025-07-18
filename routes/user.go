package routes

import (
	"github.com/gofiber/fiber/v2"
	"hospitalaria/handlers"
)

func SetupUserRoutes(app *fiber.App, jwtMiddleware fiber.Handler) {
	user := app.Group("/users", jwtMiddleware)
	user.Get("/profile", handlers.GetUserProfile) // Endpoint protegido
}