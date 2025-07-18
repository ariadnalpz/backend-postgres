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

func CreateConsultorio(c *fiber.Ctx) error {
    userID := c.Locals("user_id").(int)
    log.Printf("Solicitud recibida para userID: %d", userID)
    role := c.Locals("role").(string)
    if role != "Medico" {
        utils.LogAction(userID, "create_consultorio", "fallido", "Permiso denegado: Solo Médicos pueden crear consultorios")
        return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"error": "Permiso denegado"})
    }

    type ConsultorioInput struct {
        NumeroConsultorio string `json:"numero_consultorio"`
        Ubicacion         string `json:"ubicacion"`
        Estado            string `json:"estado,omitempty"`
        FechaActualizacion string `json:"fecha_actualizacion,omitempty"`
    }
    var input ConsultorioInput
    if err := c.BodyParser(&input); err != nil {
        utils.LogAction(userID, "create_consultorio", "fallido", "JSON inválido: "+err.Error())
        return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "JSON inválido"})
    }

    var idMedico int
    err := config.Conn.QueryRow(context.Background(), "SELECT id_medico FROM medicos WHERE id_usuario = $1", userID).Scan(&idMedico)
    if err != nil {
        utils.LogAction(userID, "create_consultorio", "fallido", "Médico no encontrado: "+err.Error())
        return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Médico no encontrado"})
    }

    log.Printf("Datos a insertar: numero_consultorio=%s, ubicacion=%s, estado=%s, fecha_actualizacion=%s, id_medico=%d", input.NumeroConsultorio, input.Ubicacion, input.Estado, input.FechaActualizacion, idMedico)
    _, err = config.Conn.Exec(context.Background(),
        "INSERT INTO consultorios (numero_consultorio, ubicacion, estado, fecha_actualizacion, id_medico) VALUES ($1, $2, $3, $4, $5)",
        input.NumeroConsultorio, input.Ubicacion, input.Estado, input.FechaActualizacion, idMedico)
    if err != nil {
        log.Printf("Error al crear consultorio: %v", err)
        utils.LogAction(userID, "create_consultorio", "fallido", "Error al crear consultorio: "+err.Error())
        return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Error al crear consultorio: " + err.Error()})
    }
    utils.LogAction(userID, "create_consultorio", "exitoso", "Consultorio creado: "+input.NumeroConsultorio)
    return c.JSON(fiber.Map{"message": "Consultorio creado", "estado": input.Estado})
}

func GetConsultorios(c *fiber.Ctx) error {
    userID := c.Locals("user_id").(int)
    log.Printf("Solicitud de lectura para userID: %d", userID)
    role := c.Locals("role").(string)
    if role != "Medico" {
        utils.LogAction(userID, "read_consultorio", "fallido", "Permiso denegado: Solo Médicos pueden ver consultorios")
        return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"error": "Permiso denegado"})
    }

    var idMedico int
    err := config.Conn.QueryRow(context.Background(), "SELECT id_medico FROM medicos WHERE id_usuario = $1", userID).Scan(&idMedico)
    if err != nil {
        utils.LogAction(userID, "read_consultorio", "fallido", "Médico no encontrado: "+err.Error())
        return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Médico no encontrado"})
    }

    rows, err := config.Conn.Query(context.Background(), "SELECT id_consultorio, numero_consultorio, ubicacion, estado, fecha_actualizacion FROM consultorios WHERE id_medico = $1", idMedico)
    if err != nil {
        log.Printf("Error en consulta de consultorios: %v", err)
        utils.LogAction(userID, "read_consultorio", "fallido", "Error al obtener consultorios: "+err.Error())
        return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Error al leer consultorios"})
    }
    defer rows.Close()
    var consultorios []struct {
        IDConsultorio    int    `json:"id_consultorio"`
        NumeroConsultorio string `json:"numero_consultorio"`
        Ubicacion        string `json:"ubicacion"`
        Estado           string `json:"estado"`
        FechaActualizacion string `json:"fecha_actualizacion"`
    }
    for rows.Next() {
        var cons struct {
            IDConsultorio    int    `json:"id_consultorio"`
            NumeroConsultorio string `json:"numero_consultorio"`
            Ubicacion        string `json:"ubicacion"`
            Estado           string `json:"estado"`
            FechaActualizacion string `json:"fecha_actualizacion"`
        }
        err := rows.Scan(&cons.IDConsultorio, &cons.NumeroConsultorio, &cons.Ubicacion, &cons.Estado, &cons.FechaActualizacion)
        if err != nil {
            log.Printf("Error al escanear fila: %v", err)
            utils.LogAction(userID, "read_consultorio", "fallido", "Error al leer consultorio: "+err.Error())
            return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Error al leer consultorios"})
        }
        consultorios = append(consultorios, cons)
    }
    if len(consultorios) == 0 {
        log.Printf("No se encontraron consultorios para id_medico: %d", idMedico)
        utils.LogAction(userID, "read_consultorio", "exitoso", "No se encontraron consultorios")
        return c.JSON([]struct{}{})
    }
    utils.LogAction(userID, "read_consultorio", "exitoso", "Consultorios leídos")
    return c.JSON(consultorios)
}

func UpdateConsultorio(c *fiber.Ctx) error {
    userID := c.Locals("user_id").(int)
    log.Printf("Solicitud de actualización para userID: %d", userID)
    role := c.Locals("role").(string)
    if role != "Medico" {
        utils.LogAction(userID, "update_consultorio", "fallido", "Permiso denegado: Solo Médicos pueden actualizar consultorios")
        return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"error": "Permiso denegado"})
    }

    type ConsultorioUpdate struct {
        IDConsultorio    int    `json:"id_consultorio"`
        NumeroConsultorio string `json:"numero_consultorio,omitempty"`
        Ubicacion        string `json:"ubicacion,omitempty"`
        Estado           string `json:"estado,omitempty"`
        FechaActualizacion string `json:"fecha_actualizacion,omitempty"`
    }
    var input ConsultorioUpdate
    if err := c.BodyParser(&input); err != nil {
        utils.LogAction(userID, "update_consultorio", "fallido", "JSON inválido: "+err.Error())
        return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "JSON inválido"})
    }

    setClause := "SET "
    args := []interface{}{input.IDConsultorio}
    paramCount := 1
    if input.NumeroConsultorio != "" {
        paramCount++
        setClause += "numero_consultorio = $" + strconv.Itoa(paramCount) + ", "
        args = append(args, input.NumeroConsultorio)
    }
    if input.Ubicacion != "" {
        paramCount++
        setClause += "ubicacion = $" + strconv.Itoa(paramCount) + ", "
        args = append(args, input.Ubicacion)
    }
    if input.Estado != "" {
        paramCount++
        setClause += "estado = $" + strconv.Itoa(paramCount) + ", "
        args = append(args, input.Estado)
    }
    if input.FechaActualizacion != "" {
        paramCount++
        setClause += "fecha_actualizacion = $" + strconv.Itoa(paramCount) + ", "
        args = append(args, input.FechaActualizacion)
    }
    setClause = strings.TrimSuffix(setClause, ", ") + " WHERE id_consultorio = $1 AND id_medico = (SELECT id_medico FROM medicos WHERE id_usuario = $2)"
    args = append(args, userID)

    result, err := config.Conn.Exec(context.Background(), "UPDATE consultorios "+setClause, args...)
    if err != nil {
        log.Printf("Error al actualizar consultorio: %v", err)
        utils.LogAction(userID, "update_consultorio", "fallido", "Error al actualizar consultorio: "+err.Error())
        return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Error al actualizar consultorio: " + err.Error()})
    }
    if result.RowsAffected() == 0 {
        utils.LogAction(userID, "update_consultorio", "fallido", "Consultorio no encontrado: ID "+strconv.Itoa(input.IDConsultorio))
        return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Consultorio no encontrado"})
    }
    utils.LogAction(userID, "update_consultorio", "exitoso", "Consultorio actualizado: ID "+strconv.Itoa(input.IDConsultorio))
    return c.JSON(fiber.Map{"message": "Consultorio actualizado"})
}

func DeleteConsultorio(c *fiber.Ctx) error {
    userID := c.Locals("user_id").(int)
    log.Printf("Solicitud de eliminación para userID: %d", userID)
    role := c.Locals("role").(string)
    if role != "Medico" {
        utils.LogAction(userID, "delete_consultorio", "fallido", "Permiso denegado: Solo Médicos pueden eliminar consultorios")
        return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"error": "Permiso denegado"})
    }

    type ConsultorioDelete struct {
        IDConsultorio int `json:"id_consultorio"`
    }
    var input ConsultorioDelete
    if err := c.BodyParser(&input); err != nil {
        utils.LogAction(userID, "delete_consultorio", "fallido", "JSON inválido: "+err.Error())
        return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "JSON inválido"})
    }

    result, err := config.Conn.Exec(context.Background(), "DELETE FROM consultorios WHERE id_consultorio = $1 AND id_medico = (SELECT id_medico FROM medicos WHERE id_usuario = $2)", input.IDConsultorio, userID)
    if err != nil {
        log.Printf("Error al eliminar consultorio: %v", err)
        utils.LogAction(userID, "delete_consultorio", "fallido", "Error al eliminar consultorio: "+err.Error())
        return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Error al eliminar consultorio: " + err.Error()})
    }
    if result.RowsAffected() == 0 {
        utils.LogAction(userID, "delete_consultorio", "fallido", "Consultorio no encontrado: ID "+strconv.Itoa(input.IDConsultorio))
        return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Consultorio no encontrado"})
    }
    utils.LogAction(userID, "delete_consultorio", "exitoso", "Consultorio eliminado: ID "+strconv.Itoa(input.IDConsultorio))
    return c.JSON(fiber.Map{"message": "Consultorio eliminado"})
}