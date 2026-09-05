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

func (r *ReservaOperativaRepository) FinalizarReserva(ctx context.Context, idReserva uint, idPrestador uint, fotoURL string) error {
	// Verificar que la reserva exista, esté asignada al prestador y esté en estado activo
	var reserva struct {
		EstadoReserva string `gorm:"column:estado_reserva"`
		IDPrestador   *uint  `gorm:"column:id_prestador"`
	}
	err := r.db.WithContext(ctx).Table("gestion.reservas").
		Select("estado_reserva, id_prestador").
		Where("id_reserva = ?", idReserva).
		Scan(&reserva).Error

	if err != nil {
		if err == gorm.ErrRecordNotFound || reserva.EstadoReserva == "" {
			return errors.New("reserva no encontrada")
		}
		return err
	}
	if reserva.EstadoReserva == "" {
		return errors.New("reserva no encontrada")
	}

	if reserva.IDPrestador == nil || *reserva.IDPrestador != idPrestador {
		return errors.New("la reserva no está asignada a este prestador")
	}

	st := reserva.EstadoReserva
	if st == "cancelada" || st == "rechazada" || st == "terminada" || st == "completada" {
		return errors.New("la reserva ya se encuentra " + st)
	}

	// Actualizar estado a terminada y guardar la URL de la foto en detalle_extra
	updates := map[string]interface{}{
		"estado_reserva": "terminada",
		"detalle_extra":  fotoURL,
	}

	result := r.db.WithContext(ctx).Table("gestion.reservas").
		Where("id_reserva = ?", idReserva).
		Updates(updates)

	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return errors.New("no se pudo finalizar la reserva")
	}

	return nil
}

