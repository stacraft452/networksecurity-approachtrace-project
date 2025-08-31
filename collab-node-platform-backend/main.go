package main

import (
	"log"

	"github.com/gin-gonic/gin"

	"collab-node-platform-backend/config"
	"collab-node-platform-backend/controller"
	"collab-node-platform-backend/db"
	"collab-node-platform-backend/middleware"
)

func main() {
	// 初始化配置和数据库
	config.InitConfig()
	db.InitDB()

	r := gin.Default()
	r.Use(middleware.CorsMiddleware())

	// 公开路由
	userGroup := r.Group("/api/user")
	{
		userGroup.POST("/login", controller.LoginHandler)
		userGroup.POST("/register", controller.RegisterHandler)
	}

	// 需认证路由
	apiGroup := r.Group("/api")
	apiGroup.Use(middleware.AuthMiddleware())
	{
		taskGroup := apiGroup.Group("/task")
		taskGroup.POST("/create", controller.CreateTaskHandler)
		taskGroup.GET("/list", controller.GetTaskListHandler)
		taskGroup.POST("/delete", controller.DeleteTaskHandler)
		taskGroup.POST("/start", controller.StartTaskHandler)
		taskGroup.POST("/finish", controller.FinishTaskHandler)

		nodeGroup := apiGroup.Group("/node")
		nodeGroup.POST("/create", controller.CreateNodeHandler)
		nodeGroup.PUT("/edit", controller.EditNodeHandler)
		nodeGroup.GET("/list", controller.GetNodeListHandler)
		nodeGroup.GET("/search", controller.SearchNodeHandler)
		nodeGroup.POST("/delete", controller.DeleteNodeHandler)
	}

	// WebSocket 路由
	r.GET(config.AppConfig.WsPath, controller.WSHandler)

	addr := "0.0.0.0:" + config.AppConfig.Port
	log.Printf("server running on %s", addr)
	err := r.Run(addr)
	if err != nil {
		log.Fatalf("服务启动失败: %v", err)
	}
}
