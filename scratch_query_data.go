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

	type UserInfo struct {
		Correo string `gorm:"column:correo"`
		Rol    string `gorm:"column:rol"`
	}

	var users []UserInfo
	db.Raw("SELECT correo, rol FROM seguridad.usuarios").Scan(&users)
	for _, u := range users {
		fmt.Printf("- Correo: %s, Rol: '%s'\n", u.Correo, u.Rol)
	}
}
