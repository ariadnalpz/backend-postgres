package pacientes

import (
	"context"
	"log"
	"strconv"
	"strings"

	"github.com/gofiber/fiber/v2"
	"hospitalaria/config"
	"hospitalaria/utils"
)

func CreateExpediente(c *fiber.Ctx) error {
    userID := c.Locals("user_id").(int)
    log.Printf("Solicitud recibida para userID: %d", userID)
    role := c.Locals("role").(string)
    if role != "Paciente" {
        utils.LogAction(userID, "create_expediente", "fallido", "Permiso denegado: Solo Pacientes pueden crear expedientes")
        return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"error": "Permiso denegado"})
    }

    type ExpedienteInput struct {
        AntecedentesMedicos string `json:"antecedentes_medicos"`
        Alergias           string `json:"alergias"`
        Tratamientos       string `json:"tratamientos"`
        FechaActualizacion string `json:"fecha_actualizacion"`
    }
    var input ExpedienteInput
    if err := c.BodyParser(&input); err != nil {
        utils.LogAction(userID, "create_expediente", "fallido", "JSON inválido: "+err.Error())
        return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "JSON inválido"})
    }

    var idPaciente int
    err := config.Conn.QueryRow(context.Background(), "SELECT id_paciente FROM pacientes WHERE id_usuario = $1", userID).Scan(&idPaciente)
    if err != nil {
        utils.LogAction(userID, "create_expediente", "fallido", "Paciente no encontrado: "+err.Error())
        return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Paciente no encontrado"})
    }

    log.Printf("Datos a insertar: id_paciente=%d, antecedentes_medicos=%s, alergias=%s, tratamientos=%s, fecha_actualizacion=%s",
        idPaciente, input.AntecedentesMedicos, input.Alergias, input.Tratamientos, input.FechaActualizacion)
    _, err = config.Conn.Exec(context.Background(),
        "INSERT INTO expedientes (id_paciente, antecedentes_medicos, alergias, tratamientos, fecha_actualizacion) VALUES ($1, $2, $3, $4, $5)",
        idPaciente, input.AntecedentesMedicos, input.Alergias, input.Tratamientos, input.FechaActualizacion)
    if err != nil {
        log.Printf("Error al crear expediente: %v", err)
        utils.LogAction(userID, "create_expediente", "fallido", "Error al crear expediente: "+err.Error())
        return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Error al crear expediente: " + err.Error()})
    }
    utils.LogAction(userID, "create_expediente", "exitoso", "Expediente creado para paciente ID "+strconv.Itoa(idPaciente))
    return c.JSON(fiber.Map{"message": "Expediente creado"})
}

func GetExpedientes(c *fiber.Ctx) error {
    userID := c.Locals("user_id").(int)
    log.Printf("Solicitud de lectura para userID: %d", userID)
    role := c.Locals("role").(string)
    if role != "Paciente" {
        utils.LogAction(userID, "read_expediente", "fallido", "Permiso denegado: Solo Pacientes pueden ver expedientes")
        return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"error": "Permiso denegado"})
    }

    var idPaciente int
    err := config.Conn.QueryRow(context.Background(), "SELECT id_paciente FROM pacientes WHERE id_usuario = $1", userID).Scan(&idPaciente)
    if err != nil {
        utils.LogAction(userID, "read_expediente", "fallido", "Paciente no encontrado: "+err.Error())
        return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Paciente no encontrado"})
    }

    row := config.Conn.QueryRow(context.Background(), "SELECT id_expediente, antecedentes_medicos, alergias, tratamientos, fecha_actualizacion FROM expedientes WHERE id_paciente = $1", idPaciente)
    var exp struct {
        IDExpediente       int    `json:"id_expediente"`
        AntecedentesMedicos string `json:"antecedentes_medicos"`
        Alergias           string `json:"alergias"`
        Tratamientos       string `json:"tratamientos"`
        FechaActualizacion string `json:"fecha_actualizacion"`
    }
    err = row.Scan(&exp.IDExpediente, &exp.AntecedentesMedicos, &exp.Alergias, &exp.Tratamientos, &exp.FechaActualizacion)
    if err != nil {
        log.Printf("Error al obtener expediente: %v", err)
        utils.LogAction(userID, "read_expediente", "fallido", "Error al obtener expediente: "+err.Error())
        return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Error al obtener expediente"})
    }
    utils.LogAction(userID, "read_expediente", "exitoso", "Expediente leído para paciente ID "+strconv.Itoa(idPaciente))
    return c.JSON(exp)
}

func UpdateExpediente(c *fiber.Ctx) error {
    userID := c.Locals("user_id").(int)
    log.Printf("Solicitud de actualización para userID: %d", userID)
    role := c.Locals("role").(string)
    if role != "Paciente" {
        utils.LogAction(userID, "update_expediente", "fallido", "Permiso denegado: Solo Pacientes pueden actualizar expedientes")
        return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"error": "Permiso denegado"})
    }

    type ExpedienteUpdate struct {
        AntecedentesMedicos string `json:"antecedentes_medicos,omitempty"`
        Alergias           string `json:"alergias,omitempty"`
        Tratamientos       string `json:"tratamientos,omitempty"`
        FechaActualizacion string `json:"fecha_actualizacion,omitempty"`
    }
    var input ExpedienteUpdate
    if err := c.BodyParser(&input); err != nil {
        utils.LogAction(userID, "update_expediente", "fallido", "JSON inválido: "+err.Error())
        return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "JSON inválido"})
    }

    var idPaciente int
    err := config.Conn.QueryRow(context.Background(), "SELECT id_paciente FROM pacientes WHERE id_usuario = $1", userID).Scan(&idPaciente)
    if err != nil {
        utils.LogAction(userID, "update_expediente", "fallido", "Paciente no encontrado: "+err.Error())
        return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Paciente no encontrado"})
    }

    setClause := "SET "
    args := []interface{}{idPaciente}
    paramCount := 1
    if input.AntecedentesMedicos != "" {
        paramCount++
        setClause += "antecedentes_medicos = $" + strconv.Itoa(paramCount) + ", "
        args = append(args, input.AntecedentesMedicos)
    }
    if input.Alergias != "" {
        paramCount++
        setClause += "alergias = $" + strconv.Itoa(paramCount) + ", "
        args = append(args, input.Alergias)
    }
    if input.Tratamientos != "" {
        paramCount++
        setClause += "tratamientos = $" + strconv.Itoa(paramCount) + ", "
        args = append(args, input.Tratamientos)
    }
    if input.FechaActualizacion != "" {
        paramCount++
        setClause += "fecha_actualizacion = $" + strconv.Itoa(paramCount) + ", "
        args = append(args, input.FechaActualizacion)
    }
    setClause = strings.TrimSuffix(setClause, ", ") + " WHERE id_paciente = $1"

    result, err := config.Conn.Exec(context.Background(), "UPDATE expedientes "+setClause, args...)
    if err != nil {
        log.Printf("Error al actualizar expediente: %v", err)
        utils.LogAction(userID, "update_expediente", "fallido", "Error al actualizar expediente: "+err.Error())
        return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Error al actualizar expediente: " + err.Error()})
    }
    if result.RowsAffected() == 0 {
        utils.LogAction(userID, "update_expediente", "fallido", "Expediente no encontrado para paciente ID "+strconv.Itoa(idPaciente))
        return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Expediente no encontrado"})
    }
    utils.LogAction(userID, "update_expediente", "exitoso", "Expediente actualizado para paciente ID "+strconv.Itoa(idPaciente))
    return c.JSON(fiber.Map{"message": "Expediente actualizado"})
}

func DeleteExpediente(c *fiber.Ctx) error {
    userID := c.Locals("user_id").(int)
    log.Printf("Solicitud de eliminación para userID: %d", userID)
    role := c.Locals("role").(string)
    if role != "Paciente" {
        utils.LogAction(userID, "delete_expediente", "fallido", "Permiso denegado: Solo Pacientes pueden eliminar expedientes")
        return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"error": "Permiso denegado"})
    }

    var idPaciente int
    err := config.Conn.QueryRow(context.Background(), "SELECT id_paciente FROM pacientes WHERE id_usuario = $1", userID).Scan(&idPaciente)
    if err != nil {
        utils.LogAction(userID, "delete_expediente", "fallido", "Paciente no encontrado: "+err.Error())
        return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Paciente no encontrado"})
    }

    result, err := config.Conn.Exec(context.Background(), "DELETE FROM expedientes WHERE id_paciente = $1", idPaciente)
    if err != nil {
        log.Printf("Error al eliminar expediente: %v", err)
        utils.LogAction(userID, "delete_expediente", "fallido", "Error al eliminar expediente: "+err.Error())
        return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Error al eliminar expediente: " + err.Error()})
    }
    if result.RowsAffected() == 0 {
        utils.LogAction(userID, "delete_expediente", "fallido", "Expediente no encontrado para paciente ID "+strconv.Itoa(idPaciente))
        return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Expediente no encontrado"})
    }
    utils.LogAction(userID, "delete_expediente", "exitoso", "Expediente eliminado para paciente ID "+strconv.Itoa(idPaciente))
    return c.JSON(fiber.Map{"message": "Expediente eliminado"})
}