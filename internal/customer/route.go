package customer

import (
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
)

func RegisterRoutes(rg *gin.RouterGroup, db *pgxpool.Pool) {
	repo := NewRepository(db)
	svc := NewService(repo)
	handler := NewHandler(svc)

	group := rg.Group("/customers")
	{
		group.POST("", handler.CreateCustomer)
		group.GET("", handler.GetCustomers)
		group.GET("/:id", handler.GetCustomerByID)
		group.PUT("/:id", handler.UpdateCustomer)
		group.DELETE("/:id", handler.DeleteCustomer)
	}
}