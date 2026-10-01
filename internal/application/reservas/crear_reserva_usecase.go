package reservas

import (
	"apiServiYa/internal/domain"
	"context"
	"errors"
	"strings"
	"time"
)

type CrearReservaUseCase struct {
	repo domain.IReservaOperativaRepository
}

func NewCrearReservaUseCase(repo domain.IReservaOperativaRepository) *CrearReservaUseCase {
	return &CrearReservaUseCase{repo: repo}
}

func (uc *CrearReservaUseCase) Ejecutar(ctx context.Context, req domain.CrearReservaRequestDTO) (*domain.Reserva, error) {
	if req.IDCliente == 0 {
		return nil, errors.New("el ID del cliente es obligatorio")
	}
	if req.IDServicio == 0 {
		return nil, errors.New("el ID del servicio es obligatorio")
	}
	if strings.TrimSpace(req.Direccion) == "" {
		return nil, errors.New("la dirección es obligatoria")
	}
	if strings.TrimSpace(req.Descripcion) == "" {
		return nil, errors.New("la descripción es obligatoria")
	}

	trimmedFecha := strings.TrimSpace(req.FechaAgenda)
	if trimmedFecha == "" {
		return nil, errors.New("la fecha de la cita es obligatoria")
	}

	// Parsear fecha con múltiples formatos soportados
	var parsedDate time.Time
	formats := []string{
		time.RFC3339Nano,
		time.RFC3339,
		"2006-01-02T15:04:05Z07:00",
		"2006-01-02T15:04:05",
		"2006-01-02 15:04:05",
		"2006-01-02",
		"02-01-2006 15:04:05",
		"02-01-2006",
		"02/01/2006 15:04:05",
		"02/01/2006",
	}
	for _, f := range formats {
		if t, err := time.Parse(f, trimmedFecha); err == nil {
			parsedDate = t
			break
		}
	}
	if parsedDate.IsZero() {
		return nil, errors.New("formato de fecha_agenda inválido. Debe ser ISO8601 o RFC3339")
	}

	nuevaReserva := &domain.Reserva{
		IDCliente:     req.IDCliente,
		IDServicio:    req.IDServicio,
		IDPrestador:   req.IDPrestador,
		Direccion:     strings.TrimSpace(req.Direccion),
		Descripcion:   strings.TrimSpace(req.Descripcion),
		FechaAgenda:   parsedDate,
		EstadoReserva: "pendiente",
	}

	return uc.repo.CrearReserva(ctx, nuevaReserva)
}
