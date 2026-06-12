package pdf

import (
	"bytes"
	"fmt"

	"apiServiYa/internal/domain"

	"github.com/go-pdf/fpdf"
)

// GenerarReporteCompletoPDF crea un documento PDF consolidado con todas las métricas.
func GenerarReporteCompletoPDF(
	adminData domain.ReporteConsolidadoMensualDTO,
	serviciosData []domain.ServicioPopularDTO,
	actividadData domain.ActividadUsuariosDTO,
) ([]byte, error) {
	pdf := fpdf.New("P", "mm", "A4", "")
	pdf.AddPage()
	
	// Título Principal
	pdf.SetFont("Arial", "B", 16)
	pdf.CellFormat(190, 10, fmt.Sprintf("Reporte General de Plataforma - %d/%d", adminData.Mes, adminData.Anio), "0", 1, "C", false, 0, "")
	pdf.Ln(10)

	// --- SECCIÓN 1: ESTADO DE RESERVAS ---
	pdf.SetFont("Arial", "B", 14)
	pdf.CellFormat(190, 10, "1. Estado de las Reservas", "0", 1, "L", false, 0, "")
	pdf.SetFont("Arial", "", 12)
	pdf.CellFormat(95, 8, fmt.Sprintf("Total Reservas: %d", adminData.TotalReservas), "0", 0, "L", false, 0, "")
	pdf.CellFormat(95, 8, fmt.Sprintf("Pendientes: %d", adminData.TotalPendientes), "0", 1, "L", false, 0, "")
	pdf.CellFormat(95, 8, fmt.Sprintf("Completadas: %d", adminData.TotalCompletadas), "0", 0, "L", false, 0, "")
	pdf.CellFormat(95, 8, fmt.Sprintf("Aceptadas: %d", adminData.TotalAceptadas), "0", 1, "L", false, 0, "")
	pdf.CellFormat(95, 8, fmt.Sprintf("Canceladas: %d", adminData.TotalCanceladas), "0", 1, "L", false, 0, "")
	pdf.Ln(5)

	// --- SECCIÓN 2: ACTIVIDAD Y PQRS ---
	pdf.SetFont("Arial", "B", 14)
	pdf.CellFormat(190, 10, "2. Actividad de Usuarios y Soporte", "0", 1, "L", false, 0, "")
	pdf.SetFont("Arial", "", 12)
	pdf.CellFormat(95, 8, fmt.Sprintf("Usuarios Nuevos: %d", actividadData.UsuariosNuevos), "0", 0, "L", false, 0, "")
	pdf.CellFormat(95, 8, fmt.Sprintf("Usuarios Activos: %d", actividadData.UsuariosActivos), "0", 1, "L", false, 0, "")
	pdf.CellFormat(95, 8, fmt.Sprintf("PQRS Abiertas: %d", adminData.PQRSAbiertas), "0", 1, "L", false, 0, "")
	pdf.Ln(5)

	// --- SECCIÓN 3: TOP PRESTADORES ---
	pdf.SetFont("Arial", "B", 14)
	pdf.CellFormat(190, 10, "3. Top Prestadores", "0", 1, "L", false, 0, "")
	
	// Cabecera Tabla Prestadores
	pdf.SetFont("Arial", "B", 12)
	pdf.SetFillColor(200, 200, 200)
	pdf.CellFormat(60, 8, "Nombre", "1", 0, "C", true, 0, "")
	pdf.CellFormat(70, 8, "Correo", "1", 0, "C", true, 0, "")
	pdf.CellFormat(30, 8, "Calificacion", "1", 0, "C", true, 0, "")
	pdf.CellFormat(30, 8, "Servicios", "1", 1, "C", true, 0, "")

	pdf.SetFont("Arial", "", 10)
	for _, p := range adminData.TopPrestadores {
		pdf.CellFormat(60, 8, p.NombrePrestador, "1", 0, "L", false, 0, "")
		pdf.CellFormat(70, 8, p.CorreoPrestador, "1", 0, "L", false, 0, "")
		pdf.CellFormat(30, 8, fmt.Sprintf("%.1f", p.Calificacion), "1", 0, "C", false, 0, "")
		pdf.CellFormat(30, 8, fmt.Sprintf("%d", p.TotalServicios), "1", 1, "C", false, 0, "")
	}
	pdf.Ln(10)

	// --- SECCIÓN 4: SERVICIOS MÁS SOLICITADOS ---
	pdf.SetFont("Arial", "B", 14)
	pdf.CellFormat(190, 10, "4. Servicios Mas Solicitados", "0", 1, "L", false, 0, "")
	
	// Cabecera Tabla Servicios
	pdf.SetFont("Arial", "B", 12)
	pdf.SetFillColor(200, 200, 200)
	pdf.CellFormat(130, 8, "Servicio", "1", 0, "C", true, 0, "")
	pdf.CellFormat(60, 8, "Solicitudes", "1", 1, "C", true, 0, "")

	pdf.SetFont("Arial", "", 10)
	for _, s := range serviciosData {
		pdf.CellFormat(130, 8, s.NombreServicio, "1", 0, "L", false, 0, "")
		pdf.CellFormat(60, 8, fmt.Sprintf("%d", s.VecesSolicitado), "1", 1, "C", false, 0, "")
	}

	// Buffer para escribir el PDF
	var buf bytes.Buffer
	if err := pdf.Output(&buf); err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}
