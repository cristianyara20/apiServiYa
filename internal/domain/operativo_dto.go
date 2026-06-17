package domain

import "context"

type PrestadorDetalladoDTO struct {
	IDPrestador          uint    `json:"id_prestador"`
	Nombre               string  `json:"nombre"`
	Apellido             string  `json:"apellido"`
	Correo               string  `json:"correo"`
	Experiencia          string  `json:"experiencia"`
	CalificacionPromedio float64 `json:"calificacion_promedio"`
	EstadoDisponibilidad string  `json:"estado_disponibilidad"`
}

type ReservaHistorialDTO struct {
	IDReserva      uint   `json:"id_reserva"`
	IDCliente      uint   `json:"id_cliente"`
	NombreCliente  string `json:"nombre_cliente"`
	IDPrestador    *uint  `json:"id_prestador"`
	NombrePrestador *string `json:"nombre_prestador"`
	IDServicio     uint   `json:"id_servicio"`
	NombreServicio string `json:"nombre_servicio"`
	FechaAgenda    string `json:"fecha_agenda"`
	Direccion      string `json:"direccion"`
	Descripcion    string `json:"descripcion"`
	EstadoReserva  string `json:"estado_reserva"`
}

type NotificationRequestDTO struct {
	IDReserva     uint   `json:"id_reserva"`
	IDCliente     uint   `json:"id_cliente"`
	IDPrestador   uint   `json:"id_prestador"`
	EstadoReserva string `json:"estado_reserva"`
}

type IPrestadorOperativoRepository interface {
	ObtenerTodosPrestadores(ctx context.Context) ([]PrestadorDetalladoDTO, error)
	ObtenerHistorialServicios(ctx context.Context) ([]ReservaHistorialDTO, error)
}
