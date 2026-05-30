package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/your-org/city-roam/backend/internal/service"
)

type VisitHandler struct {
	svc *service.VisitService
}

func NewVisitHandler(svc *service.VisitService) *VisitHandler {
	return &VisitHandler{svc: svc}
}

type anonymousUserReq struct {
	AnonymousID string `json:"anonymous_id" binding:"required"`
}

// CreateAnonymousUser handles POST /api/users/anonymous
func (h *VisitHandler) CreateAnonymousUser(c *gin.Context) {
	var req anonymousUserReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, errorResp("INVALID_PARAM", "anonymous_id is required"))
		return
	}

	user, err := h.svc.CreateAnonymousUser(c.Request.Context(), req.AnonymousID)
	if err != nil {
		writeServiceError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"user_id":         user.ID,
		"anonymous_id":    user.AnonymousID,
		"current_city_id": user.CurrentCityID,
	})
}

type freeVisitReq struct {
	UserID int64  `json:"user_id" binding:"required"`
	CityID int64  `json:"city_id" binding:"required"`
	Source string `json:"source"`
}

// CreateFreeVisit handles POST /api/visits/free
func (h *VisitHandler) CreateFreeVisit(c *gin.Context) {
	var req freeVisitReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, errorResp("INVALID_PARAM", "user_id and city_id are required"))
		return
	}

	visit, err := h.svc.CreateFreeVisit(c.Request.Context(), req.UserID, req.CityID, req.Source)
	if err != nil {
		writeServiceError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"visit_id":   visit.ID,
		"city_id":    visit.CityID,
		"visit_mode": visit.VisitMode,
	})
}
