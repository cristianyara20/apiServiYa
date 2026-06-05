package domain

import "context"

// DTO para la respuesta HTTP (JSON)
type ReporteConsolidadoMensualDTO struct {
	Mes                   int                   `json:"mes"`
	Anio                  int                   `json:"anio"`
	TotalReservas         int64                 `json:"total_reservas"`
	TotalPendientes       int64                 `json:"total_pendientes"`
	TotalAceptadas        int64                 `json:"total_aceptadas"`
	TotalCompletadas      int64                 `json:"total_completadas"`
	TotalCanceladas       int64                 `json:"total_canceladas"`
	PQRSAbiertas          int64                 `json:"pqrs_abiertas"`
	TopPrestadores        []PrestadorTopDTO     `json:"top_prestadores"`
}

type PrestadorTopDTO struct {
	IDPrestador      uint    `json:"id_prestador"`
	NombrePrestador  string  `json:"nombre_prestador"`
	CorreoPrestador  string  `json:"correo_prestador"`
	Calificacion     float64 `json:"calificacion"`
	TotalServicios   int64   `json:"total_servicios"`
}

// DTO para reporte #2: Servicios más solicitados
type ServicioPopularDTO struct {
	IDServicio      uint   `json:"id_servicio"`
	NombreServicio  string `json:"nombre_servicio"`
	Categoria       string `json:"categoria"`
	VecesSolicitado int64  `json:"veces_solicitado"`
}

// DTO para reporte #6: Actividad de usuarios
type ActividadUsuariosDTO struct {
	Mes              int   `json:"mes"`
	Anio             int   `json:"anio"`
	UsuariosNuevos   int64 `json:"usuarios_nuevos"`
	UsuariosActivos  int64 `json:"usuarios_activos"`
}

// Interfaces de Repositorios (Analítica)
type IReservaAnaliticaRepository interface {
	ContarPorEstadoYMes(ctx context.Context, estado string, mes int, anio int) (int64, error)
}

type IPrestadorAnaliticaRepository interface {
	ObtenerTopPrestadores(ctx context.Context, limite int) ([]PrestadorTopDTO, error)
}

type IPqrsAnaliticaRepository interface {
	ContarPorEstado(ctx context.Context, estado string) (int64, error)
}

type IServicioAnaliticaRepository interface {
	ObtenerServiciosMasSolicitados(ctx context.Context, mes int, anio int, limite int) ([]ServicioPopularDTO, error)
}

type IUsuarioAnaliticaRepository interface {
	ContarUsuariosNuevos(ctx context.Context, mes int, anio int) (int64, error)
	ContarUsuariosActivos(ctx context.Context, mes int, anio int) (int64, error)
}

// IReporteUnitOfWork agrupa los repositorios de reportes
type IReporteUnitOfWork interface {
	ReservasAnalitica() IReservaAnaliticaRepository
	PrestadoresAnalitica() IPrestadorAnaliticaRepository
	PqrsAnalitica() IPqrsAnaliticaRepository
	ServiciosAnalitica() IServicioAnaliticaRepository
	UsuariosAnalitica() IUsuarioAnaliticaRepository
}
