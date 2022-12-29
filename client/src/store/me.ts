import { defineStore } from 'pinia'
import api, { UserWithoutPass, UserSimple, User } from '@/lib/apis'

export const useMe = defineStore('me', {
  state: (): { me: UserWithoutPass | undefined } => ({ me: undefined }),
  getters: {
    getMe: state => state.me
  },
  actions: {
    async setMe() {
      const { data } = await api.getUserMe()
      this.me = data
    },
    async login(userData: UserSimple) {
      const { data } = await api.postLogin(userData)
      this.me = data
    },
    async logout() {
      await api.postLogout()
      this.me = undefined
    },
    async createMe(userData: User) {
      const { data } = await api.postUser(userData)
      this.me = data
    }
  }
})
