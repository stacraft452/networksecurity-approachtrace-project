package controller
import (
	"collab-node-platform-backend/service"
	"net/http"

	"github.com/gin-gonic/gin"
)

type CreateTaskRequest struct {
	TaskName string              `json:"taskName" binding:"required"`
	Members  []service.MemberDTO `json:"members"`
}
type DeleteTaskRequest struct {
	TaskID string `json:"taskId" binding:"required"`
}
type TaskStatusRequest struct {
	TaskID string `json:"taskId" binding:"required"`
}

// 开启任务接口
func StartTaskHandler(c *gin.Context) {
	var req TaskStatusRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "参数错误"})
		return
	}
	userId, ok := c.Get("userId")
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "未认证"})
		return
	}
	err := service.StartTask(req.TaskID, userId.(string))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"msg": "任务已开启"})
}
// 查询任务列表接口
func GetTaskListHandler(c *gin.Context) {
	userId, ok := c.Get("userId")
	if !ok {
		c.JSON(401, gin.H{"error": "未认证"})
		return
	}
	tasks, err := service.GetTaskList(userId.(string))
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, tasks)
}
// 结束任务接口
func FinishTaskHandler(c *gin.Context) {
	var req TaskStatusRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "参数错误"})
		return
	}
	userId, ok := c.Get("userId")
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "未认证"})
		return
	}
	err := service.FinishTask(req.TaskID, userId.(string))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"msg": "任务已结束"})
}
// 删除任务接口
func DeleteTaskHandler(c *gin.Context) {
	var req DeleteTaskRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "参数错误"})
		return
	}
	userId, ok := c.Get("userId")
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "未认证"})
		return
	}
	err := service.DeleteTask(req.TaskID, userId.(string))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"msg": "删除成功"})
}
func CreateTaskHandler(c *gin.Context) {
	var req CreateTaskRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "参数错误"})
		return
	}
	userId, ok := c.Get("userId")
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "未认证"})
		return
	}
	taskID, err := service.CreateTask(req.TaskName, userId.(string), req.Members)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"taskId": taskID})
}
