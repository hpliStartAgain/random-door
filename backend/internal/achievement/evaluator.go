package achievement

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/your-org/city-roam/backend/internal/model"
	"gorm.io/gorm"
)

// Repos holds the repository dependencies for achievement evaluation.
type Repos struct {
	DB *gorm.DB
}

// Evaluate checks all achievements and returns newly unlocked ones.
func Evaluate(ctx context.Context, userID int64, repos Repos) ([]model.Achievement, error) {
	stats, err := buildUserStats(ctx, userID, repos.DB)
	if err != nil {
		return nil, fmt.Errorf("build user stats: %w", err)
	}

	// Load all achievement definitions
	var allAchievements []model.Achievement
	if err := repos.DB.WithContext(ctx).Find(&allAchievements).Error; err != nil {
		return nil, fmt.Errorf("list achievements: %w", err)
	}

	// Load already unlocked
	var unlocked []model.UserAchievement
	repos.DB.WithContext(ctx).Where("user_id = ?", userID).Find(&unlocked)
	unlockedSet := make(map[int64]bool)
	for _, ua := range unlocked {
		unlockedSet[ua.AchievementID] = true
	}

	// Evaluate each
	var newlyUnlocked []model.Achievement
	for _, ach := range allAchievements {
		if unlockedSet[ach.ID] {
			continue
		}
		if Match(ach.RuleType, ach.RuleValue, stats) {
			// Write user_achievement (UNIQUE constraint prevents duplicates)
			ua := model.UserAchievement{
				UserID:        userID,
				AchievementID: ach.ID,
				UnlockedAt:    time.Now(),
			}
			if err := repos.DB.WithContext(ctx).Create(&ua).Error; err != nil {
				slog.Warn("create user_achievement failed (possible duplicate)", "user_id", userID, "achievement", ach.Code, "error", err)
				continue
			}
			newlyUnlocked = append(newlyUnlocked, ach)
			slog.Info("achievement unlocked", "user_id", userID, "code", ach.Code)
		}
	}

	return newlyUnlocked, nil
}

func buildUserStats(ctx context.Context, userID int64, db *gorm.DB) (UserStats, error) {
	s := UserStats{
		CheckinCityTags: make(map[string]int),
		CheckinTagAny:   make(map[string]bool),
		MaxSameDirRun:   make(map[string]int),
	}

	// Checkin count
	var checkinCount int64
	db.WithContext(ctx).Model(&model.Checkin{}).Where("user_id = ?", userID).Count(&checkinCount)
	s.CheckinCount = int(checkinCount)

	// Checkin city tags: find distinct cities user checked in, then their tags
	var checkinCityIDs []int64
	db.WithContext(ctx).Model(&model.Checkin{}).
		Where("user_id = ?", userID).
		Distinct("city_id").
		Pluck("city_id", &checkinCityIDs)

	if len(checkinCityIDs) > 0 {
		var tags []model.CityTag
		db.WithContext(ctx).Where("city_id IN ?", checkinCityIDs).Find(&tags)

		// Count distinct cities per tag
		tagCities := make(map[string]map[int64]bool)
		for _, t := range tags {
			if tagCities[t.Tag] == nil {
				tagCities[t.Tag] = make(map[int64]bool)
			}
			tagCities[t.Tag][t.CityID] = true
		}
		for tag, cities := range tagCities {
			s.CheckinCityTags[tag] = len(cities)
			s.CheckinTagAny[tag] = true
		}
	}

	// Game visit count (distinct cities in game mode)
	var gameVisitCount int64
	db.WithContext(ctx).Model(&model.CityVisit{}).
		Where("user_id = ? AND visit_mode = ?", userID, "game").
		Distinct("city_id").
		Count(&gameVisitCount)
	s.GameVisitCount = int(gameVisitCount)

	// Dice rolls stats
	var rolls []model.DiceRoll
	db.WithContext(ctx).Where("user_id = ?", userID).Order("created_at ASC").Find(&rolls)

	// Max distance
	for _, r := range rolls {
		if r.DistanceKm > s.MaxDiceDistance {
			s.MaxDiceDistance = r.DistanceKm
		}
	}

	// Max consecutive same direction
	if len(rolls) > 0 {
		currentDir := rolls[0].Direction
		currentRun := 1
		s.MaxSameDirRun[currentDir] = 1

		for i := 1; i < len(rolls); i++ {
			if rolls[i].Direction == currentDir {
				currentRun++
			} else {
				currentDir = rolls[i].Direction
				currentRun = 1
			}
			if currentRun > s.MaxSameDirRun[currentDir] {
				s.MaxSameDirRun[currentDir] = currentRun
			}
		}
	}

	return s, nil
}
