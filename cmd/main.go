package main

import (
	"queen-laundry/config"
	"queen-laundry/pkg/database"
	"queen-laundry/pkg/middleware"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"

	_ "queen-laundry/docs"
	"queen-laundry/internal/auth"
	"queen-laundry/internal/customer"
	"queen-laundry/internal/order"
	"queen-laundry/internal/payment"
	"queen-laundry/internal/report"
	"queen-laundry/internal/service"

	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// @title Queen Laundry API
// @version 1.0
// @description API untuk sistem laundry

// @host api.queenlaundry.com
// @BasePath /

// JWT CONFIG
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
func main() {
	cfg := config.LoadConfig()

	db := database.ConnectDB(cfg.DBUrl)

	r := gin.Default()

	// CORS
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"}, // bisa diganti domain nanti
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
	}))

	// Swagger
	r.GET("/docs/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// Health check
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})

	// TODO: register routes modular di sini
	// ================= PUBLIC =================
	auth.RegisterRoutes(r, db)

	// ================= PROTECTED =================
	protected := r.Group("/")
	protected.Use(middleware.JWTAuthMiddleware())

	order.RegisterRoutes(protected, db)
	customer.RegisterRoutes(protected, db)
	payment.RegisterRoutes(protected, db)
	service.RegisterRoutes(protected, db)
	report.RegisterRoutes(protected, db)
	_ = db

	r.Run(":" + cfg.Port)
}