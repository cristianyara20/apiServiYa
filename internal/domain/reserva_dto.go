package domain

import (
	"context"
	"time"
)

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

// DTO para la petición de creación de reserva
type CrearReservaRequestDTO struct {
	IDCliente   uint   `json:"id_cliente" binding:"required"`
	IDServicio  uint   `json:"id_servicio" binding:"required"`
	IDPrestador *uint  `json:"id_prestador,omitempty"`
	Direccion   string `json:"direccion" binding:"required"`
	Descripcion string `json:"descripcion" binding:"required"`
	FechaAgenda string `json:"fecha_agenda" binding:"required"`
}

// Entidad / Modelo de base de datos para la tabla gestion.reservas
type Reserva struct {
	IDReserva     uint      `gorm:"primaryKey;column:id_reserva;autoIncrement" json:"id_reserva"`
	IDCliente     uint      `gorm:"column:id_cliente" json:"id_cliente"`
	IDPrestador   *uint     `gorm:"column:id_prestador" json:"id_prestador"`
	IDServicio    uint      `gorm:"column:id_servicio" json:"id_servicio"`
	FechaAgenda   time.Time `gorm:"column:fecha_agenda" json:"fecha_agenda"`
	Direccion     string    `gorm:"column:direccion" json:"direccion"`
	Descripcion   string    `gorm:"column:descripcion" json:"descripcion"`
	EstadoReserva string    `gorm:"column:estado_reserva;default:pendiente" json:"estado_reserva"`
	DetalleExtra  *string   `gorm:"column:detalle_extra" json:"detalle_extra,omitempty"`
}

func (Reserva) TableName() string {
	return "gestion.reservas"
}

// Interface del repositorio operativo de reservas
type IReservaOperativaRepository interface {
	CrearReserva(ctx context.Context, reserva *Reserva) (*Reserva, error)
	CancelarReserva(ctx context.Context, idReserva uint, idCliente uint) error
	FinalizarReserva(ctx context.Context, idReserva uint, idPrestador uint, fotoURL string) error
}

