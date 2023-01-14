<script lang="ts" setup>
import apis, { Message } from '@/lib/apis'
import { useMe } from '@/store/me'
import { useMessages } from '@/store/message'
import { showErrorMessage } from '@/util/showErrorMessage'
import { AxiosError } from 'axios'
import { nextTick, onMounted, ref, watch } from 'vue'
import { useRoute } from 'vue-router'
import ChatInput from './components/ChatInput.vue'
import MessageList from './components/MessageList.vue'

const storeMe = useMe()
const storeMessages = useMessages()

const route = useRoute()
const roomId = route.params.id as string

const myUserName = ref(storeMe.getMe?.userName)
const otherUserName = ref('')
const messages = ref<Message[]>([])
const contentDivRef = ref<HTMLDivElement>()

const onSubmit = async () => {
  const message = storeMessages.getMessage(roomId)
  if (message) {
    const { data } = await apis.postChat(roomId, { post: message })
    storeMessages.setMessage(roomId, '')
    messages.value?.push(data)
  }
}

watch(
  messages,
  async () => {
    // scroll to bottom
    await nextTick()
    contentDivRef.value?.scrollTo(0, contentDivRef.value.scrollHeight)
  },
  { deep: true }
)

onMounted(async () => {
  try {
    const { data } = await apis.getChatMessages(roomId)
    messages.value = data.messages

    for (const message of messages.value) {
      if (message.userName !== myUserName.value) {
        otherUserName.value = message.userName
        break
      }
    }
  } catch (e: any) {
    const err: AxiosError = e
    showErrorMessage(err)
  }
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
        <el-icon size="1.5rem" class="icon"><promotion /></el-icon>
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
  flex-grow: 1;
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
