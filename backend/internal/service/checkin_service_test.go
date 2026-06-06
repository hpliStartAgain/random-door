package service

import (
	"context"
	"encoding/json"
	"errors"
	"testing"
	"time"

	aiPkg "github.com/your-org/city-roam/backend/internal/ai"
	"github.com/your-org/city-roam/backend/internal/model"
	"gorm.io/gorm"
)

type fakeCheckinUserFinder struct {
	findByID func(context.Context, int64) (*model.User, error)
}

func (r *fakeCheckinUserFinder) FindByID(ctx context.Context, id int64) (*model.User, error) {
	return r.findByID(ctx, id)
}

type fakeCheckinCityRepo struct {
	findByID      func(context.Context, int64) (*model.City, error)
	listLandmarks func(context.Context, int64) ([]model.Landmark, error)
}

func (r *fakeCheckinCityRepo) FindByID(ctx context.Context, id int64) (*model.City, error) {
	return r.findByID(ctx, id)
}
func (r *fakeCheckinCityRepo) ListLandmarks(ctx context.Context, cityID int64) ([]model.Landmark, error) {
	return r.listLandmarks(ctx, cityID)
}

type fakeCheckinStore struct {
	createAndEvaluate func(context.Context, *model.Checkin) ([]model.Achievement, error)
}

func (r *fakeCheckinStore) CreateAndEvaluate(ctx context.Context, c *model.Checkin) ([]model.Achievement, error) {
	return r.createAndEvaluate(ctx, c)
}

type fakeCheckinVisitRepo struct {
	findByID func(context.Context, int64) (*model.CityVisit, error)
}

func (r *fakeCheckinVisitRepo) FindByID(ctx context.Context, id int64) (*model.CityVisit, error) {
	return r.findByID(ctx, id)
}

type fakeImageGen struct {
	generate func(context.Context, string, string, string) ([]byte, error)
}

func (f *fakeImageGen) Generate(ctx context.Context, selfie, ref, prompt string) ([]byte, error) {
	return f.generate(ctx, selfie, ref, prompt)
}

type fakeImageStorage struct {
	saveBytes func([]byte, string, string) (string, error)
}

func (f *fakeImageStorage) SaveBytes(data []byte, subDir, ext string) (string, error) {
	return f.saveBytes(data, subDir, ext)
}

type fakeCheckinTaskRepo struct {
	create          func(context.Context, *model.AITask) error
	findByIDForUser func(context.Context, int64, int64) (*model.AITask, error)
	queueRetry      func(context.Context, int64, int64) error
	claimNext       func(context.Context, string, int) (*model.AITask, error)
	markSucceeded   func(context.Context, int64, string) error
	markFailed      func(context.Context, int64, string) error
	markRetryable   func(context.Context, int64, string) error
}

func (r *fakeCheckinTaskRepo) Create(ctx context.Context, task *model.AITask) error {
	if r.create == nil {
		return nil
	}
	return r.create(ctx, task)
}

func (r *fakeCheckinTaskRepo) FindByIDForUser(ctx context.Context, id, userID int64) (*model.AITask, error) {
	if r.findByIDForUser == nil {
		return nil, gorm.ErrRecordNotFound
	}
	return r.findByIDForUser(ctx, id, userID)
}

func (r *fakeCheckinTaskRepo) QueueRetry(ctx context.Context, id, userID int64) error {
	if r.queueRetry == nil {
		return nil
	}
	return r.queueRetry(ctx, id, userID)
}

func (r *fakeCheckinTaskRepo) ClaimNext(ctx context.Context, taskType string, maxAttempts int) (*model.AITask, error) {
	if r.claimNext == nil {
		return nil, nil
	}
	return r.claimNext(ctx, taskType, maxAttempts)
}

func (r *fakeCheckinTaskRepo) MarkSucceeded(ctx context.Context, id int64, resultURL string) error {
	if r.markSucceeded == nil {
		return nil
	}
	return r.markSucceeded(ctx, id, resultURL)
}

func (r *fakeCheckinTaskRepo) MarkFailed(ctx context.Context, id int64, message string) error {
	if r.markFailed == nil {
		return nil
	}
	return r.markFailed(ctx, id, message)
}

func (r *fakeCheckinTaskRepo) MarkRetryable(ctx context.Context, id int64, message string) error {
	if r.markRetryable == nil {
		return nil
	}
	return r.markRetryable(ctx, id, message)
}

type fakeAIUsageLimiter struct {
	incrementIfBelow func(context.Context, int64, string, time.Time, int) (int, bool, error)
	calls            int
}

func (l *fakeAIUsageLimiter) IncrementIfBelow(ctx context.Context, userID int64, usageType string, usageDate time.Time, limit int) (int, bool, error) {
	l.calls++
	if l.incrementIfBelow == nil {
		return l.calls, true, nil
	}
	return l.incrementIfBelow(ctx, userID, usageType, usageDate, limit)
}

var lm1ID int64 = 10
var lm1Image = "/static/landmarks/dayanta.png"
var lm1 = model.Landmark{ID: lm1ID, CityID: 1, Name: "大雁塔", ImageURL: &lm1Image}

func okCheckinUser(context.Context, int64) (*model.User, error) { return &model.User{ID: 1}, nil }

func newTestCheckinService(userRepo checkinUserFinder, cityRepo checkinCityRepo, store checkinStore, img imageGenerator, st imageStorage) *CheckinService {
	return &CheckinService{userRepo: userRepo, cityRepo: cityRepo, store: store, imageClient: img, storage: st}
}

func TestCheckinService_GenerateImage_Success(t *testing.T) {
	svc := newTestCheckinService(
		nil,
		&fakeCheckinCityRepo{
			findByID:      func(_ context.Context, _ int64) (*model.City, error) { return testCity, nil },
			listLandmarks: func(_ context.Context, _ int64) ([]model.Landmark, error) { return []model.Landmark{lm1}, nil },
		},
		nil,
		&fakeImageGen{generate: func(_ context.Context, _, _, _ string) ([]byte, error) { return []byte("img"), nil }},
		&fakeImageStorage{saveBytes: func(_ []byte, _, _ string) (string, error) { return "/uploads/generated/x.png", nil }},
	)
	res, err := svc.GenerateImage(context.Background(), 1, 1, lm1ID, "/tmp/selfie.jpg")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res.GeneratedImageURL != "/uploads/generated/x.png" {
		t.Fatalf("unexpected url: %s", res.GeneratedImageURL)
	}
}

func TestCheckinService_BuildImageTaskInputPrefersSceneImage(t *testing.T) {
	svc := newTestCheckinService(
		nil,
		&fakeCheckinCityRepo{
			findByID:      func(_ context.Context, _ int64) (*model.City, error) { return testCity, nil },
			listLandmarks: func(_ context.Context, _ int64) ([]model.Landmark, error) { return []model.Landmark{lm1}, nil },
		},
		nil,
		nil,
		nil,
	)

	input, err := svc.buildImageTaskInput(context.Background(), 1, 1, lm1ID, "/uploads/selfies/me.png", "/uploads/scenes/view.png")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if input.RefImagePath != "/uploads/scenes/view.png" {
		t.Fatalf("expected scene image ref, got %s", input.RefImagePath)
	}
	if input.ScenePath != "/uploads/scenes/view.png" {
		t.Fatalf("expected scene path to be recorded, got %s", input.ScenePath)
	}
}

func TestCheckinService_EnqueueImageCreatesQueuedTask(t *testing.T) {
	var created *model.AITask
	usage := &fakeAIUsageLimiter{}
	taskRepo := &fakeCheckinTaskRepo{
		create: func(_ context.Context, task *model.AITask) error {
			task.ID = 99
			created = task
			return nil
		},
	}
	svc := newTestCheckinService(
		&fakeCheckinUserFinder{findByID: okCheckinUser},
		&fakeCheckinCityRepo{
			findByID:      func(_ context.Context, _ int64) (*model.City, error) { return testCity, nil },
			listLandmarks: func(_ context.Context, _ int64) ([]model.Landmark, error) { return []model.Landmark{lm1}, nil },
		},
		nil,
		nil,
		nil,
	).WithImageTasks(taskRepo, usage, 3, 2)

	res, err := svc.EnqueueImage(context.Background(), 1, 1, lm1ID, "/uploads/selfies/me.png", "/uploads/scenes/view.png")
	if err != nil {
		t.Fatalf("EnqueueImage() error = %v", err)
	}
	if res.TaskID != 99 || res.Status != "queued" {
		t.Fatalf("EnqueueImage() = %#v, want queued task 99", res)
	}
	if usage.calls != 1 {
		t.Fatalf("usage calls = %d, want 1", usage.calls)
	}
	if created == nil || created.Type != checkinImageTaskType || created.Status != "queued" {
		t.Fatalf("created task = %#v, want queued checkin_image task", created)
	}
	var input imageTaskInput
	if err := json.Unmarshal([]byte(created.InputJSON), &input); err != nil {
		t.Fatalf("task input json = %q: %v", created.InputJSON, err)
	}
	if input.RefImagePath != "/uploads/scenes/view.png" || input.SelfiePath != "/uploads/selfies/me.png" {
		t.Fatalf("task input = %#v, want scene ref and selfie path", input)
	}
}

func TestCheckinService_EnqueueImageDoesNotConsumeQuotaWhenLandmarkInvalid(t *testing.T) {
	usage := &fakeAIUsageLimiter{}
	svc := newTestCheckinService(
		&fakeCheckinUserFinder{findByID: okCheckinUser},
		&fakeCheckinCityRepo{
			findByID:      func(_ context.Context, _ int64) (*model.City, error) { return testCity, nil },
			listLandmarks: func(_ context.Context, _ int64) ([]model.Landmark, error) { return []model.Landmark{lm1}, nil },
		},
		nil,
		nil,
		nil,
	).WithImageTasks(&fakeCheckinTaskRepo{}, usage, 3, 2)

	_, err := svc.EnqueueImage(context.Background(), 1, 1, 999, "/uploads/selfies/me.png", "")
	if !errors.Is(err, ErrInvalidParam) {
		t.Fatalf("EnqueueImage() error = %v, want ErrInvalidParam", err)
	}
	if usage.calls != 0 {
		t.Fatalf("usage calls = %d, want 0 for invalid task input", usage.calls)
	}
}

func TestCheckinService_GenerateImage_LandmarkNotFound(t *testing.T) {
	svc := newTestCheckinService(
		nil,
		&fakeCheckinCityRepo{
			findByID:      func(_ context.Context, _ int64) (*model.City, error) { return testCity, nil },
			listLandmarks: func(_ context.Context, _ int64) ([]model.Landmark, error) { return []model.Landmark{lm1}, nil },
		},
		nil,
		&fakeImageGen{},
		&fakeImageStorage{},
	)
	_, err := svc.GenerateImage(context.Background(), 1, 1, 999, "/tmp/selfie.jpg")
	if !errors.Is(err, ErrInvalidParam) {
		t.Fatalf("expected ErrInvalidParam, got %v", err)
	}
}

func TestCheckinService_GenerateImage_CityNotFound(t *testing.T) {
	svc := newTestCheckinService(
		nil,
		&fakeCheckinCityRepo{
			findByID: func(_ context.Context, _ int64) (*model.City, error) { return nil, gorm.ErrRecordNotFound },
		},
		nil,
		&fakeImageGen{},
		&fakeImageStorage{},
	)
	_, err := svc.GenerateImage(context.Background(), 1, 999, lm1ID, "/tmp/selfie.jpg")
	if !errors.Is(err, ErrNotFound) {
		t.Fatalf("expected ErrNotFound, got %v", err)
	}
}

func TestCheckinService_GenerateImage_AITimeout(t *testing.T) {
	svc := newTestCheckinService(
		nil,
		&fakeCheckinCityRepo{
			findByID:      func(_ context.Context, _ int64) (*model.City, error) { return testCity, nil },
			listLandmarks: func(_ context.Context, _ int64) ([]model.Landmark, error) { return []model.Landmark{lm1}, nil },
		},
		nil,
		&fakeImageGen{generate: func(_ context.Context, _, _, _ string) ([]byte, error) { return nil, aiPkg.ErrAITimeout }},
		&fakeImageStorage{},
	)
	_, err := svc.GenerateImage(context.Background(), 1, 1, lm1ID, "/tmp/selfie.jpg")
	if !errors.Is(err, aiPkg.ErrAITimeout) {
		t.Fatalf("expected ErrAITimeout, got %v", err)
	}
}

func TestCheckinService_Create_Success(t *testing.T) {
	svc := newTestCheckinService(
		&fakeCheckinUserFinder{findByID: okCheckinUser},
		&fakeCheckinCityRepo{findByID: func(_ context.Context, _ int64) (*model.City, error) { return testCity, nil }},
		&fakeCheckinStore{createAndEvaluate: func(_ context.Context, c *model.Checkin) ([]model.Achievement, error) {
			c.ID = 42
			return nil, nil
		}},
		nil, nil,
	)
	res, err := svc.Create(context.Background(), CreateCheckinRequest{UserID: 1, CityID: 1})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res.CheckinID != 42 {
		t.Fatalf("expected checkin_id=42, got %d", res.CheckinID)
	}
}

func TestCheckinService_Create_DefaultsCheckinModeFree(t *testing.T) {
	svc := newTestCheckinService(
		&fakeCheckinUserFinder{findByID: okCheckinUser},
		&fakeCheckinCityRepo{findByID: func(_ context.Context, _ int64) (*model.City, error) { return testCity, nil }},
		&fakeCheckinStore{createAndEvaluate: func(_ context.Context, c *model.Checkin) ([]model.Achievement, error) {
			if c.CheckinMode == nil || *c.CheckinMode != "free" {
				t.Fatalf("checkin_mode = %v, want free", c.CheckinMode)
			}
			c.ID = 43
			return nil, nil
		}},
		nil, nil,
	)
	_, err := svc.Create(context.Background(), CreateCheckinRequest{UserID: 1, CityID: 1})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestCheckinService_Create_UsesVisitMode(t *testing.T) {
	visitID := int64(100)
	svc := newTestCheckinService(
		&fakeCheckinUserFinder{findByID: okCheckinUser},
		&fakeCheckinCityRepo{findByID: func(_ context.Context, _ int64) (*model.City, error) { return testCity, nil }},
		&fakeCheckinStore{createAndEvaluate: func(_ context.Context, c *model.Checkin) ([]model.Achievement, error) {
			if c.CheckinMode == nil || *c.CheckinMode != "game" {
				t.Fatalf("checkin_mode = %v, want game", c.CheckinMode)
			}
			c.ID = 44
			return nil, nil
		}},
		nil, nil,
	).WithVisitRepo(&fakeCheckinVisitRepo{
		findByID: func(_ context.Context, _ int64) (*model.CityVisit, error) {
			return &model.CityVisit{ID: visitID, UserID: 1, CityID: 1, VisitMode: "game"}, nil
		},
	})

	_, err := svc.Create(context.Background(), CreateCheckinRequest{UserID: 1, CityID: 1, VisitID: &visitID})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestCheckinService_Create_RejectsMismatchedVisit(t *testing.T) {
	visitID := int64(100)
	svc := newTestCheckinService(
		&fakeCheckinUserFinder{findByID: okCheckinUser},
		&fakeCheckinCityRepo{findByID: func(_ context.Context, _ int64) (*model.City, error) { return testCity, nil }},
		&fakeCheckinStore{},
		nil, nil,
	).WithVisitRepo(&fakeCheckinVisitRepo{
		findByID: func(_ context.Context, _ int64) (*model.CityVisit, error) {
			return &model.CityVisit{ID: visitID, UserID: 2, CityID: 1, VisitMode: "game"}, nil
		},
	})

	_, err := svc.Create(context.Background(), CreateCheckinRequest{UserID: 1, CityID: 1, VisitID: &visitID})
	if !errors.Is(err, ErrInvalidParam) {
		t.Fatalf("expected ErrInvalidParam, got %v", err)
	}
}

func TestCheckinService_Create_StoreError(t *testing.T) {
	svc := newTestCheckinService(
		&fakeCheckinUserFinder{findByID: okCheckinUser},
		&fakeCheckinCityRepo{findByID: func(_ context.Context, _ int64) (*model.City, error) { return testCity, nil }},
		&fakeCheckinStore{createAndEvaluate: func(_ context.Context, _ *model.Checkin) ([]model.Achievement, error) {
			return nil, errors.New("db error")
		}},
		nil, nil,
	)
	_, err := svc.Create(context.Background(), CreateCheckinRequest{UserID: 1, CityID: 1})
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

func TestCheckinService_Create_UserNotFound(t *testing.T) {
	svc := newTestCheckinService(
		&fakeCheckinUserFinder{findByID: func(_ context.Context, _ int64) (*model.User, error) { return nil, gorm.ErrRecordNotFound }},
		&fakeCheckinCityRepo{},
		&fakeCheckinStore{},
		nil, nil,
	)
	_, err := svc.Create(context.Background(), CreateCheckinRequest{UserID: 999, CityID: 1})
	if !errors.Is(err, ErrNotFound) || ClientMessage(err) != "user not found" {
		t.Fatalf("expected user not found, got %v", err)
	}
}

func TestCheckinService_Create_CityNotFound(t *testing.T) {
	svc := newTestCheckinService(
		&fakeCheckinUserFinder{findByID: okCheckinUser},
		&fakeCheckinCityRepo{findByID: func(_ context.Context, _ int64) (*model.City, error) { return nil, gorm.ErrRecordNotFound }},
		&fakeCheckinStore{},
		nil, nil,
	)
	_, err := svc.Create(context.Background(), CreateCheckinRequest{UserID: 1, CityID: 999})
	if !errors.Is(err, ErrNotFound) || ClientMessage(err) != "city not found" {
		t.Fatalf("expected city not found, got %v", err)
	}
}
