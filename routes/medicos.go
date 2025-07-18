package routes

import (
	"github.com/gofiber/fiber/v2"
	"hospitalaria/handlers/medicos"
	"hospitalaria/middleware"
)

func SetupMedicoRoutes(app *fiber.App) {
	app.Put("/appointments", middleware.JWTProtected(), medicos.UpdateAppointment)
	app.Post("/consultorios", middleware.JWTProtected(), medicos.CreateConsultorio)
	app.Get("/consultorios", middleware.JWTProtected(), medicos.GetConsultorios)
	app.Put("/consultorios", middleware.JWTProtected(), medicos.UpdateConsultorio)
	app.Delete("/consultorios", middleware.JWTProtected(), medicos.DeleteConsultorio)
	app.Post("/horarios", middleware.JWTProtected(), medicos.CreateHorario)
	app.Get("/horarios", middleware.JWTProtected(), medicos.GetHorarios)
	app.Put("/horarios", middleware.JWTProtected(), medicos.UpdateHorario)
	app.Delete("/horarios", middleware.JWTProtected(), medicos.DeleteHorario)
}