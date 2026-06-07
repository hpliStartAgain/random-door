package achievement

import "testing"

func TestMatch(t *testing.T) {
	tests := []struct {
		name      string
		ruleType  string
		ruleValue string
		stats     UserStats
		want      bool
	}{
		{
			name:     "first_checkin unlocks after one checkin",
			ruleType: "first_checkin", stats: UserStats{CheckinCount: 1}, want: true,
		},
		{
			name:     "first_checkin locked with zero checkins",
			ruleType: "first_checkin", stats: UserStats{CheckinCount: 0}, want: false,
		},
		{
			name:     "checkin_count reached",
			ruleType: "checkin_count", ruleValue: "3", stats: UserStats{CheckinCount: 3}, want: true,
		},
		{
			name:     "checkin_count not reached",
			ruleType: "checkin_count", ruleValue: "5", stats: UserStats{CheckinCount: 4}, want: false,
		},
		{
			name:     "checkin_count invalid value",
			ruleType: "checkin_count", ruleValue: "abc", stats: UserStats{CheckinCount: 9}, want: false,
		},
		{
			name:     "city_tag any matched",
			ruleType: "city_tag", ruleValue: "沿海",
			stats: UserStats{VisitedTagAny: map[string]bool{"沿海": true}}, want: true,
		},
		{
			name:     "city_tag missing",
			ruleType: "city_tag", ruleValue: "大漠",
			stats: UserStats{VisitedTagAny: map[string]bool{"沿海": true}}, want: false,
		},
		{
			name:     "tag_count reached",
			ruleType: "tag_count", ruleValue: "美食:3",
			stats: UserStats{VisitedCityTags: map[string]int{"美食": 3}}, want: true,
		},
		{
			name:     "tag_count not reached",
			ruleType: "tag_count", ruleValue: "美食:3",
			stats: UserStats{VisitedCityTags: map[string]int{"美食": 2}}, want: false,
		},
		{
			name:     "tag_count malformed",
			ruleType: "tag_count", ruleValue: "美食",
			stats: UserStats{VisitedCityTags: map[string]int{"美食": 9}}, want: false,
		},
		{
			name:     "tag_count non-numeric",
			ruleType: "tag_count", ruleValue: "美食:x",
			stats: UserStats{VisitedCityTags: map[string]int{"美食": 9}}, want: false,
		},
		{
			name:     "visit_count reached",
			ruleType: "visit_count", ruleValue: "5", stats: UserStats{VisitedCityCount: 5}, want: true,
		},
		{
			name:     "visit_count not reached",
			ruleType: "visit_count", ruleValue: "5", stats: UserStats{VisitedCityCount: 4}, want: false,
		},
		{
			name:     "game_visit_count reached",
			ruleType: "game_visit_count", ruleValue: "5", stats: UserStats{GameVisitCount: 6}, want: true,
		},
		{
			name:     "game_visit_count not reached",
			ruleType: "game_visit_count", ruleValue: "5", stats: UserStats{GameVisitCount: 4}, want: false,
		},
		{
			name:     "dice_distance reached",
			ruleType: "dice_distance", ruleValue: "1200", stats: UserStats{MaxDiceDistance: 1200}, want: true,
		},
		{
			name:     "dice_distance not reached",
			ruleType: "dice_distance", ruleValue: "1200", stats: UserStats{MaxDiceDistance: 800}, want: false,
		},
		{
			name:     "dice_direction streak reached",
			ruleType: "dice_direction", ruleValue: "北:3",
			stats: UserStats{MaxSameDirRun: map[string]int{"北": 3}}, want: true,
		},
		{
			name:     "dice_direction streak short",
			ruleType: "dice_direction", ruleValue: "北:3",
			stats: UserStats{MaxSameDirRun: map[string]int{"北": 2}}, want: false,
		},
		{
			name:     "dice_direction malformed",
			ruleType: "dice_direction", ruleValue: "北",
			stats: UserStats{MaxSameDirRun: map[string]int{"北": 9}}, want: false,
		},
		{
			name:     "unknown rule type",
			ruleType: "warp_speed", ruleValue: "9", stats: UserStats{}, want: false,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			if got := Match(tc.ruleType, tc.ruleValue, tc.stats); got != tc.want {
				t.Fatalf("Match(%q,%q) = %v, want %v", tc.ruleType, tc.ruleValue, got, tc.want)
			}
		})
	}
}
