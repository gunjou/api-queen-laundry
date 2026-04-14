package service

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

// ======================= REQUEST STRUCT =======================

type CreateServiceRequest struct {
	Nama  string  `json:"nama" binding:"required"`
	Tipe  string  `json:"tipe" binding:"required"`
	Harga float64 `json:"harga" binding:"required"`
}

type UpdateServiceRequest struct {
	Nama  string  `json:"nama" binding:"required"`
	Tipe  string  `json:"tipe" binding:"required"`
	Harga float64 `json:"harga" binding:"required"`
}
// ======================= HANDLERS =======================

// CreateService godoc
// @Summary Create service
// @Description Tambah layanan laundry
// @Tags services
// @Accept json
// @Produce json
// @Param request body CreateServiceRequest true "Service Data"
// @Success 200 {object} map[string]interface{}
// @Router /services [post]
func (h *Handler) CreateService(c *gin.Context) {
	var req CreateServiceRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := h.service.CreateService(c, req.Nama, req.Tipe, req.Harga)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "service created",
	})
}

// GetServices godoc
// @Summary Get all services
// @Description Ambil semua layanan
// @Tags services
// @Produce json
// @Success 200 {array} map[string]interface{}
// @Router /services [get]
func (h *Handler) GetServices(c *gin.Context) {
	data, err := h.service.GetServices(c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, data)
}


// GetServiceByID godoc
// @Summary Get service by ID
// @Description Ambil service berdasarkan ID
// @Tags services
// @Produce json
// @Param id path int true "Service ID"
// @Success 200 {object} map[string]interface{}
// @Router /services/{id} [get]
func (h *Handler) GetServiceByID(c *gin.Context) {
	idParam := c.Param("id")

	var id int
	_, err := fmt.Sscan(idParam, &id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	data, err := h.service.GetServiceByID(c, id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "service not found"})
		return
	}

	c.JSON(http.StatusOK, data)
}


// UpdateService godoc
// @Summary Update service
// @Description Update layanan berdasarkan ID
// @Tags services
// @Accept json
// @Produce json
// @Param id path int true "Service ID"
// @Param request body UpdateServiceRequest true "Service Data"
// @Success 200 {object} map[string]interface{}
// @Router /services/{id} [put]
func (h *Handler) UpdateService(c *gin.Context) {
	idParam := c.Param("id")

	var id int
	_, err := fmt.Sscan(idParam, &id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	var req UpdateServiceRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err = h.service.UpdateService(c, id, req.Nama, req.Tipe, req.Harga)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "service updated"})
}


// DeleteService godoc
// @Summary Delete service
// @Description Soft delete service (is_active = 0)
// @Tags services
// @Param id path int true "Service ID"
// @Success 200 {object} map[string]interface{}
// @Router /services/{id} [delete]
func (h *Handler) DeleteService(c *gin.Context) {
	idParam := c.Param("id")

	var id int
	_, err := fmt.Sscan(idParam, &id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	err = h.service.DeleteService(c, id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "service deleted"})
}