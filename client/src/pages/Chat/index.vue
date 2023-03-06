<script lang="ts" setup>
import { AxiosError } from 'axios'
import { onMounted, ref, computed, reactive } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { useMessages } from '@/store/message'
import { useDraftMessages } from '@/store/draftMessage'
import { useroomNames } from '@/store/roomName'
import { useChatRooms } from '@/store/chatRoom'
import apis from '@/lib/apis'
import { showErrorMessage } from '@/util/showErrorMessage'
import ChatInput from './components/ChatInput.vue'
import MessageList from './components/MessageList.vue'

const draftMessageStore = useDraftMessages()
const messageStore = useMessages()
const storeRoom = useroomNames()
const chatRoomStore = useChatRooms()

const route = useRoute()
const roomId = route.params.id as string
const router = useRouter()

const messages = computed(() => messageStore.getMessage(roomId).messages)
const contentDivRef = ref<HTMLDivElement>()
const roominfo = computed(() => storeRoom.getRoom(roomId))
const edit = reactive({
  update: false
})
const inputData = reactive({
  roomname: ''
})

const onSubmit = async () => {
  const message = draftMessageStore.getMessage(roomId)
  if (message) {
    try {
      await messageStore.sendMessage(roomId, message)
      draftMessageStore.setMessage(roomId, '')
    } catch (e: any) {
      const err: AxiosError = e
      showErrorMessage(err)
    }
  }
}

const updateName = async () => {
  if (inputData.roomname) {
    try {
      const { data } = await apis.editRoomName(roomId, inputData.roomname)
      storeRoom.updateName(roomId, inputData.roomname)
      chatRoomStore.setName(roomId, inputData.roomname)
      messageStore.addMessage(roomId, data.latestMessage)
    } catch (e: any) {
      const err: AxiosError = e
      showErrorMessage(err)
    }
    inputData.roomname = ''
    edit.update = false
  }
}

const editName = async () => {
  if (roomId != '0') {
    edit.update = true
  }
}

const Adduser = async () => {
  try {
    router.push(`/chat/${roomId}/add`)
  } catch (e: any) {
    const err: AxiosError = e
    showErrorMessage(err)
  }
}

onMounted(async () => {
  if (messages.value.length < 20 && messageStore.getMessage(roomId).hasNext) {
    // store に情報がない時だけ初回読み込みを実行
    try {
      messageStore.setLoading(true)
      await messageStore.fetchData(roomId, 20)
    } catch (e: any) {
      const err: AxiosError = e
      showErrorMessage(err)
    } finally {
      messageStore.setLoading(false)
    }
  }

  if (!roominfo.value) {
    try {
      storeRoom.setLoading(true)
      await storeRoom.fetchRoomName(roomId)
    } catch (e: any) {
      const err: AxiosError = e
      showErrorMessage(err)
    } finally {
      storeRoom.setLoading(false)
    }
  }

  if (contentDivRef.value) {
    // scroll to bottom
    contentDivRef.value?.scrollTo(0, contentDivRef.value.scrollHeight)
  }
})
</script>

<template>
  <div class="chat-container">
    <div v-if="roomId === '0'" class="header">
      {{ roominfo }}
    </div>
    <div v-else-if="edit.update === false" class="header">
      <div class="row">
        {{ roominfo }}
        <div @click="editName" @keydown.enter="editName">
          <el-icon size="1.5rem" class="icon"><promotion /></el-icon>
        </div>
        <button class="invite-button" @click="Adduser">
          <div>招待</div>
        </button>
      </div>
    </div>
    <div v-else class="header">
      <div class="row">
        <el-input
          v-model="inputData.roomname"
          maxlength="30"
          show-word-limit
          @keyup.enter="updateName"
        />
        <div @click="updateName" @keydown.enter="updateName">
          <el-icon size="1.5rem" class="icon"><promotion /></el-icon>
        </div>
      </div>
    </div>
    <div class="content" ref="contentDivRef">
      <message-list
        :room-id="roomId"
        :messages="messages"
        :show-user-name="true"
      />
    </div>
    <div class="input-container">
      <chat-input class="input" />
      <div @click="onSubmit" @keydown.enter="onSubmit">
        <el-icon size="1.5rem" class="icon"><promotion /></el-icon>
      </div>
    </div>
  </div>
</template>

<style lang="scss" scoped>
.chat-container {
  display: flex;
  flex-direction: column;
}
.header {
  position: sticky;
  top: 0;
  font-weight: 600;
  border-bottom: 2px solid $bgcolor-primary;
  padding: 0.5rem 1rem;
}
.content {
  overflow-y: auto;
  padding: 0.5rem 1rem;
  flex-grow: 1;
}
.row {
  display: flex;
  align-items: flex-end;
  padding: 0.5rem 1rem;
}

.input-container {
  display: flex;
  align-items: flex-end;
  padding: 0.5rem 1rem;
}
.input {
  margin-right: 0.5rem;
}
.invite-button {
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
  margin-left: auto;
}
.icon {
  color: $color-secondary;
  cursor: pointer;
}
</style>
