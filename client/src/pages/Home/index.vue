<script lang="ts" setup>
import { useChatRooms } from '@/store/chatRoom'
import { showErrorMessage } from '@/util/showErrorMessage'
import { AxiosError } from 'axios'
import { ElLoading } from 'element-plus'
import { computed, onMounted, ref, watchEffect } from 'vue'
import { getRelativeTime } from './time'

const chatRoomStore = useChatRooms()

const chatRooms = computed(() => chatRoomStore.getChatList.chats)
const hasNext = computed(() => chatRoomStore.getChatList.hasNext)
const loading = computed(() => chatRoomStore.getLoading)

const loadingEle = ref<HTMLDivElement>()

const fetchData = async () => {
  try {
    chatRoomStore.setLoading(true)
    await chatRoomStore.fetchData(20)
  } catch (e: any) {
    const err: AxiosError = e
    showErrorMessage(err)
  } finally {
    chatRoomStore.setLoading(false)
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

onMounted(async () => {
  await fetchData()
  if (loadingEle.value) {
    ElLoading.service({ target: loadingEle.value })
  }
})
</script>

<template>
  <div class="chat-list-container">
    <div class="title">Chats</div>
    <div v-for="room in chatRooms" :key="room.roomId">
      <hr class="line" />
      <div class="room" @click="$router.push(`/chat/${room.roomId}`)">
        <div class="room-left">
          <div class="room-name">{{ room.name || 'Unknown' }}</div>
          <div>{{ room.latestMessage.post }}</div>
        </div>
        <div class="room-right">
          <div
            class="new-message"
            :class="{ hidden: room.newMessageCount <= 0 }"
          >
            {{ room.newMessageCount }}
          </div>
          <div>
            {{ getRelativeTime(room.latestMessage.createdAt) }}
          </div>
        </div>
      </div>
    </div>
    <div v-if="hasNext" ref="loadingEle" class="loading" />
  </div>
</template>

<style lang="scss">
.chat-list-container {
  display: flex;
  flex-direction: column;
  overflow-y: scroll;
  padding: 2rem;
}
.title {
  font-size: 1.75rem;
  font-weight: 600;
  padding: 0.5rem;
}
.line {
  border: 1px solid $bgcolor-primary;
  border-radius: 1px;
}
.room {
  display: flex;
  justify-content: space-between;
  padding: 0 0.5rem;
  word-wrap: break-word;
  cursor: pointer;
}
.room-left {
  display: flex;
  flex-direction: column;

  .room-name {
    font-size: 1.2rem;
    font-weight: 600;
  }
}
.room-right {
  display: flex;
  flex-direction: column;
  align-items: flex-end;
  font-size: 0.8rem;

  .new-message {
    background-color: $color-primary;
    border-radius: 10rem;
    color: white;
    font-weight: 600;
    padding: 0 0.5rem;
  }
  .hidden {
    visibility: hidden;
  }
}
.loading {
  height: 3rem;
}
</style>
