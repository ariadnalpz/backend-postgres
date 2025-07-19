# Changelog - Backend Postgres

## Unreleased
### @Agregado
- Pendiente de nuevas características y correcciones.

---

## [0.1.0] - 2025-07-17
### @Agregado
- Implementación inicial del backend utilizando Go y el framework Fiber.
- Autenticación segura con JWT, incluyendo:
  - Token de acceso con expiración de 10 minutos.
  - Refresh token con expiración de 24 horas para renovar el acceso.
- Integración de MFA (TOTP) para autenticación de dos factores.
- Validación de contraseñas seguras (mínimo 12 caracteres, símbolos y números).
- Almacenamiento de contraseñas hasheadas utilizando bcrypt.
- Endpoints CRUD básicos para:
  - Expedientes
  - Consultorios
  - Horarios
  - Citas
- Middleware JWT para validar permisos en rutas protegidas.
- Sistema de logging para:
  - Intentos de inicio de sesión (éxitos y fallos).
  - Operaciones CRUD.
- Configuración inicial de conexión con PostgreSQL.
- Subida del proyecto al repositorio GitHub: https://github.com/ariadnalpz/backend-postgres

### @Cambios
- Organización de la estructura del proyecto en carpetas (`config`, `handlers`, `middleware`, `models`, `routes`, `utils`).
- Ajustes en las rutas para separar autenticación y roles específicos (paciente, médico, enfermera).

### @Corregido
- Errores iniciales en la asociación de médicos con consultorios (resolución del problema "Médico no encontrado").
- Problemas en la validación de tokens en endpoints protegidos, asegurando autenticación correcta.

---

## [0.1.1] - 2025-07-18
### @Agregado
- Creación de `README.md` con descripción, instalación y estructura del proyecto.
- Creación de `CHANGELOG.md` para registrar cambios de manera estructurada.

### @Cambios
- Mejora en la documentación del proyecto.

[0.1.1]: https://github.com/ariadnalpz/backend-postgres/releases/tag/0.1.1
