package medicos

import (
	"context"
	"log"
	"strconv"
	"strings"

	"github.com/gofiber/fiber/v2"
	"hospitalaria/config"
	"hospitalaria/utils"
)

func CreateHorario(c *fiber.Ctx) error {
    userID := c.Locals("user_id").(int)
    log.Printf("Solicitud recibida para userID: %d", userID)
    role := c.Locals("role").(string)
    if role != "Medico" {
        utils.LogAction(userID, "create_horario", "fallido", "Permiso denegado: Solo Médicos pueden crear horarios")
        return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"error": "Permiso denegado"})
    }

    type HorarioInput struct {
        IDConsultorio int    `json:"id_consultorio"`
        DiaSemana     string `json:"dia_semana"`
        HoraInicio    string `json:"hora_inicio"`
        HoraFin       string `json:"hora_fin"`
        Estado        string `json:"estado,omitempty"`
    }
    var input HorarioInput
    if err := c.BodyParser(&input); err != nil {
        utils.LogAction(userID, "create_horario", "fallido", "JSON inválido: "+err.Error())
        return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "JSON inválido"})
    }

    var idMedico int
    err := config.Conn.QueryRow(context.Background(), "SELECT id_medico FROM medicos WHERE id_usuario = $1", userID).Scan(&idMedico)
    if err != nil {
        utils.LogAction(userID, "create_horario", "fallido", "Médico no encontrado: "+err.Error())
        return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Médico no encontrado"})
    }

    log.Printf("Datos a insertar: id_consultorio=%d, id_medico=%d, dia_semana=%s, hora_inicio=%s, hora_fin=%s, estado=%s",
        input.IDConsultorio, idMedico, input.DiaSemana, input.HoraInicio, input.HoraFin, input.Estado)
    _, err = config.Conn.Exec(context.Background(),
        "INSERT INTO horarios (id_consultorio, id_medico, dia_semana, hora_inicio, hora_fin, estado) VALUES ($1, $2, $3, $4, $5, $6)",
        input.IDConsultorio, idMedico, input.DiaSemana, input.HoraInicio, input.HoraFin, input.Estado)
    if err != nil {
        log.Printf("Error al crear horario: %v", err)
        utils.LogAction(userID, "create_horario", "fallido", "Error al crear horario: "+err.Error())
        return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Error al crear horario: " + err.Error()})
    }
    utils.LogAction(userID, "create_horario", "exitoso", "Horario creado para medico ID "+strconv.Itoa(idMedico))
    return c.JSON(fiber.Map{"message": "Horario creado", "estado": input.Estado})
}

func GetHorarios(c *fiber.Ctx) error {
    userID := c.Locals("user_id").(int)
    log.Printf("Solicitud de lectura para userID: %d", userID)
    role := c.Locals("role").(string)
    if role != "Medico" {
        utils.LogAction(userID, "read_horario", "fallido", "Permiso denegado: Solo Médicos pueden ver horarios")
        return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"error": "Permiso denegado"})
    }

    var idMedico int
    err := config.Conn.QueryRow(context.Background(), "SELECT id_medico FROM medicos WHERE id_usuario = $1", userID).Scan(&idMedico)
    if err != nil {
        utils.LogAction(userID, "read_horario", "fallido", "Médico no encontrado: "+err.Error())
        return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Médico no encontrado"})
    }

    rows, err := config.Conn.Query(context.Background(), "SELECT id_horario, id_consultorio, dia_semana, hora_inicio, hora_fin, estado FROM horarios WHERE id_medico = $1", idMedico)
    if err != nil {
        log.Printf("Error al obtener horarios: %v", err)
        utils.LogAction(userID, "read_horario", "fallido", "Error al obtener horarios: "+err.Error())
        return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Error al obtener horarios"})
    }
    defer rows.Close()
    var horarios []struct {
        IDHorario    int    `json:"id_horario"`
        IDConsultorio int   `json:"id_consultorio"`
        DiaSemana    string `json:"dia_semana"`
        HoraInicio   string `json:"hora_inicio"`
        HoraFin      string `json:"hora_fin"`
        Estado       string `json:"estado"`
    }
    for rows.Next() {
        var hor struct {
            IDHorario    int    `json:"id_horario"`
            IDConsultorio int   `json:"id_consultorio"`
            DiaSemana    string `json:"dia_semana"`
            HoraInicio   string `json:"hora_inicio"`
            HoraFin      string `json:"hora_fin"`
            Estado       string `json:"estado"`
        }
        err := rows.Scan(&hor.IDHorario, &hor.IDConsultorio, &hor.DiaSemana, &hor.HoraInicio, &hor.HoraFin, &hor.Estado)
        if err != nil {
            utils.LogAction(userID, "read_horario", "fallido", "Error al leer horario: "+err.Error())
            return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Error al leer horarios"})
        }
        horarios = append(horarios, hor)
    }
    utils.LogAction(userID, "read_horario", "exitoso", "Horarios leídos para medico ID "+strconv.Itoa(idMedico))
    return c.JSON(horarios)
}

func UpdateHorario(c *fiber.Ctx) error {
    userID := c.Locals("user_id").(int)
    log.Printf("Solicitud de actualización para userID: %d", userID)
    role := c.Locals("role").(string)
    if role != "Medico" {
        utils.LogAction(userID, "update_horario", "fallido", "Permiso denegado: Solo Médicos pueden actualizar horarios")
        return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"error": "Permiso denegado"})
    }

    type HorarioUpdate struct {
        IDHorario    int    `json:"id_horario"`
        IDConsultorio int   `json:"id_consultorio,omitempty"`
        DiaSemana    string `json:"dia_semana,omitempty"`
        HoraInicio   string `json:"hora_inicio,omitempty"`
        HoraFin      string `json:"hora_fin,omitempty"`
        Estado       string `json:"estado,omitempty"`
    }
    var input HorarioUpdate
    if err := c.BodyParser(&input); err != nil {
        utils.LogAction(userID, "update_horario", "fallido", "JSON inválido: "+err.Error())
        return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "JSON inválido"})
    }

    var idMedico int
    err := config.Conn.QueryRow(context.Background(), "SELECT id_medico FROM medicos WHERE id_usuario = $1", userID).Scan(&idMedico)
    if err != nil {
        utils.LogAction(userID, "update_horario", "fallido", "Médico no encontrado: "+err.Error())
        return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Médico no encontrado"})
    }

    setClause := "SET "
    args := []interface{}{input.IDHorario, idMedico}
    paramCount := 2
    if input.IDConsultorio != 0 {
        paramCount++
        setClause += "id_consultorio = $" + strconv.Itoa(paramCount) + ", "
        args = append(args, input.IDConsultorio)
    }
    if input.DiaSemana != "" {
        paramCount++
        setClause += "dia_semana = $" + strconv.Itoa(paramCount) + ", "
        args = append(args, input.DiaSemana)
    }
    if input.HoraInicio != "" {
        paramCount++
        setClause += "hora_inicio = $" + strconv.Itoa(paramCount) + ", "
        args = append(args, input.HoraInicio)
    }
    if input.HoraFin != "" {
        paramCount++
        setClause += "hora_fin = $" + strconv.Itoa(paramCount) + ", "
        args = append(args, input.HoraFin)
    }
    if input.Estado != "" {
        paramCount++
        setClause += "estado = $" + strconv.Itoa(paramCount) + ", "
        args = append(args, input.Estado)
    }
    setClause = strings.TrimSuffix(setClause, ", ") + " WHERE id_horario = $1 AND id_medico = $2"

    result, err := config.Conn.Exec(context.Background(), "UPDATE horarios "+setClause, args...)
    if err != nil {
        log.Printf("Error al actualizar horario: %v", err)
        utils.LogAction(userID, "update_horario", "fallido", "Error al actualizar horario: "+err.Error())
        return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Error al actualizar horario: " + err.Error()})
    }
    if result.RowsAffected() == 0 {
        utils.LogAction(userID, "update_horario", "fallido", "Horario no encontrado: ID "+strconv.Itoa(input.IDHorario))
        return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Horario no encontrado"})
    }
    utils.LogAction(userID, "update_horario", "exitoso", "Horario actualizado: ID "+strconv.Itoa(input.IDHorario))
    return c.JSON(fiber.Map{"message": "Horario actualizado"})
}

func DeleteHorario(c *fiber.Ctx) error {
    userID := c.Locals("user_id").(int)
    log.Printf("Solicitud de eliminación para userID: %d", userID)
    role := c.Locals("role").(string)
    if role != "Medico" {
        utils.LogAction(userID, "delete_horario", "fallido", "Permiso denegado: Solo Médicos pueden eliminar horarios")
        return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"error": "Permiso denegado"})
    }

    type HorarioDelete struct {
        IDHorario int `json:"id_horario"`
    }
    var input HorarioDelete
    if err := c.BodyParser(&input); err != nil {
        utils.LogAction(userID, "delete_horario", "fallido", "JSON inválido: "+err.Error())
        return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "JSON inválido"})
    }

    var idMedico int
    err := config.Conn.QueryRow(context.Background(), "SELECT id_medico FROM medicos WHERE id_usuario = $1", userID).Scan(&idMedico)
    if err != nil {
        utils.LogAction(userID, "delete_horario", "fallido", "Médico no encontrado: "+err.Error())
        return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Médico no encontrado"})
    }

    result, err := config.Conn.Exec(context.Background(), "DELETE FROM horarios WHERE id_horario = $1 AND id_medico = $2", input.IDHorario, idMedico)
    if err != nil {
        log.Printf("Error al eliminar horario: %v", err)
        utils.LogAction(userID, "delete_horario", "fallido", "Error al eliminar horario: "+err.Error())
        return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Error al eliminar horario: " + err.Error()})
    }
    if result.RowsAffected() == 0 {
        utils.LogAction(userID, "delete_horario", "fallido", "Horario no encontrado: ID "+strconv.Itoa(input.IDHorario))
        return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Horario no encontrado"})
    }
    utils.LogAction(userID, "delete_horario", "exitoso", "Horario eliminado: ID "+strconv.Itoa(input.IDHorario))
    return c.JSON(fiber.Map{"message": "Horario eliminado"})
}