package controller

import (
	"collab-node-platform-backend/config"
	"collab-node-platform-backend/service"
	"collab-node-platform-backend/utils"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool { return true },
}


func WSHandler(c *gin.Context) {
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		return
	}
	defer conn.Close()
	var ws *service.WsConn
	var joinedTask string
	for {
		var msg map[string]interface{}
		if err := conn.ReadJSON(&msg); err != nil {
			if joinedTask != "" && ws != nil {
				service.LeaveTaskRoom(joinedTask, ws)
			}
			break
		}
		event, _ := msg["event"].(string)
		if event == "joinTask" {
			token, _ := msg["token"].(string)
			claims, err := utils.ValidateJWT(token, config.AppConfig.JwtSecret)
			if err != nil {
				conn.WriteJSON(map[string]interface{}{"event": "error", "data": "Token无效或已过期"})
				conn.Close()
				break
			}
			ws = &service.WsConn{Conn: conn, UserID: claims.UserID}
			taskID, _ := msg["taskId"].(string)
			if taskID != "" {
				service.JoinTaskRoom(taskID, ws)
				joinedTask = taskID
			}
		}
	}
}
