package models

type User struct {
	Id_usuario     int    `json:"id_usuario"`
	Nombre         string `json:"nombre"`
	Apellido       string `json:"apellido"`
	Correo         string `json:"correo"`
	Contraseña     string `json:"contraseña,omitempty"`
	Rol            string `json:"rol"`
	Totp_secret    string `json:"totp_secret,omitempty"`
	FechaNacimiento string `json:"fecha_nacimiento,omitempty"`
	Genero         string `json:"genero,omitempty"`
	Direccion      string `json:"direccion,omitempty"`
	Especialidad   string `json:"especialidad,omitempty"`
	NumeroColegiado string `json:"numero_colegiado,omitempty"`
	Certificacion  string `json:"certificacion,omitempty"`
}