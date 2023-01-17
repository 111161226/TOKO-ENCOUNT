<script lang="ts" setup>
import { AxiosError } from 'axios'
import { useRouter } from 'vue-router'
import { ElMessage } from 'element-plus'
import { sidebarRoutes } from '@/router'
import { useMe } from '@/store/me'
import { showErrorMessage } from '@/util/showErrorMessage'

const router = useRouter()
const me = useMe()

const logout = async () => {
  try {
    await me.logout()
    ElMessage({
      message: 'ログアウトしました',
      type: 'success'
    })
    router.push({ name: 'Login' })
  } catch (e: any) {
    const err: AxiosError = e
    showErrorMessage(err)
  }
}
</script>

<template>
  <div class="sidebar-container">
    <div class="routes-container">
      <div v-for="route in sidebarRoutes" :key="route.name" class="route">
        <router-link
          v-if="route.meta && $route.meta.title"
          :to="{ name: route.name }"
          class="link"
          :class="{ 'active-link': $route.name === route.name }"
        >
          {{ route.name }}
        </router-link>
      </div>
    </div>
    <div class="logout">
      <router-link :to="{ name: 'Login' }" class="link">
        <el-button class="button" @click="logout">Logout</el-button>
      </router-link>
    </div>
  </div>
</template>

<style lang="scss" scoped>
.sidebar-container {
  display: flex;
  flex-direction: column;
  align-items: center;
  background-color: $bgcolor-primary;
  width: 20%;
  min-width: fit-content;
  max-width: 10rem;
  padding: 5rem 1.75rem;
}
.routes-container {
  display: flex;
  flex-direction: column;
  align-items: center;
  width: 100%;
}
.route {
  margin-bottom: 1.5rem;
  width: 60%;
  min-width: 5rem;
}
.link {
  text-decoration: none;
  font-size: larger;
  font-weight: 500;
  color: black;
  transition: color 0.2s;
  &:hover {
    color: $color-secondary;
  }
}
.active-link {
  color: $color-primary;
}

.logout {
  margin-top: auto;
}
.button {
  background-color: $color-primary;
  color: white;
  border-radius: 0.5rem;
}
</style>
