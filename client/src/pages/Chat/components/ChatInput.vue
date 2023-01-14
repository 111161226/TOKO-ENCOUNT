<script lang="ts" setup>
import { useMessages } from '@/store/message'
import { onMounted, ref, computed, watch } from 'vue'
import { useRoute } from 'vue-router'

const messagesStore = useMessages()
const route = useRoute()

const inputMessage = ref('')
watch(inputMessage, newMessage => {
  messagesStore.setMessage(roomId, newMessage)
})

const roomId = route.params.id as string
messagesStore.$subscribe((_, state) => {
  inputMessage.value = state.messages?.[roomId] ?? ''
})
onMounted(() => {
  const storedMessage = computed(() => messagesStore.getMessage(roomId))
  inputMessage.value = storedMessage.value ?? ''
})
</script>

<template>
  <el-input
    v-model="inputMessage"
    type="textarea"
    :autosize="{ minRows: 1, maxRows: 3 }"
    resize="none"
    placeholder="Message"
    class="input"
  />
</template>
