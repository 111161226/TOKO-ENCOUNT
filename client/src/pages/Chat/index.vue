<script lang="ts" setup>
import apis, { Message } from '@/lib/apis'
import { onMounted, ref } from 'vue'
import { useRoute } from 'vue-router'

const route = useRoute()
const messages = ref<Message[]>()

onMounted(async () => {
  const { data } = await apis.getChatMessages(route.params.id as string)
  messages.value = data.messages
})
</script>

<template>
  <div class="chat-container">
    <div v-for="message in messages" :key="message.chatId">
      {{ message.post }}
    </div>
  </div>
</template>

<style lang="scss" scoped>
.chat-container {
  display: flex;
  flex-direction: column;
}
</style>
