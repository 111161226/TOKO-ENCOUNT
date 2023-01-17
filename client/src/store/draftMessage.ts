import { defineStore } from 'pinia'

export const useDraftMessages = defineStore('draftMessages', {
  state: (): { messages: Record<string, string> | undefined } => ({
    messages: undefined
  }),
  getters: {
    getMessages: state => state.messages,
    getMessage: state => (roomId: string) => state.messages?.[roomId]
  },
  actions: {
    setMessage(roomId: string, message: string) {
      this.messages = { ...this.messages, [roomId]: message }
    }
  }
})
