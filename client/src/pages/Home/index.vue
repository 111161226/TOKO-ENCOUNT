<script lang="ts" setup>
import apis, { ChatData } from '@/lib/apis'
import { onMounted, ref } from 'vue'
import { getRelativeTime, compareFunc } from './time'

const chatRooms = ref<ChatData[]>()

onMounted(async () => {
  const { data } = await apis.getChat()
  chatRooms.value = data.chats
  chatRooms.value.sort(compareFunc)
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
          <div class="created-at">
            {{ getRelativeTime(room.latestMessage.createdAt) }}
          </div>
        </div>
      </div>
    </div>
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
</style>
