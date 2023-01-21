<script lang="ts" setup>
import apis, { UserWithoutPass } from '@/lib/apis'
import { useUsers } from '@/store/user'
import { showErrorMessage } from '@/util/showErrorMessage'
import { AxiosError } from 'axios'
import { ElLoading } from 'element-plus'
import { computed, defineProps, onMounted, ref, watchEffect } from 'vue'

const props = defineProps<{
  users: UserWithoutPass[]
  searchQuery: {
    name: string
    gender: string
    prefecture: string
  }
}>()

const userStore = useUsers()

const hasNext = computed(() => userStore.getUsers.hasNext)
const loading = computed(() => userStore.getLoading)

const loadingEle = ref<HTMLDivElement>()
const usersEle = ref<HTMLDivElement[]>([])

const onCreateRoom = async (userId: string) => {
  try {
    await apis.createChat(userId)
    // TODO: router.push to /chat/{roomId}, currently can't get roomId of the created room
  } catch (e: any) {
    const err: AxiosError = e
    showErrorMessage(err)
  }
}

const fetchData = async () => {
  try {
    userStore.setLoading(true)
    await userStore.fetchData(
      30,
      props.searchQuery.name,
      props.searchQuery.gender,
      props.searchQuery.prefecture
    )
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
    <div
      v-for="(user, index) in users"
      :key="user.userId"
      :ref="
        el => {
          usersEle[index] = el as HTMLDivElement
        }
      "
    >
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
                onCreateRoom(user.userId)
              }
            "
          >
            <div>Create Room</div>
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
