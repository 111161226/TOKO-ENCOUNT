<script lang="ts" setup>
import apis, { Message } from '@/lib/apis'
import { useMe } from '@/store/me'
import { onMounted, ref } from 'vue'
import { useRoute } from 'vue-router'

const storeMe = useMe()
const route = useRoute()

const myUserName = ref(storeMe.getMe?.userName)
const otherUserName = ref('')
const messages = ref<Message[]>()

onMounted(async () => {
  const { data } = await apis.getChatMessages(route.params.id as string)
  messages.value = data.messages

  for (const message of messages.value) {
    if (message.userName !== myUserName.value) {
      otherUserName.value = message.userName
      break
    }
  }
})
</script>

<template>
  <div class="chat-container">
    <div class="header">
      {{ otherUserName }}
    </div>
    <div class="content">
      <div class="message-list">
        <div
          v-for="message in messages"
          :key="message.chatId"
          :class="{
            message: true,
            'my-message': message.userName === myUserName,
            'other-message': message.userName !== myUserName
          }"
        >
          {{ message.post }}
        </div>
      </div>
    </div>
  </div>
</template>

<style lang="scss" scoped>
.chat-container {
  display: flex;
  flex-direction: column;
}
.header {
  position: sticky;
  top: 0;
  font-weight: 600;
  border-bottom: 3px solid $bgcolor-primary;
  padding: 0.5rem 1rem;
}
.content {
  overflow-y: auto;
  padding: 0.5rem 1rem;
}

.message-list {
  display: flex;
  flex-direction: column;
}
.message {
  white-space: pre-wrap;
  padding: 0.5rem 1rem;
  margin: 0.5rem 0;
}
.my-message {
  align-self: flex-end;
  background-color: $color-primary;
  color: white;
  border-radius: 8px;
  border-bottom-right-radius: 0;
}
.other-message {
  align-self: flex-start;
  background-color: $bgcolor-primary;
  border-radius: 8px;
  border-bottom-left-radius: 0;
}
</style>
