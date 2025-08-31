package service

import (
	"sync"

	"github.com/gorilla/websocket"
)

type WsConn struct {
	Conn   *websocket.Conn
	UserID string
}

var (
	taskRooms   = make(map[string]map[*WsConn]struct{}) // taskID -> 连接集合
	taskRoomsMu sync.RWMutex
)

// 加入房间
func JoinTaskRoom(taskID string, conn *WsConn) {
	taskRoomsMu.Lock()
	defer taskRoomsMu.Unlock()
	if taskRooms[taskID] == nil {
		taskRooms[taskID] = make(map[*WsConn]struct{})
	}
	taskRooms[taskID][conn] = struct{}{}
}

// 离开房间
func LeaveTaskRoom(taskID string, conn *WsConn) {
	taskRoomsMu.Lock()
	defer taskRoomsMu.Unlock()
	if m, ok := taskRooms[taskID]; ok {
		delete(m, conn)
		if len(m) == 0 {
			delete(taskRooms, taskID)
		}
	}
}

// 广播事件
func BroadcastToTask(taskID, event string, data interface{}) {
	taskRoomsMu.RLock()
	conns := taskRooms[taskID]
	taskRoomsMu.RUnlock()
	for c := range conns {
		msg := map[string]interface{}{"event": event, "data": data}
		_ = c.Conn.WriteJSON(msg)
	}
}
