package service

import (
	"context"
	"fmt"

	"github.com/your-org/city-roam/backend/internal/achievement"
	"github.com/your-org/city-roam/backend/internal/model"
	"github.com/your-org/city-roam/backend/internal/repository"
	"gorm.io/gorm"
)

type AchievementService struct {
	db      *gorm.DB
	achRepo *repository.AchievementRepo
}

func NewAchievementService(db *gorm.DB, achRepo *repository.AchievementRepo) *AchievementService {
	return &AchievementService{db: db, achRepo: achRepo}
}

// ensureUserExists returns a classified notFound error when the user is absent.
func (s *AchievementService) ensureUserExists(ctx context.Context, userID int64) error {
	var count int64
	if err := s.db.WithContext(ctx).Model(&model.User{}).
		Where("id = ?", userID).Count(&count).Error; err != nil {
		return fmt.Errorf("count user: %w", err)
	}
	if count == 0 {
		return notFound("user not found")
	}
	return nil
}

// Evaluate runs achievement evaluation for a user.
func (s *AchievementService) Evaluate(ctx context.Context, userID int64) ([]model.Achievement, error) {
	return achievement.Evaluate(ctx, userID, achievement.Repos{DB: s.db})
}

// WallResult is the response for the achievement wall endpoint.
type WallResult struct {
	Unlocked []UnlockedAchievement `json:"unlocked"`
	Locked   []LockedAchievement   `json:"locked"`
	Progress []ProgressItem        `json:"progress"`
}

type UnlockedAchievement struct {
	Code        string  `json:"code"`
	Name        string  `json:"name"`
	Description *string `json:"description,omitempty"`
	BadgeURL    *string `json:"badge_url,omitempty"`
	UnlockedAt  string  `json:"unlocked_at"`
}

type LockedAchievement struct {
	Code        string  `json:"code"`
	Name        string  `json:"name"`
	Description *string `json:"description,omitempty"`
	BadgeURL    *string `json:"badge_url,omitempty"`
}

type ProgressItem struct {
	Code    string `json:"code"`
	Current int    `json:"current"`
	Target  int    `json:"target"`
}

// Wall returns the full achievement wall for a user.
func (s *AchievementService) Wall(ctx context.Context, userID int64) (*WallResult, error) {
	if userID <= 0 {
		return nil, invalidParam("user_id must be a positive integer")
	}
	if err := s.ensureUserExists(ctx, userID); err != nil {
		return nil, err
	}

	allAchs, err := s.achRepo.ListAll(ctx)
	if err != nil {
		return nil, fmt.Errorf("list achievements: %w", err)
	}

	userAchs, err := s.achRepo.ListUserAchievements(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("list user achievements: %w", err)
	}

	// Build unlocked set with timestamps
	unlockedMap := make(map[int64]model.UserAchievement)
	for _, ua := range userAchs {
		unlockedMap[ua.AchievementID] = ua
	}

	// Build user stats for progress calculation
	stats, err := buildProgressStats(ctx, userID, s.db)
	if err != nil {
		return nil, fmt.Errorf("build progress stats: %w", err)
	}

	var unlocked []UnlockedAchievement
	var locked []LockedAchievement
	var progress []ProgressItem

	for _, ach := range allAchs {
		if ua, ok := unlockedMap[ach.ID]; ok {
			unlocked = append(unlocked, UnlockedAchievement{
				Code: ach.Code, Name: ach.Name, Description: ach.Description,
				BadgeURL: ach.BadgeURL, UnlockedAt: ua.UnlockedAt.Format("2006-01-02T15:04:05+08:00"),
			})
		} else {
			locked = append(locked, LockedAchievement{
				Code: ach.Code, Name: ach.Name, Description: ach.Description,
				BadgeURL: ach.BadgeURL,
			})

			// Calculate progress for quantifiable achievements
			if p := calculateProgress(ach, stats); p != nil {
				progress = append(progress, *p)
			}
		}
	}

	return &WallResult{
		Unlocked: unlocked,
		Locked:   locked,
		Progress: progress,
	}, nil
}

type progressStats struct {
	checkinCount   int
	tagCityCounts  map[string]int
	gameVisitCount int
}

func buildProgressStats(ctx context.Context, userID int64, db *gorm.DB) (progressStats, error) {
	s := progressStats{tagCityCounts: make(map[string]int)}

	var checkinCount int64
	if err := db.WithContext(ctx).Model(&model.Checkin{}).
		Where("user_id = ?", userID).Count(&checkinCount).Error; err != nil {
		return s, fmt.Errorf("count checkins: %w", err)
	}
	s.checkinCount = int(checkinCount)

	var checkinCityIDs []int64
	if err := db.WithContext(ctx).Model(&model.Checkin{}).
		Where("user_id = ?", userID).Distinct("city_id").Pluck("city_id", &checkinCityIDs).Error; err != nil {
		return s, fmt.Errorf("pluck checkin cities: %w", err)
	}

	if len(checkinCityIDs) > 0 {
		var tags []model.CityTag
		if err := db.WithContext(ctx).Where("city_id IN ?", checkinCityIDs).Find(&tags).Error; err != nil {
			return s, fmt.Errorf("find city tags: %w", err)
		}
		tagCities := make(map[string]map[int64]bool)
		for _, t := range tags {
			if tagCities[t.Tag] == nil {
				tagCities[t.Tag] = make(map[int64]bool)
			}
			tagCities[t.Tag][t.CityID] = true
		}
		for tag, cities := range tagCities {
			s.tagCityCounts[tag] = len(cities)
		}
	}

	var gameVisitCount int64
	if err := db.WithContext(ctx).Model(&model.CityVisit{}).
		Where("user_id = ? AND visit_mode = ?", userID, "game").
		Distinct("city_id").Count(&gameVisitCount).Error; err != nil {
		return s, fmt.Errorf("count game visits: %w", err)
	}
	s.gameVisitCount = int(gameVisitCount)

	return s, nil
}

func calculateProgress(ach model.Achievement, stats progressStats) *ProgressItem {
	switch ach.RuleType {
	case "checkin_count":
		target := parseIntOrZero(ach.RuleValue)
		if target > 0 {
			return &ProgressItem{Code: ach.Code, Current: stats.checkinCount, Target: target}
		}
	case "tag_count":
		parts := splitTagCount(ach.RuleValue)
		if parts != nil {
			current := stats.tagCityCounts[parts.tag]
			return &ProgressItem{Code: ach.Code, Current: current, Target: parts.count}
		}
	case "game_visit_count":
		target := parseIntOrZero(ach.RuleValue)
		if target > 0 {
			return &ProgressItem{Code: ach.Code, Current: stats.gameVisitCount, Target: target}
		}
	}
	return nil
}

type tagCountParts struct {
	tag   string
	count int
}

func splitTagCount(val string) *tagCountParts {
	for i := len(val) - 1; i >= 0; i-- {
		if val[i] == ':' {
			n := parseIntOrZero(val[i+1:])
			if n > 0 {
				return &tagCountParts{tag: val[:i], count: n}
			}
			break
		}
	}
	return nil
}

func parseIntOrZero(s string) int {
	n := 0
	for _, c := range s {
		if c >= '0' && c <= '9' {
			n = n*10 + int(c-'0')
		} else {
			return 0
		}
	}
	return n
}
