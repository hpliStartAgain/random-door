package api

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/your-org/city-roam/backend/internal/service"
)

type UserHandler struct {
	svc *service.UserService
}

func NewUserHandler(svc *service.UserService) *UserHandler {
	return &UserHandler{svc: svc}
}

type updateProfileReq struct {
	Nickname   *string `json:"nickname"`
	Age        *int    `json:"age"`
	HomeRegion *string `json:"home_region"`
}

type registerReq struct {
	UserID   *int64  `json:"user_id"`
	Username string  `json:"username" binding:"required"`
	Password string  `json:"password" binding:"required"`
	Nickname *string `json:"nickname"`
}

type loginReq struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

func (h *UserHandler) Register(c *gin.Context) {
	var req registerReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, errorResp(ErrCodeInvalidParam, "username and password are required"))
		return
	}
	result, err := h.svc.Register(c.Request.Context(), service.RegisterRequest{
		UserID:   req.UserID,
		Username: req.Username,
		Password: req.Password,
		Nickname: req.Nickname,
	})
	if err != nil {
		writeServiceError(c, err)
		return
	}
	c.JSON(http.StatusCreated, result)
}

func (h *UserHandler) Login(c *gin.Context) {
	var req loginReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, errorResp(ErrCodeInvalidParam, "username and password are required"))
		return
	}
	result, err := h.svc.Login(c.Request.Context(), service.LoginRequest{
		Username: req.Username,
		Password: req.Password,
	})
	if err != nil {
		writeServiceError(c, err)
		return
	}
	c.JSON(http.StatusOK, result)
}

func (h *UserHandler) Profile(c *gin.Context) {
	userID, err := parseUserIDParam(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, errorResp(ErrCodeInvalidParam, "user_id must be a positive integer"))
		return
	}
	profile, err := h.svc.Profile(c.Request.Context(), userID)
	if err != nil {
		writeServiceError(c, err)
		return
	}
	c.JSON(http.StatusOK, profile)
}

func (h *UserHandler) UpdateProfile(c *gin.Context) {
	userID, err := parseUserIDParam(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, errorResp(ErrCodeInvalidParam, "user_id must be a positive integer"))
		return
	}
	var req updateProfileReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, errorResp(ErrCodeInvalidParam, "invalid profile payload"))
		return
	}
	profile, err := h.svc.UpdateProfile(c.Request.Context(), userID, service.UpdateUserProfileRequest{
		Nickname:   req.Nickname,
		Age:        req.Age,
		HomeRegion: req.HomeRegion,
	})
	if err != nil {
		writeServiceError(c, err)
		return
	}
	c.JSON(http.StatusOK, profile)
}

func parseUserIDParam(c *gin.Context) (int64, error) {
	userID, err := strconv.ParseInt(c.Param("user_id"), 10, 64)
	if err != nil || userID <= 0 {
		if err == nil {
			err = errors.New("user_id must be positive")
		}
		return 0, err
	}
	return userID, nil
}
