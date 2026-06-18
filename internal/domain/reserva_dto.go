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

// Interface del repositorio operativo de reservas
type IReservaOperativaRepository interface {
	CancelarReserva(ctx context.Context, idReserva uint, idCliente uint) error
}
