package customer

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

type CreateCustomerRequest struct {
	Nama   string `json:"nama" binding:"required"`
	NoHP   string `json:"no_hp"`
	Alamat string `json:"alamat"`
}

type UpdateCustomerRequest struct {
	Nama   string `json:"nama" binding:"required"`
	NoHP   string `json:"no_hp"`
	Alamat string `json:"alamat"`
}

// ================= HANDLER =================

// CreateCustomer godoc
// @Summary Create customer
// @Description Tambah customer
// @Tags customers
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body CreateCustomerRequest true "Customer Data"
// @Success 200 {object} map[string]interface{}
// @Router /customers [post]
func (h *Handler) CreateCustomer(c *gin.Context) {
	var req CreateCustomerRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := h.service.CreateCustomer(c, req.Nama, req.NoHP, req.Alamat)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "customer created",
	})
}

// GetCustomers godoc
// @Summary Get all customers
// @Description Ambil semua customer
// @Tags customers
// @Produce json
// @Security BearerAuth
// @Success 200 {array} map[string]interface{}
// @Router /customers [get]
func (h *Handler) GetCustomers(c *gin.Context) {
	data, err := h.service.GetCustomers(c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, data)
}


// GetCustomerByID godoc
// @Summary Get customer by ID
// @Description Ambil customer berdasarkan ID
// @Tags customers
// @Produce json
// @Security BearerAuth
// @Param id path int true "Customer ID"
// @Success 200 {object} map[string]interface{}
// @Router /customers/{id} [get]
func (h *Handler) GetCustomerByID(c *gin.Context) {
	idParam := c.Param("id")

	var id int
	_, err := fmt.Sscan(idParam, &id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	data, err := h.service.GetCustomerByID(c, id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "customer not found"})
		return
	}

	c.JSON(http.StatusOK, data)
}


// UpdateCustomer godoc
// @Summary Update customer
// @Description Update data customer
// @Tags customers
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "Customer ID"
// @Param request body UpdateCustomerRequest true "Customer Data"
// @Success 200 {object} map[string]interface{}
// @Router /customers/{id} [put]
func (h *Handler) UpdateCustomer(c *gin.Context) {
	idParam := c.Param("id")

	var id int
	_, err := fmt.Sscan(idParam, &id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	var req UpdateCustomerRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err = h.service.UpdateCustomer(c, id, req.Nama, req.NoHP, req.Alamat)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "customer updated"})
}


// DeleteCustomer godoc
// @Summary Delete customer
// @Description Soft delete customer
// @Tags customers
// @Security BearerAuth
// @Param id path int true "Customer ID"
// @Success 200 {object} map[string]interface{}
// @Router /customers/{id} [delete]
func (h *Handler) DeleteCustomer(c *gin.Context) {
	idParam := c.Param("id")

	var id int
	_, err := fmt.Sscan(idParam, &id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	err = h.service.DeleteCustomer(c, id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "customer deleted"})
}