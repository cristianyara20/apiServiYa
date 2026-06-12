package http

import (
	"apiServiYa/internal/application/auth"
	"apiServiYa/internal/domain"
	"net/http"

	"github.com/gin-gonic/gin"
)

type AuthController struct {
	loginUseCase *auth.LoginUseCase
}

func NewAuthController(loginUC *auth.LoginUseCase) *AuthController {
	return &AuthController{loginUseCase: loginUC}
}

// Login godoc
// @Summary Iniciar sesión como administrador
// @Description Inicia sesión en Supabase y devuelve un JWT si el usuario es un administrador
// @Tags autenticacion
// @Accept json
// @Produce json
// @Param credenciales body domain.LoginRequestDTO true "Credenciales del administrador"
// @Success 200 {object} domain.LoginResponseDTO
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 403 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /auth/login [post]
func (ctrl *AuthController) Login(c *gin.Context) {
	var req domain.LoginRequestDTO
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Formato de petición inválido", "detalle": err.Error()})
		return
	}

	resp, err := ctrl.loginUseCase.Ejecutar(c.Request.Context(), req)
	if err != nil {
		status := http.StatusUnauthorized
		if err.Error() == "acceso denegado: no tienes permisos de administrador" {
			status = http.StatusForbidden
		} else if err.Error() == "error al contactar con el servicio de autenticación" || err.Error() == "configuración de Supabase incompleta en el servidor" {
			status = http.StatusInternalServerError
		}
		c.JSON(status, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, resp)
}
