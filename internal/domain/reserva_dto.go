package domain

import "context"

// DTO para la petición de cancelación
type CancelarReservaRequestDTO struct {
	IDCliente uint `json:"id_cliente" binding:"required"`
}

// DTO para la respuesta de cancelación
type CancelarReservaResponseDTO struct {
	IDReserva     uint   `json:"id_reserva"`
	EstadoReserva string `json:"estado_reserva"`
	Mensaje       string `json:"mensaje"`
}

// DTO para la petición de finalización con foto de evidencia
type FinalizarReservaRequestDTO struct {
	IDPrestador uint   `json:"id_prestador" binding:"required"`
	FotoURL     string `json:"foto_url" binding:"required"`
	Detalle     string `json:"detalle,omitempty"`
}

// DTO para la respuesta de finalización
type FinalizarReservaResponseDTO struct {
	IDReserva     uint   `json:"id_reserva"`
	EstadoReserva string `json:"estado_reserva"`
	FotoURL       string `json:"foto_url"`
	Mensaje       string `json:"mensaje"`
}

// Interface del repositorio operativo de reservas
type IReservaOperativaRepository interface {
	CancelarReserva(ctx context.Context, idReserva uint, idCliente uint) error
	FinalizarReserva(ctx context.Context, idReserva uint, idPrestador uint, fotoURL string) error
}

