package pacientes

import (
	"context"
	"log"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"hospitalaria/config"
	"hospitalaria/utils"
	"github.com/jackc/pgx/v4"
)

func CreateAppointment(c *fiber.Ctx) error {
    userID := c.Locals("user_id").(int)
    log.Printf("Solicitud recibida para userID: %d", userID)
    role := c.Locals("role").(string)
    if role != "Paciente" {
        utils.LogAction(userID, "create_appointment", "fallido", "Permiso denegado: Solo Pacientes pueden agendar citas")
        return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"error": "Permiso denegado"})
    }

    type AppointmentInput struct {
        IDMedico       int    `json:"id_medico"`
        FechaHora      string `json:"fecha_hora"`
        Motivo         string `json:"motivo"`
        IDConsultorio  int    `json:"id_consultorio,omitempty"`
        IDHorario      int    `json:"id_horario,omitempty"`
    }
    var input AppointmentInput
    if err := c.BodyParser(&input); err != nil {
        utils.LogAction(userID, "create_appointment", "fallido", "JSON inválido: "+err.Error())
        return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "JSON inválido"})
    }

    var idPaciente int
    err := config.Conn.QueryRow(context.Background(),
        "SELECT id_paciente FROM pacientes WHERE id_usuario = $1", userID).Scan(&idPaciente)
    if err != nil {
        log.Printf("Error al obtener id_paciente: %v", err)
        utils.LogAction(userID, "create_appointment", "fallido", "Paciente no encontrado: "+err.Error())
        return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Paciente no encontrado"})
    }

    log.Printf("Datos a insertar: id_paciente=%d, id_medico=%d, fecha_hora=%s, motivo=%s, id_consultorio=%d, id_horario=%d",
        idPaciente, input.IDMedico, input.FechaHora, input.Motivo, input.IDConsultorio, input.IDHorario)
    _, err = config.Conn.Exec(context.Background(),
        "INSERT INTO citas (id_paciente, id_medico, fecha_hora, estado, id_consultorio, id_horario, motivo) VALUES ($1, $2, $3, 'pendiente', $4, $5, $6)",
        idPaciente, input.IDMedico, input.FechaHora, input.IDConsultorio, input.IDHorario, input.Motivo)
    if err != nil {
        log.Printf("Error al crear cita: %v", err)
        utils.LogAction(userID, "create_appointment", "fallido", "Error al crear cita: "+err.Error())
        return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Error al crear cita: " + err.Error()})
    }
    utils.LogAction(userID, "create_appointment", "exitoso", "Cita agendada con medico ID "+strconv.Itoa(input.IDMedico))
    return c.JSON(fiber.Map{"message": "Cita agendada", "estado": "pendiente"})
}

func GetAppointments(c *fiber.Ctx) error {
    userID := c.Locals("user_id").(int)
    role := c.Locals("role").(string)
    var rows pgx.Rows
    var err error
    if role == "Paciente" {
        var idPaciente int
        err = config.Conn.QueryRow(context.Background(), "SELECT id_paciente FROM pacientes WHERE id_usuario = $1", userID).Scan(&idPaciente)
        if err != nil {
            utils.LogAction(userID, "read_appointment", "fallido", "Paciente no encontrado: "+err.Error())
            return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Paciente no encontrado"})
        }
        rows, err = config.Conn.Query(context.Background(), "SELECT id_cita, id_medico, fecha_hora, estado, motivo FROM citas WHERE id_paciente = $1", idPaciente)
    } else {
        utils.LogAction(userID, "read_appointment", "fallido", "Permiso denegado: Solo Pacientes pueden ver sus citas")
        return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"error": "Permiso denegado"})
    }
    if err != nil {
        log.Printf("Error al obtener citas: %v", err)
        utils.LogAction(userID, "read_appointment", "fallido", "Error al obtener citas: "+err.Error())
        return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Error al obtener citas"})
    }
    defer rows.Close()
    var appointments []struct {
        ID_cita    int    `json:"id_cita"`
        IDMedico   int    `json:"id_medico"`
        FechaHora  string `json:"fecha_hora"`
        Estado     string `json:"estado"`
        Motivo     string `json:"motivo"`
    }
    for rows.Next() {
        var app struct {
            ID_cita    int    `json:"id_cita"`
            IDMedico   int    `json:"id_medico"`
            FechaHora  string `json:"fecha_hora"`
            Estado     string `json:"estado"`
            Motivo     string `json:"motivo"`
        }
        err := rows.Scan(&app.ID_cita, &app.IDMedico, &app.FechaHora, &app.Estado, &app.Motivo)
        if err != nil {
            utils.LogAction(userID, "read_appointment", "fallido", "Error al leer cita: "+err.Error())
            return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Error al leer citas"})
        }
        appointments = append(appointments, app)
    }
    utils.LogAction(userID, "read_appointment", "exitoso", "Citas leídas para usuario "+strconv.Itoa(userID))
    return c.JSON(appointments)
}

func DeleteAppointment(c *fiber.Ctx) error {
    userID := c.Locals("user_id").(int)
    role := c.Locals("role").(string)
    if role != "Paciente" {
        utils.LogAction(userID, "delete_appointment", "fallido", "Permiso denegado: Solo Pacientes pueden cancelar citas")
        return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"error": "Permiso denegado"})
    }

    type AppointmentDelete struct {
        ID_cita int `json:"id_cita"`
    }
    var input AppointmentDelete
    if err := c.BodyParser(&input); err != nil {
        utils.LogAction(userID, "delete_appointment", "fallido", "JSON inválido: "+err.Error())
        return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "JSON inválido"})
    }

    var idPaciente int
    err := config.Conn.QueryRow(context.Background(), "SELECT id_paciente FROM pacientes WHERE id_usuario = $1", userID).Scan(&idPaciente)
    if err != nil {
        utils.LogAction(userID, "delete_appointment", "fallido", "Paciente no encontrado: "+err.Error())
        return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Paciente no encontrado"})
    }

    result, err := config.Conn.Exec(context.Background(),
        "DELETE FROM citas WHERE id_cita = $1 AND id_paciente = $2 AND estado = 'pendiente'",
        input.ID_cita, idPaciente)
    if err != nil {
        log.Printf("Error al cancelar cita: %v", err)
        utils.LogAction(userID, "delete_appointment", "fallido", "Error al cancelar cita: "+err.Error())
        return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Error al cancelar cita: " + err.Error()})
    }
    if result.RowsAffected() == 0 {
        utils.LogAction(userID, "delete_appointment", "fallido", "Cita no encontrada o no cancelable: ID "+strconv.Itoa(input.ID_cita))
        return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Cita no encontrada o no cancelable"})
    }
    utils.LogAction(userID, "delete_appointment", "exitoso", "Cita cancelada: ID "+strconv.Itoa(input.ID_cita))
    return c.JSON(fiber.Map{"message": "Cita cancelada", "estado": "cancelada"})
}