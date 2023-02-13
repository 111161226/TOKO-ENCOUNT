import { defineStore } from 'pinia'
import apis from '@/lib/apis'

export const useroomNames = defineStore('roomUsers', {
  state: (): { roomnames: Record<string, string> } => ({
    roomnames: {}
  }),
  getters: {
    getRooms: state => state.roomnames,
    getRoom: state => (roomId: string) => state.roomnames?.[roomId]
  },
  actions: {
    setRoomName(roomId: string, name: string) {
      this.roomnames = { ...this.roomnames, [roomId]: name }
    },
    async fetchRoomName(roomId: string) {
      if (!this.roomnames[roomId]) {
        const { data } = await apis.getRoomName(roomId)
        this.roomnames[roomId] = data.name
      }
    }
  }
})
