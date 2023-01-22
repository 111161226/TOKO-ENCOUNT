import { defineStore } from 'pinia'
import api, { UserWithoutPass, UserSimple, User, UserUpdate } from '@/lib/apis'

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
    async changeMeData(newData: UserUpdate) {
      if (!this.me) {
        throw new Error('not logged in')
      }
      const { data } = await api.patchUserMe(newData)
      this.me = data
    },
    async createMe(userData: User) {
      const { data } = await api.postUser(userData)
      this.me = data
    }
  }
})