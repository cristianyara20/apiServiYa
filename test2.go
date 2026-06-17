package main

import (
	"fmt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"time"
)

func main() {
	dsn := "postgresql://postgres.lvxhporsajorgckeisna:Cristianvargas2007%23@aws-1-us-east-1.pooler.supabase.com:5432/postgres?sslmode=require"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		panic(err)
	}

	type Reserva struct {
		IDReserva   int
		FechaAgenda *time.Time
	}

	var reservas []Reserva
	db.Table("gestion.reservas").Select("id_reserva, fecha_agenda").Where("estado_reserva = ?", "rechazada").Scan(&reservas)

	for _, r := range reservas {
		if r.FechaAgenda != nil {
			fmt.Printf("ID: %d, FechaAgenda: %s (Mes: %d, Año: %d)\n", r.IDReserva, r.FechaAgenda.Format("2006-01-02"), r.FechaAgenda.Month(), r.FechaAgenda.Year())
		} else {
			fmt.Printf("ID: %d, FechaAgenda: NULL\n", r.IDReserva)
		}
	}
}
