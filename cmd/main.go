package main

import (
	"queen-laundry/config"
	"queen-laundry/pkg/database"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"

	_ "queen-laundry/docs"
	"queen-laundry/internal/customer"
	"queen-laundry/internal/order"
	"queen-laundry/internal/service"

	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// @title Queen Laundry API
// @version 1.0
// @description API untuk sistem laundry
// @host localhost:5050
// @BasePath /
func main() {
	cfg := config.LoadConfig()

	db := database.ConnectDB(cfg.DBUrl)

	r := gin.Default()

	// CORS
	r.Use(cors.Default())

	// Swagger
	r.GET("/docs/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// Health check
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})

	// TODO: register routes modular di sini
	service.RegisterRoutes(r, db)
	customer.RegisterRoutes(r, db)
	order.RegisterRoutes(r, db)
	_ = db

	r.Run(":" + cfg.Port)
}