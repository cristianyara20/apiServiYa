package repository

import (
	"apiServiYa/internal/domain"
	"context"
	"gorm.io/gorm"
)

type CalificacionRepository struct {
	db *gorm.DB
}

func NewCalificacionRepository(db *gorm.DB) domain.ICalificacionRepository {
	return &CalificacionRepository{db: db}
}

func (r *CalificacionRepository) ObtenerTodas(ctx context.Context) ([]domain.CalificacionDTO, error) {
	var resultados []domain.CalificacionDTO
	
	query := `
		SELECT 
			c.id_calificacion,
			c.id_reserva,
			c.puntuacion,
			c.comentario,
			c.fecha_calificacion,
			COALESCE(u_cli.nombre, '') as nombre_cliente,
			COALESCE(u_pres.nombre, '') as nombre_prestador,
			COALESCE(s.nombre_servicio, '') as nombre_servicio
		FROM gestion.calificaciones c
		JOIN gestion.reservas res ON c.id_reserva = res.id_reserva
		LEFT JOIN seguridad.usuarios u_cli ON res.id_cliente = u_cli.id_usuario
		LEFT JOIN seguridad.usuarios u_pres ON res.id_prestador = u_pres.id_usuario
		LEFT JOIN gestion.servicios s ON res.id_servicio = s.id_servicio
		ORDER BY c.fecha_calificacion DESC
	`
	
	err := r.db.WithContext(ctx).Raw(query).Scan(&resultados).Error
	return resultados, err
}
