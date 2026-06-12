package auth

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"

	"apiServiYa/internal/domain"
)

type LoginUseCase struct {
	repo domain.IAuthRepository
}

func NewLoginUseCase(repo domain.IAuthRepository) *LoginUseCase {
	return &LoginUseCase{repo: repo}
}

type supabaseAuthResponse struct {
	AccessToken string `json:"access_token"`
	TokenType   string `json:"token_type"`
	ExpiresIn   int    `json:"expires_in"`
}

type supabaseAuthError struct {
	ErrorDescription string `json:"error_description"`
}

func (uc *LoginUseCase) Ejecutar(ctx context.Context, req domain.LoginRequestDTO) (*domain.LoginResponseDTO, error) {
	supabaseURL := os.Getenv("SUPABASE_URL")
	anonKey := os.Getenv("SUPABASE_ANON_KEY")

	if supabaseURL == "" || anonKey == "" {
		return nil, errors.New("configuración de Supabase incompleta en el servidor")
	}

	// 1. Llamar a Supabase Auth para verificar correo y contraseña
	authEndpoint := fmt.Sprintf("%s/auth/v1/token?grant_type=password", supabaseURL)
	body, _ := json.Marshal(map[string]string{
		"email":    req.Correo,
		"password": req.Password,
	})

	httpReq, err := http.NewRequestWithContext(ctx, "POST", authEndpoint, bytes.NewBuffer(body))
	if err != nil {
		return nil, err
	}
	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("apikey", anonKey)

	client := &http.Client{}
	resp, err := client.Do(httpReq)
	if err != nil {
		return nil, errors.New("error al contactar con el servicio de autenticación")
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		var supErr supabaseAuthError
		json.NewDecoder(resp.Body).Decode(&supErr)
		return nil, fmt.Errorf("credenciales inválidas: %s", supErr.ErrorDescription)
	}

	var supResp supabaseAuthResponse
	if err := json.NewDecoder(resp.Body).Decode(&supResp); err != nil {
		return nil, errors.New("error al procesar la respuesta de autenticación")
	}

	// 2. Verificar si es administrador en nuestra base de datos
	esAdmin, err := uc.repo.EsAdministrador(ctx, req.Correo)
	if err != nil {
		return nil, errors.New("error al verificar permisos del usuario")
	}

	if !esAdmin {
		return nil, errors.New("acceso denegado: no tienes permisos de administrador")
	}

	// 3. Devolver nuestro token de Supabase directo
	return &domain.LoginResponseDTO{
		AccessToken: supResp.AccessToken,
		TokenType:   supResp.TokenType,
		ExpiresIn:   supResp.ExpiresIn,
	}, nil
}
