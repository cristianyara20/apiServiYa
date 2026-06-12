package main

import (
	"fmt"
	"os"
	"log"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"github.com/joho/godotenv"
)

func main() {
	godotenv.Load(".env")
	db, err := gorm.Open(postgres.Open(os.Getenv("DATABASE_URL")), &gorm.Config{})
	if err != nil {
		log.Fatal(err)
	}

	type ColumnInfo struct {
		ColumnName string `gorm:"column:column_name"`
		DataType   string `gorm:"column:data_type"`
	}

	var cols []ColumnInfo
	db.Raw("SELECT column_name, data_type FROM information_schema.columns WHERE table_schema = 'seguridad' AND table_name = 'usuarios'").Scan(&cols)
	for _, c := range cols {
		fmt.Printf("- %s (%s)\n", c.ColumnName, c.DataType)
	}
}
