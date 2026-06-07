package model

import "time"

type GuessAnswer struct {
	ID            int64     `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
	ChallengeCode string    `gorm:"column:challenge_code;type:varchar(16);not null;index:idx_guess_answers_challenge_time,priority:1" json:"challenge_code"`
	AnswerText    string    `gorm:"column:answer_text;type:varchar(64);not null" json:"answer_text"`
	IsCorrect     bool      `gorm:"column:is_correct;not null" json:"is_correct"`
	CreatedAt     time.Time `gorm:"column:created_at;autoCreateTime;index:idx_guess_answers_challenge_time,priority:2" json:"created_at"`
}

func (GuessAnswer) TableName() string { return "guess_answers" }
