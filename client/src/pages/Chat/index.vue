<script lang="ts" setup>
import apis, { Message } from '@/lib/apis'
import { useMe } from '@/store/me'
import { useMessages } from '@/store/message'
import { onMounted, onUpdated, ref } from 'vue'
import { useRoute } from 'vue-router'
import ChatInput from './components/ChatInput.vue'
import MessageList from './components/MessageList.vue'

const storeMe = useMe()
const storeMessages = useMessages()

const route = useRoute()
const roomId = route.params.id as string

const myUserName = ref(storeMe.getMe?.userName)
const otherUserName = ref('')
const messages = ref<Message[]>()
const contentDivRef = ref<HTMLDivElement>()

const onSubmit = async () => {
  const message = storeMessages.getMessage(roomId)?.message
  if (message) {
    await apis.postChat(roomId, { post: message })
    storeMessages.setMessage(roomId, '')
  }
}

onMounted(async () => {
  const { data } = await apis.getChatMessages(roomId)
  messages.value = data.messages

  for (const message of messages.value) {
    if (message.userName !== myUserName.value) {
      otherUserName.value = message.userName
      break
    }
  }
})
onUpdated(() => {
  // scroll to bottom of messages
  contentDivRef.value?.scrollTo(0, contentDivRef.value.scrollHeight)
})
</script>

<template>
  <div class="chat-container">
    <div class="header">
      {{ otherUserName }}
    </div>
    <div class="content" ref="contentDivRef">
      <message-list :messages="messages" :my-user-name="myUserName" />
    </div>
    <div class="input-container">
      <chat-input class="input" />
      <div @click="onSubmit" @keydown.enter="onSubmit">
        <el-icon size="1.5rem" class="icon"><Promotion /></el-icon>
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
  border-bottom: 2px solid $bgcolor-primary;
  padding: 0.5rem 1rem;
}
.content {
  overflow-y: auto;
  padding: 0.5rem 1rem;
}

.input-container {
  display: flex;
  align-items: flex-end;
  padding: 0.5rem 1rem;
}
.input {
  margin-right: 0.5rem;
}
.icon {
  color: $color-secondary;
  cursor: pointer;
}
</style>
