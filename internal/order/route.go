package order

import (
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
)

func RegisterRoutes(r *gin.Engine, db *pgxpool.Pool) {
	repo := NewRepository(db)
	svc := NewService(repo)
	handler := NewHandler(svc)

	group := r.Group("/orders")
	{
		group.POST("", handler.CreateOrder)
		group.GET("", handler.GetOrders)
		group.GET("/:id", handler.GetOrderByID)
		group.PUT("/:id", handler.UpdateOrder)
		group.DELETE("/:id", handler.DeleteOrder)
		group.PUT("/:id/status", handler.UpdateOrderStatus)
		group.PUT("/:id/payment", handler.UpdateOrderPayment)
		group.GET("/active", handler.GetActiveOrders)
		group.GET("/completed", handler.GetCompletedOrders)
		group.GET("/summary", handler.GetOrderSummary)
	}
}