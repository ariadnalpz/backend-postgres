package handlers

import (
	"context"
	"encoding/base64"
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
	"github.com/pquerna/otp/totp"
	"github.com/skip2/go-qrcode"
	"golang.org/x/crypto/bcrypt"
	"hospitalaria/config"
	"hospitalaria/models"
	"hospitalaria/utils"
)

func CheckPasswordStrength(password string) (bool, string) {
	if len(password) < 12 {
		return false, "La contraseña debe tener al menos 12 caracteres"
	}
	if !strings.ContainsAny(password, "!@#$%^&*()") {
		return false, "La contraseña debe incluir símbolos"
	}
	if !strings.ContainsAny(password, "0123456789") {
		return false, "La contraseña debe incluir números"
	}
	return true, "Contraseña segura"
}

func GenerateTokens(userID int, role string) (string, string, error) {
	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": userID,
		"role":    role,
		"exp":     time.Now().Add(10 * time.Minute).Unix(),
	})
	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": userID,
		"role":    role, // Añadido para consistencia
		"exp":     time.Now().Add(24 * time.Hour).Unix(),
	})
	access, err := accessToken.SignedString([]byte(os.Getenv("JWT_SECRET")))
	if err != nil {
		return "", "", err
	}
	refresh, err := refreshToken.SignedString([]byte(os.Getenv("JWT_SECRET")))
	if err != nil {
		return "", "", err
	}
	return access, refresh, nil
}

func CreateUser(c *fiber.Ctx) error {
	type UserInput struct {
		Password        string `json:"password"`
		Nombre          string `json:"nombre"`
		Apellido        string `json:"apellido"`
		Correo          string `json:"correo"`
		Rol             string `json:"rol"`
		FechaNacimiento string `json:"fecha_nacimiento,omitempty"`
		Genero          string `json:"genero,omitempty"`
		Direccion       string `json:"direccion,omitempty"`
		Especialidad    string `json:"especialidad,omitempty"`
		NumeroColegiado string `json:"numero_colegiado,omitempty"`
		Certificacion   string `json:"certificacion,omitempty"`
	}
	var input UserInput
	if err := c.BodyParser(&input); err != nil {
		utils.LogAction(0, "create_user", "fallido", "JSON inválido: "+err.Error())
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "JSON inválido"})
	}

	isStrong, message := CheckPasswordStrength(input.Password)
	if !isStrong {
		utils.LogAction(0, "create_user", "fallido", "Contraseña débil: "+message)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": message})
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
	if err != nil {
		log.Printf("Error al hashear contraseña: %v", err)
		utils.LogAction(0, "create_user", "fallido", "Error al hashear contraseña: "+err.Error())
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Error al procesar contraseña"})
	}

	key, err := totp.Generate(totp.GenerateOpts{
		Issuer:      "MyHospitalApp",
		AccountName: input.Correo,
	})
	if err != nil {
		log.Printf("Error al generar secreto TOTP: %v", err)
		utils.LogAction(0, "create_user", "fallido", "Error al generar TOTP: "+err.Error())
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Error al generar código TOTP"})
	}

	var userID int
	user := models.User{
		Nombre:      input.Nombre,
		Apellido:    input.Apellido,
		Correo:      input.Correo,
		Contraseña:  string(hash),
		Rol:         input.Rol,
		Totp_secret: key.Secret(),
	}
	err = config.Conn.QueryRow(context.Background(),
		"INSERT INTO usuarios (nombre, apellido, correo, contraseña, rol, totp_secret) VALUES ($1, $2, $3, $4, $5, $6) RETURNING id_usuario",
		user.Nombre, user.Apellido, user.Correo, user.Contraseña, user.Rol, user.Totp_secret).Scan(&userID)
	if err != nil {
		log.Printf("Error al insertar usuario: %v", err)
		utils.LogAction(0, "create_user", "fallido", "Error al insertar usuario: "+err.Error())
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Error al guardar usuario"})
	}
	utils.LogAction(userID, "create_user", "exitoso", "Usuario creado con ID "+strconv.Itoa(userID))

	if user.Rol == "Paciente" {
		_, err = config.Conn.Exec(context.Background(),
			"INSERT INTO pacientes (id_usuario, fecha_nacimiento, genero, direccion) VALUES ($1, $2, $3, $4)",
			userID, input.FechaNacimiento, input.Genero, input.Direccion)
		if err != nil {
			log.Printf("Error al insertar paciente: %v", err)
			utils.LogAction(userID, "create_paciente", "fallido", "Error al insertar datos de paciente: "+err.Error())
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Error al guardar datos de rol: " + err.Error()})
		}
		utils.LogAction(userID, "create_paciente", "exitoso", "Datos de paciente creados para usuario "+strconv.Itoa(userID))
	} else if user.Rol == "Medico" {
		_, err = config.Conn.Exec(context.Background(),
			"INSERT INTO medicos (id_usuario, especialidad, numero_colegiado) VALUES ($1, $2, $3)",
			userID, input.Especialidad, input.NumeroColegiado)
		if err != nil {
			log.Printf("Error al insertar medico: %v", err)
			utils.LogAction(userID, "create_medico", "fallido", "Error al insertar datos de medico: "+err.Error())
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Error al guardar datos de rol: " + err.Error()})
		}
		utils.LogAction(userID, "create_medico", "exitoso", "Datos de medico creados para usuario "+strconv.Itoa(userID))
	} else if user.Rol == "Enfermero" {
		_, err = config.Conn.Exec(context.Background(),
			"INSERT INTO enfermeras (id_usuario, certificacion) VALUES ($1, $2)",
			userID, input.Certificacion)
		if err != nil {
			log.Printf("Error al insertar enfermero: %v", err)
			utils.LogAction(userID, "create_enfermero", "fallido", "Error al insertar datos de enfermero: "+err.Error())
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Error al guardar datos de rol: " + err.Error()})
		}
		utils.LogAction(userID, "create_enfermero", "exitoso", "Datos de enfermero creados para usuario "+strconv.Itoa(userID))
	}

	qrCode, err := qrcode.Encode(key.URL(), qrcode.Medium, 256)
	if err != nil {
		log.Printf("Error al generar código QR: %v", err)
		utils.LogAction(userID, "create_user", "fallido", "Error al generar QR: "+err.Error())
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Error al generar código QR"})
	}
	qrBase64 := base64.StdEncoding.EncodeToString(qrCode)

	return c.JSON(fiber.Map{
		"id_usuario":  userID,
		"nombre":      user.Nombre,
		"correo":      user.Correo,
		"rol":         user.Rol,
		"totp_secret": user.Totp_secret,
		"totp_qr":     "data:image/png;base64," + qrBase64,
	})
}

func Login(c *fiber.Ctx) error {
	var input struct {
		Correo   string `json:"correo"`
		Password string `json:"password"`
		TOTPCode string `json:"totp_code"`
	}
	if err := c.BodyParser(&input); err != nil {
		utils.LogAction(0, "login", "fallido", "JSON inválido: "+err.Error())
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "JSON inválido"})
	}

	var user models.User
	err := config.Conn.QueryRow(context.Background(),
		"SELECT id_usuario, contraseña, rol, totp_secret FROM usuarios WHERE correo = $1", input.Correo).Scan(
		&user.Id_usuario, &user.Contraseña, &user.Rol, &user.Totp_secret)
	if err != nil {
		utils.LogAction(0, "login", "fallido", "Correo no encontrado o error en consulta: "+err.Error())
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Credenciales inválidas"})
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Contraseña), []byte(input.Password)); err != nil {
		utils.LogAction(user.Id_usuario, "login", "fallido", "Contraseña incorrecta para "+input.Correo)
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Credenciales inválidas"})
	}
	utils.LogAction(user.Id_usuario, "login", "exitoso", "Contraseña validada para "+input.Correo)

	if !totp.Validate(input.TOTPCode, user.Totp_secret) {
		utils.LogAction(user.Id_usuario, "login", "fallido", "Código TOTP inválido para "+input.Correo)
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Código TOTP inválido"})
	}
	utils.LogAction(user.Id_usuario, "login", "exitoso", "Código TOTP validado para "+input.Correo)

	accessToken, refreshToken, err := GenerateTokens(user.Id_usuario, user.Rol)
	if err != nil {
		log.Printf("Error al generar tokens: %v", err)
		utils.LogAction(user.Id_usuario, "login", "fallido", "Error al generar tokens: "+err.Error())
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Error al generar tokens"})
	}

	utils.LogAction(user.Id_usuario, "login", "exitoso", "Inicio de sesión exitoso para "+input.Correo)

	return c.JSON(fiber.Map{
		"access_token":  accessToken,
		"refresh_token": refreshToken,
	})
}

func GetUserProfile(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(int)
	var user models.User
	err := config.Conn.QueryRow(context.Background(),
		"SELECT id_usuario, nombre, apellido, correo, rol FROM usuarios WHERE id_usuario = $1", userID).Scan(
		&user.Id_usuario, &user.Nombre, &user.Apellido, &user.Correo, &user.Rol)
	if err != nil {
		log.Printf("Error al obtener perfil: %v", err)
		utils.LogAction(userID, "read_user", "fallido", "Error al obtener perfil: "+err.Error())
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Usuario no encontrado"})
	}
	utils.LogAction(userID, "read_user", "exitoso", "Perfil leído para usuario "+strconv.Itoa(userID))
	return c.JSON(user)
}

func RefreshToken(c *fiber.Ctx) error {
	type RefreshTokenInput struct {
		RefreshToken string `json:"refresh_token"`
	}

	var input RefreshTokenInput
	if err := c.BodyParser(&input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "JSON inválido"})
	}

	// Parsear y validar el refresh token
	token, err := jwt.Parse(input.RefreshToken, func(token *jwt.Token) (interface{}, error) {
		return []byte(os.Getenv("JWT_SECRET")), nil
	})

	if err != nil || !token.Valid {
		log.Printf("Refresh token inválido: %v", err)
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Token de refresco inválido"})
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Claim inválido"})
	}

	// Extraer user_id y verificar expiración
	userIDFloat, ok := claims["user_id"].(float64)
	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "user_id no encontrado en token"})
	}
	userID := int(userIDFloat)

	expiry, ok := claims["exp"].(float64)
	if !ok || time.Unix(int64(expiry), 0).Before(time.Now()) {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Token de refresco expirado"})
	}

	// Obtener el rol desde el token (ya que no se guarda en refresh_token originalmente, lo simulamos)
	role, ok := claims["role"].(string)
	if !ok {
		// Si no está en el refresh_token, consultarlo desde la base de datos
		var userRole string
		err := config.Conn.QueryRow(context.Background(), "SELECT rol FROM usuarios WHERE id_usuario = $1", userID).Scan(&userRole)
		if err != nil {
			log.Printf("Error al obtener rol para userID %d: %v", userID, err)
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Error al obtener rol"})
		}
		role = userRole
	}

	// Generar nuevo access token
	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": userID,
		"role":    role,
		"exp":     time.Now().Add(10 * time.Minute).Unix(),
	})

	accessTokenStr, err := accessToken.SignedString([]byte(os.Getenv("JWT_SECRET")))
	if err != nil {
		log.Printf("Error al firmar nuevo access token: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Error al generar token"})
	}

	// Opcional: Generar nuevo refresh token para rotación
	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": userID,
		"exp":     time.Now().Add(24 * time.Hour).Unix(),
	})
	refreshTokenStr, err := refreshToken.SignedString([]byte(os.Getenv("JWT_SECRET")))
	if err != nil {
		log.Printf("Error al firmar nuevo refresh token: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Error al generar token de refresco"})
	}

	utils.LogAction(userID, "refresh_token", "exitoso", "Token renovado")
	return c.JSON(fiber.Map{
		"access_token":  accessTokenStr,
		"refresh_token": refreshTokenStr, // Opcional, para rotación
	})
}
