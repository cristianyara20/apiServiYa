package reportes

import (
	"apiServiYa/internal/domain"
	"context"
)

type ObtenerPqrsUseCase struct {
	repo domain.IPqrRepository
}

func NewObtenerPqrsUseCase(repo domain.IPqrRepository) *ObtenerPqrsUseCase {
	return &ObtenerPqrsUseCase{repo: repo}
}

func (uc *ObtenerPqrsUseCase) Ejecutar(ctx context.Context) ([]domain.PqrDTO, error) {
	return uc.repo.ObtenerTodasPqrs(ctx)
}
