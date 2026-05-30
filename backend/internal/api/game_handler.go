package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/your-org/city-roam/backend/internal/service"
)

type GameHandler struct {
	svc *service.GameService
}

func NewGameHandler(svc *service.GameService) *GameHandler {
	return &GameHandler{svc: svc}
}

type initReq struct {
	UserID int64   `json:"user_id" binding:"required"`
	Lat    float64 `json:"lat"`
	Lng    float64 `json:"lng"`
}

// Init handles POST /api/game/init
func (h *GameHandler) Init(c *gin.Context) {
	var req initReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, errorResp("INVALID_PARAM", "user_id is required"))
		return
	}

	result, err := h.svc.Init(c.Request.Context(), req.UserID, req.Lat, req.Lng)
	if err != nil {
		writeServiceError(c, err)
		return
	}
	c.JSON(http.StatusOK, result)
}

type rollReq struct {
	UserID     int64   `json:"user_id" binding:"required"`
	FromCityID int64   `json:"from_city_id" binding:"required"`
	Lat        float64 `json:"lat" binding:"required"`
	Lng        float64 `json:"lng" binding:"required"`
}

// Roll handles POST /api/game/roll
func (h *GameHandler) Roll(c *gin.Context) {
	var req rollReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, errorResp("INVALID_PARAM", "user_id, from_city_id, lat, lng are required"))
		return
	}

	result, err := h.svc.Roll(c.Request.Context(), req.UserID, req.FromCityID, req.Lat, req.Lng)
	if err != nil {
		writeServiceError(c, err)
		return
	}
	c.JSON(http.StatusOK, result)
}
