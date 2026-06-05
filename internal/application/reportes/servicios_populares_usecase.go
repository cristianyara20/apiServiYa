package reportes

import (
	"apiServiYa/internal/domain"
	"context"
)

type ObtenerServiciosPopularesUseCase struct {
	uow domain.IReporteUnitOfWork
}

func NewObtenerServiciosPopularesUseCase(uow domain.IReporteUnitOfWork) *ObtenerServiciosPopularesUseCase {
	return &ObtenerServiciosPopularesUseCase{uow: uow}
}

func (uc *ObtenerServiciosPopularesUseCase) Ejecutar(ctx context.Context, mes int, anio int) ([]domain.ServicioPopularDTO, error) {
	// Devuelve los 10 servicios más solicitados del mes
	return uc.uow.ServiciosAnalitica().ObtenerServiciosMasSolicitados(ctx, mes, anio, 10)
}
