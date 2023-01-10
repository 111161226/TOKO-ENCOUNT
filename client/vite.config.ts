import { defineConfig } from 'vite'
import vue from '@vitejs/plugin-vue'
import path from 'path'

const devHost = 'localhost:3050'

// https://vitejs.dev/config/
export default defineConfig({
  resolve: {
    alias: {
      '@': path.resolve(__dirname, './src').replace(/\\/g, '/')
    }
  },
  server: {
    port: 3000,
    proxy: {
      '/api': {
        target: `http://${devHost}`,
        changeOrigin: true
      }
    }
  },
  css: {
    preprocessorOptions: {
      scss: {
        additionalData: `@import "@/variables.scss";`
      }
    }
  },
  plugins: [vue()]
})
