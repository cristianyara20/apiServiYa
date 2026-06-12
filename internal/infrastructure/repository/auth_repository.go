package repository

import (
	"context"
	"gorm.io/gorm"
)

type AuthRepository struct {
	db *gorm.DB
}

func NewAuthRepository(db *gorm.DB) *AuthRepository {
	return &AuthRepository{db: db}
}

func (r *AuthRepository) EsAdministrador(ctx context.Context, correo string) (bool, error) {
	var rol string
	err := r.db.WithContext(ctx).Table("seguridad.usuarios").Select("rol").Where("correo = ?", correo).Scan(&rol).Error
	if err != nil {
		return false, err
	}
	// Asumimos que los administradores tienen el rol 'ADMIN' o 'administrador'
	// Ajusta esto si tu rol en BD se llama diferente
	if rol == "ADMIN" || rol == "administrador" || rol == "admin" {
		return true, nil
	}
	return false, nil
}
