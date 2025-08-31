
let ws = null

export function connectWS(taskId, onEvent) {
  if (ws) ws.close()
  const wsProtocol = window.location.protocol === 'https:' ? 'wss' : 'ws'
  const wsUrl = `${wsProtocol}://${window.location.hostname}:8087/ws`
  ws = new WebSocket(wsUrl)
  ws.onopen = () => {
    const token = document.cookie.match(/token=([^;]+)/)?.[1]
    // 只有连接建立后才发送
    if (ws && ws.readyState === 1) {
      ws.send(JSON.stringify({ event: 'joinTask', taskId, token }))
    }
  }
  ws.onmessage = (e) => {
    try {
      const msg = JSON.parse(e.data)
      if (onEvent) onEvent(msg.event, msg.data)
    } catch {}
  }
  ws.onclose = () => { ws = null }
}

export function disconnectWS() {
  if (ws) ws.close()
  ws = null
}
