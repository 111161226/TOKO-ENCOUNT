import { defineStore } from 'pinia'
import apis, { ChatList, Message } from '@/lib/apis'
import { reorderArray } from '@/util/reorderArray'

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
    },
    /** websocket 以外からは呼ばない */
    catchNewMessage(roomId: string, message: Message) {
      const idx = this.chatList.chats.findIndex(c => c.roomId === roomId)
      if (idx === -1) {
        // room が見つからない場合
        this.chatList.chats.unshift({
          roomId,
          name: message.userName,
          latestMessage: message,
          newMessageCount: 1
        })
        return
      }

      const prevData = this.chatList.chats[idx]
      this.chatList.chats[idx] = {
        ...prevData,
        latestMessage: message,
        newMessageCount: prevData.newMessageCount + 1
      }
      this.chatList.chats = reorderArray(this.chatList.chats, idx, 0)
    }
  }
})
