package repository

import (
	"apiServiYa/internal/domain"
	"context"
	"errors"
	"gorm.io/gorm"
)

type ReservaOperativaRepository struct {
	db *gorm.DB
}

func NewReservaOperativaRepository(db *gorm.DB) domain.IReservaOperativaRepository {
	return &ReservaOperativaRepository{db: db}
}

func (r *ReservaOperativaRepository) CancelarReserva(ctx context.Context, idReserva uint, idCliente uint) error {
	// Verificar que la reserva exista, pertenezca al cliente y esté en estado cancelable
	var estado string
	err := r.db.WithContext(ctx).Table("gestion.reservas").
		Select("estado_reserva").
		Where("id_reserva = ? AND id_cliente = ?", idReserva, idCliente).
		Scan(&estado).Error

	if err != nil {
		if err == gorm.ErrRecordNotFound || estado == "" {
			return errors.New("reserva no encontrada o no pertenece al cliente")
		}
		return err
	}
	if estado == "" {
		return errors.New("reserva no encontrada o no pertenece al cliente")
	}

	if estado != "pendiente" && estado != "aceptada" {
		return errors.New("solo se pueden cancelar reservas en estado pendiente o aceptada")
	}

	// Actualizar estado
	result := r.db.WithContext(ctx).Table("gestion.reservas").
		Where("id_reserva = ? AND id_cliente = ?", idReserva, idCliente).
		Update("estado_reserva", "cancelada")

	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return errors.New("no se pudo cancelar la reserva")
	}

	return nil
}
