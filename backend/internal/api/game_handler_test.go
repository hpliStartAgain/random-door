package api

import (
	"context"
	"errors"
	"net/http"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/your-org/city-roam/backend/internal/model"
	"github.com/your-org/city-roam/backend/internal/service"
	"gorm.io/gorm"
)

type gameUserFinder struct {
	findByID func(context.Context, int64) (*model.User, error)
}

func (r *gameUserFinder) FindByID(ctx context.Context, id int64) (*model.User, error) {
	return r.findByID(ctx, id)
}

type gameCityRepo struct {
	listAll  func(context.Context) ([]model.City, error)
	findByID func(context.Context, int64) (*model.City, error)
}

func (r *gameCityRepo) ListAll(ctx context.Context) ([]model.City, error) { return r.listAll(ctx) }
func (r *gameCityRepo) FindByID(ctx context.Context, id int64) (*model.City, error) {
	return r.findByID(ctx, id)
}

type gameVisitReader struct{}

func (gameVisitReader) ListVisitedCityIDs(context.Context, int64) ([]int64, error) {
	return nil, nil
}

type gameStore struct {
	create func(context.Context, *model.DiceRoll, *model.CityVisit) error
}

func (r *gameStore) CreateRollWithVisit(ctx context.Context, roll *model.DiceRoll, visit *model.CityVisit) error {
	return r.create(ctx, roll, visit)
}

func newGameRouter(svc *service.GameService) *gin.Engine {
	gin.SetMode(gin.TestMode)
	h := NewGameHandler(svc)
	r := gin.New()
	r.POST("/api/game/init", h.Init)
	r.POST("/api/game/roll", h.Roll)
	return r
}

var gameRouterCities = []model.City{
	{ID: 1, Name: "北京", Province: "北京", Lat: 39.9042, Lng: 116.4074},
	{ID: 3, Name: "西安", Province: "陕西", Lat: 34.3416, Lng: 108.9398},
}

func TestGameHandlersSuccess(t *testing.T) {
	cityRepo := &gameCityRepo{
		listAll:  func(context.Context) ([]model.City, error) { return gameRouterCities, nil },
		findByID: func(context.Context, int64) (*model.City, error) { return &model.City{ID: 1}, nil },
	}
	store := &gameStore{create: func(_ context.Context, roll *model.DiceRoll, visit *model.CityVisit) error {
		roll.ID = 2001
		visit.DiceRollID = &roll.ID
		visit.ID = 1002
		return nil
	}}
	svc := service.NewGameService(
		&gameUserFinder{findByID: func(context.Context, int64) (*model.User, error) { return &model.User{ID: 1}, nil }},
		cityRepo, gameVisitReader{}, store,
	).WithRandSource(&handlerSeqSource{values: []int{5, 4}})

	router := newGameRouter(svc)

	init := assertTask3Response(t, router, http.MethodPost, "/api/game/init", `{"user_id":1}`, http.StatusOK, "")
	near, ok := init["nearest_city"].(map[string]any)
	if !ok || near["id"] != float64(1) || near["name"] != "北京" {
		t.Fatalf("init response = %#v, want nearest Beijing", init)
	}

	roll := assertTask3Response(t, router, http.MethodPost, "/api/game/roll", `{"user_id":1,"from_city_id":1,"lat":39.9042,"lng":116.4074}`, http.StatusOK, "")
	if roll["visit_id"] != float64(1002) || roll["dice_roll_id"] != float64(2001) ||
		roll["direction"] != "西南" || roll["distance_km"] != float64(800) {
		t.Fatalf("roll response = %#v, want contract fields", roll)
	}
	target, ok := roll["target_city"].(map[string]any)
	if !ok || target["id"] == float64(1) {
		t.Fatalf("roll target = %#v, want city other than origin", roll)
	}
}

func TestGameHandlersErrors(t *testing.T) {
	okUserRepo := &gameUserFinder{findByID: func(context.Context, int64) (*model.User, error) { return &model.User{ID: 1}, nil }}
	okCityRepo := &gameCityRepo{
		listAll:  func(context.Context) ([]model.City, error) { return gameRouterCities, nil },
		findByID: func(context.Context, int64) (*model.City, error) { return &model.City{ID: 1}, nil },
	}

	tests := []struct {
		name        string
		svc         *service.GameService
		path        string
		body        string
		wantStatus  int
		wantCode    string
		wantMessage string
	}{
		{
			name: "init requires user", svc: service.NewGameService(okUserRepo, okCityRepo, gameVisitReader{}, nil),
			path: "/api/game/init", body: `{}`,
			wantStatus: http.StatusBadRequest, wantCode: "INVALID_PARAM", wantMessage: "user_id is required",
		},
		{
			name: "init invalid coords", svc: service.NewGameService(okUserRepo, okCityRepo, gameVisitReader{}, nil),
			path: "/api/game/init", body: `{"user_id":1,"lat":200,"lng":10}`,
			wantStatus: http.StatusBadRequest, wantCode: "INVALID_PARAM", wantMessage: "lat must be between -90 and 90",
		},
		{
			name: "init missing user", svc: service.NewGameService(
				&gameUserFinder{findByID: func(context.Context, int64) (*model.User, error) { return nil, gorm.ErrRecordNotFound }},
				okCityRepo, gameVisitReader{}, nil),
			path: "/api/game/init", body: `{"user_id":1}`,
			wantStatus: http.StatusNotFound, wantCode: "NOT_FOUND", wantMessage: "user not found",
		},
		{
			name: "roll requires fields", svc: service.NewGameService(okUserRepo, okCityRepo, gameVisitReader{}, nil),
			path: "/api/game/roll", body: `{"user_id":1}`,
			wantStatus: http.StatusBadRequest, wantCode: "INVALID_PARAM", wantMessage: "user_id, from_city_id, lat, lng are required",
		},
		{
			name: "roll missing from city", svc: service.NewGameService(okUserRepo,
				&gameCityRepo{findByID: func(context.Context, int64) (*model.City, error) { return nil, gorm.ErrRecordNotFound }},
				gameVisitReader{}, nil),
			path: "/api/game/roll", body: `{"user_id":1,"from_city_id":99,"lat":39.9,"lng":116.4}`,
			wantStatus: http.StatusNotFound, wantCode: "NOT_FOUND", wantMessage: "from city not found",
		},
		{
			name: "roll internal error masked", svc: service.NewGameService(okUserRepo, okCityRepo, gameVisitReader{},
				&gameStore{create: func(context.Context, *model.DiceRoll, *model.CityVisit) error {
					return errors.New("mysql password secret")
				}}).WithRandSource(&handlerSeqSource{values: []int{0, 0}}),
			path: "/api/game/roll", body: `{"user_id":1,"from_city_id":1,"lat":39.9,"lng":116.4}`,
			wantStatus: http.StatusInternalServerError, wantCode: "INTERNAL_ERROR", wantMessage: "internal server error",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			router := newGameRouter(tc.svc)
			body := assertTask3Response(t, router, http.MethodPost, tc.path, tc.body, tc.wantStatus, tc.wantCode)
			if body["error"].(map[string]any)["message"] != tc.wantMessage {
				t.Fatalf("response = %#v, want message %q", body, tc.wantMessage)
			}
			if strings.Contains(mustJSON(t, body), "mysql password secret") {
				t.Fatalf("response leaked internal error: %#v", body)
			}
		})
	}
}

type handlerSeqSource struct {
	values []int
	idx    int
}

func (s *handlerSeqSource) Intn(int) int {
	v := s.values[s.idx%len(s.values)]
	s.idx++
	return v
}
