import { createRouter, createWebHistory } from 'vue-router'
import { useUserStore } from './store/user'

import Login from './pages/Login.vue'
import Register from './pages/Register.vue'
import TaskList from './pages/TaskList.vue'
import TaskDetail from './pages/TaskDetail.vue'

const routes = [
  { path: '/login', component: Login },
  { path: '/register', component: Register },
  { path: '/', component: TaskList },
  { path: '/task/:id', component: TaskDetail, props: true }
]

const router = createRouter({
  history: createWebHistory(),
  routes
})


function getCookie(name) {
  return document.cookie.split('; ').reduce((r, v) => {
    const parts = v.split('=')
    return parts[0] === name ? decodeURIComponent(parts[1]) : r
  }, '')
}

router.beforeEach((to, from, next) => {
  const token = getCookie('token')
  const whiteList = ['/login', '/register']
  if (!token && !whiteList.includes(to.path)) {
    return next({ path: '/login', replace: true })
  }
  if (token && to.path === '/login') {
    return next({ path: '/', replace: true })
  }
  next()
})

export default router
