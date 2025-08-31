package model

import "time"

type Node struct {
	NodeID       string    `gorm:"column:node_id;type:char(36);primaryKey" json:"nodeId"`
	TaskID       string    `gorm:"column:task_id;type:char(36);not null;index:idx_task_id" json:"taskId"`
	ParentNodeID *string   `gorm:"column:parent_node_id;type:char(36);default:null" json:"parentNodeId"`
	NodeName     string    `gorm:"column:node_name;type:varchar(128);not null" json:"nodeName"`
	NodeContent  string    `gorm:"column:node_content;type:text" json:"nodeContent"`
	Site         string    `gorm:"column:site;type:varchar(255);default:null" json:"site"`
	Result       string    `gorm:"column:result;type:text;default:null" json:"result"`
	NextStep     string    `gorm:"column:next_step;type:text;default:null" json:"nextStep"`
	CreatorID    string    `gorm:"column:creator_id;type:char(36);not null" json:"creatorId"`
	Version      int       `gorm:"column:version;type:int;default:1" json:"version"` // 乐观锁
	CreatedAt    time.Time `gorm:"column:created_at;autoCreateTime" json:"createdAt"`
	UpdatedAt    time.Time `gorm:"column:updated_at;autoUpdateTime" json:"updatedAt"`
}

func (Node) TableName() string {
	return "nodes"
}
