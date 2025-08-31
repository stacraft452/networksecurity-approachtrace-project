import { defineConfig } from 'vite'
import vue from '@vitejs/plugin-vue'

export default defineConfig({
  plugins: [vue()],
  server: {
    host: '0.0.0.0',
    port: 5173,
    proxy: {
      '/api': 'http://39.99.41.108:8087',
      '/ws': {
        target: 'ws://39.99.41.108:8087',
        ws: true
      }
    }
  }
})
