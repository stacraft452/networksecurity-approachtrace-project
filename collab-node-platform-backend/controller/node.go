package controller
import (
	"collab-node-platform-backend/service"
	"net/http"

	"github.com/gin-gonic/gin"
)

type CreateNodeRequest struct {
	TaskID       string `json:"taskId" binding:"required"`
	ParentNodeID string `json:"parentNodeId"`
	NodeName     string `json:"nodeName" binding:"required"`
	NodeContent  string `json:"nodeContent"`
	Site         string `json:"site"`
	Result       string `json:"result"`
	NextStep     string `json:"nextStep"`
}

type EditNodeRequest struct {
	NodeID      string `json:"nodeId" binding:"required"`
	NodeName    string `json:"nodeName" binding:"required"`
	NodeContent string `json:"nodeContent"`
	Site        string `json:"site"`
	Result      string `json:"result"`
	NextStep    string `json:"nextStep"`
	Version     int    `json:"version" binding:"required"`
}

type DeleteNodeRequest struct {
	NodeID string `json:"nodeId" binding:"required"`
}

// 删除节点接口
func DeleteNodeHandler(c *gin.Context) {
	var req DeleteNodeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "参数错误"})
		return
	}
	userId, ok := c.Get("userId")
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "未认证"})
		return
	}
	err := service.DeleteNode(req.NodeID, userId.(string))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"msg": "删除成功"})
}

func CreateNodeHandler(c *gin.Context) {
	var req CreateNodeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "参数错误"})
		return
	}
	userId, _ := c.Get("userId")
	node, err := service.CreateNode(req.TaskID, req.ParentNodeID, req.NodeName, req.NodeContent, req.Site, req.Result, req.NextStep, userId.(string))
	if err != nil {
		c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, node)
}

func EditNodeHandler(c *gin.Context) {
	var req EditNodeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "参数错误"})
		return
	}
	userId, _ := c.Get("userId")
	node, err := service.EditNode(req.NodeID, req.NodeName, req.NodeContent, req.Site, req.Result, req.NextStep, userId.(string), req.Version)
	if err != nil {
		c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, node)
}

func GetNodeListHandler(c *gin.Context) {
	taskId := c.Query("taskId")
	if taskId == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "缺少taskId"})
		return
	}
	userId, _ := c.Get("userId")
	nodes, err := service.GetNodeList(taskId, userId.(string))
	if err != nil {
		c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, nodes)
}

func SearchNodeHandler(c *gin.Context) {
	taskId := c.Query("taskId")
	keyword := c.Query("keyword")
	if taskId == "" || keyword == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "缺少参数"})
		return
	}
	userId, _ := c.Get("userId")
	nodes, err := service.SearchNode(taskId, keyword, userId.(string))
	if err != nil {
		c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, nodes)
}
