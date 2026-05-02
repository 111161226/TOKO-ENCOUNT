import { defineConfig } from 'vite'
import vue from '@vitejs/plugin-vue'
import path from 'path'

const devHost = import.meta.env.VITE_API_HOST 

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
        target: `https://${devHost}`,
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
