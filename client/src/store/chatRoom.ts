import apis, { ChatList } from '@/lib/apis'
import { defineStore } from 'pinia'

export const useChatRooms = defineStore('chatRooms', {
  state: (): { chatList: ChatList; loading: boolean } => ({
    chatList: {
      chats: [],
      hasNext: false
    },
    loading: false
  }),
  getters: {
    getChatList: state => state.chatList,
    getLoading: state => state.loading
  },
  actions: {
    async fetchData(limit: number) {
      const prevData = this.chatList
      const { data } = await apis.getChat(limit, prevData.chats.length)

      this.chatList = {
        chats: [...prevData.chats, ...data.chats],
        hasNext: data.hasNext
      }
    },
    setLoading(loading: boolean) {
      this.loading = loading
    }
  }
})
