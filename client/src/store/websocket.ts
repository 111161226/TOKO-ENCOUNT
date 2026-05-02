import { defineStore } from 'pinia'
import { ElMessage } from 'element-plus'
import { Message } from '@/lib/apis'
import { useMessages } from './message'
import { useChatRooms } from './chatRoom'

interface EventData {
  type: string
  body: {
    roomId: string
    message: Message
  }
}

const isProduction = import.meta.env.PROD;
const host = import.meta.env.VITE_API_HOST || 'localhost:3050';

// 本番なら wss://、ローカルなら ws:// を自動選択
const protocol = isProduction ? 'wss' : 'ws';

const ENTRY_POINT = `${protocol}://${host}/api/ws`;

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
      const roomStore = useChatRooms()
      socket.onmessage = function (event) {
        const data: EventData = JSON.parse(event.data)
        messageStore.catchNewMessage(data.body.roomId, data.body.message)
        roomStore.catchNewMessage(data.body.roomId, data.body.message)
        ElMessage({
          message: `新着メッセージ\n${data.body.message.userName}: ${data.body.message.post}`
        })
      }
      this.socket = socket
    }
  }
})
