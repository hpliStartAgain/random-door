package api

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/your-org/city-roam/backend/internal/model"
	"github.com/your-org/city-roam/backend/internal/service"
)

type handlerProfileRepo struct {
	user *model.User
}

func (r *handlerProfileRepo) FindByID(context.Context, int64) (*model.User, error) {
	return r.user, nil
}

func (r *handlerProfileRepo) UpdateProfile(_ context.Context, _ int64, fields map[string]any) error {
	if v, ok := fields["nickname"].(string); ok {
		r.user.Nickname = &v
	}
	if v, ok := fields["age"].(int); ok {
		r.user.Age = &v
	}
	if v, ok := fields["home_region"].(string); ok {
		r.user.HomeRegion = &v
	}
	return nil
}

func TestUserHandlerProfile(t *testing.T) {
	gin.SetMode(gin.TestMode)
	repo := &handlerProfileRepo{user: &model.User{ID: 1, AnonymousID: handlerTestAnonymousID}}
	handler := NewUserHandler(service.NewUserService(repo))
	router := gin.New()
	router.GET("/api/users/:user_id/profile", handler.Profile)
	router.PATCH("/api/users/:user_id/profile", handler.UpdateProfile)

	body := requestJSON(t, router, http.MethodPatch, "/api/users/1/profile", `{"nickname":"游客","age":28,"home_region":"广东"}`, http.StatusOK)
	if body["nickname"] != "游客" || body["age"] != float64(28) || body["home_region"] != "广东" {
		t.Fatalf("PATCH profile = %#v, want updated fields", body)
	}

	body = requestJSON(t, router, http.MethodGet, "/api/users/1/profile", "", http.StatusOK)
	if body["user_id"] != float64(1) || body["anonymous_id"] != handlerTestAnonymousID {
		t.Fatalf("GET profile = %#v, want profile fields", body)
	}
}

func requestJSON(t *testing.T, router *gin.Engine, method, path, requestBody string, wantStatus int) map[string]any {
	t.Helper()
	req := httptest.NewRequest(method, path, bytes.NewBufferString(requestBody))
	if requestBody != "" {
		req.Header.Set("Content-Type", "application/json")
	}
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)
	if rec.Code != wantStatus {
		t.Fatalf("%s %s status = %d body=%s, want %d", method, path, rec.Code, rec.Body.String(), wantStatus)
	}
	var body map[string]any
	if err := json.Unmarshal(rec.Body.Bytes(), &body); err != nil {
		t.Fatalf("decode response: %v", err)
	}
	return body
}
