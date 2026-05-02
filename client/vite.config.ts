import { defineConfig, loadEnv } from 'vite' // loadEnv を追加
import vue from '@vitejs/plugin-vue'
import path from 'path'

export default defineConfig(({ mode }) => {
  // 現在の作業ディレクトリにある環境変数を読み込む
  // 第3引数を '' にすると、VITE_ 以外の変数も読み込めます
  const env = loadEnv(mode, process.cwd(), '');
  const devHost = env.VITE_API_HOST || 'localhost:3050';

  return {
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
  }
})