package main

import (
	"log"
	"os"
	"time"

	"apiServiYa/internal/application/auth"
	"apiServiYa/internal/application/reportes"
	"apiServiYa/internal/infrastructure/repository"
	presentation_http "apiServiYa/internal/presentation/http"
	"apiServiYa/internal/presentation/http/middleware"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	// Swagger
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"apiServiYa/docs"
)

// @title API REST ServiYa (Reportes)
// @version 2.0
// @description API de Reportes generada con Onion Architecture, Unit of Work y Repository Pattern
// @host localhost:8080
// @BasePath /api/v1
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
func main() {
	// Cargar variables de entorno desde .env (solo en local, en Render no existe y eso está bien)
	if err := godotenv.Load(); err != nil {
		log.Println("ℹ️ No se encontró archivo .env, usando variables de entorno del sistema")
	}

	// Configurar Swagger dinámicamente según el entorno
	if os.Getenv("RENDER") != "" || os.Getenv("PORT") != "" {
		docs.SwaggerInfo.Host = ""
		docs.SwaggerInfo.Schemes = []string{"https", "http"}
	} else {
		docs.SwaggerInfo.Host = "localhost:8080"
		docs.SwaggerInfo.Schemes = []string{"http", "https"}
	}

	// 1. Conexión a la Base de Datos (Supabase)
	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		log.Fatal("❌ DATABASE_URL no está configurada. Agrega un archivo .env o configura la variable de entorno.")
	}
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Println("⚠️ Advertencia: No se pudo conectar a la base de datos (ignora esto si es prueba de compilación)")
	}

	// 2. Inicialización de Repositorios e Infraestructura
	reporteUoW := repository.NewReporteUnitOfWork(db)
	authRepo := repository.NewAuthRepository(db)

	// 3. Inicialización de Casos de Uso (Capa de Aplicación)
	reporteAdminUseCase := reportes.NewGenerarReporteConsolidadoUseCase(reporteUoW)
	serviciosPopularesUseCase := reportes.NewObtenerServiciosPopularesUseCase(reporteUoW)
	actividadUsuariosUseCase := reportes.NewObtenerActividadUsuariosUseCase(reporteUoW)
	loginUseCase := auth.NewLoginUseCase(authRepo)

	// 4. Inicialización de Controladores (Capa de Presentación)
	reportesController := presentation_http.NewReportesController(
		reporteAdminUseCase,
		serviciosPopularesUseCase,
		actividadUsuariosUseCase,
	)
	authController := presentation_http.NewAuthController(loginUseCase)

	// 5. Configuración de Gin Router
	router := gin.Default()

	// Middleware de CORS (gin-contrib/cors)
	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:3000", "https://tu-dominio-vercel.com"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	// Endpoints Públicos (Autenticación)
	authGroup := router.Group("/api/v1/auth")
	{
		authGroup.POST("/login", authController.Login)
	}

	// Endpoints Protegidos (Reportes)
	reportesGroup := router.Group("/api/v1/reportes")
	reportesGroup.Use(middleware.AuthMiddleware())
	{
		reportesGroup.GET("/admin", reportesController.ObtenerReporteAdmin)
		reportesGroup.GET("/admin/pdf", reportesController.DescargarReporteAdminPDF)
		reportesGroup.GET("/servicios-populares", reportesController.ObtenerServiciosPopulares)
		reportesGroup.GET("/actividad-usuarios", reportesController.ObtenerActividadUsuarios)
	}

	// Ruta de Swagger UI
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// Redirección amigable: /docs manda directo a Swagger
	router.GET("/docs", func(c *gin.Context) {
		c.Redirect(301, "/swagger/index.html")
	})

	// Endpoint raíz para verificar que la API está corriendo
	router.GET("/", func(c *gin.Context) {
		c.String(200, "API DE SERVIYA CORRIENDO EN LENGIAJE GO CON ORM GORM")
	})

	// 6. Arrancar servidor leyendo puerto dinámico (requerido por Render)
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	
	log.Println("=========================================================")
	log.Println("🚀 API REST ServiYa iniciada correctamente")
	log.Println("=========================================================")
	log.Printf("📖 Documentación:   http://localhost:%s/docs\n", port)
	log.Printf("🌐 Estado API:      http://localhost:%s/\n", port)
	log.Printf("📄 Descargar PDF:   http://localhost:%s/api/v1/reportes/admin/pdf?mes=6&anio=2026\n", port)
	log.Println("=========================================================")

	if err := router.Run(":" + port); err != nil {
		log.Fatal("Error al iniciar el servidor:", err)
	}
}

