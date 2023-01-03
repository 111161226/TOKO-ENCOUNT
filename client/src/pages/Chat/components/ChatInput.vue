<script lang="ts" setup>
import { useMessages } from '@/store/message'
import { onMounted, ref, computed } from 'vue'
import { onBeforeRouteLeave, useRoute } from 'vue-router'

const messagesStore = useMessages()
const route = useRoute()

const inputMessage = ref('')

const roomId = route.params.id as string
onMounted(() => {
  const storedMessage = computed(() => messagesStore.getMessage(roomId))
  inputMessage.value = storedMessage.value?.message ?? ''
})
onBeforeRouteLeave(() => {
  messagesStore.setMessage(roomId, inputMessage.value)
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
  ></el-input>
</template>
