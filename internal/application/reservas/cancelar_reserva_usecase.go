package reservas

import (
	"apiServiYa/internal/domain"
	"context"
)

type CancelarReservaUseCase struct {
	repo domain.IReservaOperativaRepository
}

func NewCancelarReservaUseCase(repo domain.IReservaOperativaRepository) *CancelarReservaUseCase {
	return &CancelarReservaUseCase{repo: repo}
}

func (uc *CancelarReservaUseCase) Ejecutar(ctx context.Context, idReserva uint, idCliente uint) (*domain.CancelarReservaResponseDTO, error) {
	err := uc.repo.CancelarReserva(ctx, idReserva, idCliente)
	if err != nil {
		return nil, err
	}

	return &domain.CancelarReservaResponseDTO{
		IDReserva:     idReserva,
		EstadoReserva: "cancelada",
		Mensaje:       "Reserva cancelada exitosamente",
	}, nil
}
