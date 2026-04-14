package service

import (
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
)

func RegisterRoutes(r *gin.Engine, db *pgxpool.Pool) {
	repo := NewRepository(db)
	svc := NewService(repo)
	handler := NewHandler(svc)

	group := r.Group("/services")
	{
		group.POST("", handler.CreateService)
		group.GET("", handler.GetServices)
		group.GET("/:id", handler.GetServiceByID)
		group.PUT("/:id", handler.UpdateService)
		group.DELETE("/:id", handler.DeleteService)
	}
}