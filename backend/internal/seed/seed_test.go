package seed

import (
	"path/filepath"
	"runtime"
	"strings"
	"testing"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func TestLoadCatalog(t *testing.T) {
	catalog, err := LoadCatalog(testSeedDir(t))
	if err != nil {
		t.Fatalf("LoadCatalog() error = %v", err)
	}

	if got := len(catalog.Cities); got < MinCityCount || got > MaxCityCount {
		t.Fatalf("city count = %d, want between %d and %d", got, MinCityCount, MaxCityCount)
	}
	if got := len(catalog.Achievements); got != 11 {
		t.Fatalf("achievement count = %d, want 11", got)
	}

	expectedCities := map[string]bool{
		"北京": false, "西安": false, "南京": false, "杭州": false,
		"成都": false, "广州": false, "哈尔滨": false, "苏州": false,
		"大理": false, "厦门": false, "长沙": false, "拉萨": false,
	}
	for _, city := range catalog.Cities {
		if _, ok := expectedCities[city.Name]; ok {
			expectedCities[city.Name] = true
		}
		if len(city.Landmarks) < 1 || len(city.Landmarks) > 2 {
			t.Fatalf("city %q landmark count = %d", city.Name, len(city.Landmarks))
		}
		if len(city.Foods) < 1 || len(city.Foods) > 2 {
			t.Fatalf("city %q food count = %d", city.Name, len(city.Foods))
		}
		if len(city.Characters) != 1 {
			t.Fatalf("city %q character count = %d", city.Name, len(city.Characters))
		}
	}
	for city, found := range expectedCities {
		if !found {
			t.Fatalf("expected city %q is missing", city)
		}
	}
}

func TestValidateCatalogRejectsDuplicateCity(t *testing.T) {
	catalog := mustLoadCatalog(t)
	catalog.Cities[1].Name = catalog.Cities[0].Name

	err := ValidateCatalog(catalog)
	if err == nil || !strings.Contains(err.Error(), "duplicate city name") {
		t.Fatalf("ValidateCatalog() error = %v, want duplicate city name", err)
	}
}

func TestValidateCatalogRejectsImpossibleTagAchievement(t *testing.T) {
	catalog := mustLoadCatalog(t)
	for i := range catalog.Achievements {
		if catalog.Achievements[i].RuleType == "tag_count" {
			catalog.Achievements[i].RuleValue = "spicy_food:99"
			break
		}
	}

	err := ValidateCatalog(catalog)
	if err == nil || !strings.Contains(err.Error(), "seed has") {
		t.Fatalf("ValidateCatalog() error = %v, want impossible tag_count", err)
	}
}

func TestValidateCatalogRejectsNonCompliantPrompt(t *testing.T) {
	catalog := mustLoadCatalog(t)
	catalog.Cities[0].Characters[0].Prompt = "只介绍城市风景。"

	err := ValidateCatalog(catalog)
	if err == nil || !strings.Contains(err.Error(), "compliance reminders") {
		t.Fatalf("ValidateCatalog() error = %v, want compliance reminder error", err)
	}
}

func TestUpsertCatalogGeneratesIdempotentStatements(t *testing.T) {
	catalog := mustLoadCatalog(t)
	db, err := gorm.Open(mysql.New(mysql.Config{
		DSN:                       "city_roam:city_roam@tcp(localhost:3306)/city_roam?charset=utf8mb4&parseTime=True&loc=Local",
		SkipInitializeWithVersion: true,
	}), &gorm.Config{
		DryRun:                 true,
		DisableAutomaticPing:   true,
		SkipDefaultTransaction: true,
	})
	if err != nil {
		t.Fatalf("gorm.Open() error = %v", err)
	}

	var statements []string
	if err := db.Callback().Create().After("gorm:create").Register("test:capture_upserts", func(tx *gorm.DB) {
		statements = append(statements, tx.Statement.SQL.String())
	}); err != nil {
		t.Fatalf("register create callback: %v", err)
	}

	if err := upsertCatalog(db, catalog); err != nil {
		t.Fatalf("upsertCatalog() error = %v", err)
	}
	if len(statements) == 0 {
		t.Fatal("upsertCatalog() generated no create statements")
	}
	for _, statement := range statements {
		if !strings.Contains(statement, "ON DUPLICATE KEY UPDATE") {
			t.Errorf("statement is not idempotent: %s", statement)
		}
	}
}

func mustLoadCatalog(t *testing.T) Catalog {
	t.Helper()
	catalog, err := LoadCatalog(testSeedDir(t))
	if err != nil {
		t.Fatalf("LoadCatalog() error = %v", err)
	}
	return catalog
}

func testSeedDir(t *testing.T) string {
	t.Helper()
	_, filename, _, ok := runtime.Caller(0)
	if !ok {
		t.Fatal("runtime.Caller() failed")
	}
	return filepath.Join(filepath.Dir(filename), "..", "..", "data", "seed")
}
