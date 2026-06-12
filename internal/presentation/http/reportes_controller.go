package http

import (
	"apiServiYa/internal/application/reportes"
	"apiServiYa/internal/infrastructure/pdf"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type ReportesController struct {
	adminUseCase     *reportes.GenerarReporteConsolidadoUseCase
	serviciosUseCase *reportes.ObtenerServiciosPopularesUseCase
	actividadUseCase *reportes.ObtenerActividadUsuariosUseCase
}

func NewReportesController(
	adminUC *reportes.GenerarReporteConsolidadoUseCase,
	serviciosUC *reportes.ObtenerServiciosPopularesUseCase,
	actividadUC *reportes.ObtenerActividadUsuariosUseCase,
) *ReportesController {
	return &ReportesController{
		adminUseCase:     adminUC,
		serviciosUseCase: serviciosUC,
		actividadUseCase: actividadUC,
	}
}

// ObtenerReporteAdmin godoc
// @Summary Obtiene un reporte consolidado para el administrador
// @Description Devuelve el top de prestadores, cantidad de pqrs abiertas y estadísticas de reservas para un mes y año
// @Tags reportes
// @Produce json
// @Param mes query int true "Mes (1-12)"
// @Param anio query int true "Año (ej. 2026)"
// @Success 200 {object} domain.ReporteConsolidadoMensualDTO
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /reportes/admin [get]
// @Security BearerAuth
func (ctrl *ReportesController) ObtenerReporteAdmin(c *gin.Context) {
	mes, anio, ok := ctrl.parseMesAnio(c)
	if !ok {
		return
	}

	reporteDTO, err := ctrl.adminUseCase.Ejecutar(c.Request.Context(), mes, anio)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Fallo al generar el reporte", "detalle": err.Error()})
		return
	}

	c.JSON(http.StatusOK, reporteDTO)
}

// ObtenerServiciosPopulares godoc
// @Summary Obtiene los servicios más solicitados del mes
// @Description Devuelve los 10 servicios con más reservas en el mes y año indicados, incluyendo nombre y categoría
// @Tags reportes
// @Produce json
// @Param mes query int true "Mes (1-12)"
// @Param anio query int true "Año (ej. 2026)"
// @Success 200 {array} domain.ServicioPopularDTO
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /reportes/servicios-populares [get]
// @Security BearerAuth
func (ctrl *ReportesController) ObtenerServiciosPopulares(c *gin.Context) {
	mes, anio, ok := ctrl.parseMesAnio(c)
	if !ok {
		return
	}

	resultado, err := ctrl.serviciosUseCase.Ejecutar(c.Request.Context(), mes, anio)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Fallo al obtener servicios populares", "detalle": err.Error()})
		return
	}

	c.JSON(http.StatusOK, resultado)
}

// ObtenerActividadUsuarios godoc
// @Summary Obtiene la actividad de usuarios del mes
// @Description Devuelve cuántos usuarios nuevos se registraron y cuántos usuarios activos (con al menos una reserva) hubo en el mes
// @Tags reportes
// @Produce json
// @Param mes query int true "Mes (1-12)"
// @Param anio query int true "Año (ej. 2026)"
// @Success 200 {object} domain.ActividadUsuariosDTO
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /reportes/actividad-usuarios [get]
// @Security BearerAuth
func (ctrl *ReportesController) ObtenerActividadUsuarios(c *gin.Context) {
	mes, anio, ok := ctrl.parseMesAnio(c)
	if !ok {
		return
	}

	resultado, err := ctrl.actividadUseCase.Ejecutar(c.Request.Context(), mes, anio)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Fallo al obtener actividad de usuarios", "detalle": err.Error()})
		return
	}

	c.JSON(http.StatusOK, resultado)
}

// parseMesAnio es un helper privado que parsea los query params mes y anio
func (ctrl *ReportesController) parseMesAnio(c *gin.Context) (int, int, bool) {
	mesStr := c.Query("mes")
	anioStr := c.Query("anio")

	mes, err1 := strconv.Atoi(mesStr)
	anio, err2 := strconv.Atoi(anioStr)

	if err1 != nil || err2 != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Debe enviar mes y anio como parámetros válidos"})
		return 0, 0, false
	}
	return mes, anio, true
}

// DescargarReporteAdminPDF godoc
// @Summary Descarga el reporte general en PDF
// @Description Genera y devuelve un documento PDF con los datos de administrador, actividad y servicios
// @Tags reportes
// @Produce application/pdf
// @Param mes query int true "Mes (1-12)"
// @Param anio query int true "Año (ej. 2026)"
// @Success 200 {file} file "reporte.pdf"
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /reportes/admin/pdf [get]
// @Security BearerAuth
func (ctrl *ReportesController) DescargarReporteAdminPDF(c *gin.Context) {
	mes, anio, ok := ctrl.parseMesAnio(c)
	if !ok {
		return
	}

	ctx := c.Request.Context()

	// Obtener datos consolidados del admin
	adminData, err := ctrl.adminUseCase.Ejecutar(ctx, mes, anio)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Fallo al generar datos del admin", "detalle": err.Error()})
		return
	}

	// Obtener servicios populares
	serviciosData, err := ctrl.serviciosUseCase.Ejecutar(ctx, mes, anio)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Fallo al obtener servicios populares", "detalle": err.Error()})
		return
	}

	// Obtener actividad de usuarios
	actividadData, err := ctrl.actividadUseCase.Ejecutar(ctx, mes, anio)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Fallo al obtener actividad de usuarios", "detalle": err.Error()})
		return
	}

	// Generar el PDF
	pdfBytes, err := pdf.GenerarReporteCompletoPDF(*adminData, serviciosData, *actividadData)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Fallo al generar el archivo PDF", "detalle": err.Error()})
		return
	}

	// Configurar headers para descarga
	c.Header("Content-Disposition", fmt.Sprintf(`attachment; filename="reporte_general_%d_%d.pdf"`, mes, anio))
	c.Header("Content-Type", "application/pdf")
	c.Header("Content-Length", strconv.Itoa(len(pdfBytes)))

	// Enviar los bytes del PDF
	c.Data(http.StatusOK, "application/pdf", pdfBytes)
}
