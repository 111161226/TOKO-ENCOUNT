<script lang="ts" setup>
import { ref } from 'vue'

const name = ref('')
const gender = ref('')
const prefecture = ref('')
</script>

<template>
  <div>
    <div class="search-bar">
      <el-input v-model="name" placeholder="Name" class="input" />
      <el-select v-model="gender" class="input"></el-select>
      <el-select v-model="prefecture" class="input"></el-select>
      <button class="search-button">
        <el-icon class="search-icon"><search /></el-icon>
        <div>Search</div>
      </button>
    </div>
  </div>
</template>

<style lang="scss" scoped>
.search-bar {
  display: flex;
  flex-direction: row;
  align-items: center;
  padding: 2rem 4rem;
}
.input {
  margin-right: 1rem;
}
.search-button {
  display: flex;
  align-items: center;
  justify-content: center;
  background-color: $color-primary;
  color: white;
  border: none;
  border-radius: 0.25rem;
  padding: 0.5rem 0.75rem;
  cursor: pointer;
  transition: background-color 0.2s ease-in-out;
  &:hover {
    background-color: $color-secondary;
  }
}
.search-icon {
  margin-right: 0.5rem;
}
</style>
