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

type IAuthRepository interface {
	EsAdministrador(ctx context.Context, correo string) (bool, error)
}
