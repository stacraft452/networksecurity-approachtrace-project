package model

import "time"

type Task struct {
	TaskID    string    `gorm:"column:task_id;type:char(36);primaryKey" json:"taskId"`
	TaskName  string    `gorm:"column:task_name;type:varchar(128);not null" json:"taskName"`
	CreatorID string    `gorm:"column:creator_id;type:char(36);not null" json:"creatorId"`
	Status    int8      `gorm:"column:status;type:tinyint;default:0" json:"status"`
	CreatedAt time.Time `gorm:"column:created_at;autoCreateTime" json:"createdAt"`
	UpdatedAt time.Time `gorm:"column:updated_at;autoUpdateTime" json:"updatedAt"`
}

func (Task) TableName() string {
	return "tasks"
}
