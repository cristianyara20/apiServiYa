package middleware

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
)

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			log.Println("Middleware: Falta el token de autorización en el header")
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Falta el token de autorización"})
			return
		}

		// Permitir tanto "Bearer <token>" como solo "<token>"
		var tokenString string
		if strings.HasPrefix(authHeader, "Bearer ") {
			tokenString = strings.TrimPrefix(authHeader, "Bearer ")
		} else {
			tokenString = authHeader
		}

		// Validar el token consultando a Supabase
		supabaseURL := os.Getenv("SUPABASE_URL")
		anonKey := os.Getenv("SUPABASE_ANON_KEY")

		if supabaseURL == "" || anonKey == "" {
			log.Println("Middleware: Configuración de Supabase incompleta")
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Configuración de Supabase incompleta en el servidor"})
			return
		}

		req, err := http.NewRequest("GET", fmt.Sprintf("%s/auth/v1/user", supabaseURL), nil)
		if err != nil {
			log.Println("Middleware: Error interno al crear request", err)
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Error interno del servidor"})
			return
		}

		req.Header.Set("Authorization", "Bearer "+tokenString)
		req.Header.Set("apikey", anonKey)

		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			log.Println("Middleware: Error al contactar con Supabase", err)
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Token inválido o expirado"})
			return
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			bodyBytes, _ := io.ReadAll(resp.Body)
			log.Printf("Middleware: Supabase rechazó el token. Status: %d, Body: %s\n", resp.StatusCode, string(bodyBytes))
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Token inválido o expirado"})
			return
		}

		// Si llegamos aquí, Supabase confirmó que el token es válido
		c.Next()
	}
}
