package service

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"time"

	"github.com/your-org/city-roam/backend/internal/ai"
	"github.com/your-org/city-roam/backend/internal/model"
	"gorm.io/gorm"
)

type checkinUserFinder interface {
	FindByID(ctx context.Context, id int64) (*model.User, error)
}

type checkinCityRepo interface {
	FindByID(ctx context.Context, id int64) (*model.City, error)
	ListLandmarks(ctx context.Context, cityID int64) ([]model.Landmark, error)
}

type checkinVisitRepo interface {
	FindByID(ctx context.Context, id int64) (*model.CityVisit, error)
}

// checkinStore writes a checkin and evaluates achievements in a single transaction.
type checkinStore interface {
	CreateAndEvaluate(ctx context.Context, checkin *model.Checkin) ([]model.Achievement, error)
}

type imageGenerator interface {
	Generate(ctx context.Context, selfiePath, refImagePath, prompt string) ([]byte, error)
}

type imageStorage interface {
	SaveBytes(data []byte, subDir, ext string) (string, error)
}

type checkinTaskRepo interface {
	Create(ctx context.Context, task *model.AITask) error
	FindByIDForUser(ctx context.Context, id, userID int64) (*model.AITask, error)
	QueueRetry(ctx context.Context, id, userID int64) error
	ClaimNext(ctx context.Context, taskType string, maxAttempts int) (*model.AITask, error)
	MarkSucceeded(ctx context.Context, id int64, resultURL string) error
	MarkFailed(ctx context.Context, id int64, message string) error
	MarkRetryable(ctx context.Context, id int64, message string) error
}

type CheckinService struct {
	userRepo        checkinUserFinder
	cityRepo        checkinCityRepo
	visitRepo       checkinVisitRepo
	store           checkinStore
	imageClient     imageGenerator
	storage         imageStorage
	taskRepo        checkinTaskRepo
	usageLimiter    aiUsageLimiter
	imageDailyLimit int
	maxTaskAttempts int
}

func NewCheckinService(
	userRepo checkinUserFinder,
	cityRepo checkinCityRepo,
	store checkinStore,
	imageClient imageGenerator,
	storage imageStorage,
) *CheckinService {
	return &CheckinService{
		userRepo: userRepo, cityRepo: cityRepo, store: store,
		imageClient: imageClient, storage: storage,
	}
}

func (s *CheckinService) WithImageTasks(taskRepo checkinTaskRepo, limiter aiUsageLimiter, imageDailyLimit, maxTaskAttempts int) *CheckinService {
	s.taskRepo = taskRepo
	s.usageLimiter = limiter
	s.imageDailyLimit = imageDailyLimit
	if maxTaskAttempts <= 0 {
		maxTaskAttempts = 3
	}
	s.maxTaskAttempts = maxTaskAttempts
	return s
}

func (s *CheckinService) WithVisitRepo(visitRepo checkinVisitRepo) *CheckinService {
	s.visitRepo = visitRepo
	return s
}

// GenerateImageResult is the response from image generation.
type GenerateImageResult struct {
	Status            string `json:"status"`
	GeneratedImageURL string `json:"generated_image_url"`
}

type EnqueueImageResult struct {
	TaskID int64  `json:"task_id"`
	Status string `json:"status"`
}

type ImageTaskStatus struct {
	TaskID    int64   `json:"task_id"`
	Status    string  `json:"status"`
	ResultURL *string `json:"result_url"`
	Error     *string `json:"error"`
	Attempts  int     `json:"attempts"`
	CreatedAt string  `json:"created_at"`
	UpdatedAt string  `json:"updated_at"`
}

type imageTaskInput struct {
	UserID       int64  `json:"user_id"`
	CityID       int64  `json:"city_id"`
	LandmarkID   int64  `json:"landmark_id"`
	SelfiePath   string `json:"selfie_path"`
	ScenePath    string `json:"scene_path,omitempty"`
	CityName     string `json:"city_name"`
	LandmarkName string `json:"landmark_name"`
	RefImagePath string `json:"ref_image_path"`
	Prompt       string `json:"prompt"`
}

const checkinImageTaskType = "checkin_image"

func (s *CheckinService) EnqueueImage(ctx context.Context, userID, cityID, landmarkID int64, selfiePath, scenePath string) (*EnqueueImageResult, error) {
	if s.taskRepo == nil {
		return nil, fmt.Errorf("image task repository is not configured")
	}
	if _, err := s.userRepo.FindByID(ctx, userID); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, notFound("user not found")
		}
		return nil, fmt.Errorf("find user: %w", err)
	}

	input, err := s.buildImageTaskInput(ctx, userID, cityID, landmarkID, selfiePath, scenePath)
	if err != nil {
		return nil, err
	}
	if s.usageLimiter != nil && s.imageDailyLimit > 0 {
		_, allowed, err := s.usageLimiter.IncrementIfBelow(ctx, userID, "image", time.Now(), s.imageDailyLimit)
		if err != nil {
			return nil, fmt.Errorf("increment image usage: %w", err)
		}
		if !allowed {
			return nil, quotaExceeded("daily image quota exceeded")
		}
	}

	payload, err := json.Marshal(input)
	if err != nil {
		return nil, fmt.Errorf("marshal image task input: %w", err)
	}

	task := &model.AITask{
		UserID:    userID,
		Type:      checkinImageTaskType,
		Status:    "queued",
		InputJSON: string(payload),
	}
	if err := s.taskRepo.Create(ctx, task); err != nil {
		return nil, fmt.Errorf("create image task: %w", err)
	}
	slog.Info("image task queued", "task_id", task.ID, "user_id", userID, "city_id", cityID)
	return &EnqueueImageResult{TaskID: task.ID, Status: task.Status}, nil
}

func (s *CheckinService) GetImageTask(ctx context.Context, userID, taskID int64) (*ImageTaskStatus, error) {
	if taskID <= 0 || userID <= 0 {
		return nil, invalidParam("task_id and user_id must be positive integers")
	}
	task, err := s.taskRepo.FindByIDForUser(ctx, taskID, userID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, notFound("image task not found")
		}
		return nil, fmt.Errorf("find image task: %w", err)
	}
	return imageTaskStatus(task), nil
}

func (s *CheckinService) RetryImageTask(ctx context.Context, userID, taskID int64) (*EnqueueImageResult, error) {
	if taskID <= 0 || userID <= 0 {
		return nil, invalidParam("task_id and user_id must be positive integers")
	}
	if err := s.taskRepo.QueueRetry(ctx, taskID, userID); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, notFound("image task not found or not retryable")
		}
		return nil, fmt.Errorf("queue image task retry: %w", err)
	}
	task, err := s.taskRepo.FindByIDForUser(ctx, taskID, userID)
	if err != nil {
		return nil, fmt.Errorf("reload image task: %w", err)
	}
	return &EnqueueImageResult{TaskID: task.ID, Status: task.Status}, nil
}

func imageTaskStatus(task *model.AITask) *ImageTaskStatus {
	return &ImageTaskStatus{
		TaskID:    task.ID,
		Status:    task.Status,
		ResultURL: task.ResultURL,
		Error:     task.Error,
		Attempts:  task.Attempts,
		CreatedAt: task.CreatedAt.Format(time.RFC3339),
		UpdatedAt: task.UpdatedAt.Format(time.RFC3339),
	}
}

// GenerateImage processes selfie + landmark reference to create a cyber check-in photo.
func (s *CheckinService) GenerateImage(ctx context.Context, userID, cityID, landmarkID int64, selfiePath string) (*GenerateImageResult, error) {
	input, err := s.buildImageTaskInput(ctx, userID, cityID, landmarkID, selfiePath, "")
	if err != nil {
		return nil, err
	}

	generatedURL, err := s.generateImageFromInput(ctx, input)
	if err != nil {
		return nil, err
	}

	slog.Info("image generated", "user_id", userID, "url", generatedURL)
	return &GenerateImageResult{
		Status:            "success",
		GeneratedImageURL: generatedURL,
	}, nil
}

func (s *CheckinService) buildImageTaskInput(ctx context.Context, userID, cityID, landmarkID int64, selfiePath, scenePath string) (imageTaskInput, error) {
	city, err := s.cityRepo.FindByID(ctx, cityID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return imageTaskInput{}, notFound("city not found")
		}
		return imageTaskInput{}, fmt.Errorf("find city: %w", err)
	}

	landmarks, _ := s.cityRepo.ListLandmarks(ctx, cityID)
	var landmarkName string
	var refImagePath string
	for _, l := range landmarks {
		if l.ID == landmarkID {
			landmarkName = l.Name
			if l.ImageURL != nil {
				refImagePath = *l.ImageURL
			}
			break
		}
	}
	if landmarkName == "" {
		return imageTaskInput{}, invalidParam("landmark not found")
	}
	if scenePath != "" {
		refImagePath = scenePath
	}

	prompt := ai.BuildImagePrompt(city.Name, landmarkName)
	return imageTaskInput{
		UserID:       userID,
		CityID:       cityID,
		LandmarkID:   landmarkID,
		SelfiePath:   selfiePath,
		ScenePath:    scenePath,
		CityName:     city.Name,
		LandmarkName: landmarkName,
		RefImagePath: refImagePath,
		Prompt:       prompt,
	}, nil
}

func (s *CheckinService) generateImageFromInput(ctx context.Context, input imageTaskInput) (string, error) {
	imgBytes, err := s.imageClient.Generate(ctx, input.SelfiePath, input.RefImagePath, input.Prompt)
	if err != nil {
		slog.Error("image generation failed", "error", err, "user_id", input.UserID)
		return "", err
	}

	generatedURL, err := s.storage.SaveBytes(imgBytes, "generated", ".png")
	if err != nil {
		return "", fmt.Errorf("save generated image: %w", err)
	}
	return generatedURL, nil
}

func (s *CheckinService) StartImageWorker(ctx context.Context, interval time.Duration, concurrency int) {
	if s.taskRepo == nil || s.imageClient == nil || s.storage == nil {
		slog.Warn("image worker disabled: dependencies are not configured")
		return
	}
	if interval <= 0 {
		interval = 2 * time.Second
	}
	if concurrency <= 0 {
		concurrency = 1
	}
	for i := 0; i < concurrency; i++ {
		workerID := i + 1
		go s.imageWorkerLoop(ctx, workerID, interval)
	}
}

func (s *CheckinService) imageWorkerLoop(ctx context.Context, workerID int, interval time.Duration) {
	ticker := time.NewTicker(interval)
	defer ticker.Stop()
	for {
		select {
		case <-ctx.Done():
			return
		default:
			s.processOneImageTask(ctx, workerID)
		}

		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
		}
	}
}

func (s *CheckinService) processOneImageTask(ctx context.Context, workerID int) {
	task, err := s.taskRepo.ClaimNext(ctx, checkinImageTaskType, s.maxTaskAttempts)
	if err != nil {
		slog.Error("claim image task failed", "error", err, "worker_id", workerID)
		return
	}
	if task == nil {
		return
	}

	var input imageTaskInput
	if err := json.Unmarshal([]byte(task.InputJSON), &input); err != nil {
		_ = s.taskRepo.MarkFailed(ctx, task.ID, "invalid task input")
		slog.Error("invalid image task input", "task_id", task.ID, "error", err)
		return
	}

	resultURL, err := s.generateImageFromInput(ctx, input)
	if err == nil {
		if markErr := s.taskRepo.MarkSucceeded(ctx, task.ID, resultURL); markErr != nil {
			slog.Error("mark image task succeeded failed", "task_id", task.ID, "error", markErr)
		}
		slog.Info("image task completed", "task_id", task.ID, "user_id", task.UserID, "url", resultURL)
		return
	}

	message := "image generation failed"
	if errors.Is(err, ai.ErrAITimeout) {
		message = "image generation timeout"
	}
	if task.Attempts >= s.maxTaskAttempts {
		if markErr := s.taskRepo.MarkFailed(ctx, task.ID, message); markErr != nil {
			slog.Error("mark image task failed failed", "task_id", task.ID, "error", markErr)
		}
	} else {
		if markErr := s.taskRepo.MarkRetryable(ctx, task.ID, message); markErr != nil {
			slog.Error("mark image task retryable failed", "task_id", task.ID, "error", markErr)
		}
	}
}

// CreateCheckinRequest is the input for creating a check-in.
type CreateCheckinRequest struct {
	UserID            int64   `json:"user_id"`
	CityID            int64   `json:"city_id"`
	LandmarkID        *int64  `json:"landmark_id,omitempty"`
	VisitID           *int64  `json:"visit_id,omitempty"`
	GeneratedImageURL *string `json:"generated_image_url,omitempty"`
}

// CreateCheckinResult is the response from creating a check-in.
type CreateCheckinResult struct {
	CheckinID            int64              `json:"checkin_id"`
	UnlockedAchievements []AchievementBrief `json:"unlocked_achievements"`
}

type AchievementBrief struct {
	Code        string  `json:"code"`
	Name        string  `json:"name"`
	Description *string `json:"description,omitempty"`
}

// Create records a check-in and evaluates achievements within a single transaction.
func (s *CheckinService) Create(ctx context.Context, req CreateCheckinRequest) (*CreateCheckinResult, error) {
	// checkins has no FK constraints, so verify references explicitly.
	if _, err := s.userRepo.FindByID(ctx, req.UserID); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, notFound("user not found")
		}
		return nil, fmt.Errorf("find user: %w", err)
	}
	if _, err := s.cityRepo.FindByID(ctx, req.CityID); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, notFound("city not found")
		}
		return nil, fmt.Errorf("find city: %w", err)
	}

	checkinMode, err := s.checkinMode(ctx, req)
	if err != nil {
		return nil, err
	}

	checkin := &model.Checkin{
		UserID:            req.UserID,
		CityID:            req.CityID,
		LandmarkID:        req.LandmarkID,
		VisitID:           req.VisitID,
		GeneratedImageURL: req.GeneratedImageURL,
		CheckinMode:       &checkinMode,
	}

	newAchs, err := s.store.CreateAndEvaluate(ctx, checkin)
	if err != nil {
		return nil, fmt.Errorf("create checkin: %w", err)
	}

	achBriefs := briefAchievements(newAchs)

	slog.Info("checkin created", "user_id", req.UserID, "city_id", req.CityID,
		"checkin_id", checkin.ID, "new_achievements", len(newAchs))

	return &CreateCheckinResult{
		CheckinID:            checkin.ID,
		UnlockedAchievements: achBriefs,
	}, nil
}

func (s *CheckinService) checkinMode(ctx context.Context, req CreateCheckinRequest) (string, error) {
	if req.VisitID == nil {
		return "free", nil
	}
	if s.visitRepo == nil {
		return "free", nil
	}
	visit, err := s.visitRepo.FindByID(ctx, *req.VisitID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return "", notFound("visit not found")
		}
		return "", fmt.Errorf("find visit: %w", err)
	}
	if visit.UserID != req.UserID || visit.CityID != req.CityID {
		return "", invalidParam("visit does not match user or city")
	}
	if visit.VisitMode == "game" {
		return "game", nil
	}
	return "free", nil
}
