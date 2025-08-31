<template>
  <div class="register-container">
    <el-card class="register-card">
      <h2>用户注册</h2>
      <el-form :model="form" :rules="rules" ref="formRef" label-width="80px">
        <el-form-item label="用户名" prop="username">
          <el-input v-model="form.username" autocomplete="off" />
        </el-form-item>
        <el-form-item label="密码" prop="password">
          <el-input v-model="form.password" type="password" autocomplete="off" />
        </el-form-item>
        <el-form-item label="确认密码" prop="confirm">
          <el-input v-model="form.confirm" type="password" autocomplete="off" />
        </el-form-item>
        <el-form-item>
          <el-button type="primary" @click="onRegister">注册</el-button>
          <el-button type="text" @click="$router.push('/login')">返回登录</el-button>
        </el-form-item>
      </el-form>
    </el-card>
  </div>
</template>

<script setup>
import { ref } from 'vue'
import { useRouter } from 'vue-router'
import { ElMessage } from 'element-plus'
import { register } from '../api'

const router = useRouter()
const formRef = ref()
const form = ref({ username: '', password: '', confirm: '' })
const rules = {
  username: [ { required: true, message: '请输入用户名', trigger: 'blur' } ],
  password: [ { required: true, message: '请输入密码', trigger: 'blur' } ],
  confirm: [
    { required: true, message: '请确认密码', trigger: 'blur' },
    { validator: (rule, value, callback) => {
        if (value !== form.value.password) {
          callback(new Error('两次密码不一致'))
        } else {
          callback()
        }
      }, trigger: 'blur' }
  ]
}

const onRegister = () => {
  formRef.value.validate(async valid => {
    if (!valid) return
    try {
      await register(form.value.username, form.value.password)
      ElMessage.success('注册成功，请登录')
      router.push('/login')
    } catch (e) {}
  })
}
</script>

<style scoped>
.register-container {
  display: flex;
  justify-content: center;
  align-items: center;
  height: 100vh;
  background: #f5f6fa;
}
.register-card {
  width: 400px;
}
</style>
