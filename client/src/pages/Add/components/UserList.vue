<script lang="ts" setup>
import apis, { UserWithoutPass } from '@/lib/apis'
import { useUsers } from '@/store/user'
import { useMessages } from '@/store/message'
import { showErrorMessage } from '@/util/showErrorMessage'
import { AxiosError } from 'axios'
import { ElLoading } from 'element-plus'
import { computed, defineProps, onMounted, ref, watchEffect } from 'vue'
import { useRouter, useRoute } from 'vue-router'

const props = defineProps<{
  users: UserWithoutPass[]
  searchQuery: {
    name: string
  }
}>()

const router = useRouter()
const userStore = useUsers()
const messageStore = useMessages()
const route = useRoute()
const roomId = route.params.id as string

const hasNext = computed(() => userStore.getUsers.hasNext)
const loading = computed(() => userStore.getLoading)

const loadingEle = ref<HTMLDivElement>()

const onAddRoom = async (userId: string) => {
  try {
    const { data } = await apis.addChat(roomId, userId)
    messageStore.addMessage(roomId, data.latestMessage)
    router.push(`/chat/${data.roomId}`)
  } catch (e: any) {
    const err: AxiosError = e
    showErrorMessage(err)
  }
}

const fetchData = async () => {
  try {
    userStore.setLoading(true)
    await userStore.fetchAddMembers(30, props.searchQuery.name, roomId)
  } catch (e: any) {
    const err: AxiosError = e
    showErrorMessage(err)
  } finally {
    userStore.setLoading(false)
  }
}

watchEffect(onCleanup => {
  const { value } = loadingEle
  if (!value || !hasNext.value || loading.value) return

  const observer = new IntersectionObserver(
    async entries => {
      for (const entry of entries) {
        if (!entry.isIntersecting) {
          return
        }
      }

      await fetchData()
    },
    { threshold: 0.5 }
  )
  observer.observe(value)

  onCleanup(() => {
    observer.disconnect()
  })
})

onMounted(() => {
  if (loadingEle.value) {
    ElLoading.service({ target: loadingEle.value })
  }
})
</script>

<template>
  <div class="user-list-container">
    <div v-for="(user, index) in users" :key="user.userId">
      <hr v-if="index !== 0" class="line" />
      <div class="user">
        <div class="user-left">
          <div class="username">{{ user.userName }}</div>
          <div class="user-info">
            <div class="gender">{{ user.gender }}</div>
            <div>{{ user.prefect }}</div>
          </div>
        </div>
        <div>
          <button
            class="button"
            @click="
              () => {
                onAddRoom(user.userId)
              }
            "
          >
            <div>Add Room member</div>
          </button>
        </div>
      </div>
    </div>
    <div v-if="hasNext" ref="loadingEle" class="loading"></div>
  </div>
</template>

<style lang="scss" scoped>
.user-list-container {
  display: flex;
  flex-direction: column;
  padding: 1rem 4rem;
}
.user {
  display: flex;
  flex-direction: row;
  justify-content: space-between;
  align-items: center;
  padding: 0.5rem 0;
}
.line {
  border: 1px solid $bgcolor-primary;
  border-radius: 50rem;
  margin: 0;
}
.user-left {
  display: flex;
  flex-direction: column;
}
.username {
  font-size: 1.25rem;
  font-weight: 600;
}
.user-info {
  display: flex;
  font-size: 0.9rem;
}
.gender {
  margin-right: 0.5rem;
}
.button {
  background-color: $color-primary;
  color: white;
  border: none;
  border-radius: 0.25rem;
  white-space: nowrap;
  padding: 0.25rem 0.5rem;
  cursor: pointer;
  transition: background-color 0.2s;
  &:hover {
    background-color: $color-secondary;
  }
}
.loading {
  height: 3rem;
}
</style>
