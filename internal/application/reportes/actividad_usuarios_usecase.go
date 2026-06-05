package reportes

import (
	"apiServiYa/internal/domain"
	"context"
	"sync"
)

type ObtenerActividadUsuariosUseCase struct {
	uow domain.IReporteUnitOfWork
}

func NewObtenerActividadUsuariosUseCase(uow domain.IReporteUnitOfWork) *ObtenerActividadUsuariosUseCase {
	return &ObtenerActividadUsuariosUseCase{uow: uow}
}

func (uc *ObtenerActividadUsuariosUseCase) Ejecutar(ctx context.Context, mes int, anio int) (*domain.ActividadUsuariosDTO, error) {
	dto := &domain.ActividadUsuariosDTO{
		Mes:  mes,
		Anio: anio,
	}

	// Ejecutamos las 2 consultas en paralelo con Goroutines
	var wg sync.WaitGroup
	var errNuevos, errActivos error

	wg.Add(1)
	go func() {
		defer wg.Done()
		dto.UsuariosNuevos, errNuevos = uc.uow.UsuariosAnalitica().ContarUsuariosNuevos(ctx, mes, anio)
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		dto.UsuariosActivos, errActivos = uc.uow.UsuariosAnalitica().ContarUsuariosActivos(ctx, mes, anio)
	}()

	wg.Wait()

	if errNuevos != nil {
		return nil, errNuevos
	}
	if errActivos != nil {
		return nil, errActivos
	}

	return dto, nil
}
