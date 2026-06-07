package main

import (
	"context"
	"flag"
	"fmt"
	"io"
	"log"
	"sort"
	"strings"
	"time"

	"github.com/joho/godotenv"
	"github.com/your-org/city-roam/backend/internal/config"
	"github.com/your-org/city-roam/backend/internal/model"
	"github.com/your-org/city-roam/backend/internal/seed"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type catalogCounts struct {
	Cities       int64
	CityTags     int64
	Landmarks    int64
	Foods        int64
	Characters   int64
	Achievements int64
}

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}

func run() error {
	mode := flag.String("mode", "audit", "audit, bootstrap, or sync")
	confirmOverwrite := flag.Bool("confirm-overwrite", false, "required with -mode sync because it overwrites matching catalog rows")
	flag.Parse()

	_ = godotenv.Load("../.env")
	cfg := loadConfig()

	db, err := openDB(cfg)
	if err != nil {
		return err
	}
	sqlDB, err := db.DB()
	if err != nil {
		return err
	}
	defer sqlDB.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	switch strings.ToLower(strings.TrimSpace(*mode)) {
	case "audit":
		catalog, err := seed.LoadCatalog(cfg.Server.SeedDir)
		if err != nil {
			return err
		}
		return auditCatalog(ctx, db, catalog)
	case string(seed.ImportModeBootstrap):
		if err := seed.Bootstrap(ctx, db, cfg.Server.SeedDir); err != nil {
			return err
		}
		fmt.Printf("seed bootstrap completed from %s\n", cfg.Server.SeedDir)
		return nil
	case string(seed.ImportModeSync):
		if !*confirmOverwrite {
			return fmt.Errorf("sync mode overwrites matching catalog rows; rerun with -confirm-overwrite after reviewing -mode audit")
		}
		if err := seed.Sync(ctx, db, cfg.Server.SeedDir); err != nil {
			return err
		}
		fmt.Printf("seed sync completed from %s\n", cfg.Server.SeedDir)
		return nil
	default:
		return fmt.Errorf("unsupported mode %q (want audit, bootstrap, or sync)", *mode)
	}
}

func loadConfig() *config.Config {
	logOutput := log.Writer()
	log.SetOutput(io.Discard)
	cfg := config.Load()
	log.SetOutput(logOutput)
	return cfg
}

func openDB(cfg *config.Config) (*gorm.DB, error) {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		cfg.DB.User, cfg.DB.Password, cfg.DB.Host, cfg.DB.Port, cfg.DB.Name)
	return gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
}

func auditCatalog(ctx context.Context, db *gorm.DB, catalog seed.Catalog) error {
	dbCounts, err := countDBCatalog(ctx, db)
	if err != nil {
		return err
	}
	seedCounts := countSeedCatalog(catalog)

	dbCityNames, err := listDBCityNames(ctx, db)
	if err != nil {
		return err
	}
	seedCityNames := make([]string, 0, len(catalog.Cities))
	for _, city := range catalog.Cities {
		seedCityNames = append(seedCityNames, city.Name)
	}

	dbCitySet := stringSet(dbCityNames)
	seedCitySet := stringSet(seedCityNames)
	missingSeedCities := missingFrom(seedCityNames, dbCitySet)
	dbOnlyCities := missingFrom(dbCityNames, seedCitySet)
	overlap := len(seedCityNames) - len(missingSeedCities)

	fmt.Printf("database counts: cities=%d city_tags=%d landmarks=%d foods=%d characters=%d achievements=%d\n",
		dbCounts.Cities, dbCounts.CityTags, dbCounts.Landmarks, dbCounts.Foods, dbCounts.Characters, dbCounts.Achievements)
	fmt.Printf("seed counts: cities=%d city_tags=%d landmarks=%d foods=%d characters=%d achievements=%d\n",
		seedCounts.Cities, seedCounts.CityTags, seedCounts.Landmarks, seedCounts.Foods, seedCounts.Characters, seedCounts.Achievements)
	fmt.Printf("seed city overlap: %d/%d already exist in database\n", overlap, len(seedCityNames))
	printNameList("seed cities missing from database", missingSeedCities)
	printNameList("database-only cities", dbOnlyCities)
	naturalKeyOverlaps, err := compareNaturalKeys(ctx, db, catalog)
	if err != nil {
		return err
	}
	for _, overlap := range naturalKeyOverlaps {
		fmt.Printf("seed %s overlap: %d/%d already exist in database\n", overlap.Label, overlap.Existing, overlap.Total)
		printNameList("seed "+overlap.Label+" missing from database", overlap.Missing)
	}
	fmt.Println("recommended: keep SEED_MODE=off for managed databases; use -mode bootstrap only for a reviewed one-time backfill")
	fmt.Println("dangerous: -mode sync overwrites rows matching seed natural keys and requires -confirm-overwrite")
	return nil
}

func countDBCatalog(ctx context.Context, db *gorm.DB) (catalogCounts, error) {
	var counts catalogCounts
	if err := countModel(ctx, db, &model.City{}, &counts.Cities); err != nil {
		return catalogCounts{}, err
	}
	if err := countModel(ctx, db, &model.CityTag{}, &counts.CityTags); err != nil {
		return catalogCounts{}, err
	}
	if err := countModel(ctx, db, &model.Landmark{}, &counts.Landmarks); err != nil {
		return catalogCounts{}, err
	}
	if err := countModel(ctx, db, &model.Food{}, &counts.Foods); err != nil {
		return catalogCounts{}, err
	}
	if err := countModel(ctx, db, &model.Character{}, &counts.Characters); err != nil {
		return catalogCounts{}, err
	}
	if err := countModel(ctx, db, &model.Achievement{}, &counts.Achievements); err != nil {
		return catalogCounts{}, err
	}
	return counts, nil
}

func countModel(ctx context.Context, db *gorm.DB, row any, dest *int64) error {
	return db.WithContext(ctx).Model(row).Count(dest).Error
}

func countSeedCatalog(catalog seed.Catalog) catalogCounts {
	counts := catalogCounts{
		Cities:       int64(len(catalog.Cities)),
		Achievements: int64(len(catalog.Achievements)),
	}
	for _, city := range catalog.Cities {
		counts.CityTags += int64(len(city.Tags))
		counts.Landmarks += int64(len(city.Landmarks))
		counts.Foods += int64(len(city.Foods))
		counts.Characters += int64(len(city.Characters))
	}
	return counts
}

func listDBCityNames(ctx context.Context, db *gorm.DB) ([]string, error) {
	var cities []model.City
	if err := db.WithContext(ctx).Find(&cities).Error; err != nil {
		return nil, err
	}
	names := make([]string, 0, len(cities))
	for _, city := range cities {
		names = append(names, city.Name)
	}
	sort.Strings(names)
	return names, nil
}

func stringSet(values []string) map[string]bool {
	set := make(map[string]bool, len(values))
	for _, value := range values {
		set[value] = true
	}
	return set
}

func missingFrom(values []string, allowed map[string]bool) []string {
	var missing []string
	for _, value := range values {
		if !allowed[value] {
			missing = append(missing, value)
		}
	}
	sort.Strings(missing)
	return missing
}

func printNameList(label string, names []string) {
	if len(names) == 0 {
		fmt.Printf("%s: none\n", label)
		return
	}
	shown := names
	if len(shown) > 20 {
		shown = shown[:20]
	}
	suffix := ""
	if len(names) > len(shown) {
		suffix = fmt.Sprintf(" ... (+%d more)", len(names)-len(shown))
	}
	fmt.Printf("%s (%d): %s%s\n", label, len(names), strings.Join(shown, ", "), suffix)
}

type naturalKeyOverlap struct {
	Label    string
	Existing int
	Total    int
	Missing  []string
}

func compareNaturalKeys(ctx context.Context, db *gorm.DB, catalog seed.Catalog) ([]naturalKeyOverlap, error) {
	dbKeys, err := dbNaturalKeys(ctx, db)
	if err != nil {
		return nil, err
	}
	seedKeys := seedNaturalKeys(catalog)

	labels := []string{"city_tags", "landmarks", "foods", "characters", "achievements"}
	overlaps := make([]naturalKeyOverlap, 0, len(labels))
	for _, label := range labels {
		missing := missingFrom(seedKeys[label], dbKeys[label])
		overlaps = append(overlaps, naturalKeyOverlap{
			Label:    label,
			Existing: len(seedKeys[label]) - len(missing),
			Total:    len(seedKeys[label]),
			Missing:  missing,
		})
	}
	return overlaps, nil
}

func seedNaturalKeys(catalog seed.Catalog) map[string][]string {
	keys := map[string][]string{
		"city_tags":    {},
		"landmarks":    {},
		"foods":        {},
		"characters":   {},
		"achievements": {},
	}
	for _, city := range catalog.Cities {
		for _, tag := range city.Tags {
			keys["city_tags"] = append(keys["city_tags"], scopedKey(city.Name, tag))
		}
		for _, landmark := range city.Landmarks {
			keys["landmarks"] = append(keys["landmarks"], scopedKey(city.Name, landmark.Name))
		}
		for _, food := range city.Foods {
			keys["foods"] = append(keys["foods"], scopedKey(city.Name, food.Name))
		}
		for _, character := range city.Characters {
			keys["characters"] = append(keys["characters"], scopedKey(city.Name, character.Name))
		}
	}
	for _, achievement := range catalog.Achievements {
		keys["achievements"] = append(keys["achievements"], achievement.Code)
	}
	for label := range keys {
		sort.Strings(keys[label])
	}
	return keys
}

func dbNaturalKeys(ctx context.Context, db *gorm.DB) (map[string]map[string]bool, error) {
	var cities []model.City
	if err := db.WithContext(ctx).Find(&cities).Error; err != nil {
		return nil, err
	}
	cityNamesByID := make(map[int64]string, len(cities))
	for _, city := range cities {
		cityNamesByID[city.ID] = city.Name
	}

	keys := map[string]map[string]bool{
		"city_tags":    {},
		"landmarks":    {},
		"foods":        {},
		"characters":   {},
		"achievements": {},
	}

	var tags []model.CityTag
	if err := db.WithContext(ctx).Find(&tags).Error; err != nil {
		return nil, err
	}
	for _, tag := range tags {
		if cityName := cityNamesByID[tag.CityID]; cityName != "" {
			keys["city_tags"][scopedKey(cityName, tag.Tag)] = true
		}
	}

	var landmarks []model.Landmark
	if err := db.WithContext(ctx).Find(&landmarks).Error; err != nil {
		return nil, err
	}
	for _, landmark := range landmarks {
		if cityName := cityNamesByID[landmark.CityID]; cityName != "" {
			keys["landmarks"][scopedKey(cityName, landmark.Name)] = true
		}
	}

	var foods []model.Food
	if err := db.WithContext(ctx).Find(&foods).Error; err != nil {
		return nil, err
	}
	for _, food := range foods {
		if cityName := cityNamesByID[food.CityID]; cityName != "" {
			keys["foods"][scopedKey(cityName, food.Name)] = true
		}
	}

	var characters []model.Character
	if err := db.WithContext(ctx).Find(&characters).Error; err != nil {
		return nil, err
	}
	for _, character := range characters {
		if cityName := cityNamesByID[character.CityID]; cityName != "" {
			keys["characters"][scopedKey(cityName, character.Name)] = true
		}
	}

	var achievements []model.Achievement
	if err := db.WithContext(ctx).Find(&achievements).Error; err != nil {
		return nil, err
	}
	for _, achievement := range achievements {
		keys["achievements"][achievement.Code] = true
	}
	return keys, nil
}

func scopedKey(scope string, name string) string {
	return scope + "/" + name
}
