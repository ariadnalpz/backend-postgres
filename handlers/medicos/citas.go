package medicos

import (
	"context"
	"log"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"hospitalaria/config"
	"hospitalaria/utils"
)

func UpdateAppointment(c *fiber.Ctx) error {
    userID := c.Locals("user_id").(int)
    role := c.Locals("role").(string)
    if role != "Medico" {
        utils.LogAction(userID, "update_appointment", "fallido", "Permiso denegado: Solo Médicos pueden aceptar citas")
        return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"error": "Permiso denegado"})
    }

    type AppointmentUpdate struct {
        ID_cita int `json:"id_cita"`
    }
    var input AppointmentUpdate
    if err := c.BodyParser(&input); err != nil {
        utils.LogAction(userID, "update_appointment", "fallido", "JSON inválido: "+err.Error())
        return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "JSON inválido"})
    }

    var idMedico int
    err := config.Conn.QueryRow(context.Background(), "SELECT id_medico FROM medicos WHERE id_usuario = $1", userID).Scan(&idMedico)
    if err != nil {
        utils.LogAction(userID, "update_appointment", "fallido", "Médico no encontrado: "+err.Error())
        return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Médico no encontrado"})
    }

    result, err := config.Conn.Exec(context.Background(),
        "UPDATE citas SET estado = 'aceptada' WHERE id_cita = $1 AND id_medico = $2 AND estado = 'pendiente'",
        input.ID_cita, idMedico)
    if err != nil {
        log.Printf("Error al actualizar cita: %v", err)
        utils.LogAction(userID, "update_appointment", "fallido", "Error al aceptar cita: "+err.Error())
        return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Error al aceptar cita: " + err.Error()})
    }
    if result.RowsAffected() == 0 {
        utils.LogAction(userID, "update_appointment", "fallido", "Cita no encontrada o ya procesada: ID "+strconv.Itoa(input.ID_cita))
        return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Cita no encontrada o ya procesada"})
    }
    utils.LogAction(userID, "update_appointment", "exitoso", "Cita aceptada: ID "+strconv.Itoa(input.ID_cita))
    return c.JSON(fiber.Map{"message": "Cita aceptada", "estado": "aceptada"})
}