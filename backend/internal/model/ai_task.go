package model

import "time"

type AITask struct {
	ID        int64     `gorm:"column:id;primaryKey;autoIncrement" json:"task_id"`
	UserID    int64     `gorm:"column:user_id;not null;index:idx_ai_tasks_user_time,priority:1" json:"user_id"`
	Type      string    `gorm:"column:type;type:varchar(32);not null;index:idx_ai_tasks_type_status,priority:1" json:"type"`
	Status    string    `gorm:"column:status;type:varchar(32);not null;index:idx_ai_tasks_status_time,priority:1;index:idx_ai_tasks_type_status,priority:2" json:"status"`
	InputJSON string    `gorm:"column:input_json;type:json;not null" json:"-"`
	ResultURL *string   `gorm:"column:result_url;type:varchar(512)" json:"result_url,omitempty"`
	Error     *string   `gorm:"column:error;type:text" json:"error,omitempty"`
	Attempts  int       `gorm:"column:attempts;not null;default:0" json:"attempts"`
	CreatedAt time.Time `gorm:"column:created_at;autoCreateTime;index:idx_ai_tasks_user_time,priority:2" json:"created_at"`
	UpdatedAt time.Time `gorm:"column:updated_at;autoUpdateTime;index:idx_ai_tasks_status_time,priority:2" json:"updated_at"`
}

func (AITask) TableName() string { return "ai_tasks" }
