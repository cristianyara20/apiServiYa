package reportes

import (
	"apiServiYa/internal/domain"
	"context"
)

type ObtenerHistorialServiciosUseCase struct {
	repo domain.IPrestadorOperativoRepository
}

func NewObtenerHistorialServiciosUseCase(repo domain.IPrestadorOperativoRepository) *ObtenerHistorialServiciosUseCase {
	return &ObtenerHistorialServiciosUseCase{repo: repo}
}

func (uc *ObtenerHistorialServiciosUseCase) Ejecutar(ctx context.Context) ([]domain.ReservaHistorialDTO, error) {
	return uc.repo.ObtenerHistorialServicios(ctx)
}
