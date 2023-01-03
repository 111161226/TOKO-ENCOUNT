import { defineStore } from 'pinia'

export interface MessageState {
  roomId: string
  message: string
}

export const useMessages = defineStore('messages', {
  state: (): { messages: MessageState[] | undefined } => ({
    messages: undefined
  }),
  getters: {
    getMessages: state => state.messages,
    getMessage: state => (roomId: string) =>
      state.messages?.find(messageState => messageState.roomId === roomId)
  },
  actions: {
    setMessage(roomId: string, message: string) {
      const index = this.messages?.findIndex(
        messageState => messageState.roomId === roomId
      )
      if (index !== undefined) {
        this.messages?.splice(index, 1, { roomId, message })
      } else {
        if (this.messages !== undefined) {
          this.messages.push({ roomId, message })
        } else {
          this.messages = [{ roomId, message }]
        }
      }
    }
  }
})
