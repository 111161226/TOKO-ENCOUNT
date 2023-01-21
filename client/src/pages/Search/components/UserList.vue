<script lang="ts" setup>
import apis, { UserWithoutPass } from '@/lib/apis'
import { showErrorMessage } from '@/util/showErrorMessage'
import { AxiosError } from 'axios'
import { defineProps } from 'vue'

defineProps<{
  users: UserWithoutPass[]
}>()

const onCreateRoom = async (userId: string) => {
  try {
    await apis.createChat(userId)
    // TODO: router.push to chat page, currently can't get roomId of the created room
  } catch (e: any) {
    const err: AxiosError = e
    showErrorMessage(err)
  }
}
</script>

<template>
  <div class="user-list-container">
    <div v-for="(user, index) in users" :key="user.userId">
      <hr v-if="index !== 0" class="line" />
      <div class="user">
        <div class="user-left">
          <div class="username">{{ user.userName }}</div>
          <div class="user-info">
            <div class="gender">{{ user.gender }}</div>
            <div>{{ user.prefect }}</div>
          </div>
        </div>
        <div>
          <button
            class="button"
            @click="
              () => {
                onCreateRoom(user.userId)
              }
            "
          >
            <div>Create Room</div>
          </button>
        </div>
      </div>
    </div>
  </div>
</template>

<style lang="scss" scoped>
.user-list-container {
  display: flex;
  flex-direction: column;
  padding: 1rem 4rem;
}
.user {
  display: flex;
  flex-direction: row;
  justify-content: space-between;
  align-items: center;
  padding: 0.5rem 0;
}
.line {
  border: 1px solid $bgcolor-primary;
  border-radius: 50rem;
  margin: 0;
}
.user-left {
  display: flex;
  flex-direction: column;
}
.username {
  font-size: 1.25rem;
  font-weight: 600;
}
.user-info {
  display: flex;
  font-size: 0.9rem;
}
.gender {
  margin-right: 0.5rem;
}
.button {
  background-color: $color-primary;
  color: white;
  border: none;
  border-radius: 0.25rem;
  white-space: nowrap;
  padding: 0.25rem 0.5rem;
  cursor: pointer;
  transition: background-color 0.2s;
  &:hover {
    background-color: $color-secondary;
  }
}
</style>
