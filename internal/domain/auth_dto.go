package domain

import "context"

type LoginRequestDTO struct {
	Correo   string `json:"correo" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

type LoginResponseDTO struct {
	AccessToken string `json:"access_token"`
	TokenType   string `json:"token_type"`
	ExpiresIn   int    `json:"expires_in"`
}

type RegisterRequestDTO struct {
	Correo          string `json:"correo" binding:"required,email"`
	Password        string `json:"password" binding:"required,min=6"`
	Nombre          string `json:"nombre" binding:"required"`
	Apellido        string `json:"apellido" binding:"required"`
	FechaNacimiento string `json:"fecha_nacimiento"`
	Rol             string `json:"rol"`
}

type RegisterResponseDTO struct {
	Mensaje string `json:"mensaje"`
	AuthID  string `json:"auth_id,omitempty"`
	Correo  string `json:"correo"`
}

type UsuarioDB struct {
	IDUsuario       int     `gorm:"column:id_usuario;primaryKey;autoIncrement" json:"id_usuario"`
	AuthID          string  `gorm:"column:auth_id" json:"auth_id"`
	Correo          string  `gorm:"column:correo" json:"correo"`
	Nombre          string  `gorm:"column:nombre" json:"nombre"`
	Apellido        string  `gorm:"column:apellido" json:"apellido"`
	FechaNacimiento *string `gorm:"column:fecha_nacimiento" json:"fecha_nacimiento"`
	Rol             string  `gorm:"column:rol" json:"rol"`
}

type IAuthRepository interface {
	EsAdministrador(ctx context.Context, correo string) (bool, error)
	ExisteCorreo(ctx context.Context, correo string) (bool, error)
	CrearUsuario(ctx context.Context, usuario *UsuarioDB) error
}

