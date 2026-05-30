package api

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/your-org/city-roam/backend/internal/service"
)

type CityHandler struct {
	svc *service.CityService
}

func NewCityHandler(svc *service.CityService) *CityHandler {
	return &CityHandler{svc: svc}
}

// List handles GET /api/cities
func (h *CityHandler) List(c *gin.Context) {
	cities, err := h.svc.List(c.Request.Context())
	if err != nil {
		writeServiceError(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"cities": cities})
}

// Detail handles GET /api/cities/:city_id
func (h *CityHandler) Detail(c *gin.Context) {
	cityID, err := strconv.ParseInt(c.Param("city_id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, errorResp("INVALID_PARAM", "city_id must be a positive integer"))
		return
	}

	detail, err := h.svc.Detail(c.Request.Context(), cityID)
	if err != nil {
		writeServiceError(c, err)
		return
	}
	c.JSON(http.StatusOK, detail)
}
