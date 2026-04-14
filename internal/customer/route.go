package customer

import (
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
)

func RegisterRoutes(r *gin.Engine, db *pgxpool.Pool) {
	repo := NewRepository(db)
	svc := NewService(repo)
	handler := NewHandler(svc)

	group := r.Group("/customers")
	{
		group.POST("", handler.CreateCustomer)
		group.GET("", handler.GetCustomers)
		group.GET("/:id", handler.GetCustomerByID)
		group.PUT("/:id", handler.UpdateCustomer)
		group.DELETE("/:id", handler.DeleteCustomer)
	}
}