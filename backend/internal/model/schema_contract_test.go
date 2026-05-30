package model

import (
	"os"
	"path/filepath"
	"reflect"
	"runtime"
	"strings"
	"testing"
)

func TestSchemaDefinesTwelveTablesAndRequiredIndexes(t *testing.T) {
	schema := readSchema(t)
	if got := strings.Count(schema, "CREATE TABLE IF NOT EXISTS "); got != 12 {
		t.Fatalf("table count = %d, want 12", got)
	}

	required := []string{
		"UNIQUE KEY uk_anonymous_id (anonymous_id)",
		"KEY idx_user_current_city (current_city_id)",
		"UNIQUE KEY uk_city_name (name)",
		"UNIQUE KEY uk_ct_city_tag (city_id, tag)",
		"UNIQUE KEY uk_lm_city_name (city_id, name)",
		"UNIQUE KEY uk_food_city_name (city_id, name)",
		"UNIQUE KEY uk_char_city_name (city_id, name)",
		"KEY idx_cv_user_city (user_id, city_id)",
		"KEY idx_cv_user_time (user_id, created_at)",
		"KEY idx_cv_city (city_id)",
		"KEY idx_cv_from_city (from_city_id)",
		"KEY idx_cv_dice_roll (dice_roll_id)",
		"KEY idx_dr_user_time (user_id, created_at)",
		"KEY idx_dr_from_city (from_city_id)",
		"KEY idx_dr_to_city (to_city_id)",
		"KEY idx_ck_user_time (user_id, created_at)",
		"KEY idx_ck_user_city (user_id, city_id)",
		"KEY idx_ck_city (city_id)",
		"KEY idx_ck_landmark (landmark_id)",
		"KEY idx_ck_visit (visit_id)",
		"UNIQUE KEY uk_ach_code (code)",
		"UNIQUE KEY uk_ua_user_ach (user_id, achievement_id)",
		"KEY idx_ua_achievement (achievement_id)",
		"KEY idx_cm_user_char_time (user_id, character_id, created_at)",
		"KEY idx_cm_city (city_id)",
		"KEY idx_cm_character (character_id)",
	}
	for _, item := range required {
		if !strings.Contains(schema, item) {
			t.Errorf("schema is missing %q", item)
		}
	}
}

func TestModelTableNames(t *testing.T) {
	tables := map[string]string{
		(User{}).TableName():            "users",
		(City{}).TableName():            "cities",
		(CityTag{}).TableName():         "city_tags",
		(Landmark{}).TableName():        "landmarks",
		(Food{}).TableName():            "foods",
		(Character{}).TableName():       "characters",
		(CityVisit{}).TableName():       "city_visits",
		(DiceRoll{}).TableName():        "dice_rolls",
		(Checkin{}).TableName():         "checkins",
		(Achievement{}).TableName():     "achievements",
		(UserAchievement{}).TableName(): "user_achievements",
		(ChatMessage{}).TableName():     "chat_messages",
	}
	if len(tables) != 12 {
		t.Fatalf("unique model table names = %d, want 12", len(tables))
	}
	for got, want := range tables {
		if got != want {
			t.Errorf("table name = %q, want %q", got, want)
		}
	}
}

func TestModelIndexTagsMatchSchemaContract(t *testing.T) {
	tests := []struct {
		model any
		field string
		want  []string
	}{
		{User{}, "CurrentCityID", []string{"index:idx_user_current_city"}},
		{City{}, "Name", []string{"uniqueIndex:uk_city_name"}},
		{CityTag{}, "CityID", []string{"index:idx_ct_city", "uniqueIndex:uk_ct_city_tag,priority:1"}},
		{CityTag{}, "Tag", []string{"index:idx_ct_tag", "uniqueIndex:uk_ct_city_tag,priority:2"}},
		{Landmark{}, "CityID", []string{"index:idx_lm_city", "uniqueIndex:uk_lm_city_name,priority:1"}},
		{Landmark{}, "Name", []string{"uniqueIndex:uk_lm_city_name,priority:2"}},
		{Food{}, "CityID", []string{"index:idx_food_city", "uniqueIndex:uk_food_city_name,priority:1"}},
		{Food{}, "Name", []string{"uniqueIndex:uk_food_city_name,priority:2"}},
		{Character{}, "CityID", []string{"index:idx_char_city", "uniqueIndex:uk_char_city_name,priority:1"}},
		{Character{}, "Name", []string{"uniqueIndex:uk_char_city_name,priority:2"}},
		{CityVisit{}, "CityID", []string{"index:idx_cv_city"}},
		{CityVisit{}, "FromCityID", []string{"index:idx_cv_from_city"}},
		{CityVisit{}, "DiceRollID", []string{"index:idx_cv_dice_roll"}},
		{DiceRoll{}, "FromCityID", []string{"index:idx_dr_from_city"}},
		{DiceRoll{}, "ToCityID", []string{"index:idx_dr_to_city"}},
		{Checkin{}, "CityID", []string{"index:idx_ck_city"}},
		{Checkin{}, "LandmarkID", []string{"index:idx_ck_landmark"}},
		{Checkin{}, "VisitID", []string{"index:idx_ck_visit"}},
		{UserAchievement{}, "AchievementID", []string{"index:idx_ua_achievement"}},
		{ChatMessage{}, "CityID", []string{"index:idx_cm_city"}},
		{ChatMessage{}, "CharacterID", []string{"index:idx_cm_character"}},
	}

	for _, tt := range tests {
		modelType := reflect.TypeOf(tt.model)
		field, ok := modelType.FieldByName(tt.field)
		if !ok {
			t.Errorf("%s.%s does not exist", modelType.Name(), tt.field)
			continue
		}
		tag := field.Tag.Get("gorm")
		for _, want := range tt.want {
			if !strings.Contains(tag, want) {
				t.Errorf("%s.%s gorm tag %q does not contain %q", modelType.Name(), tt.field, tag, want)
			}
		}
	}
}

func readSchema(t *testing.T) string {
	t.Helper()
	_, filename, _, ok := runtime.Caller(0)
	if !ok {
		t.Fatal("runtime.Caller() failed")
	}
	path := filepath.Join(filepath.Dir(filename), "..", "..", "migrations", "schema.sql")
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read schema: %v", err)
	}
	return string(data)
}
