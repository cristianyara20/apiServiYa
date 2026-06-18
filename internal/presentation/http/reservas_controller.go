package http

import (
	"apiServiYa/internal/application/reservas"
	"apiServiYa/internal/domain"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type ReservasController struct {
	cancelarUC *reservas.CancelarReservaUseCase
}

func NewReservasController(cancelarUC *reservas.CancelarReservaUseCase) *ReservasController {
	return &ReservasController{cancelarUC: cancelarUC}
}

// CancelarReserva godoc
// @Summary Cancela una reserva existente
// @Description Permite a un cliente cancelar una de sus reservas siempre y cuando esté en estado pendiente o aceptada.
// @Tags reservas
// @Accept json
// @Produce json
// @Param id path int true "ID de la Reserva"
// @Param body body domain.CancelarReservaRequestDTO true "Datos de cancelación (ID del Cliente)"
// @Success 200 {object} domain.CancelarReservaResponseDTO
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /reservas/{id}/cancelar [put]
// @Security BearerAuth
func (ctrl *ReservasController) CancelarReserva(c *gin.Context) {
	idParam := c.Param("id")
	idReserva, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID de reserva inválido"})
		return
	}

	var req domain.CancelarReservaRequestDTO
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Cuerpo de la petición inválido o incompleto", "detalle": err.Error()})
		return
	}

	res, err := ctrl.cancelarUC.Ejecutar(c.Request.Context(), uint(idReserva), req.IDCliente)
	if err != nil {
		if err.Error() == "reserva no encontrada o no pertenece al cliente" {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		if err.Error() == "solo se pueden cancelar reservas en estado pendiente o aceptada" {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Fallo al cancelar reserva", "detalle": err.Error()})
		return
	}

	c.JSON(http.StatusOK, res)
}
