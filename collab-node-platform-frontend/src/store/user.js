import { defineStore } from 'pinia'

function setCookie(name, value, days = 7) {
  const expires = new Date(Date.now() + days * 864e5).toUTCString()
  document.cookie = `${name}=${encodeURIComponent(value)}; expires=${expires}; path=/`
}
function getCookie(name) {
  return document.cookie.split('; ').reduce((r, v) => {
    const parts = v.split('=')
    return parts[0] === name ? decodeURIComponent(parts[1]) : r
  }, '')
}
function delCookie(name) {
  document.cookie = `${name}=; expires=Thu, 01 Jan 1970 00:00:00 GMT; path=/`
}

export const useUserStore = defineStore('user', {
  state: () => ({
    token: getCookie('token') || '',
    userId: getCookie('userId') || '',
    username: getCookie('username') || ''
  }),
  actions: {
    setUser(token, userId, username) {
      this.token = token
      this.userId = userId
      this.username = username
      setCookie('token', token)
      setCookie('userId', userId)
      setCookie('username', username)
    },
    logout() {
      this.token = ''
      this.userId = ''
      this.username = ''
      delCookie('token')
      delCookie('userId')
      delCookie('username')
    }
  }
})
