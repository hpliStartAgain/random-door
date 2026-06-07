package seed

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/your-org/city-roam/backend/internal/model"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

const (
	MinCityCount = 12
	MaxCityCount = 100
)

var (
	allowedCharacterTypes = map[string]bool{
		"history": true,
		"culture": true,
		"symbol":  true,
	}
	allowedDirections = map[string]bool{
		"北":  true,
		"东北": true,
		"东":  true,
		"东南": true,
		"南":  true,
		"西南": true,
		"西":  true,
		"西北": true,
	}
)

type Catalog struct {
	Cities       []City
	Achievements []Achievement
}

type City struct {
	Name               string      `json:"name"`
	Province           string      `json:"province"`
	Lat                float64     `json:"lat"`
	Lng                float64     `json:"lng"`
	Intro              string      `json:"intro"`
	CoverImageURL      string      `json:"cover_image_url"`
	DialectSample      string      `json:"dialect_sample"`
	DialectExplanation string      `json:"dialect_explanation"`
	Tags               []string    `json:"tags"`
	Landmarks          []POI       `json:"landmarks"`
	Foods              []POI       `json:"foods"`
	Characters         []Character `json:"characters"`
}

type POI struct {
	Name        string `json:"name"`
	ImageURL    string `json:"image_url"`
	Description string `json:"description"`
}

type Character struct {
	Name          string `json:"name"`
	CharacterType string `json:"character_type"`
	AvatarURL     string `json:"avatar_url"`
	Persona       string `json:"persona"`
	DialectStyle  string `json:"dialect_style"`
	Prompt        string `json:"prompt"`
	RoleTitle     string `json:"role_title"`
	LifeSpan      string `json:"life_span"`
	IntroQuote    string `json:"intro_quote"`
}

type Achievement struct {
	Code        string `json:"code"`
	Name        string `json:"name"`
	Description string `json:"description"`
	RuleType    string `json:"rule_type"`
	RuleValue   string `json:"rule_value"`
	BadgeURL    string `json:"badge_url"`
}

// Load validates the JSON catalog and upserts it atomically.
func Load(ctx context.Context, db *gorm.DB, dir string) error {
	catalog, err := LoadCatalog(dir)
	if err != nil {
		return err
	}

	return db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		return upsertCatalog(tx, catalog)
	})
}

// LoadCatalog reads and validates the JSON files without touching the database.
func LoadCatalog(dir string) (Catalog, error) {
	var catalog Catalog
	if err := readJSON(filepath.Join(dir, "cities.json"), &catalog.Cities); err != nil {
		return Catalog{}, err
	}
	if err := readJSON(filepath.Join(dir, "achievements.json"), &catalog.Achievements); err != nil {
		return Catalog{}, err
	}
	if err := ValidateCatalog(catalog); err != nil {
		return Catalog{}, fmt.Errorf("validate seed catalog: %w", err)
	}
	return catalog, nil
}

func readJSON(path string, dest any) error {
	data, err := os.ReadFile(path)
	if err != nil {
		return fmt.Errorf("read %s: %w", path, err)
	}
	if err := json.Unmarshal(data, dest); err != nil {
		return fmt.Errorf("parse %s: %w", path, err)
	}
	return nil
}

// ValidateCatalog rejects incomplete or internally inconsistent demo data.
func ValidateCatalog(catalog Catalog) error {
	if len(catalog.Cities) < MinCityCount || len(catalog.Cities) > MaxCityCount {
		return fmt.Errorf("expected %d-%d cities, got %d", MinCityCount, MaxCityCount, len(catalog.Cities))
	}
	if len(catalog.Achievements) == 0 {
		return fmt.Errorf("achievements must not be empty")
	}

	cityNames := make(map[string]bool)
	tagCityCount := make(map[string]int)
	for i, city := range catalog.Cities {
		if err := validateCity(city); err != nil {
			return fmt.Errorf("city[%d]: %w", i, err)
		}
		if cityNames[city.Name] {
			return fmt.Errorf("duplicate city name %q", city.Name)
		}
		cityNames[city.Name] = true

		for _, tag := range city.Tags {
			tagCityCount[tag]++
		}
	}

	achievementCodes := make(map[string]bool)
	hasFirstCheckin := false
	hasGameAchievement := false
	for i, achievement := range catalog.Achievements {
		if achievementCodes[achievement.Code] {
			return fmt.Errorf("duplicate achievement code %q", achievement.Code)
		}
		achievementCodes[achievement.Code] = true

		if err := validateAchievement(achievement, tagCityCount); err != nil {
			return fmt.Errorf("achievement[%d]: %w", i, err)
		}
		hasFirstCheckin = hasFirstCheckin || achievement.RuleType == "first_checkin"
		hasGameAchievement = hasGameAchievement ||
			achievement.RuleType == "game_visit_count" ||
			achievement.RuleType == "dice_direction" ||
			achievement.RuleType == "dice_distance"
	}
	if !hasFirstCheckin {
		return fmt.Errorf("first_checkin achievement is required")
	}
	if !hasGameAchievement {
		return fmt.Errorf("at least one game achievement is required")
	}
	return nil
}

func validateCity(city City) error {
	if err := requireText("name", city.Name); err != nil {
		return err
	}
	if err := requireText("province", city.Province); err != nil {
		return err
	}
	if city.Lat < -90 || city.Lat > 90 || city.Lng < -180 || city.Lng > 180 {
		return fmt.Errorf("coordinates out of range for %q", city.Name)
	}
	if err := requireText("intro", city.Intro); err != nil {
		return err
	}
	if err := requireStaticURL("cover_image_url", city.CoverImageURL); err != nil {
		return err
	}
	if err := requireText("dialect_sample", city.DialectSample); err != nil {
		return err
	}
	if err := requireText("dialect_explanation", city.DialectExplanation); err != nil {
		return err
	}
	if len(city.Tags) == 0 {
		return fmt.Errorf("tags must not be empty for %q", city.Name)
	}
	if err := uniqueStrings("tag", city.Tags); err != nil {
		return fmt.Errorf("%s in city %q", err, city.Name)
	}
	if len(city.Landmarks) < 1 || len(city.Landmarks) > 2 {
		return fmt.Errorf("city %q must have 1 or 2 landmarks", city.Name)
	}
	if err := validatePOIs("landmark", city.Name, city.Landmarks); err != nil {
		return err
	}
	if len(city.Foods) < 1 || len(city.Foods) > 2 {
		return fmt.Errorf("city %q must have 1 or 2 foods", city.Name)
	}
	if err := validatePOIs("food", city.Name, city.Foods); err != nil {
		return err
	}
	if len(city.Characters) != 1 {
		return fmt.Errorf("city %q must have exactly 1 character", city.Name)
	}

	character := city.Characters[0]
	if err := requireText("character name", character.Name); err != nil {
		return err
	}
	if !allowedCharacterTypes[character.CharacterType] {
		return fmt.Errorf("unsupported character_type %q in city %q", character.CharacterType, city.Name)
	}
	if err := requireStaticURL("avatar_url", character.AvatarURL); err != nil {
		return err
	}
	if err := requireText("persona", character.Persona); err != nil {
		return err
	}
	if err := requireText("dialect_style", character.DialectStyle); err != nil {
		return err
	}
	if err := requireText("prompt", character.Prompt); err != nil {
		return err
	}
	if !strings.Contains(character.Prompt, "不声称真实复活") || !strings.Contains(character.Prompt, "不编史") {
		return fmt.Errorf("character prompt in city %q must contain compliance reminders", city.Name)
	}
	return nil
}

func validatePOIs(kind, cityName string, pois []POI) error {
	names := make([]string, 0, len(pois))
	for _, poi := range pois {
		if err := requireText(kind+" name", poi.Name); err != nil {
			return err
		}
		if err := requireStaticURL(kind+" image_url", poi.ImageURL); err != nil {
			return err
		}
		if err := requireText(kind+" description", poi.Description); err != nil {
			return err
		}
		names = append(names, poi.Name)
	}
	if err := uniqueStrings(kind, names); err != nil {
		return fmt.Errorf("%s in city %q", err, cityName)
	}
	return nil
}

func validateAchievement(achievement Achievement, tagCityCount map[string]int) error {
	if err := requireText("achievement code", achievement.Code); err != nil {
		return err
	}
	if err := requireText("achievement name", achievement.Name); err != nil {
		return err
	}
	if err := requireText("achievement description", achievement.Description); err != nil {
		return err
	}
	if err := requireStaticURL("achievement badge_url", achievement.BadgeURL); err != nil {
		return err
	}

	switch achievement.RuleType {
	case "first_checkin":
		if achievement.RuleValue != "" {
			return fmt.Errorf("first_checkin rule_value must be empty")
		}
	case "checkin_count", "game_visit_count", "dice_distance":
		if _, err := parsePositiveInt(achievement.RuleValue); err != nil {
			return fmt.Errorf("%s: %w", achievement.RuleType, err)
		}
	case "city_tag":
		if tagCityCount[achievement.RuleValue] == 0 {
			return fmt.Errorf("city_tag %q is not present in city seed", achievement.RuleValue)
		}
	case "tag_count":
		tag, count, err := parsePair(achievement.RuleValue)
		if err != nil {
			return fmt.Errorf("tag_count: %w", err)
		}
		if tagCityCount[tag] < count {
			return fmt.Errorf("tag_count %q requires %d cities, seed has %d", tag, count, tagCityCount[tag])
		}
	case "dice_direction":
		direction, _, err := parsePair(achievement.RuleValue)
		if err != nil {
			return fmt.Errorf("dice_direction: %w", err)
		}
		if !allowedDirections[direction] {
			return fmt.Errorf("unsupported dice direction %q", direction)
		}
	default:
		return fmt.Errorf("unsupported rule_type %q", achievement.RuleType)
	}
	return nil
}

func requireText(field, value string) error {
	if strings.TrimSpace(value) == "" {
		return fmt.Errorf("%s must not be empty", field)
	}
	return nil
}

func requireStaticURL(field, value string) error {
	if !strings.HasPrefix(value, "/static/") {
		return fmt.Errorf("%s must start with /static/", field)
	}
	return nil
}

func uniqueStrings(kind string, values []string) error {
	seen := make(map[string]bool)
	for _, value := range values {
		if strings.TrimSpace(value) == "" {
			return fmt.Errorf("%s must not be empty", kind)
		}
		if seen[value] {
			return fmt.Errorf("duplicate %s %q", kind, value)
		}
		seen[value] = true
	}
	return nil
}

func parsePair(value string) (string, int, error) {
	parts := strings.SplitN(value, ":", 2)
	if len(parts) != 2 || strings.TrimSpace(parts[0]) == "" {
		return "", 0, fmt.Errorf("expected name:count, got %q", value)
	}
	count, err := parsePositiveInt(parts[1])
	if err != nil {
		return "", 0, err
	}
	return parts[0], count, nil
}

func parsePositiveInt(value string) (int, error) {
	n, err := strconv.Atoi(value)
	if err != nil || n <= 0 {
		return 0, fmt.Errorf("expected positive integer, got %q", value)
	}
	return n, nil
}

func upsertCatalog(tx *gorm.DB, catalog Catalog) error {
	for _, citySeed := range catalog.Cities {
		if err := upsertCity(tx, citySeed); err != nil {
			return err
		}
	}
	for _, achievementSeed := range catalog.Achievements {
		if err := upsertAchievement(tx, achievementSeed); err != nil {
			return err
		}
	}
	return nil
}

func upsertCity(tx *gorm.DB, source City) error {
	existingCity, err := findExistingCity(tx, source.Name)
	if err != nil {
		return fmt.Errorf("load existing city %q: %w", source.Name, err)
	}
	city := model.City{
		Name:               source.Name,
		Province:           source.Province,
		Lat:                source.Lat,
		Lng:                source.Lng,
		Intro:              ptr(source.Intro),
		CoverImageURL:      seedAssetURL(existingCity.CoverImageURL, source.CoverImageURL),
		DialectSample:      ptr(source.DialectSample),
		DialectExplanation: ptr(source.DialectExplanation),
	}
	if err := tx.Clauses(clause.OnConflict{
		Columns: []clause.Column{{Name: "name"}},
		DoUpdates: clause.AssignmentColumns([]string{
			"province", "lat", "lng", "intro", "cover_image_url",
			"dialect_sample", "dialect_explanation", "updated_at",
		}),
	}).Create(&city).Error; err != nil {
		return fmt.Errorf("upsert city %q: %w", source.Name, err)
	}
	if err := tx.Where("name = ?", source.Name).First(&city).Error; err != nil {
		return fmt.Errorf("reload city %q: %w", source.Name, err)
	}

	for _, tag := range source.Tags {
		row := model.CityTag{CityID: city.ID, Tag: tag}
		if err := tx.Clauses(clause.OnConflict{DoNothing: true}).Create(&row).Error; err != nil {
			return fmt.Errorf("upsert tag %q for city %q: %w", tag, source.Name, err)
		}
	}
	for _, landmark := range source.Landmarks {
		existingLandmark, err := findExistingLandmark(tx, city.ID, landmark.Name)
		if err != nil {
			return fmt.Errorf("load existing landmark %q for city %q: %w", landmark.Name, source.Name, err)
		}
		row := model.Landmark{
			CityID:      city.ID,
			Name:        landmark.Name,
			ImageURL:    seedAssetURL(existingLandmark.ImageURL, landmark.ImageURL),
			Description: ptr(landmark.Description),
		}
		if err := tx.Clauses(clause.OnConflict{
			Columns:   []clause.Column{{Name: "city_id"}, {Name: "name"}},
			DoUpdates: clause.AssignmentColumns([]string{"image_url", "description"}),
		}).Create(&row).Error; err != nil {
			return fmt.Errorf("upsert landmark %q for city %q: %w", landmark.Name, source.Name, err)
		}
	}
	for _, food := range source.Foods {
		existingFood, err := findExistingFood(tx, city.ID, food.Name)
		if err != nil {
			return fmt.Errorf("load existing food %q for city %q: %w", food.Name, source.Name, err)
		}
		row := model.Food{
			CityID:      city.ID,
			Name:        food.Name,
			ImageURL:    seedAssetURL(existingFood.ImageURL, food.ImageURL),
			Description: ptr(food.Description),
		}
		if err := tx.Clauses(clause.OnConflict{
			Columns:   []clause.Column{{Name: "city_id"}, {Name: "name"}},
			DoUpdates: clause.AssignmentColumns([]string{"image_url", "description"}),
		}).Create(&row).Error; err != nil {
			return fmt.Errorf("upsert food %q for city %q: %w", food.Name, source.Name, err)
		}
	}
	for _, character := range source.Characters {
		existingCharacter, err := findExistingCharacter(tx, city.ID, character.Name)
		if err != nil {
			return fmt.Errorf("load existing character %q for city %q: %w", character.Name, source.Name, err)
		}
		row := model.Character{
			CityID:        city.ID,
			Name:          character.Name,
			CharacterType: character.CharacterType,
			AvatarURL:     seedAssetURL(existingCharacter.AvatarURL, character.AvatarURL),
			Persona:       character.Persona,
			DialectStyle:  ptr(character.DialectStyle),
			Prompt:        character.Prompt,
			RoleTitle:     nilIfEmpty(character.RoleTitle),
			LifeSpan:      nilIfEmpty(character.LifeSpan),
			IntroQuote:    nilIfEmpty(character.IntroQuote),
		}
		if err := tx.Clauses(clause.OnConflict{
			Columns: []clause.Column{{Name: "city_id"}, {Name: "name"}},
			DoUpdates: clause.AssignmentColumns([]string{
				"character_type", "avatar_url", "persona", "dialect_style", "prompt",
				"role_title", "life_span", "intro_quote",
			}),
		}).Create(&row).Error; err != nil {
			return fmt.Errorf("upsert character %q for city %q: %w", character.Name, source.Name, err)
		}
	}
	return nil
}

func findExistingCity(tx *gorm.DB, name string) (model.City, error) {
	var row model.City
	err := tx.Where("name = ?", name).First(&row).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return model.City{}, nil
	}
	return row, err
}

func findExistingLandmark(tx *gorm.DB, cityID int64, name string) (model.Landmark, error) {
	var row model.Landmark
	err := tx.Where("city_id = ? AND name = ?", cityID, name).First(&row).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return model.Landmark{}, nil
	}
	return row, err
}

func findExistingFood(tx *gorm.DB, cityID int64, name string) (model.Food, error) {
	var row model.Food
	err := tx.Where("city_id = ? AND name = ?", cityID, name).First(&row).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return model.Food{}, nil
	}
	return row, err
}

func findExistingCharacter(tx *gorm.DB, cityID int64, name string) (model.Character, error) {
	var row model.Character
	err := tx.Where("city_id = ? AND name = ?", cityID, name).First(&row).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return model.Character{}, nil
	}
	return row, err
}

func upsertAchievement(tx *gorm.DB, source Achievement) error {
	row := model.Achievement{
		Code:        source.Code,
		Name:        source.Name,
		Description: ptr(source.Description),
		RuleType:    source.RuleType,
		RuleValue:   source.RuleValue,
		BadgeURL:    ptr(source.BadgeURL),
	}
	if err := tx.Clauses(clause.OnConflict{
		Columns: []clause.Column{{Name: "code"}},
		DoUpdates: clause.AssignmentColumns([]string{
			"name", "description", "rule_type", "rule_value", "badge_url",
		}),
	}).Create(&row).Error; err != nil {
		return fmt.Errorf("upsert achievement %q: %w", source.Code, err)
	}
	return nil
}

func seedAssetURL(existing *string, seeded string) *string {
	if existing != nil && strings.HasPrefix(*existing, "/uploads/") {
		preserved := *existing
		return &preserved
	}
	return ptr(seeded)
}

func ptr(value string) *string {
	return &value
}

func nilIfEmpty(s string) *string {
	if s == "" {
		return nil
	}
	return &s
}
