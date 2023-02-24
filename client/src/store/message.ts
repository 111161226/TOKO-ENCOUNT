import { defineStore } from 'pinia'
import api, { Message, MessageList } from '@/lib/apis'

export const useMessages = defineStore('messages', {
  state: (): { messages: Record<string, MessageList>; loading: boolean } => ({
    messages: {},
    loading: false
  }),
  getters: {
    getMessage: state => (roomId: string) =>
      state.messages[roomId] ?? { hasNext: true, messages: [] },
    getLoading: state => () => state.loading
  },
  actions: {
    async fetchData(roomId: string, limit: number) {
      if (!this.messages[roomId]) {
        const { data } = await api.getChatMessages(roomId, limit, 0)
        this.messages = { ...this.messages, [roomId]: data }
        return
      }

      const prevData = this.messages[roomId]
      const { data } = await api.getChatMessages(
        roomId,
        limit,
        prevData.messages.length
      )
      this.messages = {
        ...this.messages,
        [roomId]: {
          ...data,
          messages: prevData.messages.concat(data.messages)
        }
      }
    },
    async sendMessage(roomId: string, message: string) {
      const { data } = await api.postChat(roomId, {
        post: message
      })
      this.messages[roomId].messages.unshift(data)
    },
    addMessage(roomId: string, message: Message) {
      this.messages[roomId].messages.unshift(message)
    },
    setLoading(value: boolean) {
      this.loading = value
    },
    /** websocket 以外では呼び出さない */
    catchNewMessage(roomId: string, message: Message) {
      if (!this.messages[roomId]) {
        this.messages = {
          ...this.messages,
          [roomId]: { hasNext: false, messages: [message] }
        }
        return
      }

      this.messages[roomId].messages.unshift(message)
    }
  }
})
