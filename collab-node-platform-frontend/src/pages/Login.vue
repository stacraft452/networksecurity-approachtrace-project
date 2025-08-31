<template>
  <div class="login-container">
    <el-card class="login-card">
      <h2>多人协作节点平台登录</h2>
      <el-form :model="form" @submit.prevent="onLogin">
        <el-form-item label="用户名">
          <el-input v-model="form.username" autocomplete="username" />
        </el-form-item>
        <el-form-item label="密码">
          <el-input v-model="form.password" type="password" autocomplete="current-password" />
        </el-form-item>
        <el-form-item>
          <el-button type="primary" @click="onLogin" :loading="loading">登录</el-button>
          <el-button type="text" @click="$router.push('/register')">注册账号</el-button>
        </el-form-item>
      </el-form>
    </el-card>
  </div>
</template>

<script setup>
import { ref } from 'vue'
import { useRouter } from 'vue-router'
import { useUserStore } from '../store/user'
import { login } from '../api'
import { ElMessage } from 'element-plus'

const form = ref({ username: '', password: '' })
const loading = ref(false)
const router = useRouter()
const userStore = useUserStore()

const onLogin = async () => {
  if (!form.value.username || !form.value.password) {
    ElMessage.error('请输入用户名和密码')
    return
  }
  loading.value = true
  try {
    const res = await login(form.value.username, form.value.password)
    userStore.setUser(res.token, res.userId, form.value.username)
    ElMessage.success('登录成功')
    router.push('/')
  } catch (e) {}
  loading.value = false
}
</script>

<style scoped>
.login-container {
  display: flex;
  height: 100vh;
  align-items: center;
  justify-content: center;
  background: #f5f6fa;
}
.login-card {
  width: 350px;
  padding: 30px 20px 20px 20px;
}
</style>
