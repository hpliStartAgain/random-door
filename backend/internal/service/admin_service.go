package service

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/your-org/city-roam/backend/internal/model"
	"gorm.io/gorm"
)

type adminRepo interface {
	ListCities(ctx context.Context) ([]model.City, error)
	FindCityByID(ctx context.Context, id int64) (*model.City, error)
	ListTags(ctx context.Context, cityID int64) ([]model.CityTag, error)
	ListAllTags(ctx context.Context) ([]model.CityTag, error)
	ListLandmarks(ctx context.Context, cityID int64) ([]model.Landmark, error)
	ListFoods(ctx context.Context, cityID int64) ([]model.Food, error)
	ListCharacters(ctx context.Context, cityID int64) ([]model.Character, error)
	CreateLandmark(ctx context.Context, row *model.Landmark) error
	CreateFood(ctx context.Context, row *model.Food) error
	CreateCharacter(ctx context.Context, row *model.Character) error
	ListAchievements(ctx context.Context) ([]model.Achievement, error)
	CreateAchievement(ctx context.Context, row *model.Achievement) error
	UpdateAchievement(ctx context.Context, id int64, fields map[string]any) error
	DeleteAchievement(ctx context.Context, id int64) error
	RenameTag(ctx context.Context, oldTag, newTag string) error
	DeleteTag(ctx context.Context, tag string) error
	UpdateCity(ctx context.Context, id int64, fields map[string]any, tags *[]string) error
	UpdateLandmark(ctx context.Context, id int64, fields map[string]any) error
	UpdateFood(ctx context.Context, id int64, fields map[string]any) error
	UpdateCharacter(ctx context.Context, id int64, fields map[string]any) error
	DeleteLandmark(ctx context.Context, id int64) error
	DeleteFood(ctx context.Context, id int64) error
	DeleteCharacter(ctx context.Context, id int64) error
}

type bytesStorage interface {
	SaveBytes(data []byte, subDir, ext string) (string, error)
}

type AdminService struct {
	repo    adminRepo
	storage bytesStorage
	client  *http.Client
}

func NewAdminService(repo adminRepo, storage bytesStorage) *AdminService {
	return &AdminService{
		repo:    repo,
		storage: storage,
		client:  &http.Client{Timeout: 15 * time.Second},
	}
}

type CatalogCoverage struct {
	TotalCities    int            `json:"total_cities"`
	CompleteCities int            `json:"complete_cities"`
	Items          []CoverageItem `json:"items"`
}

type CoverageItem struct {
	CityID         int64    `json:"city_id"`
	CityName       string   `json:"city_name"`
	HasCoverImage  bool     `json:"has_cover_image"`
	TagCount       int      `json:"tag_count"`
	LandmarkCount  int      `json:"landmark_count"`
	FoodCount      int      `json:"food_count"`
	CharacterCount int      `json:"character_count"`
	MissingFields  []string `json:"missing_fields"`
}

type TagListResult struct {
	Tags []TagItem `json:"tags"`
}

type TagItem struct {
	Tag       string `json:"tag"`
	CityCount int    `json:"city_count"`
}

type RenameTagRequest struct {
	Tag string `json:"tag"`
}

type AchievementListResult struct {
	Achievements []model.Achievement `json:"achievements"`
}

type CreateAchievementRequest struct {
	Code        string  `json:"code"`
	Name        string  `json:"name"`
	Description *string `json:"description"`
	RuleType    string  `json:"rule_type"`
	RuleValue   string  `json:"rule_value"`
	BadgeURL    *string `json:"badge_url"`
}

type UpdateAchievementRequest struct {
	Code        *string `json:"code"`
	Name        *string `json:"name"`
	Description *string `json:"description"`
	RuleType    *string `json:"rule_type"`
	RuleValue   *string `json:"rule_value"`
	BadgeURL    *string `json:"badge_url"`
}

func (s *AdminService) Coverage(ctx context.Context) (*CatalogCoverage, error) {
	cities, err := s.repo.ListCities(ctx)
	if err != nil {
		return nil, fmt.Errorf("list cities: %w", err)
	}
	result := &CatalogCoverage{TotalCities: len(cities), Items: make([]CoverageItem, 0, len(cities))}
	for _, city := range cities {
		tags, err := s.repo.ListTags(ctx, city.ID)
		if err != nil {
			return nil, fmt.Errorf("list tags: %w", err)
		}
		landmarks, err := s.repo.ListLandmarks(ctx, city.ID)
		if err != nil {
			return nil, fmt.Errorf("list landmarks: %w", err)
		}
		foods, err := s.repo.ListFoods(ctx, city.ID)
		if err != nil {
			return nil, fmt.Errorf("list foods: %w", err)
		}
		characters, err := s.repo.ListCharacters(ctx, city.ID)
		if err != nil {
			return nil, fmt.Errorf("list characters: %w", err)
		}

		item := CoverageItem{
			CityID:         city.ID,
			CityName:       city.Name,
			HasCoverImage:  isLocalAsset(city.CoverImageURL),
			TagCount:       len(tags),
			LandmarkCount:  len(landmarks),
			FoodCount:      len(foods),
			CharacterCount: len(characters),
		}
		item.MissingFields = coverageMissingFields(city, tags, landmarks, foods, characters)
		if len(item.MissingFields) == 0 {
			result.CompleteCities++
		}
		result.Items = append(result.Items, item)
	}
	return result, nil
}

func (s *AdminService) ListTags(ctx context.Context) (*TagListResult, error) {
	rows, err := s.repo.ListAllTags(ctx)
	if err != nil {
		return nil, fmt.Errorf("list all tags: %w", err)
	}
	counts := make(map[string]map[int64]bool)
	for _, row := range rows {
		tag := strings.TrimSpace(row.Tag)
		if tag == "" {
			continue
		}
		if counts[tag] == nil {
			counts[tag] = make(map[int64]bool)
		}
		counts[tag][row.CityID] = true
	}
	tags := make([]TagItem, 0, len(counts))
	for tag, cities := range counts {
		tags = append(tags, TagItem{Tag: tag, CityCount: len(cities)})
	}
	sort.Slice(tags, func(i, j int) bool {
		if tags[i].CityCount == tags[j].CityCount {
			return tags[i].Tag < tags[j].Tag
		}
		return tags[i].CityCount > tags[j].CityCount
	})
	return &TagListResult{Tags: tags}, nil
}

func (s *AdminService) RenameTag(ctx context.Context, oldTag string, req RenameTagRequest) error {
	oldTag, err := normalizeTag(oldTag)
	if err != nil {
		return err
	}
	newTag, err := normalizeTag(req.Tag)
	if err != nil {
		return err
	}
	if oldTag == newTag {
		return nil
	}
	if err := s.repo.RenameTag(ctx, oldTag, newTag); err != nil {
		return classifyAdminRepoError(err, "tag not found")
	}
	return nil
}

func (s *AdminService) DeleteTag(ctx context.Context, tag string) error {
	tag, err := normalizeTag(tag)
	if err != nil {
		return err
	}
	achievements, err := s.repo.ListAchievements(ctx)
	if err != nil {
		return fmt.Errorf("list achievements: %w", err)
	}
	for _, ach := range achievements {
		if achievementUsesTag(ach, tag) {
			return conflict("tag is used by achievements")
		}
	}
	if err := s.repo.DeleteTag(ctx, tag); err != nil {
		return classifyAdminRepoError(err, "tag not found")
	}
	return nil
}

func (s *AdminService) ListAchievements(ctx context.Context) (*AchievementListResult, error) {
	rows, err := s.repo.ListAchievements(ctx)
	if err != nil {
		return nil, fmt.Errorf("list achievements: %w", err)
	}
	return &AchievementListResult{Achievements: rows}, nil
}

func (s *AdminService) CreateAchievement(ctx context.Context, req CreateAchievementRequest) (*model.Achievement, error) {
	row, err := createAchievementRow(req)
	if err != nil {
		return nil, err
	}
	if err := s.repo.CreateAchievement(ctx, row); err != nil {
		return nil, classifyCreateError(err, "achievement already exists")
	}
	return row, nil
}

func (s *AdminService) UpdateAchievement(ctx context.Context, id int64, req UpdateAchievementRequest) error {
	if id <= 0 {
		return invalidParam("achievement_id must be a positive integer")
	}
	fields, err := achievementUpdateFields(req)
	if err != nil {
		return err
	}
	if err := s.repo.UpdateAchievement(ctx, id, fields); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return notFound("achievement not found")
		}
		return classifyCreateError(err, "achievement already exists")
	}
	return nil
}

func (s *AdminService) DeleteAchievement(ctx context.Context, id int64) error {
	if id <= 0 {
		return invalidParam("achievement_id must be a positive integer")
	}
	if err := s.repo.DeleteAchievement(ctx, id); err != nil {
		return classifyAdminRepoError(err, "achievement not found")
	}
	return nil
}

type UpdateCityRequest struct {
	Name               *string   `json:"name"`
	Province           *string   `json:"province"`
	Lat                *float64  `json:"lat"`
	Lng                *float64  `json:"lng"`
	Intro              *string   `json:"intro"`
	CoverImageURL      *string   `json:"cover_image_url"`
	DialectSample      *string   `json:"dialect_sample"`
	DialectExplanation *string   `json:"dialect_explanation"`
	Tags               *[]string `json:"tags"`
}

type UpdatePOIRequest struct {
	Name        *string  `json:"name"`
	Lat         *float64 `json:"lat"`
	Lng         *float64 `json:"lng"`
	ImageURL    *string  `json:"image_url"`
	Description *string  `json:"description"`
}

type CreatePOIRequest struct {
	Name        string   `json:"name"`
	Lat         *float64 `json:"lat"`
	Lng         *float64 `json:"lng"`
	ImageURL    *string  `json:"image_url"`
	Description *string  `json:"description"`
}

type UpdateCharacterRequest struct {
	Name          *string `json:"name"`
	CharacterType *string `json:"character_type"`
	AvatarURL     *string `json:"avatar_url"`
	Persona       *string `json:"persona"`
	DialectStyle  *string `json:"dialect_style"`
	Prompt        *string `json:"prompt"`
}

type CreateCharacterRequest struct {
	Name          string  `json:"name"`
	CharacterType *string `json:"character_type"`
	AvatarURL     *string `json:"avatar_url"`
	Persona       *string `json:"persona"`
	DialectStyle  *string `json:"dialect_style"`
	Prompt        *string `json:"prompt"`
}

func (s *AdminService) UpdateCity(ctx context.Context, id int64, req UpdateCityRequest) error {
	if id <= 0 {
		return invalidParam("city_id must be a positive integer")
	}
	fields := map[string]any{}
	addStringField(fields, "name", req.Name)
	addStringField(fields, "province", req.Province)
	addStringField(fields, "intro", req.Intro)
	addStringField(fields, "dialect_sample", req.DialectSample)
	addStringField(fields, "dialect_explanation", req.DialectExplanation)
	if req.Lat != nil {
		if *req.Lat < -90 || *req.Lat > 90 {
			return invalidParam("lat out of range")
		}
		fields["lat"] = *req.Lat
	}
	if req.Lng != nil {
		if *req.Lng < -180 || *req.Lng > 180 {
			return invalidParam("lng out of range")
		}
		fields["lng"] = *req.Lng
	}
	if req.CoverImageURL != nil {
		if err := validateLocalImageURL(*req.CoverImageURL); err != nil {
			return err
		}
		fields["cover_image_url"] = *req.CoverImageURL
	}
	if req.Tags != nil {
		if len(*req.Tags) == 0 {
			return invalidParam("tags must not be empty")
		}
		for _, tag := range *req.Tags {
			if strings.TrimSpace(tag) == "" {
				return invalidParam("tags must not contain empty values")
			}
		}
	}
	if err := s.repo.UpdateCity(ctx, id, fields, req.Tags); err != nil {
		return classifyAdminRepoError(err, "city not found")
	}
	return nil
}

func (s *AdminService) CreateLandmark(ctx context.Context, cityID int64, req CreatePOIRequest) (*model.Landmark, error) {
	fields, err := createLandmarkFields(req)
	if err != nil {
		return nil, err
	}
	if err := s.ensureCity(ctx, cityID); err != nil {
		return nil, err
	}
	row := &model.Landmark{
		CityID: cityID, Name: fields.name, Lat: fields.lat, Lng: fields.lng,
		ImageURL: fields.imageURL, Description: fields.description,
	}
	if err := s.repo.CreateLandmark(ctx, row); err != nil {
		return nil, classifyCreateError(err, "landmark already exists")
	}
	return row, nil
}

func (s *AdminService) UpdateLandmark(ctx context.Context, id int64, req UpdatePOIRequest) error {
	fields, err := landmarkFields(req)
	if err != nil {
		return err
	}
	if err := s.repo.UpdateLandmark(ctx, id, fields); err != nil {
		return classifyAdminRepoError(err, "landmark not found")
	}
	return nil
}

func (s *AdminService) DeleteLandmark(ctx context.Context, id int64) error {
	if id <= 0 {
		return invalidParam("landmark_id must be a positive integer")
	}
	if err := s.repo.DeleteLandmark(ctx, id); err != nil {
		return classifyAdminRepoError(err, "landmark not found")
	}
	return nil
}

func (s *AdminService) CreateFood(ctx context.Context, cityID int64, req CreatePOIRequest) (*model.Food, error) {
	fields, err := createPOIFields(req)
	if err != nil {
		return nil, err
	}
	if err := s.ensureCity(ctx, cityID); err != nil {
		return nil, err
	}
	row := &model.Food{CityID: cityID, Name: fields.name, ImageURL: fields.imageURL, Description: fields.description}
	if err := s.repo.CreateFood(ctx, row); err != nil {
		return nil, classifyCreateError(err, "food already exists")
	}
	return row, nil
}

func (s *AdminService) UpdateFood(ctx context.Context, id int64, req UpdatePOIRequest) error {
	fields, err := poiFields(req)
	if err != nil {
		return err
	}
	if err := s.repo.UpdateFood(ctx, id, fields); err != nil {
		return classifyAdminRepoError(err, "food not found")
	}
	return nil
}

func (s *AdminService) DeleteFood(ctx context.Context, id int64) error {
	if id <= 0 {
		return invalidParam("food_id must be a positive integer")
	}
	if err := s.repo.DeleteFood(ctx, id); err != nil {
		return classifyAdminRepoError(err, "food not found")
	}
	return nil
}

func (s *AdminService) CreateCharacter(ctx context.Context, cityID int64, req CreateCharacterRequest) (*model.Character, error) {
	name := strings.TrimSpace(req.Name)
	if name == "" {
		return nil, invalidParam("name is required")
	}
	if err := s.ensureCity(ctx, cityID); err != nil {
		return nil, err
	}
	characterType := "history"
	if req.CharacterType != nil {
		switch *req.CharacterType {
		case "history", "culture", "symbol":
			characterType = *req.CharacterType
		default:
			return nil, invalidParam("unsupported character_type")
		}
	}
	if req.AvatarURL != nil {
		if err := validateLocalImageURL(*req.AvatarURL); err != nil {
			return nil, err
		}
	}
	persona := fmt.Sprintf("%s 的城市文化导览角色。", name)
	if req.Persona != nil && strings.TrimSpace(*req.Persona) != "" {
		persona = strings.TrimSpace(*req.Persona)
	}
	prompt := defaultCharacterPrompt(name)
	if req.Prompt != nil && strings.TrimSpace(*req.Prompt) != "" {
		prompt = strings.TrimSpace(*req.Prompt)
		if !strings.Contains(prompt, "不声称真实复活") || !strings.Contains(prompt, "不编史") {
			return nil, invalidParam("prompt must contain compliance reminders")
		}
	}
	row := &model.Character{
		CityID:        cityID,
		Name:          name,
		CharacterType: characterType,
		AvatarURL:     req.AvatarURL,
		Persona:       persona,
		DialectStyle:  req.DialectStyle,
		Prompt:        prompt,
	}
	if err := s.repo.CreateCharacter(ctx, row); err != nil {
		return nil, classifyCreateError(err, "character already exists")
	}
	return row, nil
}

func (s *AdminService) UpdateCharacter(ctx context.Context, id int64, req UpdateCharacterRequest) error {
	fields := map[string]any{}
	addStringField(fields, "name", req.Name)
	if req.CharacterType != nil {
		switch *req.CharacterType {
		case "history", "culture", "symbol":
			fields["character_type"] = *req.CharacterType
		default:
			return invalidParam("unsupported character_type")
		}
	}
	if req.AvatarURL != nil {
		if err := validateLocalImageURL(*req.AvatarURL); err != nil {
			return err
		}
		fields["avatar_url"] = *req.AvatarURL
	}
	addStringField(fields, "persona", req.Persona)
	addStringField(fields, "dialect_style", req.DialectStyle)
	if req.Prompt != nil {
		if !strings.Contains(*req.Prompt, "不声称真实复活") || !strings.Contains(*req.Prompt, "不编史") {
			return invalidParam("prompt must contain compliance reminders")
		}
		fields["prompt"] = *req.Prompt
	}
	if err := s.repo.UpdateCharacter(ctx, id, fields); err != nil {
		return classifyAdminRepoError(err, "character not found")
	}
	return nil
}

func (s *AdminService) DeleteCharacter(ctx context.Context, id int64) error {
	if id <= 0 {
		return invalidParam("character_id must be a positive integer")
	}
	if err := s.repo.DeleteCharacter(ctx, id); err != nil {
		return classifyAdminRepoError(err, "character not found")
	}
	return nil
}

func (s *AdminService) ImportImageURL(ctx context.Context, remoteURL string) (string, error) {
	parsed, err := url.Parse(remoteURL)
	if err != nil || (parsed.Scheme != "http" && parsed.Scheme != "https") {
		return "", invalidParam("url must be http or https")
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, remoteURL, nil)
	if err != nil {
		return "", fmt.Errorf("create import request: %w", err)
	}
	resp, err := s.client.Do(req)
	if err != nil {
		return "", fmt.Errorf("download image: %w", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return "", invalidParam("image url is not reachable")
	}
	data, err := io.ReadAll(io.LimitReader(resp.Body, 6<<20))
	if err != nil {
		return "", fmt.Errorf("read image body: %w", err)
	}
	if len(data) > 5<<20 {
		return "", invalidParam("remote image exceeds 5MB")
	}
	contentType := http.DetectContentType(data)
	if !allowedImageContentType(contentType) {
		return "", invalidParam("remote image type not supported")
	}
	ext := strings.ToLower(filepath.Ext(parsed.Path))
	if ext == "" {
		ext = extensionFromContentType(contentType)
	}
	if ext == "" {
		ext = extensionFromContentType(resp.Header.Get("Content-Type"))
	}
	if !allowedImageExt(ext) {
		return "", invalidParam("remote image type not supported")
	}
	localURL, err := s.storage.SaveBytes(data, "admin_imports", ext)
	if err != nil {
		return "", fmt.Errorf("save imported image: %w", err)
	}
	return localURL, nil
}

func (s *AdminService) ensureCity(ctx context.Context, cityID int64) error {
	if cityID <= 0 {
		return invalidParam("city_id must be a positive integer")
	}
	if _, err := s.repo.FindCityByID(ctx, cityID); err != nil {
		return classifyAdminRepoError(err, "city not found")
	}
	return nil
}

type poiCreateFields struct {
	name        string
	lat         *float64
	lng         *float64
	imageURL    *string
	description *string
}

func createLandmarkFields(req CreatePOIRequest) (poiCreateFields, error) {
	fields, err := createPOIFields(req)
	if err != nil {
		return poiCreateFields{}, err
	}
	if err := validateOptionalCoordinates(req.Lat, req.Lng); err != nil {
		return poiCreateFields{}, err
	}
	fields.lat = req.Lat
	fields.lng = req.Lng
	return fields, nil
}

func createPOIFields(req CreatePOIRequest) (poiCreateFields, error) {
	name := strings.TrimSpace(req.Name)
	if name == "" {
		return poiCreateFields{}, invalidParam("name is required")
	}
	if req.ImageURL != nil {
		if err := validateLocalImageURL(*req.ImageURL); err != nil {
			return poiCreateFields{}, err
		}
	}
	return poiCreateFields{name: name, imageURL: req.ImageURL, description: req.Description}, nil
}

func createAchievementRow(req CreateAchievementRequest) (*model.Achievement, error) {
	code, err := normalizeCode(req.Code)
	if err != nil {
		return nil, err
	}
	name, err := normalizeRequiredText("name", req.Name, 128)
	if err != nil {
		return nil, err
	}
	ruleType := strings.TrimSpace(req.RuleType)
	ruleValue := strings.TrimSpace(req.RuleValue)
	if err := validateAchievementRule(ruleType, ruleValue); err != nil {
		return nil, err
	}
	if req.BadgeURL != nil && strings.TrimSpace(*req.BadgeURL) != "" {
		if err := validateLocalImageURL(*req.BadgeURL); err != nil {
			return nil, err
		}
	}
	return &model.Achievement{
		Code:        code,
		Name:        name,
		Description: trimOptionalText(req.Description),
		RuleType:    ruleType,
		RuleValue:   ruleValue,
		BadgeURL:    trimOptionalText(req.BadgeURL),
	}, nil
}

func achievementUpdateFields(req UpdateAchievementRequest) (map[string]any, error) {
	fields := map[string]any{}
	if req.Code != nil {
		code, err := normalizeCode(*req.Code)
		if err != nil {
			return nil, err
		}
		fields["code"] = code
	}
	if req.Name != nil {
		name, err := normalizeRequiredText("name", *req.Name, 128)
		if err != nil {
			return nil, err
		}
		fields["name"] = name
	}
	if req.Description != nil {
		fields["description"] = strings.TrimSpace(*req.Description)
	}
	if req.BadgeURL != nil {
		badgeURL := strings.TrimSpace(*req.BadgeURL)
		if badgeURL != "" {
			if err := validateLocalImageURL(badgeURL); err != nil {
				return nil, err
			}
		}
		fields["badge_url"] = badgeURL
	}
	if req.RuleType != nil || req.RuleValue != nil {
		if req.RuleType == nil || req.RuleValue == nil {
			return nil, invalidParam("rule_type and rule_value must be updated together")
		}
		ruleType := strings.TrimSpace(*req.RuleType)
		ruleValue := strings.TrimSpace(*req.RuleValue)
		if err := validateAchievementRule(ruleType, ruleValue); err != nil {
			return nil, err
		}
		fields["rule_type"] = ruleType
		fields["rule_value"] = ruleValue
	}
	return fields, nil
}

func normalizeTag(value string) (string, error) {
	value = strings.TrimSpace(value)
	if value == "" {
		return "", invalidParam("tag must not be empty")
	}
	if len([]rune(value)) > 64 {
		return "", invalidParam("tag must be at most 64 characters")
	}
	return value, nil
}

func normalizeCode(value string) (string, error) {
	value = strings.TrimSpace(value)
	if value == "" {
		return "", invalidParam("code is required")
	}
	for _, r := range value {
		if (r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z') || (r >= '0' && r <= '9') || r == '_' || r == '-' {
			continue
		}
		return "", invalidParam("code must contain only letters, numbers, underscores, or hyphens")
	}
	if len(value) > 64 {
		return "", invalidParam("code must be at most 64 characters")
	}
	return value, nil
}

func normalizeRequiredText(field, value string, max int) (string, error) {
	value = strings.TrimSpace(value)
	if value == "" {
		return "", invalidParam(field + " is required")
	}
	if len([]rune(value)) > max {
		return "", invalidParam(fmt.Sprintf("%s must be at most %d characters", field, max))
	}
	return value, nil
}

func trimOptionalText(value *string) *string {
	if value == nil {
		return nil
	}
	trimmed := strings.TrimSpace(*value)
	if trimmed == "" {
		return nil
	}
	return &trimmed
}

func validateAchievementRule(ruleType, ruleValue string) error {
	switch ruleType {
	case "first_checkin":
		if ruleValue != "" {
			return invalidParam("first_checkin rule_value must be empty")
		}
	case "checkin_count", "visit_count", "game_visit_count", "dice_distance":
		n, err := strconv.Atoi(ruleValue)
		if err != nil || n <= 0 {
			return invalidParam(ruleType + " rule_value must be a positive integer")
		}
	case "city_tag":
		if _, err := normalizeTag(ruleValue); err != nil {
			return err
		}
	case "tag_count":
		tag, count, ok := splitRulePair(ruleValue)
		if !ok || tag == "" || count <= 0 {
			return invalidParam("tag_count rule_value must be tag:count")
		}
	case "dice_direction":
		direction, count, ok := splitRulePair(ruleValue)
		if !ok || count <= 0 {
			return invalidParam("dice_direction rule_value must be direction:count")
		}
		if !isDiceDirection(direction) {
			return invalidParam("unsupported dice direction")
		}
	default:
		return invalidParam("unsupported rule_type")
	}
	return nil
}

func splitRulePair(value string) (string, int, bool) {
	for i := len(value) - 1; i >= 0; i-- {
		if value[i] != ':' {
			continue
		}
		n, err := strconv.Atoi(value[i+1:])
		if err != nil {
			return "", 0, false
		}
		return strings.TrimSpace(value[:i]), n, true
	}
	return "", 0, false
}

func isDiceDirection(value string) bool {
	switch value {
	case "北", "东北", "东", "东南", "南", "西南", "西", "西北":
		return true
	default:
		return false
	}
}

func achievementUsesTag(ach model.Achievement, tag string) bool {
	switch ach.RuleType {
	case "city_tag":
		return ach.RuleValue == tag
	case "tag_count":
		ruleTag, _, ok := splitRulePair(ach.RuleValue)
		return ok && ruleTag == tag
	default:
		return false
	}
}

func coverageMissingFields(city model.City, tags []model.CityTag, landmarks []model.Landmark, foods []model.Food, characters []model.Character) []string {
	missing := []string{}
	if strings.TrimSpace(city.Name) == "" {
		missing = append(missing, "name")
	}
	if strings.TrimSpace(city.Province) == "" {
		missing = append(missing, "province")
	}
	if city.Intro == nil || strings.TrimSpace(*city.Intro) == "" {
		missing = append(missing, "intro")
	}
	if !isLocalAsset(city.CoverImageURL) {
		missing = append(missing, "cover_image_url")
	}
	if city.DialectSample == nil || strings.TrimSpace(*city.DialectSample) == "" {
		missing = append(missing, "dialect_sample")
	}
	if city.DialectExplanation == nil || strings.TrimSpace(*city.DialectExplanation) == "" {
		missing = append(missing, "dialect_explanation")
	}
	if len(tags) == 0 {
		missing = append(missing, "tags")
	}
	if len(landmarks) < 1 || len(landmarks) > 2 {
		missing = append(missing, "landmarks")
	}
	for _, landmark := range landmarks {
		if !isLocalAsset(landmark.ImageURL) || strings.TrimSpace(landmark.Name) == "" ||
			isBlankPtr(landmark.Description) || landmark.Lat == nil || landmark.Lng == nil {
			missing = append(missing, "landmark:"+landmark.Name)
			break
		}
	}
	if len(foods) < 1 || len(foods) > 2 {
		missing = append(missing, "foods")
	}
	for _, food := range foods {
		if !isLocalAsset(food.ImageURL) || strings.TrimSpace(food.Name) == "" || isBlankPtr(food.Description) {
			missing = append(missing, "food:"+food.Name)
			break
		}
	}
	if len(characters) != 1 {
		missing = append(missing, "characters")
	}
	for _, character := range characters {
		if !isLocalAsset(character.AvatarURL) || strings.TrimSpace(character.Name) == "" ||
			strings.TrimSpace(character.Persona) == "" || strings.TrimSpace(character.Prompt) == "" {
			missing = append(missing, "character:"+character.Name)
			break
		}
	}
	return missing
}

func poiFields(req UpdatePOIRequest) (map[string]any, error) {
	fields := map[string]any{}
	addStringField(fields, "name", req.Name)
	addStringField(fields, "description", req.Description)
	if req.ImageURL != nil {
		if err := validateLocalImageURL(*req.ImageURL); err != nil {
			return nil, err
		}
		fields["image_url"] = *req.ImageURL
	}
	return fields, nil
}

func landmarkFields(req UpdatePOIRequest) (map[string]any, error) {
	fields, err := poiFields(req)
	if err != nil {
		return nil, err
	}
	if err := validateOptionalCoordinates(req.Lat, req.Lng); err != nil {
		return nil, err
	}
	if req.Lat != nil {
		fields["lat"] = *req.Lat
	}
	if req.Lng != nil {
		fields["lng"] = *req.Lng
	}
	return fields, nil
}

func validateOptionalCoordinates(lat, lng *float64) error {
	if lat == nil && lng == nil {
		return nil
	}
	if lat == nil || lng == nil {
		return invalidParam("lat and lng must be provided together")
	}
	if *lat < -90 || *lat > 90 {
		return invalidParam("lat out of range")
	}
	if *lng < -180 || *lng > 180 {
		return invalidParam("lng out of range")
	}
	return nil
}

func addStringField(fields map[string]any, key string, value *string) {
	if value != nil {
		fields[key] = *value
	}
}

func validateLocalImageURL(value string) error {
	if !strings.HasPrefix(value, "/static/") && !strings.HasPrefix(value, "/uploads/") {
		return invalidParam("image url must be a local /static or /uploads path")
	}
	return nil
}

func isLocalAsset(value *string) bool {
	return value != nil && (strings.HasPrefix(*value, "/static/") || strings.HasPrefix(*value, "/uploads/"))
}

func isBlankPtr(value *string) bool {
	return value == nil || strings.TrimSpace(*value) == ""
}

func classifyAdminRepoError(err error, message string) error {
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return notFound(message)
	}
	return err
}

func classifyCreateError(err error, message string) error {
	if strings.Contains(strings.ToLower(err.Error()), "duplicate") {
		return conflict(message)
	}
	return err
}

func defaultCharacterPrompt(name string) string {
	return fmt.Sprintf("你在和用户进行角色扮演的游戏，你扮演的人物是%s。不声称真实复活，不编史。回答控制在150字以内。", name)
}

func extensionFromContentType(contentType string) string {
	switch strings.ToLower(strings.TrimSpace(strings.Split(contentType, ";")[0])) {
	case "image/jpeg":
		return ".jpg"
	case "image/png":
		return ".png"
	case "image/webp":
		return ".webp"
	default:
		return ""
	}
}

func allowedImageContentType(contentType string) bool {
	switch strings.ToLower(strings.TrimSpace(strings.Split(contentType, ";")[0])) {
	case "image/jpeg", "image/png", "image/webp":
		return true
	default:
		return false
	}
}

func allowedImageExt(ext string) bool {
	switch ext {
	case ".jpg", ".jpeg", ".png", ".webp":
		return true
	default:
		return false
	}
}
