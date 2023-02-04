import { defineStore } from 'pinia'

interface RoomUser {
  userId: string
  userName: string
}

export const useroomUsers = defineStore('roomUsers', {
  state: (): { users: Record<string, RoomUser[]>; roomnames: Record<string, string> } => ({
    users: {},
    roomnames: {}
  }),
  getters: {
    getUsers: state => state.users,
    getUser: state => (roomId: string) => state.users?.[roomId],
    getRooms: state => state.roomnames,
    getRoom: state => (roomId: string) => state.roomnames?.[roomId]
  },
  actions: {
    setUser(roomId: string, userName: string, userId: string) {
      if (!this.users[roomId]) {
        this.users = {
          ...this.users,
          [roomId]: [{ userId: userId, userName: userName }]
        }
        return
      }
      const prevData = this.users[roomId]
      this.users = {
        ...this.users,
        [roomId]: prevData.concat({ userId: userId, userName: userName })
      }
    },
    setRoomName(roomId: string, name: string) {
      this.roomnames = { ...this.roomnames, [roomId]: name }
    }
  }
})
