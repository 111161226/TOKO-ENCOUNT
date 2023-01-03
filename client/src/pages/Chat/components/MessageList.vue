<script lang="ts" setup>
import { Message } from '@/lib/apis'

defineProps<{
  messages: Message[] | undefined
  myUserName: string | undefined
}>()
</script>

<template>
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
</template>

<style lang="scss" scoped>
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
