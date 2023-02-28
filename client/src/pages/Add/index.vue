<script lang="ts" setup>
import { computed, ref } from 'vue'
import UserList from './components/UserList.vue'
import { AxiosError } from 'axios'
import { showErrorMessage } from '@/util/showErrorMessage'
import { useUsers } from '@/store/user'
import { useRoute } from 'vue-router'

const userStore = useUsers()
const route = useRoute()
const roomId = route.params.id as string

const input = ref({
  name: ''
})

const result = computed(() => userStore.getUsers.users)

const onSearch = async () => {
  try {
    userStore.initializeUsers()
    userStore.setLoading(true)
    await userStore.fetchAddMembers(30, input.value.name, roomId)
  } catch (e: any) {
    const err: AxiosError = e
    showErrorMessage(err)
  } finally {
    userStore.setLoading(false)
  }
}
</script>

<template>
  <div class="search-container">
    <div class="search-bar">
      <el-input v-model="input.name" placeholder="名前" class="input" />
      <button class="search-button" @click="onSearch">
        <el-icon class="search-icon"><search /></el-icon>
        <div>Search</div>
      </button>
    </div>
    <div class="content">
      <user-list :users="result" :search-query="input" />
    </div>
  </div>
</template>

<style lang="scss" scoped>
.search-container {
  display: flex;
  flex-direction: column;
}
.search-bar {
  display: flex;
  flex-direction: row;
  align-items: center;
  padding: 1.5rem 3rem;
}
.input {
  margin-right: 1rem;
}
.search-button {
  display: flex;
  align-items: center;
  background-color: $color-primary;
  color: white;
  border: none;
  border-radius: 0.25rem;
  padding: 0.5rem 0.75rem;
  cursor: pointer;
  transition: background-color 0.2s ease-in-out;
  &:hover {
    background-color: $color-secondary;
  }
}
.search-icon {
  margin-right: 0.25rem;
}
.content {
  overflow-y: auto;
}
</style>
