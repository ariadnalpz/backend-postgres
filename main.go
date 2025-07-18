package main

import (
	"log"

	"github.com/gofiber/fiber/v2"
	"hospitalaria/config"
	"hospitalaria/routes"
	"hospitalaria/utils"
)

func main() {
	utils.LoadEnv()

	if err := config.InitDatabase(); err != nil {
		log.Fatal("No se pudo iniciar la base de datos:", err)
	}

	defer config.Conn.Close()

	app := fiber.New()
	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("¡Bienvenido al backend del Sistema de Citas y Reportes del Hospital!")
	})

	// Registra todas las rutas
	routes.SetupAuthRoutes(app)
	routes.SetupPacienteRoutes(app)
	routes.SetupMedicoRoutes(app)
	routes.SetupEnfermeraRoutes(app)

	log.Fatal(app.Listen(":3000"))
}