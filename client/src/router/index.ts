import { AxiosError } from 'axios'
import {
  createRouter,
  createWebHistory,
  RouteLocation,
  RouteRecordRaw
} from 'vue-router'
import { ElMessage } from 'element-plus'
import Layout from '@/layout/index.vue'
import { useMe } from '@/store/me'
import { showErrorMessage } from '@/util/showErrorMessage'

interface RouteMeta {
  title?: string
  isPublic?: boolean
}

type IRouteRecordRaw = RouteRecordRaw & {
  meta?: RouteMeta
  children?: IRouteRecordRaw[]
}

type IRoute = Omit<RouteLocation, 'meta'> & {
  meta?: RouteMeta
}

export const sidebarRoutes: IRouteRecordRaw[] = [
  {
    path: '',
    name: 'Home',
    component: () => import('@/pages/Home/index.vue'),
    meta: { title: 'Home' }
  },
  {
    path: '/search',
    name: 'Search',
    component: () => import('@/pages/Search/index.vue'),
    meta: { title: 'Search' }
  },
  {
    path: '/profile',
    name: 'Profile',
    component: () => import('@/pages/EditAccount/index.vue'),
    meta: { title: 'Profile' }
  }
]

const privateRoutes: IRouteRecordRaw[] = [
  {
    path: '/chat/:id',
    name: 'Chat',
    component: () => import('@/pages/Chat/index.vue'),
    meta: { title: 'Chat' }
  },
  {
    path: '/chat/:id/add',
    name: 'Search Member to add to private chat',
    component: () => import('@/pages/Add/index.vue'),
    meta: { title: 'Search Member to add to private chat' }
  }
]

const constantRouts: IRouteRecordRaw[] = [
  {
    path: '/',
    component: Layout,
    children: sidebarRoutes.concat(privateRoutes)
  }
]

const publicRoutes: IRouteRecordRaw[] = [
  {
    path: '/login',
    name: 'Login',
    component: () => import('@/pages/Login/index.vue')
  },
  {
    path: '/signup',
    name: 'CreateAccount',
    component: () => import('@/pages/CreateAccount/index.vue')
  },
  {
    path: '/:pathMatch(.*)',
    name: 'NotFound',
    component: () => import('@/pages/404/index.vue')
  }
].map((route: IRouteRecordRaw) => {
  return {
    ...route,
    meta: { isPublic: true }
  }
})

const router = createRouter({
  history: createWebHistory(),
  routes: constantRouts.concat(publicRoutes)
})

router.beforeEach(async (to: IRoute, _, next) => {
  const meStore = useMe()

  if (to.meta && to.meta.isPublic) {
    // 行き先が未ログインユーザー向けページの場合
    if (meStore.getMe) {
      // ログイン済みの場合、トップページへリダイレクト
      ElMessage({
        message: 'ログイン済みです',
        type: 'info'
      })
      next({ name: 'Home' })
      return
    } else {
      try {
        await meStore.setMe()
        ElMessage({
          message: 'ログイン済みです',
          type: 'info'
        })
        next({ name: 'Home' })
        return
      } catch (e: any) {
        const err: AxiosError = e
        if (err.response && err.response.status === 401) {
          // 未ログインの場合のエラーは無視
          next()
        } else {
          showErrorMessage(err)
        }
        return
      }
    }
  }

  if (meStore.getMe) {
    next()
    return
  }

  try {
    await meStore.setMe()
    next()
    return
  } catch (e: any) {
    const err: AxiosError = e
    if (err.response && err.response.status === 401) {
      ElMessage({
        message: '未ログインでした',
        type: 'error'
      })
    } else {
      showErrorMessage(err)
    }
    next({ name: 'Login' })
    return
  }
})

export default router
