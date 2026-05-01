package payment

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	service *Service
}

func NewHandler(service *Service) *Handler {
	return &Handler{service: service}
}

// ================= REQUEST =================
type CreatePaymentRequest struct {
	IdOrder int     `json:"id_order" binding:"required"`
	Jumlah  float64 `json:"jumlah" binding:"required"`
	Metode  string  `json:"metode" binding:"required"`
}

// ================= HANDLER =================
// CreatePayment godoc
// @Summary Create payment
// @Tags payments
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body CreatePaymentRequest true "Payment Data"
// @Success 200 {object} map[string]interface{}
// @Router /payments [post]
func (h *Handler) CreatePayment(c *gin.Context) {
	var req CreatePaymentRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := h.service.CreatePayment(c, req.IdOrder, req.Jumlah, req.Metode)
	if err != nil {
		if err.Error() == "invalid payment method" {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "payment created"})
}


// GetPayments godoc
// @Summary Get all payments
// @Tags payments
// @Produce json
// @Security BearerAuth
// @Success 200 {array} map[string]interface{}
// @Router /payments [get]
func (h *Handler) GetPayments(c *gin.Context) {
	data, err := h.service.GetPayments(c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, data)
}


// GetPaymentByID godoc
// @Summary Get payment by ID
// @Tags payments
// @Security BearerAuth
// @Param id path int true "Payment ID"
// @Success 200 {object} map[string]interface{}
// @Router /payments/{id} [get]
func (h *Handler) GetPaymentByID(c *gin.Context) {
	var id int
	_, err := fmt.Sscan(c.Param("id"), &id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	data, err := h.service.GetPaymentByID(c, id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "payment not found"})
		return
	}

	c.JSON(http.StatusOK, data)
}