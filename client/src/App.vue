<script lang="ts" setup>
import { watchEffect } from 'vue'
import { useMe } from '@/store/me'
import { useWebSocket } from '@/store/websocket'

const meStore = useMe()
const wsStore = useWebSocket()

watchEffect(() => {
  if (meStore.getMe && !wsStore.getSocket) {
    wsStore.connect()
  }
})
</script>

<template>
  <div id="app">
    <router-view />
  </div>
</template>
