package repository

import (
	"apiServiYa/internal/domain"
	"context"

	"gorm.io/gorm"
)

// Implementación del UnitOfWork de Lectura para Reportes
type ReporteUnitOfWork struct {
	db *gorm.DB
}

func NewReporteUnitOfWork(db *gorm.DB) *ReporteUnitOfWork {
	return &ReporteUnitOfWork{db: db}
}

func (u *ReporteUnitOfWork) ReservasAnalitica() domain.IReservaAnaliticaRepository {
	return &ReservaAnaliticaRepository{db: u.db}
}

func (u *ReporteUnitOfWork) PrestadoresAnalitica() domain.IPrestadorAnaliticaRepository {
	return &PrestadorAnaliticaRepository{db: u.db}
}

func (u *ReporteUnitOfWork) PqrsAnalitica() domain.IPqrsAnaliticaRepository {
	return &PqrsAnaliticaRepository{db: u.db}
}

func (u *ReporteUnitOfWork) ServiciosAnalitica() domain.IServicioAnaliticaRepository {
	return &ServicioAnaliticaRepository{db: u.db}
}

func (u *ReporteUnitOfWork) UsuariosAnalitica() domain.IUsuarioAnaliticaRepository {
	return &UsuarioAnaliticaRepository{db: u.db}
}


// --- Implementaciones de los Repositorios de Reportes ---

type ReservaAnaliticaRepository struct {
	db *gorm.DB
}

func (r *ReservaAnaliticaRepository) ContarPorEstadoYMes(ctx context.Context, estado string, mes int, anio int) (int64, error) {
	var count int64
	// Filtra por estado_reserva (pendiente, aceptada, rechazada, terminada) y por mes/año
	err := r.db.WithContext(ctx).Table("gestion.reservas").
		Where("estado_reserva = ? AND EXTRACT(MONTH FROM fecha_agenda) = ? AND EXTRACT(YEAR FROM fecha_agenda) = ?", estado, mes, anio).
		Count(&count).Error
	return count, err
}


type PrestadorAnaliticaRepository struct {
	db *gorm.DB
}

func (r *PrestadorAnaliticaRepository) ObtenerTopPrestadores(ctx context.Context, limite int) ([]domain.PrestadorTopDTO, error) {
	var resultados []domain.PrestadorTopDTO
	query := `
		SELECT 
			p.id_prestador,
			CONCAT(u.nombre, ' ', u.apellido) as nombre_prestador,
			u.correo as correo_prestador,
			COALESCE(AVG(c.puntuacion), p.calificacion_promedio) as calificacion,
			COUNT(r.id_reserva) as total_servicios
		FROM gestion.prestadores p
		JOIN seguridad.usuarios u ON u.id_usuario = p.id_prestador
		JOIN gestion.reservas r ON r.id_prestador = p.id_prestador AND r.estado_reserva = 'terminada'
		LEFT JOIN gestion.calificaciones c ON c.id_reserva = r.id_reserva
		GROUP BY p.id_prestador, u.nombre, u.apellido, u.correo, p.calificacion_promedio
		ORDER BY total_servicios DESC, calificacion DESC
		LIMIT ?
	`
	err := r.db.WithContext(ctx).Raw(query, limite).Scan(&resultados).Error
	return resultados, err
}


type PqrsAnaliticaRepository struct {
	db *gorm.DB
}

func (r *PqrsAnaliticaRepository) ContarPorEstado(ctx context.Context, estado string) (int64, error) {
	var count int64
	// En tu esquema, la columna de PQRS se llama estado_pqr
	err := r.db.WithContext(ctx).Table("soporte.pqrs").Where("estado_pqr = ?", estado).Count(&count).Error
	return count, err
}


// --- NUEVO: Repositorio de Servicios Populares (#2) ---

type ServicioAnaliticaRepository struct {
	db *gorm.DB
}

func (r *ServicioAnaliticaRepository) ObtenerServiciosMasSolicitados(ctx context.Context, mes int, anio int, limite int) ([]domain.ServicioPopularDTO, error) {
	var resultados []domain.ServicioPopularDTO
	query := `
		SELECT s.id_servicio, s.nombre_servicio, s.categoria, COUNT(r.id_reserva) as veces_solicitado
		FROM gestion.servicios s
		JOIN gestion.reservas r ON r.id_servicio = s.id_servicio
		WHERE EXTRACT(MONTH FROM r.fecha_agenda) = ? AND EXTRACT(YEAR FROM r.fecha_agenda) = ?
		GROUP BY s.id_servicio, s.nombre_servicio, s.categoria
		ORDER BY veces_solicitado DESC
		LIMIT ?
	`
	err := r.db.WithContext(ctx).Raw(query, mes, anio, limite).Scan(&resultados).Error
	return resultados, err
}


// --- NUEVO: Repositorio de Actividad de Usuarios (#6) ---

type UsuarioAnaliticaRepository struct {
	db *gorm.DB
}

func (r *UsuarioAnaliticaRepository) ContarUsuariosNuevos(ctx context.Context, mes int, anio int) (int64, error) {
	var count int64
	// Cuenta cuántos clientes se registraron en el mes dado
	err := r.db.WithContext(ctx).Table("gestion.clientes").
		Where("EXTRACT(MONTH FROM fecha_registro) = ? AND EXTRACT(YEAR FROM fecha_registro) = ?", mes, anio).
		Count(&count).Error
	return count, err
}

func (r *UsuarioAnaliticaRepository) ContarUsuariosActivos(ctx context.Context, mes int, anio int) (int64, error) {
	var count int64
	// Cuenta cuántos clientes únicos tienen al menos una reserva en el mes dado
	query := `
		SELECT COUNT(DISTINCT r.id_cliente)
		FROM gestion.reservas r
		WHERE EXTRACT(MONTH FROM r.fecha_agenda) = ? AND EXTRACT(YEAR FROM r.fecha_agenda) = ?
	`
	err := r.db.WithContext(ctx).Raw(query, mes, anio).Scan(&count).Error
	return count, err
}

