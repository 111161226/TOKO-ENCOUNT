<script lang="ts" setup>
import Sidebar from './Sidebar.vue'
</script>

<template>
  <div class="root-container">
    <sidebar />
    <router-view />
  </div>
</template>

<style lang="scss" scoped>
.root-container {
  display: flex;
  flex-direction: row;
  height: 100%;
}
</style>
