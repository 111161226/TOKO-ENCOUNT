import { createApp } from 'vue'
import './style.scss'
import App from './App.vue'

import ElementPlus from 'element-plus'
import ja from 'element-plus/dist/locale/ja.mjs'
import 'element-plus/dist/index.css'
import { createPinia } from 'pinia'

import router from '@/router'

const app = createApp(App)
app.use(router)
app.use(ElementPlus, { locale: ja })
app.use(createPinia())
app.mount('#app')
