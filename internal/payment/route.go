package payment

import (
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
)

func RegisterRoutes(rg *gin.RouterGroup, db *pgxpool.Pool) {
	repo := NewRepository(db)
	svc := NewService(repo)
	handler := NewHandler(svc)

	group := rg.Group("/payments")
	{
		group.POST("", handler.CreatePayment)
		group.GET("", handler.GetPayments)
		group.GET("/:id", handler.GetPaymentByID)
	}
}