package routes

import (
	"github.com/gofiber/fiber/v2"
	"hospitalaria/handlers/pacientes"
	"hospitalaria/middleware"
)

func SetupPacienteRoutes(app *fiber.App) {
	app.Post("/appointments", middleware.JWTProtected(), pacientes.CreateAppointment)
	app.Get("/appointments", middleware.JWTProtected(), pacientes.GetAppointments)
	app.Delete("/appointments", middleware.JWTProtected(), pacientes.DeleteAppointment)
	app.Post("/expedientes", middleware.JWTProtected(), pacientes.CreateExpediente)
	app.Get("/expedientes", middleware.JWTProtected(), pacientes.GetExpedientes)
	app.Put("/expedientes", middleware.JWTProtected(), pacientes.UpdateExpediente)
	app.Delete("/expedientes", middleware.JWTProtected(), pacientes.DeleteExpediente)
}