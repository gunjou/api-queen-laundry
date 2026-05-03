package report

import (
	"queen-laundry/pkg/middleware"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
)

func RegisterRoutes(rg *gin.RouterGroup, db *pgxpool.Pool) {

	repo := NewRepository(db)
	service := NewService(repo)
	handler := NewHandler(service)

	group := rg.Group("/reports")
	group.Use(middleware.JWTAuthMiddleware())
	{
		group.GET("/revenue", handler.GetRevenueReport)
	}
}