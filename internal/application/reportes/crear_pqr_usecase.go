package reportes

import (
	"apiServiYa/internal/domain"
	"context"
	"errors"
	"strings"
)

type CrearPqrUseCase struct {
	repo domain.IPqrRepository
}

func NewCrearPqrUseCase(repo domain.IPqrRepository) *CrearPqrUseCase {
	return &CrearPqrUseCase{repo: repo}
}

func (uc *CrearPqrUseCase) Ejecutar(ctx context.Context, req domain.CrearPqrRequestDTO) (*domain.CrearPqrResponseDTO, error) {
	if req.IDCliente == 0 {
		return nil, errors.New("el ID del cliente es obligatorio")
	}
	if req.IDReserva == 0 {
		return nil, errors.New("el ID de la reserva es obligatorio")
	}
	if strings.TrimSpace(req.TipoPqr) == "" {
		return nil, errors.New("el tipo de PQR es obligatorio")
	}
	if strings.TrimSpace(req.Descripcion) == "" {
		return nil, errors.New("la descripción es obligatoria")
	}

	req.TipoPqr = strings.TrimSpace(req.TipoPqr)
	req.Descripcion = strings.TrimSpace(req.Descripcion)

	return uc.repo.CrearPqr(ctx, &req)
}
