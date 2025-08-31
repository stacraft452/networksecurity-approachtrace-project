package model

import "time"

type OperationLog struct {
	ID               int64     `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
	TaskID           string    `gorm:"column:task_id;type:char(36);not null;index:idx_task_id" json:"taskId"`
	NodeID           string    `gorm:"column:node_id;type:char(36);not null;index:idx_node_id" json:"nodeId"`
	UserID           string    `gorm:"column:user_id;type:char(36);not null;index:idx_user_id" json:"userId"`
	OperationType    string    `gorm:"column:operation_type;type:enum('create','edit','delete');not null" json:"operationType"`
	OperationContent string    `gorm:"column:operation_content;type:text" json:"operationContent"`
	CreatedAt        time.Time `gorm:"column:created_at;autoCreateTime" json:"createdAt"`
}

func (OperationLog) TableName() string {
	return "operation_logs"
}
