package enfermeras

import (
	"context"
	"log"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"hospitalaria/config"
	"hospitalaria/utils"
)

func AssignConsulta(c *fiber.Ctx) error {
    userID := c.Locals("user_id").(int)
    log.Printf("Solicitud recibida para userID: %d", userID)
    role := c.Locals("role").(string)
    if role != "Enfermero" {
        utils.LogAction(userID, "assign_consulta", "fallido", "Permiso denegado: Solo Enfermeras pueden asignar consultas")
        return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"error": "Permiso denegado"})
    }

    type ConsultaInput struct {
        IDCita       int    `json:"id_cita"`
        IDPaciente   int    `json:"id_paciente"`
        IDMedico     int    `json:"id_medico"`
        FechaHora    string `json:"fecha_hora"`
        Diagnostico  string `json:"diagnostico,omitempty"`
        Estado       string `json:"estado"`
    }
    var input ConsultaInput
    if err := c.BodyParser(&input); err != nil {
        utils.LogAction(userID, "assign_consulta", "fallido", "JSON inválido: "+err.Error())
        return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "JSON inválido"})
    }

    var idEnfermera int
    err := config.Conn.QueryRow(context.Background(), "SELECT id_enfermera FROM enfermeras WHERE id_usuario = $1", userID).Scan(&idEnfermera)
    if err != nil {
        utils.LogAction(userID, "assign_consulta", "fallido", "Enfermera no encontrada: "+err.Error())
        return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Enfermera no encontrada"})
    }

    log.Printf("Datos a insertar: id_cita=%d, id_paciente=%d, id_medico=%d, id_enfermera=%d, fecha_hora=%s, diagnostico=%s, estado=%s",
        input.IDCita, input.IDPaciente, input.IDMedico, idEnfermera, input.FechaHora, input.Diagnostico, input.Estado)
    _, err = config.Conn.Exec(context.Background(),
        "INSERT INTO consultas (id_cita, id_paciente, id_medico, id_enfermera, fecha_hora, diagnostico, estado) VALUES ($1, $2, $3, $4, $5, $6, $7)",
        input.IDCita, input.IDPaciente, input.IDMedico, idEnfermera, input.FechaHora, input.Diagnostico, input.Estado)
    if err != nil {
        log.Printf("Error al asignar consulta: %v", err)
        utils.LogAction(userID, "assign_consulta", "fallido", "Error al asignar consulta: "+err.Error())
        return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Error al asignar consulta: " + err.Error()})
    }
    utils.LogAction(userID, "assign_consulta", "exitoso", "Consulta asignada para enfermera ID "+strconv.Itoa(idEnfermera))
    return c.JSON(fiber.Map{"message": "Consulta asignada", "estado": input.Estado})
}