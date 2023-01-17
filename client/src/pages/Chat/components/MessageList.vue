<script lang="ts" setup>
import { AxiosError } from 'axios'
import { ref, watchEffect, onMounted, computed } from 'vue'
import { ElLoading } from 'element-plus'
import { Message } from '@/lib/apis'
import { showErrorMessage } from '@/util/showErrorMessage'
import { useMe } from '@/store/me'
import { useMessages } from '@/store/message'

const props = defineProps<{
  roomId: string
  messages: Message[]
  showUserName: boolean
}>()

const meStore = useMe()
const messageStore = useMessages()

const myUserId = computed(() => meStore.getMe?.userId)
const hasNext = computed(() => messageStore.getMessage(props.roomId).hasNext)
const loading = computed(() => messageStore.getLoading())

const loadingEle = ref<HTMLDivElement>()
const messagesEle = ref<HTMLDivElement[]>([])

const fetchData = async () => {
  try {
    messageStore.setLoading(true)
    await messageStore.fetchData(props.roomId, 20)
  } catch (e: any) {
    const err: AxiosError = e
    showErrorMessage(err)
  } finally {
    messageStore.setLoading(false)
  }
}

watchEffect(onCleanup => {
  const { value } = loadingEle
  if (!value || !hasNext.value || loading.value) {
    return
  }

  const observer = new IntersectionObserver(
    async entries => {
      for (const entry of entries) {
        if (!entry.isIntersecting) {
          return
        }
      }

      const oldestEle = messagesEle.value.slice(-1)[0]
      await fetchData()
      oldestEle.scrollTo({ top: 100 })
    },
    {
      threshold: 0.5
    }
  )
  observer.observe(value)

  onCleanup(() => observer.disconnect())
})

onMounted(async () => {
  if (loadingEle.value) {
    ElLoading.service({ target: loadingEle.value })
  }
})
</script>

<template>
  <div class="message-list">
    <div
      v-for="(message, index) in messages"
      :key="message.chatId"
      :ref="
        el => {
          messagesEle[index] = el as HTMLDivElement
        }
      "
    >
      <div
        v-if="showUserName && message.postUserId !== myUserId"
        class="user-name"
      >
        {{ message.userName }}
      </div>
      <div
        class="message"
        :class="[
          message.postUserId === myUserId ? 'my-message' : 'other-message',
          showUserName ? 'show-user-name' : ''
        ]"
      >
        {{ message.post }}
      </div>
    </div>
    <div v-if="hasNext" ref="loadingEle" class="loading" />
  </div>
</template>

<style lang="scss" scoped>
.message-list {
  display: flex;
  flex-direction: column-reverse;
}
.message {
  white-space: pre-wrap;
  padding: 0.5rem 1rem;
  margin: 0.5rem 0;
  border-radius: 8px;
  width: fit-content;
  max-width: 80%;

  &.show-user-name {
    margin-top: 0.25rem;
  }

  &.my-message {
    align-self: flex-end;
    background-color: $color-primary;
    color: white;
    border-bottom-right-radius: 0;
  }

  &.other-message {
    align-self: flex-start;
    background-color: $bgcolor-primary;
    border-bottom-left-radius: 0;
  }
}

.user-name {
  margin-top: 0.5rem;
}

.loading {
  height: 3rem;
}
</style>
