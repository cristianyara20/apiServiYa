package main

import (
	"fmt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	dsn := "postgresql://postgres.lvxhporsajorgckeisna:Cristianvargas2007%23@aws-1-us-east-1.pooler.supabase.com:5432/postgres?sslmode=require"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		panic(err)
	}

	type Result struct {
		EstadoReserva string
		Count         int
	}

	var results []Result
	db.Table("gestion.reservas").Select("estado_reserva, count(*) as count").Group("estado_reserva").Scan(&results)

	for _, r := range results {
		fmt.Printf("Estado: %s, Count: %d\n", r.EstadoReserva, r.Count)
	}
}
