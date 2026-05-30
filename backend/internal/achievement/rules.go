package achievement

import (
	"log/slog"
	"strconv"
	"strings"
)

// UserStats is a snapshot of a user's activity, built once per evaluation.
type UserStats struct {
	CheckinCount    int
	CheckinCityTags map[string]int  // tag -> distinct city count
	CheckinTagAny   map[string]bool // whether user checked in any city with this tag
	GameVisitCount  int             // distinct cities in game mode
	MaxDiceDistance int             // max single dice distance
	MaxSameDirRun   map[string]int  // direction -> max consecutive count
}

// Match checks if a single achievement rule is satisfied.
func Match(ruleType, ruleValue string, s UserStats) bool {
	switch ruleType {
	case "first_checkin":
		return s.CheckinCount >= 1

	case "checkin_count":
		n, err := strconv.Atoi(ruleValue)
		if err != nil {
			slog.Warn("invalid rule_value for checkin_count", "value", ruleValue)
			return false
		}
		return s.CheckinCount >= n

	case "city_tag":
		return s.CheckinTagAny[ruleValue]

	case "tag_count":
		parts := strings.SplitN(ruleValue, ":", 2)
		if len(parts) != 2 {
			slog.Warn("invalid rule_value for tag_count", "value", ruleValue)
			return false
		}
		tag := parts[0]
		n, err := strconv.Atoi(parts[1])
		if err != nil {
			slog.Warn("invalid count in tag_count", "value", ruleValue)
			return false
		}
		return s.CheckinCityTags[tag] >= n

	case "game_visit_count":
		n, err := strconv.Atoi(ruleValue)
		if err != nil {
			slog.Warn("invalid rule_value for game_visit_count", "value", ruleValue)
			return false
		}
		return s.GameVisitCount >= n

	case "dice_distance":
		n, err := strconv.Atoi(ruleValue)
		if err != nil {
			slog.Warn("invalid rule_value for dice_distance", "value", ruleValue)
			return false
		}
		return s.MaxDiceDistance >= n

	case "dice_direction":
		parts := strings.SplitN(ruleValue, ":", 2)
		if len(parts) != 2 {
			slog.Warn("invalid rule_value for dice_direction", "value", ruleValue)
			return false
		}
		dir := parts[0]
		n, err := strconv.Atoi(parts[1])
		if err != nil {
			slog.Warn("invalid count in dice_direction", "value", ruleValue)
			return false
		}
		return s.MaxSameDirRun[dir] >= n

	default:
		slog.Warn("unknown rule_type", "type", ruleType)
		return false
	}
}
