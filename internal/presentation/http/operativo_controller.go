package http

import (
	"apiServiYa/internal/application/reportes"
	"apiServiYa/internal/domain"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type OperativoController struct {
	prestadoresUC *reportes.ObtenerPrestadoresOperativosUseCase
	historialUC   *reportes.ObtenerHistorialServiciosUseCase
	pqrsUC        *reportes.ObtenerPqrsUseCase
	responderUC   *reportes.ResponderPqrUseCase
}

func NewOperativoController(
	prestadoresUC *reportes.ObtenerPrestadoresOperativosUseCase,
	historialUC *reportes.ObtenerHistorialServiciosUseCase,
	pqrsUC *reportes.ObtenerPqrsUseCase,
	responderUC *reportes.ResponderPqrUseCase,
) *OperativoController {
	return &OperativoController{
		prestadoresUC: prestadoresUC,
		historialUC:   historialUC,
		pqrsUC:        pqrsUC,
		responderUC:   responderUC,
	}
}

// ObtenerPrestadores godoc
// @Summary Obtiene la lista de prestadores con su estado de disponibilidad en tiempo real
// @Description Devuelve la lista completa de prestadores de la base de datos, con nombre, correo, experiencia, calificación promedio y estado (disponible, ocupado, inactivo).
// @Tags operativo
// @Produce json
// @Success 200 {array} domain.PrestadorDetalladoDTO
// @Failure 401 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /operativo/prestadores [get]
// @Security BearerAuth
func (ctrl *OperativoController) ObtenerPrestadores(c *gin.Context) {
	resultados, err := ctrl.prestadoresUC.Ejecutar(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Fallo al obtener prestadores", "detalle": err.Error()})
		return
	}
	c.JSON(http.StatusOK, resultados)
}

// ObtenerHistorialServicios godoc
// @Summary Obtiene el historial completo de servicios/reservas prestadas en la plataforma
// @Description Devuelve todas las reservas del sistema incluyendo los nombres de cliente, servicio, fecha agendada, dirección y estado del servicio.
// @Tags operativo
// @Produce json
// @Success 200 {array} domain.ReservaHistorialDTO
// @Failure 401 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /operativo/historial-servicios [get]
// @Security BearerAuth
func (ctrl *OperativoController) ObtenerHistorialServicios(c *gin.Context) {
	resultados, err := ctrl.historialUC.Ejecutar(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Fallo al obtener historial de servicios", "detalle": err.Error()})
		return
	}
	c.JSON(http.StatusOK, resultados)
}

// InformarAceptacion godoc
// @Summary Informa al cliente que un prestador aceptó su reserva
// @Description Endpoint que simula o registra la notificación del prestador al cliente de que se aceptó la reserva. Por ahora retorna un mensaje exitoso de confirmación de notificación.
// @Tags operativo
// @Accept json
// @Produce json
// @Param body body domain.NotificationRequestDTO true "Datos de la reserva"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Router /operativo/notificar-aceptacion [post]
// @Security BearerAuth
func (ctrl *OperativoController) InformarAceptacion(c *gin.Context) {
	var req domain.NotificationRequestDTO
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Formato de petición inválido", "detalle": err.Error()})
		return
	}

	// Aquí se integraría el envío de notificaciones por WebSockets, Email, SMS, Firebase, etc.
	// Por ahora devolvemos la respuesta exitosa confirmando que se informó al cliente con ID especificado.
	c.JSON(http.StatusOK, gin.H{
		"mensaje":     "Cliente notificado con éxito de que su servicio fue aceptado",
		"id_cliente":  strconvFormat(req.IDCliente),
		"id_reserva":  strconvFormat(req.IDReserva),
		"id_prestador": strconvFormat(req.IDPrestador),
		"estado":      req.EstadoReserva,
	})
}

// ObtenerPqrs godoc
// @Summary Obtiene todas las PQRs del sistema con el nombre del cliente
// @Description Devuelve la lista completa de PQRs incluyendo datos del cliente, tipo, estado, descripción, y respuesta del administrador si existe.
// @Tags operativo
// @Produce json
// @Success 200 {array} domain.PqrDTO
// @Failure 401 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /operativo/pqrs [get]
// @Security BearerAuth
func (ctrl *OperativoController) ObtenerPqrs(c *gin.Context) {
	resultados, err := ctrl.pqrsUC.Ejecutar(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Fallo al obtener PQRs", "detalle": err.Error()})
		return
	}
	c.JSON(http.StatusOK, resultados)
}

// ResponderPqr godoc
// @Summary Permite al administrador responder una PQR
// @Description Recibe el ID de la PQR y la respuesta del administrador, actualiza el estado a 'Cerrado' y registra la fecha de respuesta.
// @Tags operativo
// @Accept json
// @Produce json
// @Param body body domain.ResponderPqrRequestDTO true "ID de la PQR y respuesta del admin"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /operativo/pqrs/responder [post]
// @Security BearerAuth
func (ctrl *OperativoController) ResponderPqr(c *gin.Context) {
	var req domain.ResponderPqrRequestDTO
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Formato de petición inválido", "detalle": err.Error()})
		return
	}

	if err := ctrl.responderUC.Ejecutar(c.Request.Context(), req.IDPqr, req.RespuestaAdmin); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "No se pudo responder la PQR", "detalle": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"mensaje": "PQR respondida exitosamente",
		"id_pqr":  strconvFormat(req.IDPqr),
		"estado":  "Cerrado",
	})
}

func strconvFormat(val uint) string {
	return strconv.FormatUint(uint64(val), 10)
}
