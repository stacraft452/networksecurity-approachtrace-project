package service

import (
	"collab-node-platform-backend/db"
	"collab-node-platform-backend/model"
	"collab-node-platform-backend/utils"
	"errors"
	"strings"
)

// 权限校验辅助
func checkNodePermission(taskID, userID string, minRole int8) (bool, error) {
	var member model.TaskMember
	err := db.DB.Where("task_id = ? AND user_id = ?", taskID, userID).First(&member).Error
	if err != nil {
		return false, errors.New("无任务成员权限")
	}
	if member.Role > minRole {
		return false, errors.New("无操作权限")
	}
	return true, nil
}

// 删除节点（仅末尾节点可删，所有成员可操作）
// 删除节点及其所有子节点
func DeleteNode(nodeID, userID string) error {
	var node model.Node
	if err := db.DB.Where("node_id = ?", nodeID).First(&node).Error; err != nil {
		return errors.New("节点不存在")
	}
	// 权限校验：任务成员
	ok, err := checkNodePermission(node.TaskID, userID, 2)
	if !ok {
		return err
	}
	// 递归查找所有子节点
	var allIds []string
	var findDescendants func(id string)
	findDescendants = func(id string) {
		allIds = append(allIds, id)
		var children []model.Node
		db.DB.Where("parent_node_id = ?", id).Find(&children)
		for _, child := range children {
			findDescendants(child.NodeID)
		}
	}
	findDescendants(nodeID)
	// 批量删除
	if err := db.DB.Delete(&model.Node{}, "node_id IN (?)", allIds).Error; err != nil {
		return err
	}
	return nil
}

// 创建节点
func CreateNode(taskID, parentNodeID, nodeName, nodeContent, site, result, nextStep, creatorID string) (*model.Node, error) {
	ok, err := checkNodePermission(taskID, creatorID, 1)
	if !ok {
		return nil, err
	}
	node := &model.Node{
		NodeID:       utils.GenerateUUID(),
		TaskID:       taskID,
		ParentNodeID: nil,
		NodeName:     nodeName,
		NodeContent:  nodeContent,
		Site:         site,
		Result:       result,
		NextStep:     nextStep,
		CreatorID:    creatorID,
		Version:      1,
	}
	if parentNodeID != "" {
		node.ParentNodeID = &parentNodeID
	}
	if err := db.DB.Create(node).Error; err != nil {
		return nil, err
	}
	// 记录操作日志
	log := model.OperationLog{
		TaskID:           taskID,
		NodeID:           node.NodeID,
		UserID:           creatorID,
		OperationType:    "create",
		OperationContent: nodeContent,
	}
	_ = db.DB.Create(&log)
	// WebSocket 广播
	BroadcastNodeEvent(taskID, "newNode", node)
	return node, nil
}

// 编辑节点（乐观锁）
func EditNode(nodeID, nodeName, nodeContent, site, result, nextStep, editorID string, version int) (*model.Node, error) {
	var node model.Node
	if err := db.DB.Where("node_id = ?", nodeID).First(&node).Error; err != nil {
		return nil, errors.New("节点不存在")
	}
	// 权限校验：创建者或任务owner/editor
	ok, _ := checkNodePermission(node.TaskID, editorID, 1)
	if !ok && node.CreatorID != editorID {
		return nil, errors.New("无编辑权限")
	}
	// 乐观锁
	update := map[string]interface{}{
		"node_name":    nodeName,
		"node_content": nodeContent,
		"site":         site,
		"result":       result,
		"next_step":    nextStep,
		"version":      version + 1,
	}
	tx := db.DB.Model(&model.Node{}).Where("node_id = ? AND version = ?", nodeID, version).Updates(update)
	if tx.RowsAffected == 0 {
		return nil, errors.New("并发冲突，请刷新后重试")
	}
	// 记录操作日志
	log := model.OperationLog{
		TaskID:           node.TaskID,
		NodeID:           nodeID,
		UserID:           editorID,
		OperationType:    "edit",
		OperationContent: nodeContent,
	}
	_ = db.DB.Create(&log)
	// 返回最新节点
	_ = db.DB.Where("node_id = ?", nodeID).First(&node)
	// WebSocket 广播
	BroadcastNodeEvent(node.TaskID, "editNode", &node)
	return &node, nil
}

// 广播节点事件
func BroadcastNodeEvent(taskID, event string, node *model.Node) {
	BroadcastToTask(taskID, event, node)
}

// 查询节点列表
type NodeWithCreatorName struct {
	model.Node
	CreatorUsername string `json:"creatorUsername"`
}

func GetNodeList(taskID, userID string) ([]NodeWithCreatorName, error) {
	ok, err := checkNodePermission(taskID, userID, 2)
	if !ok {
		return nil, err
	}
	var nodes []model.Node
	err = db.DB.Where("task_id = ?", taskID).Find(&nodes).Error
	if err != nil {
		return nil, err
	}
	// 批量查用户名
	userIdSet := make(map[string]struct{})
	for _, n := range nodes {
		userIdSet[n.CreatorID] = struct{}{}
	}
	userIds := make([]string, 0, len(userIdSet))
	for id := range userIdSet {
		userIds = append(userIds, id)
	}
	var users []model.User
	db.DB.Where("user_id IN (?)", userIds).Find(&users)
	userMap := make(map[string]string)
	for _, u := range users {
		userMap[u.UserID] = u.Username
	}
	var result []NodeWithCreatorName
	for _, n := range nodes {
		result = append(result, NodeWithCreatorName{
			Node:            n,
			CreatorUsername: userMap[n.CreatorID],
		})
	}
	return result, nil
}

// 搜索节点
func SearchNode(taskID, keyword, userID string) ([]model.Node, error) {
	ok, err := checkNodePermission(taskID, userID, 2)
	if !ok {
		return nil, err
	}
	var nodes []model.Node
	kw := "%" + strings.TrimSpace(keyword) + "%"
	err = db.DB.Where("task_id = ? AND (node_name LIKE ? OR node_content LIKE ?)", taskID, kw, kw).Find(&nodes).Error
	return nodes, err
}
