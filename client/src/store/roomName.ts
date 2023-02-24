import { defineStore } from 'pinia'
import apis from '@/lib/apis'

export const useroomNames = defineStore('roomUsers', {
  state: (): { roomnames: Record<string, string>; loading: boolean } => ({
    roomnames: {},
    loading: false
  }),
  getters: {
    getRooms: state => state.roomnames,
    getRoom: state => (roomId: string) => state.roomnames?.[roomId]
  },
  actions: {
    setRoomName(roomId: string, name: string) {
      this.roomnames = { ...this.roomnames, [roomId]: name }
    },
    setLoading(value: boolean) {
      this.loading = value
    },
    async fetchRoomName(roomId: string) {
      const { data } = await apis.getRoomName(roomId)
      this.roomnames = { ...this.roomnames, [roomId]: data.roomName }
    },
    async updateName(roomId: string, name: string) {
      this.roomnames = { ...this.roomnames, [roomId]: name }
    }
  }
})
