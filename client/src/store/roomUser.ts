import { defineStore } from 'pinia'

interface RoomUser {
    userId: string
    userName: string
}

export const useroomUsers = defineStore('roomUsers', {
  state: (): { users: Record<string, RoomUser>} => ({
    users: {}
  }),
  getters: {
    getUsers: state => state.users,
    getUser: state => (roomId: string) => state.users?.[roomId]
  },
  actions: {
    setUser(roomId: string, userName: string, userId: string) {
      this.users = { ...this.users, 
                    [roomId]: {userId: userId, userName: userName} 
                }
    }
  }
})