package reservas

import (
	"apiServiYa/internal/domain"
	"context"
	"errors"
	"strings"
)

type FinalizarReservaUseCase struct {
	repo domain.IReservaOperativaRepository
}

func NewFinalizarReservaUseCase(repo domain.IReservaOperativaRepository) *FinalizarReservaUseCase {
	return &FinalizarReservaUseCase{repo: repo}
}

func (uc *FinalizarReservaUseCase) Ejecutar(ctx context.Context, idReserva uint, idPrestador uint, fotoURL string) (*domain.FinalizarReservaResponseDTO, error) {
	if idReserva == 0 {
		return nil, errors.New("ID de reserva inválido")
	}
	if idPrestador == 0 {
		return nil, errors.New("ID de prestador inválido")
	}
	trimmedFoto := strings.TrimSpace(fotoURL)
	if trimmedFoto == "" {
		return nil, errors.New("la foto de evidencia es obligatoria para finalizar el servicio")
	}

	err := uc.repo.FinalizarReserva(ctx, idReserva, idPrestador, trimmedFoto)
	if err != nil {
		return nil, err
	}

	return &domain.FinalizarReservaResponseDTO{
		IDReserva:     idReserva,
		EstadoReserva: "terminada",
		FotoURL:       trimmedFoto,
		Mensaje:       "Servicio finalizado con éxito con evidencia fotográfica",
	}, nil
}
