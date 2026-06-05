package reportes

import (
	"apiServiYa/internal/domain"
	"context"
	"sync"
)

type GenerarReporteConsolidadoUseCase struct {
	uow domain.IReporteUnitOfWork
}

func NewGenerarReporteConsolidadoUseCase(uow domain.IReporteUnitOfWork) *GenerarReporteConsolidadoUseCase {
	return &GenerarReporteConsolidadoUseCase{
		uow: uow,
	}
}

func (uc *GenerarReporteConsolidadoUseCase) Ejecutar(ctx context.Context, mes int, anio int) (*domain.ReporteConsolidadoMensualDTO, error) {
	// DTO Final
	reporte := &domain.ReporteConsolidadoMensualDTO{
		Mes:  mes,
		Anio: anio,
	}

	// Usamos WaitGroup para ejecutar las consultas concurrentemente
	var wg sync.WaitGroup
	var errPendientes, errAceptadas, errCompletadas, errCanceladas, errPQRS, errTop error

	// 1. Contar pendientes
	wg.Add(1)
	go func() {
		defer wg.Done()
		reporte.TotalPendientes, errPendientes = uc.uow.ReservasAnalitica().ContarPorEstadoYMes(ctx, "pendiente", mes, anio)
	}()

	// 2. Contar aceptadas
	wg.Add(1)
	go func() {
		defer wg.Done()
		reporte.TotalAceptadas, errAceptadas = uc.uow.ReservasAnalitica().ContarPorEstadoYMes(ctx, "aceptada", mes, anio)
	}()

	// 3. Contar terminadas (completadas)
	wg.Add(1)
	go func() {
		defer wg.Done()
		reporte.TotalCompletadas, errCompletadas = uc.uow.ReservasAnalitica().ContarPorEstadoYMes(ctx, "terminada", mes, anio)
	}()

	// 4. Contar rechazadas (canceladas)
	wg.Add(1)
	go func() {
		defer wg.Done()
		reporte.TotalCanceladas, errCanceladas = uc.uow.ReservasAnalitica().ContarPorEstadoYMes(ctx, "rechazada", mes, anio)
	}()

	// 5. Contar PQRS Abiertas
	wg.Add(1)
	go func() {
		defer wg.Done()
		reporte.PQRSAbiertas, errPQRS = uc.uow.PqrsAnalitica().ContarPorEstado(ctx, "abierta")
	}()

	// 6. Obtener Top Prestadores
	wg.Add(1)
	go func() {
		defer wg.Done()
		reporte.TopPrestadores, errTop = uc.uow.PrestadoresAnalitica().ObtenerTopPrestadores(ctx, 3)
	}()

	// Esperamos a que todas las rutinas (goroutines) terminen
	wg.Wait()

	// Calcular total de reservas del mes (suma de todos los estados)
	reporte.TotalReservas = reporte.TotalPendientes + reporte.TotalAceptadas + reporte.TotalCompletadas + reporte.TotalCanceladas

	// Validación de errores
	for _, err := range []error{errPendientes, errAceptadas, errCompletadas, errCanceladas, errPQRS, errTop} {
		if err != nil {
			return nil, err
		}
	}

	return reporte, nil
}
