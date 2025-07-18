package routes

import (
	"github.com/gofiber/fiber/v2"
	"hospitalaria/handlers/enfermeras"
	"hospitalaria/middleware"
)

func SetupEnfermeraRoutes(app *fiber.App) {
	app.Post("/consultas", middleware.JWTProtected(), enfermeras.AssignConsulta)
}