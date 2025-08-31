package model

import "time"

type TaskMember struct {
	ID        int64     `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
	TaskID    string    `gorm:"column:task_id;type:char(36);not null;index:idx_task_id" json:"taskId"`
	UserID    string    `gorm:"column:user_id;type:char(36);not null;index:idx_user_id" json:"userId"`
	Role      int8      `gorm:"column:role;type:tinyint;not null" json:"role"`
	CreatedAt time.Time `gorm:"column:created_at;autoCreateTime" json:"createdAt"`
}

func (TaskMember) TableName() string {
	return "task_members"
}
