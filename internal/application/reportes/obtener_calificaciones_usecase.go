package reportes

import (
	"apiServiYa/internal/domain"
	"context"
)

type ObtenerCalificacionesUseCase struct {
	repo domain.ICalificacionRepository
}

func NewObtenerCalificacionesUseCase(repo domain.ICalificacionRepository) *ObtenerCalificacionesUseCase {
	return &ObtenerCalificacionesUseCase{repo: repo}
}

func (uc *ObtenerCalificacionesUseCase) Ejecutar(ctx context.Context) ([]domain.CalificacionDTO, error) {
	return uc.repo.ObtenerTodas(ctx)
}
