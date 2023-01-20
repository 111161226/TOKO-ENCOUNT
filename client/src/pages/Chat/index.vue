<script lang="ts" setup>
import { AxiosError } from 'axios'
import { onMounted, ref, computed } from 'vue'
import { useRoute } from 'vue-router'
import { useMe } from '@/store/me'
import { useMessages } from '@/store/message'
import { useDraftMessages } from '@/store/draftMessage'
import { showErrorMessage } from '@/util/showErrorMessage'
import ChatInput from './components/ChatInput.vue'
import MessageList from './components/MessageList.vue'

const storeMe = useMe()
const draftMessageStore = useDraftMessages()
const messageStore = useMessages()

const route = useRoute()
const roomId = route.params.id as string

const otherUserName = ref('')
const messages = computed(() => messageStore.getMessage(roomId).messages)
const contentDivRef = ref<HTMLDivElement>()

const onSubmit = async () => {
  const message = draftMessageStore.getMessage(roomId)
  if (message) {
    try {
      await messageStore.sendMessage(roomId, message)
      draftMessageStore.setMessage(roomId, '')
    } catch (e: any) {
      const err: AxiosError = e
      showErrorMessage(err)
    }
  }
}

onMounted(async () => {
  try {
    messageStore.setLoading(true)
    await messageStore.fetchData(roomId, 20)

    if (roomId == '0') {
      otherUserName.value = '全体チャット'
    } else {
      for (const message of messages.value) {
        if (message.postUserId !== storeMe.getMe?.userId) {
          otherUserName.value = message.userName
          break
        }
      }
    }
  } catch (e: any) {
    const err: AxiosError = e
    showErrorMessage(err)
  } finally {
    messageStore.setLoading(false)
  }

  if (contentDivRef.value) {
    // scroll to bottom
    contentDivRef.value?.scrollTo(0, contentDivRef.value.scrollHeight)
  }
})
</script>

<template>
  <div class="chat-container">
    <div class="header">
      {{ otherUserName }}
    </div>
    <div class="content" ref="contentDivRef">
      <message-list
        :room-id="roomId"
        :messages="messages"
        :show-user-name="roomId === '0'"
      />
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
