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

	type TableInfo struct {
		TableSchema string `gorm:"column:table_schema"`
		TableName   string `gorm:"column:table_name"`
	}

	var tables []TableInfo
	db.Raw("SELECT table_schema, table_name FROM information_schema.tables WHERE table_schema IN ('gestion', 'seguridad', 'soporte')").Scan(&tables)
	for _, t := range tables {
		fmt.Printf("Tabla: %s.%s\n", t.TableSchema, t.TableName)
		type ColumnInfo struct {
			ColumnName string `gorm:"column:column_name"`
			DataType   string `gorm:"column:data_type"`
		}
		var cols []ColumnInfo
		db.Raw("SELECT column_name, data_type FROM information_schema.columns WHERE table_schema = ? AND table_name = ?", t.TableSchema, t.TableName).Scan(&cols)
		for _, c := range cols {
			fmt.Printf("  - %s (%s)\n", c.ColumnName, c.DataType)
		}
	}
}
