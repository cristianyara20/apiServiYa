package auth

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"time"

	"apiServiYa/internal/domain"
)

type RegisterUseCase struct {
	repo domain.IAuthRepository
}

func NewRegisterUseCase(repo domain.IAuthRepository) *RegisterUseCase {
	return &RegisterUseCase{repo: repo}
}

type supabaseUserResponse struct {
	ID    string `json:"id"`
	Email string `json:"email"`
	User  *struct {
		ID    string `json:"id"`
		Email string `json:"email"`
	} `json:"user,omitempty"`
}

type supabaseErrorResponse struct {
	Message          string `json:"message"`
	Msg              string `json:"msg"`
	ErrorDescription string `json:"error_description"`
	Error            string `json:"error"`
}

func (uc *RegisterUseCase) Ejecutar(ctx context.Context, req domain.RegisterRequestDTO) (*domain.RegisterResponseDTO, error) {
	supabaseURL := strings.TrimRight(os.Getenv("SUPABASE_URL"), "/")
	anonKey := os.Getenv("SUPABASE_ANON_KEY")
	serviceKey := os.Getenv("SUPABASE_SERVICE_ROLE_KEY")

	if supabaseURL == "" || (anonKey == "" && serviceKey == "") {
		return nil, errors.New("configuración de Supabase incompleta en el servidor")
	}

	req.Correo = strings.TrimSpace(strings.ToLower(req.Correo))
	req.Nombre = strings.TrimSpace(req.Nombre)
	req.Apellido = strings.TrimSpace(req.Apellido)
	if req.Rol == "" {
		req.Rol = "usuario"
	}

	// 1. Validar si el correo ya existe en base de datos
	existe, err := uc.repo.ExisteCorreo(ctx, req.Correo)
	if err != nil {
		return nil, fmt.Errorf("error al verificar existencia del correo: %w", err)
	}
	if existe {
		return nil, errors.New("el correo electrónico ya se encuentra registrado")
	}

	// 2. Normalizar fecha de nacimiento a formato estándar YYYY-MM-DD
	fechaNacimientoPtr := parsearFechaNacimiento(req.FechaNacimiento)


	var authID string

	// 3. Registrar en Supabase Auth
	// Prioridad A: Usar Supabase Admin API (con Service Role Key) para confirmar el email inmediatamente
	if serviceKey != "" {
		adminEndpoint := fmt.Sprintf("%s/auth/v1/admin/users", supabaseURL)
		adminBody := map[string]interface{}{
			"email":         req.Correo,
			"password":      req.Password,
			"email_confirm": true,
			"user_metadata": map[string]interface{}{
				"nombre":   req.Nombre,
				"apellido": req.Apellido,
				"rol":      req.Rol,
			},
		}
		if fechaNacimientoPtr != nil {
			adminBody["user_metadata"].(map[string]interface{})["fecha_nacimiento"] = *fechaNacimientoPtr
		}

		bodyBytes, _ := json.Marshal(adminBody)
		httpReq, reqErr := http.NewRequestWithContext(ctx, "POST", adminEndpoint, bytes.NewBuffer(bodyBytes))
		if reqErr == nil {
			httpReq.Header.Set("Content-Type", "application/json")
			httpReq.Header.Set("apikey", serviceKey)
			httpReq.Header.Set("Authorization", "Bearer "+serviceKey)

			client := &http.Client{Timeout: 10 * time.Second}
			resp, doErr := client.Do(httpReq)
			if doErr == nil {
				defer resp.Body.Close()
				respBytes, _ := io.ReadAll(resp.Body)
				if resp.StatusCode == http.StatusOK || resp.StatusCode == http.StatusCreated {
					var sUser supabaseUserResponse
					if err := json.Unmarshal(respBytes, &sUser); err == nil {
						if sUser.ID != "" {
							authID = sUser.ID
						} else if sUser.User != nil && sUser.User.ID != "" {
							authID = sUser.User.ID
						}
					}
				}
			}
		}
	}

	// Prioridad B: Si no se obtuvo authID con Admin, usar endpoint estándar /auth/v1/signup con Anon Key
	if authID == "" {
		signupEndpoint := fmt.Sprintf("%s/auth/v1/signup", supabaseURL)
		signupBody := map[string]interface{}{
			"email":    req.Correo,
			"password": req.Password,
			"data": map[string]interface{}{
				"nombre":   req.Nombre,
				"apellido": req.Apellido,
				"rol":      req.Rol,
			},
		}
		if fechaNacimientoPtr != nil {
			signupBody["data"].(map[string]interface{})["fecha_nacimiento"] = *fechaNacimientoPtr
		}

		bodyBytes, _ := json.Marshal(signupBody)
		httpReq, reqErr := http.NewRequestWithContext(ctx, "POST", signupEndpoint, bytes.NewBuffer(bodyBytes))
		if reqErr != nil {
			return nil, fmt.Errorf("error al preparar petición de registro: %w", reqErr)
		}
		httpReq.Header.Set("Content-Type", "application/json")
		httpReq.Header.Set("apikey", anonKey)

		client := &http.Client{Timeout: 10 * time.Second}
		resp, doErr := client.Do(httpReq)
		if doErr != nil {
			return nil, errors.New("error al comunicarse con el servicio de autenticación")
		}
		defer resp.Body.Close()

		respBytes, _ := io.ReadAll(resp.Body)

		if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
			var supErr supabaseErrorResponse
			_ = json.Unmarshal(respBytes, &supErr)
			msg := supErr.Message
			if msg == "" {
				msg = supErr.ErrorDescription
			}
			if msg == "" {
				msg = supErr.Msg
			}
			if msg == "" {
				msg = supErr.Error
			}
			if msg == "" {
				msg = fmt.Sprintf("código de estado %d", resp.StatusCode)
			}
			return nil, fmt.Errorf("error en autenticación: %s", msg)
		}

		var sUser supabaseUserResponse
		if err := json.Unmarshal(respBytes, &sUser); err == nil {
			if sUser.ID != "" {
				authID = sUser.ID
			} else if sUser.User != nil && sUser.User.ID != "" {
				authID = sUser.User.ID
			}
		}
	}

	if authID == "" {
		return nil, errors.New("no se pudo obtener el identificador de usuario de autenticación")
	}

	// 4. Guardar en Base de Datos (seguridad.usuarios y gestion.clientes)
	usuarioDB := &domain.UsuarioDB{
		AuthID:          authID,
		Correo:          req.Correo,
		Nombre:          req.Nombre,
		Apellido:        req.Apellido,
		FechaNacimiento: fechaNacimientoPtr,
		Rol:             req.Rol,
	}

	if err := uc.repo.CrearUsuario(ctx, usuarioDB); err != nil {
		return nil, fmt.Errorf("error al registrar perfil de usuario en base de datos: %w", err)
	}

	return &domain.RegisterResponseDTO{
		Mensaje: "Usuario registrado con éxito",
		AuthID:  authID,
		Correo:  req.Correo,
	}, nil
}

func parsearFechaNacimiento(raw string) *string {
	raw = strings.TrimSpace(raw)
	if raw == "" || raw == "--" || raw == "null" {
		return nil
	}

	// Formatos comunes: YYYY-MM-DD, DD-MM-YYYY, MM-DD-YYYY, etc.
	formatos := []string{
		"2006-01-02",
		"02-01-2006",
		"01-02-2006",
		"2006/01/02",
		"02/01/2006",
		"01/02/2006",
		"2006.01.02",
		"02.01.2006",
		time.RFC3339,
		"2006-01-02 15:04:05",
	}

	for _, f := range formatos {
		if t, err := time.Parse(f, raw); err == nil {
			formatted := t.Format("2006-01-02")
			return &formatted
		}
	}
	return nil
}

