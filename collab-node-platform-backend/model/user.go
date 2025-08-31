package model

import "time"

type User struct {
	UserID    string    `gorm:"column:user_id;type:char(36);primaryKey" json:"userId"`
	Username  string    `gorm:"column:username;type:varchar(64);uniqueIndex;not null" json:"username"`
	Password  string    `gorm:"column:password;type:varchar(255);not null" json:"-"` // 不序列化
	CreatedAt time.Time `gorm:"column:created_at;autoCreateTime" json:"createdAt"`
	UpdatedAt time.Time `gorm:"column:updated_at;autoUpdateTime" json:"updatedAt"`
}

func (User) TableName() string {
	return "users"
}
