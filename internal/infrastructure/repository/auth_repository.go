package repository

import (
	"apiServiYa/internal/domain"
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

func (r *AuthRepository) ExisteCorreo(ctx context.Context, correo string) (bool, error) {
	var count int64
	err := r.db.WithContext(ctx).Table("seguridad.usuarios").Where("correo = ?", correo).Count(&count).Error
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

func (r *AuthRepository) CrearUsuario(ctx context.Context, usuario *domain.UsuarioDB) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// 1. Verificar si ya existe en seguridad.usuarios por auth_id o correo
		var existing domain.UsuarioDB
		err := tx.Table("seguridad.usuarios").Where("auth_id = ? OR correo = ?", usuario.AuthID, usuario.Correo).First(&existing).Error
		if err == nil {
			usuario.IDUsuario = existing.IDUsuario
			updates := map[string]interface{}{}
			if existing.Nombre == "" && usuario.Nombre != "" {
				updates["nombre"] = usuario.Nombre
			}
			if existing.Apellido == "" && usuario.Apellido != "" {
				updates["apellido"] = usuario.Apellido
			}
			if usuario.FechaNacimiento != nil {
				updates["fecha_nacimiento"] = usuario.FechaNacimiento
			}

			if len(updates) > 0 {
				if err := tx.Table("seguridad.usuarios").Where("id_usuario = ?", existing.IDUsuario).Updates(updates).Error; err != nil {
					return err
				}
			}
		} else {
			if err := tx.Table("seguridad.usuarios").Create(usuario).Error; err != nil {
				return err
			}
		}

		// 2. Si el rol es 'usuario' (o está vacío), asegurar que exista registro en gestion.clientes
		if usuario.Rol == "" || usuario.Rol == "usuario" {
			type ClienteDB struct {
				IDCliente int    `gorm:"column:id_cliente;primaryKey"`
				AuthID    string `gorm:"column:auth_id"`
			}
			cliente := ClienteDB{
				IDCliente: usuario.IDUsuario,
				AuthID:    usuario.AuthID,
			}
			var existingCliente ClienteDB
			err := tx.Table("gestion.clientes").Where("id_cliente = ? OR auth_id = ?", usuario.IDUsuario, usuario.AuthID).First(&existingCliente).Error
			if err != nil {
				if err := tx.Table("gestion.clientes").Create(&cliente).Error; err != nil {
					return err
				}
			}
		}

		return nil
	})
}
