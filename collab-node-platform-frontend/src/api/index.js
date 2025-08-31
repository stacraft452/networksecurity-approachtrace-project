export const getTaskList = () =>
  request.get('/task/list')
export const startTask = (taskId) =>
  request.post('/task/start', { taskId })

export const finishTask = (taskId) =>
  request.post('/task/finish', { taskId })
export const deleteNode = (nodeId) =>
  request.post('/node/delete', { nodeId })
export const deleteTask = (taskId) =>
  request.post('/task/delete', { taskId })
import request from './request'

export const login = (username, password) =>
  request.post('/user/login', { username, password })

export const register = (username, password) =>
  request.post('/user/register', { username, password })

export const createTask = (taskName, members) =>
  request.post('/task/create', { taskName, members })

export const getNodeList = (taskId) =>
  request.get('/node/list', { params: { taskId } })

export const createNode = (data) =>
  request.post('/node/create', data)

export const editNode = (data) =>
  request.put('/node/edit', data)

export const searchNode = (taskId, keyword) =>
  request.get('/node/search', { params: { taskId, keyword } })
