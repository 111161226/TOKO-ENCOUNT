import { defineStore } from 'pinia'
import { ElMessage } from 'element-plus'
import { Message } from '@/lib/apis'
import { useMessages } from './message'

interface EventData {
  type: string
  body: {
    roomId: string
    message: Message
  }
}

const ENTRY_POINT = 'ws://localhost:3050/api/ws'

export const useWebSocket = defineStore('websocket', {
  state: (): { socket: WebSocket | undefined } => ({
    socket: undefined
  }),
  getters: {
    getSocket: state => state.socket
  },
  actions: {
    connect() {
      const socket = new WebSocket(ENTRY_POINT)
      const messageStore = useMessages()
      socket.addEventListener(
        'message',
        function (event: MessageEvent<EventData>) {
          const { body } = event.data
          messageStore.catchNewMessage(body.roomId, body.message)
          ElMessage({
            message: `新着メッセージ\n${body.message.userName}: ${body.message.post}`
          })
        }
      )
      this.socket = socket
    }
  }
})
