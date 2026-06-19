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

// --- PQR DTOs ---

type PqrDTO struct {
	IDPqr           uint    `json:"id_pqr"`
	IDCliente       uint    `json:"id_cliente"`
	NombreCliente   string  `json:"nombre_cliente"`
	IDReserva       *uint   `json:"id_reserva"`
	TipoPqr         string  `json:"tipo_pqr"`
	Descripcion     string  `json:"descripcion"`
	EstadoPqr       string  `json:"estado_pqr"`
	FechaPqr        string  `json:"fecha_pqr"`
	RespuestaAdmin  *string `json:"respuesta_admin"`
	FechaRespuesta  *string `json:"fecha_respuesta"`
}

type ResponderPqrRequestDTO struct {
	IDPqr          uint   `json:"id_pqr" binding:"required"`
	RespuestaAdmin string `json:"respuesta_admin" binding:"required"`
}

type IPqrRepository interface {
	ObtenerTodasPqrs(ctx context.Context) ([]PqrDTO, error)
	ResponderPqr(ctx context.Context, idPqr uint, respuesta string) error
}
