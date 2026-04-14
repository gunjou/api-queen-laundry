package payment

import (
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
)

func RegisterRoutes(r *gin.Engine, db *pgxpool.Pool) {
	repo := NewRepository(db)
	svc := NewService(repo)
	handler := NewHandler(svc)

	group := r.Group("/payments")
	{
		group.POST("", handler.CreatePayment)
		group.GET("", handler.GetPayments)
		group.GET("/:id", handler.GetPaymentByID)
	}
}