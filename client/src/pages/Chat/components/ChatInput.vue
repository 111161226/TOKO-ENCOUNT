<script lang="ts" setup>
import { onMounted, ref, computed, watch } from 'vue'
import { useRoute } from 'vue-router'
import { useDraftMessages } from '@/store/draftMessage'

const draftMessageStore = useDraftMessages()
const route = useRoute()

const inputMessage = ref('')
watch(inputMessage, newMessage => {
  draftMessageStore.setMessage(roomId, newMessage)
})

const roomId = route.params.id as string
draftMessageStore.$subscribe((_, state) => {
  inputMessage.value = state.messages?.[roomId] ?? ''
})
onMounted(() => {
  const storedMessage = computed(() => draftMessageStore.getMessage(roomId))
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
