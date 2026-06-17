package repository

import (
	"apiServiYa/internal/domain"
	"context"

	"gorm.io/gorm"
)

type PrestadorOperativoRepository struct {
	db *gorm.DB
}

func NewPrestadorOperativoRepository(db *gorm.DB) *PrestadorOperativoRepository {
	return &PrestadorOperativoRepository{db: db}
}

func (r *PrestadorOperativoRepository) ObtenerTodosPrestadores(ctx context.Context) ([]domain.PrestadorDetalladoDTO, error) {
	var resultados []domain.PrestadorDetalladoDTO
	query := `
		SELECT 
			p.id_prestador,
			u.nombre,
			u.apellido,
			u.correo,
			COALESCE(p.experiencia, '') as experiencia,
			COALESCE(p.calificacion_promedio, 5.0) as calificacion_promedio,
			COALESCE(p.estado_disponibilidad, 'disponible') as estado_disponibilidad
		FROM gestion.prestadores p
		JOIN seguridad.usuarios u ON u.id_usuario = p.id_prestador
		ORDER BY p.id_prestador ASC
	`
	err := r.db.WithContext(ctx).Raw(query).Scan(&resultados).Error
	return resultados, err
}

func (r *PrestadorOperativoRepository) ObtenerHistorialServicios(ctx context.Context) ([]domain.ReservaHistorialDTO, error) {
	var resultados []domain.ReservaHistorialDTO
	query := `
		SELECT 
			r.id_reserva,
			r.id_cliente,
			CONCAT(uc.nombre, ' ', uc.apellido) as nombre_cliente,
			r.id_prestador,
			CONCAT(up.nombre, ' ', up.apellido) as nombre_prestador,
			r.id_servicio,
			s.nombre_servicio,
			TO_CHAR(r.fecha_agenda, 'YYYY-MM-DD HH24:MI:SS') as fecha_agenda,
			COALESCE(r.direccion, '') as direccion,
			COALESCE(r.descripcion, '') as descripcion,
			r.estado_reserva
		FROM gestion.reservas r
		JOIN seguridad.usuarios uc ON uc.id_usuario = r.id_cliente
		LEFT JOIN seguridad.usuarios up ON up.id_usuario = r.id_prestador
		JOIN gestion.servicios s ON s.id_servicio = r.id_servicio
		ORDER BY r.fecha_agenda DESC
	`
	err := r.db.WithContext(ctx).Raw(query).Scan(&resultados).Error
	return resultados, err
}
