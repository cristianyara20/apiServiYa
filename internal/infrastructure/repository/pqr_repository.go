package repository

import (
	"apiServiYa/internal/domain"
	"context"
	"time"

	"gorm.io/gorm"
)

type PqrRepository struct {
	db *gorm.DB
}

func NewPqrRepository(db *gorm.DB) *PqrRepository {
	return &PqrRepository{db: db}
}

func (r *PqrRepository) ObtenerTodasPqrs(ctx context.Context) ([]domain.PqrDTO, error) {
	var resultados []domain.PqrDTO
	query := `
		SELECT 
			p.id_pqr,
			p.id_cliente,
			CONCAT(u.nombre, ' ', u.apellido) as nombre_cliente,
			p.id_reserva,
			p.tipo_pqr,
			COALESCE(p.descripcion, '') as descripcion,
			COALESCE(p.estado_pqr, 'Abierto') as estado_pqr,
			p.respuesta_admin,
			CASE WHEN p.fecha_respuesta IS NOT NULL 
				THEN TO_CHAR(p.fecha_respuesta, 'YYYY-MM-DD HH24:MI:SS') 
				ELSE NULL 
			END as fecha_respuesta
		FROM soporte.pqrs p
		JOIN seguridad.usuarios u ON u.id_usuario = p.id_cliente
		ORDER BY p.id_pqr DESC
	`
	err := r.db.WithContext(ctx).Raw(query).Scan(&resultados).Error
	return resultados, err
}

func (r *PqrRepository) ResponderPqr(ctx context.Context, idPqr uint, respuesta string) error {
	query := `
		UPDATE soporte.pqrs 
		SET respuesta_admin = ?, 
			estado_pqr = 'Cerrado', 
			fecha_respuesta = ?
		WHERE id_pqr = ?
	`
	return r.db.WithContext(ctx).Exec(query, respuesta, time.Now(), idPqr).Error
}

func (r *PqrRepository) CrearPqr(ctx context.Context, req *domain.CrearPqrRequestDTO) (*domain.CrearPqrResponseDTO, error) {
	var idPqr uint
	query := `
		INSERT INTO soporte.pqrs (id_cliente, id_reserva, tipo_pqr, descripcion, estado_pqr)
		VALUES (?, ?, ?, ?, 'Abierto')
		RETURNING id_pqr
	`
	err := r.db.WithContext(ctx).Raw(query, req.IDCliente, req.IDReserva, req.TipoPqr, req.Descripcion).Scan(&idPqr).Error
	if err != nil {
		return nil, err
	}

	return &domain.CrearPqrResponseDTO{
		IDPqr:       idPqr,
		IDCliente:   req.IDCliente,
		IDReserva:   req.IDReserva,
		TipoPqr:     req.TipoPqr,
		Descripcion: req.Descripcion,
		EstadoPqr:   "Abierto",
		Mensaje:     "PQR radicada exitosamente",
	}, nil
}

