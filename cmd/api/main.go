package main

import (
	"log"
	"os"

	"apiServiYa/internal/application/reportes"
	"apiServiYa/internal/infrastructure/repository"
	presentation_http "apiServiYa/internal/presentation/http"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	// Swagger
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	_ "apiServiYa/docs"
)

// @title API REST ServiYa (Reportes)
// @version 2.0
// @description API de Reportes generada con Onion Architecture, Unit of Work y Repository Pattern
// @host localhost:8080
// @BasePath /api/v1
func main() {
	// 1. Conexión a la Base de Datos (Supabase)
	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		dsn = "postgresql://postgres:Cristianvargas2007%23@db.lvxhporsajorgckeisna.supabase.co:5432/postgres?sslmode=require"
	}
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Println("⚠️ Advertencia: No se pudo conectar a la base de datos (ignora esto si es prueba de compilación)")
	}

	// 2. Inicialización de Repositorios e Infraestructura
	reporteUoW := repository.NewReporteUnitOfWork(db)

	// 3. Inicialización de Casos de Uso (Capa de Aplicación)
	reporteAdminUseCase := reportes.NewGenerarReporteConsolidadoUseCase(reporteUoW)
	serviciosPopularesUseCase := reportes.NewObtenerServiciosPopularesUseCase(reporteUoW)
	actividadUsuariosUseCase := reportes.NewObtenerActividadUsuariosUseCase(reporteUoW)

	// 4. Inicialización de Controladores (Capa de Presentación)
	reportesController := presentation_http.NewReportesController(
		reporteAdminUseCase,
		serviciosPopularesUseCase,
		actividadUsuariosUseCase,
	)

	// 5. Configuración de Gin Router
	router := gin.Default()

	// Middleware de CORS
	router.Use(func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	})

	// Endpoints de Reportes
	api := router.Group("/api/v1")
	{
		api.GET("/reportes/admin", reportesController.ObtenerReporteAdmin)
		api.GET("/reportes/servicios-populares", reportesController.ObtenerServiciosPopulares)
		api.GET("/reportes/actividad-usuarios", reportesController.ObtenerActividadUsuarios)
	}

	// Ruta de Swagger UI
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// Redirección amigable: /docs manda directo a Swagger
	router.GET("/docs", func(c *gin.Context) {
		c.Redirect(301, "/swagger/index.html")
	})

	// 6. Arrancar servidor leyendo puerto dinámico (requerido por Render)
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	log.Println("🚀 Servidor corriendo en el puerto " + port)
	if err := router.Run(":" + port); err != nil {
		log.Fatal("Error al iniciar el servidor:", err)
	}
}
