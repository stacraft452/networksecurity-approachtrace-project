import axios from 'axios'
import { useUserStore } from '../store/user'
import { ElMessage } from 'element-plus'

const instance = axios.create({
  baseURL: '/api', // 让Vite代理生效，生产环境由nginx等反代
  timeout: 10000
})

instance.interceptors.request.use(config => {
  const userStore = useUserStore()
  if (userStore.token) {
    config.headers.Authorization = `Bearer ${userStore.token}`
  }
  return config
})

instance.interceptors.response.use(
  res => res.data,
  err => {
    if (err.response && err.response.data && err.response.data.error) {
      ElMessage.error(err.response.data.error)
    }
    return Promise.reject(err)
  }
)

export default instance

function getCookie(name) {
  return document.cookie.split('; ').reduce((r, v) => {
    const parts = v.split('=')
    return parts[0] === name ? decodeURIComponent(parts[1]) : r
  }, '')
}

instance.interceptors.request.use(config => {
  const token = getCookie('token')
  if (token) {
    config.headers.Authorization = `Bearer ${token}`
  }
  return config
})
