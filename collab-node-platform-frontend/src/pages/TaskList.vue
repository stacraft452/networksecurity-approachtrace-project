<template>
  <div class="task-list-page">
    <el-row justify="space-between" align="middle" style="margin-bottom: 20px;">
      <el-col><h2>我的任务</h2></el-col>
      <el-col>
        <el-button type="primary" @click="showCreate = true">新建任务</el-button>
        <el-button style="margin-left:10px;" @click="logout">退出登录</el-button>
      </el-col>
    </el-row>
    <el-table :data="tasks" style="width: 100%">
      <el-table-column prop="taskName" label="任务名称" />
  <el-table-column prop="status" label="状态" :formatter="statusFormatter" />
      <el-table-column label="操作">
        <template #default="{ row }">
          <el-button type="primary" size="small" @click="goDetail(row)">进入</el-button>
          <el-button
            v-if="row.creatorId === userStore.userId && row.status === 0"
            type="success"
            size="small"
            style="margin-left:8px;"
            @click="onStartTask(row)"
          >开启</el-button>
          <el-button
            v-if="row.creatorId === userStore.userId && row.status === 1"
            type="warning"
            size="small"
            style="margin-left:8px;"
            @click="onFinishTask(row)"
          >结束</el-button>
          <el-button
            v-if="row.creatorId === userStore.userId"
            type="danger"
            size="small"
            style="margin-left:8px;"
            @click="onDeleteTask(row)"
          >删除</el-button>
        </template>
      </el-table-column>
    </el-table>
    <el-dialog v-model="showCreate" title="新建任务">
      <el-form :model="createForm">
        <el-form-item label="任务名称">
          <el-input v-model="createForm.taskName" />
        </el-form-item>
        <el-form-item label="协作成员(用户名,逗号分隔)">
          <el-input v-model="createForm.members" placeholder="user1,user2" />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="showCreate = false">取消</el-button>
        <el-button type="primary" @click="onCreate">创建</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup>
// 任务状态格式化
function statusFormatter(row) {
  if (row.status === 0) return '任务未正式开启'
  if (row.status === 1) return '正式开启'
  if (row.status === 2) return '已结束'
  return '未知状态'
}
const logout = () => {
  userStore.logout()
  router.replace('/login')
}
import { ref, onMounted, onUnmounted } from 'vue'
import { useRouter } from 'vue-router'
import { createTask, deleteTask, startTask, finishTask, getTaskList } from '../api'
import { connectWS, disconnectWS } from '../utils/ws'
const onStartTask = async (row) => {
  if (!row.taskId) return
  try {
    await startTask(row.taskId)
    ElMessage.success('任务已开启')
    fetchTasks()
  } catch (e) {}
}

const onFinishTask = async (row) => {
  if (!row.taskId) return
  try {
    await finishTask(row.taskId)
    ElMessage.success('任务已结束')
    // 本地模拟：直接修改状态为2
    const t = tasks.value.find(t => t.taskId === row.taskId)
    if (t) t.status = 2
    fetchTasks()
  } catch (e) {}
}
import { ElMessage } from 'element-plus'
import { useUserStore } from '../store/user'

const tasks = ref([])
const showCreate = ref(false)
const createForm = ref({ taskName: '', members: '' })
const router = useRouter()
const userStore = useUserStore()

const fetchTasks = async () => {
  try {
    const res = await getTaskList()
    tasks.value = res
  } catch (e) {
    tasks.value = []
  }
}

const onDeleteTask = async (row) => {
  if (!row.taskId) return
  try {
    await deleteTask(row.taskId)
    ElMessage.success('删除成功')
    fetchTasks()
  } catch (e) {}
}

const goDetail = row => {
  router.push(`/task/${row.taskId}`)
}

const onCreate = async () => {
  if (!createForm.value.taskName) {
    ElMessage.error('请输入任务名称')
    return
  }
  const members = createForm.value.members
    .split(',')
    .map(u => u.trim())
    .filter(Boolean)
    .map(username => ({ username, role: 1 }))
  try {
    const res = await createTask(createForm.value.taskName, members)
    ElMessage.success('创建成功')
    // 本地模拟任务列表
    const newTask = { taskId: res.taskId, taskName: createForm.value.taskName, status: 0 }
    const arr = JSON.parse(localStorage.getItem('tasks') || '[]')
    arr.push(newTask)
    localStorage.setItem('tasks', JSON.stringify(arr))
    fetchTasks()
    showCreate.value = false
    createForm.value = { taskName: '', members: '' }
  } catch (e) {}
}

onMounted(() => {
  fetchTasks()
  connectWS('all', (event, data) => {
    if (["newTask", "deleteTask", "updateTask"].includes(event)) {
      fetchTasks()
    }
  })
})
onUnmounted(() => {
  disconnectWS()
})
</script>

<style scoped>
.task-list-page {
  max-width: 800px;
  margin: 40px auto;
  background: #fff;
  padding: 30px;
  border-radius: 8px;
  box-shadow: 0 2px 8px #eee;
}
</style>
