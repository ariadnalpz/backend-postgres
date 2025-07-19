# Backend - Hospitalaria

## Descripción

Este es el **backend** de un sistema de gestión hospitalaria desarrollado en **Go** utilizando el framework **Fiber**. Ofrece funcionalidades **CRUD** para:

- Expedientes
- Consultorios
- Horarios
- Citas
- Pacientes
- Medicos
- Enfermeras

Además, cuenta con autenticación segura mediante:

- **JWT** (token de acceso de 10 min, renovable con refresh token de 24 horas)
- **MFA (TOTP)** para autenticación de dos factores
- Contraseñas **hasheadas** con bcrypt

---

## Requisitos

- [Go 1.18+](https://go.dev/)
- [PostgreSQL](https://www.postgresql.org/)
- Archivo `.env` con las siguientes variables:

```env
DB_HOST=
DB_PORT=
DB_USER=
DB_PASSWORD=
DB_NAME=
JWT_SECRET=
```

---

## Instalación

1. **Clonar el repositorio:**

   ```bash
   git clone https://github.com/ariadnalpz/backend-postgres.git
   cd backend-postgres
   ```

2. **Instalar dependencias:**

   ```bash
   go mod tidy
   ```

(Configura las variables de entorno en un archivo .env)

3. **Ejecutar el proyecto:**

   ```bash
   go run main.go
   ```

El servidor se iniciará en [http://localhost:3000](http://localhost:3000).

---

## Endpoints

*Registro:* POST /register - Crea un nuevo usuario (Paciente, Médico, Enfermero).
*Login:* POST /login - Autenticación con contraseña y TOTP.
*Refresh Token:* POST /refresh-token - Renueva el access_token con un refresh_token.
*Perfil:* GET /profile - Obtiene el perfil del usuario autenticado (requiere token).
*Rutas protegidas:* Accede a /paciente, /medico, /enfermera con un access_token válido (ejemplo: GET /medico/consultorios con header `Authorization: Bearer <token>`).

---

## Estructura del Proyecto

   backend-hospitalaria/
   ├── config/              # Configuración de la base de datos y conexión
   │   └── db.go
   ├── handlers/            # Lógica de negocio y endpoints
   │   ├── auth.go
   │   └── medicos/
   ├── middleware/          # Middlewares (ej. validación JWT)
   │   └── jwt.go
   ├── models/              # Estructuras de datos (ej. modelos de usuario)
   │   └── user.go
   ├── routes/              # Definición de rutas por rol
   │   ├── auth.go
   │   ├── paciente.go
   │   ├── medico.go
   │   └── enfermera.go
   ├── utils/               # Utilidades como logging
   │   └── logging.go
   ├── main.go              # Punto de entrada de la aplicación
   ├── go.mod / go.sum      # Gestión de dependencias
   ├── .gitignore           # Archivos ignorados por Git
   ├── README.md            # Este archivo
   └── CHANGELOG.md         # Historial de cambios

---

## Contribuciones

*Crea una rama:* git checkout -b feature/nueva-funcionalidad.

*Haz commits:* git commit -m "Descripción".

*Sube y crea un Pull Request:* git push origin feature/nueva-funcionalidad.