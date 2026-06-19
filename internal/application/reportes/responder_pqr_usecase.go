package reportes

import (
	"apiServiYa/internal/domain"
	"context"
	"errors"
)

type ResponderPqrUseCase struct {
	repo domain.IPqrRepository
}

func NewResponderPqrUseCase(repo domain.IPqrRepository) *ResponderPqrUseCase {
	return &ResponderPqrUseCase{repo: repo}
}

func (uc *ResponderPqrUseCase) Ejecutar(ctx context.Context, idPqr uint, respuesta string) error {
	if respuesta == "" {
		return errors.New("la respuesta no puede estar vacía")
	}
	return uc.repo.ResponderPqr(ctx, idPqr, respuesta)
}
