package report

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	service *Service
}

func NewHandler(service *Service) *Handler {
	return &Handler{
		service: service,
	}
}

// GetRevenueReport godoc
// @Summary Get revenue report
// @Description Get revenue report by type (weekly, monthly, yearly)
// @Tags reports
// @Produce json
// @Security BearerAuth
// @Param type query string true "Report Type" Enums(weekly,monthly,yearly)
// @Param month query int false "Month (required for monthly report)"
// @Param year query int false "Year (required for monthly/yearly report)"
// @Success 200 {array} map[string]interface{}
// @Router /reports/revenue [get]
func (h *Handler) GetRevenueReport(c *gin.Context) {

	reportType := c.Query("type")
	month := c.Query("month")
	year := c.Query("year")

	// ================= VALIDASI TYPE =================
	validType := map[string]bool{
		"weekly":  true,
		"monthly": true,
		"yearly":  true,
	}

	if !validType[reportType] {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid report type",
		})
		return
	}

	// ================= VALIDASI MONTHLY =================
	if reportType == "monthly" {
		if month == "" || year == "" {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "month and year are required for monthly report",
			})
			return
		}
	}

	// ================= VALIDASI YEARLY =================
	if reportType == "yearly" {
		if year == "" {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "year is required for yearly report",
			})
			return
		}
	}

	data, err := h.service.GetRevenueReport(
		c,
		reportType,
		month,
		year,
	)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, data)
}