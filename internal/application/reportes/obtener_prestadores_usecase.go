package reportes

import (
	"apiServiYa/internal/domain"
	"context"
)

type ObtenerPrestadoresOperativosUseCase struct {
	repo domain.IPrestadorOperativoRepository
}

func NewObtenerPrestadoresOperativosUseCase(repo domain.IPrestadorOperativoRepository) *ObtenerPrestadoresOperativosUseCase {
	return &ObtenerPrestadoresOperativosUseCase{repo: repo}
}

func (uc *ObtenerPrestadoresOperativosUseCase) Ejecutar(ctx context.Context) ([]domain.PrestadorDetalladoDTO, error) {
	return uc.repo.ObtenerTodosPrestadores(ctx)
}
