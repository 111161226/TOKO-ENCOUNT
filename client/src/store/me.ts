import { defineStore } from 'pinia'
import api, { UserWithoutPass, UserSimple, User, UserUpdate } from '@/lib/apis'

export const useMe = defineStore('me', {
  state: (): { me: UserWithoutPass | undefined; old: boolean | undefined} => ({
     me: undefined,
     old: false 
    }),
  getters: {
    getMe: state => state.me,
    getOld: state => state.old
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
      this.old = undefined
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
    },
    async deleteMe() {
      await api.deleteMe()
      this.me = undefined
    },
    async checkpresentMe(userData: UserSimple) {
      const { data } = await api.getUserPresent(userData)
      this.old = data.old
    },
    async restoreMe(userData: UserSimple) {
      await api.restoreUser(userData)
      this.login(userData)
      this.old = false
    }
  }
})
