package routes

import (
	"github.com/gofiber/fiber/v2"
	"hospitalaria/handlers"
)

func SetupAuthRoutes(app *fiber.App) {
	app.Post("/register", handlers.CreateUser)
	app.Post("/login", handlers.Login)
	app.Post("/refresh-token", handlers.RefreshToken) // Nuevo endpoint para refresh token
	app.Get("/profile", handlers.GetUserProfile)
}