package domain

import "context"

type CalificacionDTO struct {
	IDCalificacion    uint   `json:"id_calificacion"`
	IDReserva         uint   `json:"id_reserva"`
	Puntuacion        int    `json:"puntuacion"`
	Comentario        string `json:"comentario"`
	FechaCalificacion string `json:"fecha_calificacion"`
	// Datos extra para hacerlo más rico
	NombreCliente    string `json:"nombre_cliente,omitempty"`
	NombrePrestador  string `json:"nombre_prestador,omitempty"`
	NombreServicio   string `json:"nombre_servicio,omitempty"`
}

type ICalificacionRepository interface {
	ObtenerTodas(ctx context.Context) ([]CalificacionDTO, error)
}
