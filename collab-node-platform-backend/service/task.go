package service

import (
	"collab-node-platform-backend/db"
	"collab-node-platform-backend/model"
	"collab-node-platform-backend/utils"
	"errors"
)

type MemberDTO struct {
	Username string `json:"username"`
	Role     int8   `json:"role"`
}

// 删除任务（仅创建者可删）
func DeleteTask(taskID, userID string) error {
	var task model.Task
	if err := db.DB.Where("task_id = ?", taskID).First(&task).Error; err != nil {
		return errors.New("任务不存在")
	}
	if task.CreatorID != userID {
		return errors.New("只有任务创建者可以删除任务")
	}
	// 事务删除任务、成员、节点
	tx := db.DB.Begin()
	if err := tx.Where("task_id = ?", taskID).Delete(&model.TaskMember{}).Error; err != nil {
		tx.Rollback()
		return err
	}
	if err := tx.Where("task_id = ?", taskID).Delete(&model.Node{}).Error; err != nil {
		tx.Rollback()
		return err
	}
	if err := tx.Where("task_id = ?", taskID).Delete(&model.Task{}).Error; err != nil {
		tx.Rollback()
		return err
	}
	tx.Commit()
	// WebSocket 广播任务删除
	BroadcastToTask(taskID, "deleteTask", map[string]interface{}{"taskId": taskID})
	return nil
}

// 查询用户参与的任务列表
func GetTaskList(userID string) ([]model.Task, error) {
	// 查询用户参与的所有任务id
	var taskIds []string
	err := db.DB.Model(&model.TaskMember{}).Where("user_id = ?", userID).Pluck("task_id", &taskIds).Error
	if err != nil {
		return nil, err
	}
	if len(taskIds) == 0 {
		return []model.Task{}, nil
	}
	var tasks []model.Task
	err = db.DB.Where("task_id IN (?)", taskIds).Find(&tasks).Error
	return tasks, err
}

// 开启任务（仅创建者可操作）
func StartTask(taskID, userID string) error {
	var task model.Task
	if err := db.DB.Where("task_id = ?", taskID).First(&task).Error; err != nil {
		return errors.New("任务不存在")
	}
	if task.CreatorID != userID {
		return errors.New("只有任务创建者可以开启任务")
	}
	if task.Status != 0 {
		return errors.New("任务已开启或已结束")
	}
	err := db.DB.Model(&task).Update("status", 1).Error
	if err == nil {
		// WebSocket 广播任务状态变更
		BroadcastToTask(taskID, "updateTask", map[string]interface{}{"taskId": taskID, "status": 1})
	}
	return err
}

// 结束任务（仅创建者可操作）
func FinishTask(taskID, userID string) error {
	var task model.Task
	if err := db.DB.Where("task_id = ?", taskID).First(&task).Error; err != nil {
		return errors.New("任务不存在")
	}
	if task.CreatorID != userID {
		return errors.New("只有任务创建者可以结束任务")
	}
	if task.Status != 1 {
		return errors.New("任务未处于进行中，无法结束")
	}
	err := db.DB.Model(&task).Update("status", 2).Error
	if err == nil {
		// WebSocket 广播任务状态变更
		BroadcastToTask(taskID, "updateTask", map[string]interface{}{"taskId": taskID, "status": 2})
	}
	return err
}

func CreateTask(taskName string, creatorID string, members []MemberDTO) (string, error) {
	tx := db.DB.Begin()
	if tx.Error != nil {
		return "", tx.Error
	}
	taskID := utils.GenerateUUID()
	// 创建任务
	task := model.Task{
		TaskID:    taskID,
		TaskName:  taskName,
		CreatorID: creatorID,
		Status:    0,
	}
	if err := tx.Create(&task).Error; err != nil {
		tx.Rollback()
		return "", err
	}
	// 创建成员（包含创建者），防止重复
	memberMap := make(map[string]struct{})
	var taskMembers []model.TaskMember
	// 创建者
	memberMap[creatorID] = struct{}{}
	taskMembers = append(taskMembers, model.TaskMember{
		TaskID: taskID, UserID: creatorID, Role: 0,
	})
	// 其他成员
	for _, m := range members {
		var user model.User
		if err := tx.Where("username = ?", m.Username).First(&user).Error; err != nil {
			tx.Rollback()
			return "", errors.New("成员 " + m.Username + " 不存在")
		}
		if _, exists := memberMap[user.UserID]; exists {
			continue // 已添加过该成员，跳过
		}
		memberMap[user.UserID] = struct{}{}
		taskMembers = append(taskMembers, model.TaskMember{
			TaskID: taskID, UserID: user.UserID, Role: m.Role,
		})
	}
	if len(taskMembers) > 0 {
		if err := tx.Create(&taskMembers).Error; err != nil {
			tx.Rollback()
			return "", err
		}
	}
	// 创建根节点
	rootNode := model.Node{
		NodeID:       utils.GenerateUUID(),
		TaskID:       taskID,
		ParentNodeID: nil,
		NodeName:     taskName,
		NodeContent:  "",
		CreatorID:    creatorID,
		Version:      1,
	}
	if err := tx.Create(&rootNode).Error; err != nil {
		tx.Rollback()
		return "", err
	}
	tx.Commit()
	// WebSocket 广播新任务
	BroadcastToTask(taskID, "newTask", map[string]interface{}{"taskId": taskID, "taskName": taskName, "creatorId": creatorID})
	return taskID, nil
}
