package order

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

type CreateOrderRequest struct {
	IdCustomer    int     `json:"id_customer"`
	IdService     int     `json:"id_service" binding:"required"`
	Berat         float64 `json:"berat" binding:"required"`
	Catatan       string  `json:"catatan"`
	LangsungBayar bool    `json:"langsung_bayar"`
	Metode        *string `json:"metode"`
}

type UpdateOrderRequest struct {
	Berat   float64 `json:"berat"`
	Catatan string  `json:"catatan"`
}

type UpdateOrderStatusRequest struct {
	Status string `json:"status" binding:"required"`
}

type UpdateOrderPaymentRequest struct {
	Jumlah  float64 `json:"jumlah" binding:"required"`
	Metode  string  `json:"metode" binding:"required"`
}

// ================= HANDLER =================

// CreateOrder godoc
// @Summary Create order
// @Tags orders
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body CreateOrderRequest true "Order Data"
// @Success 200 {object} map[string]interface{}
// @Router /orders [post]
func (h *Handler) CreateOrder(c *gin.Context) {
	var req CreateOrderRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// VALIDASI LOGIC
	if req.LangsungBayar && req.Metode == nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "metode wajib diisi jika langsung_bayar = true",
		})
		return
	}

	err := h.service.CreateOrder(
		c,
		req.IdCustomer,
		req.IdService,
		req.Berat,
		req.Catatan,
		req.Metode,
		req.LangsungBayar,
	)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "order created",
	})
}

// GetOrders godoc
// @Summary Get all orders
// @Tags orders
// @Produce json
// @Security BearerAuth
// @Success 200 {array} map[string]interface{}
// @Router /orders [get]
func (h *Handler) GetOrders(c *gin.Context) {
	data, err := h.service.GetOrders(c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, data)
}

// GetOrderByID godoc
// @Summary Get order by ID
// @Tags orders
// @Security BearerAuth
// @Param id path int true "Order ID"
// @Success 200 {object} map[string]interface{}
// @Router /orders/{id} [get]
func (h *Handler) GetOrderByID(c *gin.Context) {
	var id int
	_, err := fmt.Sscan(c.Param("id"), &id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	data, err := h.service.GetOrderByID(c, id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "order not found"})
		return
	}

	c.JSON(http.StatusOK, data)
}

// UpdateOrder godoc
// @Summary Update order
// @Tags orders
// @Accept json
// @Security BearerAuth
// @Param id path int true "Order ID"
// @Param request body UpdateOrderRequest true "Order Data"
// @Success 200 {object} map[string]interface{}
// @Router /orders/{id} [put]
func (h *Handler) UpdateOrder(c *gin.Context) {
	var id int
	_, err := fmt.Sscan(c.Param("id"), &id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	var req UpdateOrderRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err = h.service.UpdateOrder(c, id, req.Berat, req.Catatan)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "order updated"})
}

// DeleteOrder godoc
// @Summary Delete order
// @Tags orders
// @Security BearerAuth
// @Param id path int true "Order ID"
// @Success 200 {object} map[string]interface{}
// @Router /orders/{id} [delete]
func (h *Handler) DeleteOrder(c *gin.Context) {
	var id int
	_, err := fmt.Sscan(c.Param("id"), &id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	err = h.service.DeleteOrder(c, id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "order deleted"})
}

// UpdateOrderStatus godoc
// @Summary Update order status
// @Description Update status order (DITERIMA, DIPROSES, SELESAI, DIAMBIL)
// @Tags orders
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "Order ID"
// @Param request body UpdateOrderStatusRequest true "Status Data"
// @Success 200 {object} map[string]interface{}
// @Router /orders/{id}/status [put]
func (h *Handler) UpdateOrderStatus(c *gin.Context) {
	var id int
	_, err := fmt.Sscan(c.Param("id"), &id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	var req UpdateOrderStatusRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err = h.service.UpdateOrderStatus(c, id, req.Status)
	if err != nil {
		if err.Error() == "invalid status" {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "order status updated",
	})
}

// UpdateOrderPayment godoc
// @Summary Update order payment
// @Description Input pembayaran dan update status jadi SUDAH_BAYAR
// @Tags orders
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "Order ID"
// @Param request body UpdateOrderPaymentRequest true "Payment Data"
// @Success 200 {object} map[string]interface{}
// @Router /orders/{id}/payment [put]
func (h *Handler) UpdateOrderPayment(c *gin.Context) {
	var id int
	_, err := fmt.Sscan(c.Param("id"), &id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	var req UpdateOrderPaymentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err = h.service.UpdateOrderPayment(c, id, req.Jumlah, req.Metode)
	if err != nil {
		if err.Error() == "invalid payment method" {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "payment success",
	})
}


// GetActiveOrders godoc
// @Summary Get active orders
// @Tags orders
// @Produce json
// @Security BearerAuth
// @Success 200 {array} map[string]interface{}
// @Router /orders/active [get]
func (h *Handler) GetActiveOrders(c *gin.Context) {
	data, err := h.service.GetActiveOrders(c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, data)
}

// GetCompletedOrders godoc
// @Summary Get completed orders
// @Tags orders
// @Produce json
// @Security BearerAuth
// @Success 200 {array} map[string]interface{}
// @Router /orders/completed [get]
func (h *Handler) GetCompletedOrders(c *gin.Context) {
	data, err := h.service.GetCompletedOrders(c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, data)
}

// GetOrderSummary godoc
// @Summary Get order summary
// @Tags orders
// @Produce json
// @Security BearerAuth
// @Success 200 {object} map[string]interface{}
// @Router /orders/summary [get]
func (h *Handler) GetOrderSummary(c *gin.Context) {
	data, err := h.service.GetOrderSummary(c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, data)
}