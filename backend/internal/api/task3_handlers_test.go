package api

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/your-org/city-roam/backend/internal/model"
	"github.com/your-org/city-roam/backend/internal/service"
	"gorm.io/gorm"
)

const handlerTestAnonymousID = "550e8400-e29b-41d4-a716-446655440000"

type task3CityRepository struct {
	listAll        func(context.Context) ([]model.City, error)
	findByID       func(context.Context, int64) (*model.City, error)
	listTags       func(context.Context, int64) ([]model.CityTag, error)
	listLandmarks  func(context.Context, int64) ([]model.Landmark, error)
	listFoods      func(context.Context, int64) ([]model.Food, error)
	listCharacters func(context.Context, int64) ([]model.Character, error)
}

func (r *task3CityRepository) ListAll(ctx context.Context) ([]model.City, error) {
	if r.listAll == nil {
		return nil, errors.New("unexpected ListAll call")
	}
	return r.listAll(ctx)
}

func (r *task3CityRepository) FindByID(ctx context.Context, id int64) (*model.City, error) {
	if r.findByID == nil {
		return nil, errors.New("unexpected FindByID call")
	}
	return r.findByID(ctx, id)
}

func (r *task3CityRepository) ListTags(ctx context.Context, cityID int64) ([]model.CityTag, error) {
	if r.listTags == nil {
		return nil, errors.New("unexpected ListTags call")
	}
	return r.listTags(ctx, cityID)
}

func (r *task3CityRepository) ListLandmarks(ctx context.Context, cityID int64) ([]model.Landmark, error) {
	if r.listLandmarks == nil {
		return nil, errors.New("unexpected ListLandmarks call")
	}
	return r.listLandmarks(ctx, cityID)
}

func (r *task3CityRepository) ListFoods(ctx context.Context, cityID int64) ([]model.Food, error) {
	if r.listFoods == nil {
		return nil, errors.New("unexpected ListFoods call")
	}
	return r.listFoods(ctx, cityID)
}

func (r *task3CityRepository) ListCharacters(ctx context.Context, cityID int64) ([]model.Character, error) {
	if r.listCharacters == nil {
		return nil, errors.New("unexpected ListCharacters call")
	}
	return r.listCharacters(ctx, cityID)
}

type task3UserRepository struct {
	findByAnonymousID func(context.Context, string) (*model.User, error)
	findByID          func(context.Context, int64) (*model.User, error)
	create            func(context.Context, *model.User) error
}

func (r *task3UserRepository) FindByAnonymousID(ctx context.Context, anonymousID string) (*model.User, error) {
	if r.findByAnonymousID == nil {
		return nil, errors.New("unexpected FindByAnonymousID call")
	}
	return r.findByAnonymousID(ctx, anonymousID)
}

func (r *task3UserRepository) FindByID(ctx context.Context, id int64) (*model.User, error) {
	if r.findByID == nil {
		return nil, errors.New("unexpected FindByID call")
	}
	return r.findByID(ctx, id)
}

func (r *task3UserRepository) Create(ctx context.Context, user *model.User) error {
	if r.create == nil {
		return errors.New("unexpected Create call")
	}
	return r.create(ctx, user)
}

type task3VisitRepository struct {
	create func(context.Context, *model.CityVisit) error
}

func (r *task3VisitRepository) Create(ctx context.Context, visit *model.CityVisit) error {
	if r.create == nil {
		return errors.New("unexpected Create call")
	}
	return r.create(ctx, visit)
}

func TestTask3HandlersSuccess(t *testing.T) {
	var createdVisit *model.CityVisit
	cityRepo := successfulTask3CityRepository()
	userRepo := &task3UserRepository{
		findByAnonymousID: func(context.Context, string) (*model.User, error) {
			return &model.User{ID: 1, AnonymousID: handlerTestAnonymousID}, nil
		},
		findByID: func(context.Context, int64) (*model.User, error) {
			return &model.User{ID: 1, AnonymousID: handlerTestAnonymousID}, nil
		},
	}
	visitRepo := &task3VisitRepository{create: func(_ context.Context, visit *model.CityVisit) error {
		createdVisit = visit
		visit.ID = 1001
		return nil
	}}
	router := newTask3Router(cityRepo, userRepo, visitRepo)

	anonymous := assertTask3Response(t, router, http.MethodPost, "/api/users/anonymous", `{"anonymous_id":"`+handlerTestAnonymousID+`"}`, http.StatusOK, "")
	if anonymous["user_id"] != float64(1) || anonymous["anonymous_id"] != handlerTestAnonymousID || anonymous["current_city_id"] != nil {
		t.Fatalf("anonymous response = %#v, want contract fields", anonymous)
	}

	cities := assertTask3Response(t, router, http.MethodGet, "/api/cities", "", http.StatusOK, "")
	cityItems, ok := cities["cities"].([]any)
	if !ok || len(cityItems) != 1 || cityItems[0].(map[string]any)["name"] != "西安" {
		t.Fatalf("cities response = %#v, want city list", cities)
	}

	detail := assertTask3Response(t, router, http.MethodGet, "/api/cities/3", "", http.StatusOK, "")
	if detail["name"] != "西安" || strings.Contains(mustJSON(t, detail), "persona") || strings.Contains(mustJSON(t, detail), "prompt") {
		t.Fatalf("detail response = %#v, want public city detail", detail)
	}

	visit := assertTask3Response(t, router, http.MethodPost, "/api/visits/free", `{"user_id":1,"city_id":3}`, http.StatusOK, "")
	if visit["visit_id"] != float64(1001) || visit["city_id"] != float64(3) || visit["visit_mode"] != "free" {
		t.Fatalf("visit response = %#v, want free visit contract fields", visit)
	}

	if createdVisit == nil || createdVisit.VisitMode != "free" || createdVisit.Source == nil || *createdVisit.Source != "map_click" {
		t.Fatalf("created visit = %#v, want default free map_click visit", createdVisit)
	}
}

func TestTask3HandlersErrors(t *testing.T) {
	tests := []struct {
		name        string
		router      *gin.Engine
		method      string
		path        string
		body        string
		wantStatus  int
		wantCode    string
		wantMessage string
	}{
		{
			name: "anonymous id must be uuid", router: newTask3Router(nil, nil, nil),
			method: http.MethodPost, path: "/api/users/anonymous", body: `{"anonymous_id":"not-a-uuid"}`,
			wantStatus: http.StatusBadRequest, wantCode: "INVALID_PARAM", wantMessage: "anonymous_id must be a UUID",
		},
		{
			name: "anonymous id is required", router: newTask3Router(nil, nil, nil),
			method: http.MethodPost, path: "/api/users/anonymous", body: `{}`,
			wantStatus: http.StatusBadRequest, wantCode: "INVALID_PARAM", wantMessage: "anonymous_id is required",
		},
		{
			name: "city id must be integer", router: newTask3Router(nil, nil, nil),
			method: http.MethodGet, path: "/api/cities/nope",
			wantStatus: http.StatusBadRequest, wantCode: "INVALID_PARAM", wantMessage: "city_id must be a positive integer",
		},
		{
			name: "city id must be positive", router: newTask3Router(nil, nil, nil),
			method: http.MethodGet, path: "/api/cities/0",
			wantStatus: http.StatusBadRequest, wantCode: "INVALID_PARAM", wantMessage: "city_id must be a positive integer",
		},
		{
			name: "city detail missing", router: newTask3Router(&task3CityRepository{
				findByID: func(context.Context, int64) (*model.City, error) { return nil, gorm.ErrRecordNotFound },
			}, nil, nil),
			method: http.MethodGet, path: "/api/cities/404",
			wantStatus: http.StatusNotFound, wantCode: "NOT_FOUND", wantMessage: "city not found",
		},
		{
			name: "free visit invalid source", router: newTask3Router(nil, nil, nil),
			method: http.MethodPost, path: "/api/visits/free", body: `{"user_id":1,"city_id":3,"source":"dice_roll"}`,
			wantStatus: http.StatusBadRequest, wantCode: "INVALID_PARAM", wantMessage: "source must be map_click or search",
		},
		{
			name: "free visit ids are required", router: newTask3Router(nil, nil, nil),
			method: http.MethodPost, path: "/api/visits/free", body: `{}`,
			wantStatus: http.StatusBadRequest, wantCode: "INVALID_PARAM", wantMessage: "user_id and city_id are required",
		},
		{
			name: "free visit missing user", router: newTask3Router(nil, &task3UserRepository{
				findByID: func(context.Context, int64) (*model.User, error) { return nil, gorm.ErrRecordNotFound },
			}, nil),
			method: http.MethodPost, path: "/api/visits/free", body: `{"user_id":1,"city_id":3}`,
			wantStatus: http.StatusNotFound, wantCode: "NOT_FOUND", wantMessage: "user not found",
		},
		{
			name: "free visit missing city", router: newTask3Router(&task3CityRepository{
				findByID: func(context.Context, int64) (*model.City, error) { return nil, gorm.ErrRecordNotFound },
			}, &task3UserRepository{
				findByID: func(context.Context, int64) (*model.User, error) { return &model.User{ID: 1}, nil },
			}, nil),
			method: http.MethodPost, path: "/api/visits/free", body: `{"user_id":1,"city_id":3}`,
			wantStatus: http.StatusNotFound, wantCode: "NOT_FOUND", wantMessage: "city not found",
		},
		{
			name: "internal error is masked", router: newTask3Router(&task3CityRepository{
				listAll: func(context.Context) ([]model.City, error) { return nil, errors.New("mysql password secret") },
			}, nil, nil),
			method: http.MethodGet, path: "/api/cities",
			wantStatus: http.StatusInternalServerError, wantCode: "INTERNAL_ERROR", wantMessage: "internal server error",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			body := assertTask3Response(t, tc.router, tc.method, tc.path, tc.body, tc.wantStatus, tc.wantCode)
			if body["error"].(map[string]any)["message"] != tc.wantMessage {
				t.Fatalf("response = %#v, want message %q", body, tc.wantMessage)
			}
			if strings.Contains(mustJSON(t, body), "mysql password secret") {
				t.Fatalf("response leaked internal error: %#v", body)
			}
		})
	}
}

func newTask3Router(cityRepo service.CityRepository, userRepo service.UserRepository, visitRepo service.VisitRepository) *gin.Engine {
	gin.SetMode(gin.TestMode)
	cityHandler := NewCityHandler(service.NewCityService(cityRepo))
	visitHandler := NewVisitHandler(service.NewVisitService(userRepo, cityRepo, visitRepo))
	router := gin.New()
	router.POST("/api/users/anonymous", visitHandler.CreateAnonymousUser)
	router.GET("/api/cities", cityHandler.List)
	router.GET("/api/cities/:city_id", cityHandler.Detail)
	router.POST("/api/visits/free", visitHandler.CreateFreeVisit)
	return router
}

func successfulTask3CityRepository() *task3CityRepository {
	return &task3CityRepository{
		listAll: func(context.Context) ([]model.City, error) {
			return []model.City{{ID: 3, Name: "西安", Province: "陕西", Lat: 34.3416, Lng: 108.9398}}, nil
		},
		findByID: func(context.Context, int64) (*model.City, error) {
			return &model.City{ID: 3, Name: "西安", Province: "陕西", Lat: 34.3416, Lng: 108.9398}, nil
		},
		listTags: func(context.Context, int64) ([]model.CityTag, error) {
			return []model.CityTag{{CityID: 3, Tag: "ancient_capital"}}, nil
		},
		listLandmarks: func(context.Context, int64) ([]model.Landmark, error) {
			return []model.Landmark{{ID: 12, CityID: 3, Name: "兵马俑"}}, nil
		},
		listFoods: func(context.Context, int64) ([]model.Food, error) {
			return []model.Food{{ID: 5, CityID: 3, Name: "肉夹馍"}}, nil
		},
		listCharacters: func(context.Context, int64) ([]model.Character, error) {
			return []model.Character{{ID: 8, CityID: 3, Name: "唐代长安书生", CharacterType: "history"}}, nil
		},
	}
}

func assertTask3Response(t *testing.T, router *gin.Engine, method, path, requestBody string, wantStatus int, wantCode string) map[string]any {
	t.Helper()
	request := httptest.NewRequest(method, path, bytes.NewBufferString(requestBody))
	if requestBody != "" {
		request.Header.Set("Content-Type", "application/json")
	}
	recorder := httptest.NewRecorder()
	router.ServeHTTP(recorder, request)
	if recorder.Code != wantStatus {
		t.Fatalf("%s %s status = %d, body = %s; want %d", method, path, recorder.Code, recorder.Body.String(), wantStatus)
	}
	var body map[string]any
	if err := json.Unmarshal(recorder.Body.Bytes(), &body); err != nil {
		t.Fatalf("decode response %q: %v", recorder.Body.String(), err)
	}
	if wantCode != "" {
		errorBody, ok := body["error"].(map[string]any)
		if !ok || errorBody["code"] != wantCode {
			t.Fatalf("response = %#v, want error code %q", body, wantCode)
		}
	}
	return body
}

func mustJSON(t *testing.T, value any) string {
	t.Helper()
	payload, err := json.Marshal(value)
	if err != nil {
		t.Fatalf("json.Marshal() error = %v", err)
	}
	return string(payload)
}
