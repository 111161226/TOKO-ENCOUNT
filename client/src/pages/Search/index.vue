<script lang="ts" setup>
import { ref } from 'vue'
import { prefectures } from '@/util/prefectures'
import apis, { UserWithoutPass } from '@/lib/apis'
import UserList from './components/UserList.vue'
import { AxiosError } from 'axios'
import { showErrorMessage } from '@/util/showErrorMessage'

const input = ref({
  name: '',
  gender: '',
  prefecture: ''
})
const result = ref<UserWithoutPass[]>([])

const onSearch = async () => {
  try {
    const { data } = await apis.getUsers(
      50,
      0,
      input.value.name,
      input.value.gender,
      input.value.prefecture
    )
    result.value = data.users
  } catch (e: any) {
    const err: AxiosError = e
    showErrorMessage(err)
  }
}
</script>

<template>
  <div class="search-container">
    <div class="search-bar">
      <el-input v-model="input.name" placeholder="名前" class="input" />
      <el-select
        v-model="input.gender"
        placeholder="性別"
        class="input"
        clearable
      >
        <el-option label="male" value="male" />
        <el-option label="female" value="female" />
      </el-select>
      <el-select
        v-model="input.prefecture"
        placeholder="都道府県"
        clearable
        class="input"
      >
        <el-option
          v-for="prefecture in prefectures"
          :key="prefecture.value"
          :label="prefecture.label"
          :value="prefecture.value"
        />
      </el-select>
      <button class="search-button" @click="onSearch">
        <el-icon class="search-icon"><search /></el-icon>
        <div>Search</div>
      </button>
    </div>
    <div class="content">
      <user-list :users="result" />
    </div>
  </div>
</template>

<style lang="scss" scoped>
.search-container {
  display: flex;
  flex-direction: column;
}
.search-bar {
  display: flex;
  flex-direction: row;
  align-items: center;
  padding: 1.5rem 3rem;
}
.input {
  margin-right: 1rem;
}
.search-button {
  display: flex;
  align-items: center;
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
  margin-right: 0.25rem;
}
.content {
  overflow-y: auto;
}
</style>
