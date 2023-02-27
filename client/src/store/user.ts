import apis, { UserList } from '@/lib/apis'
import { defineStore } from 'pinia'

export const useUsers = defineStore('users', {
  state: (): { users: UserList; loading: boolean } => ({
    users: { users: [], hasNext: false },
    loading: false
  }),
  getters: {
    getUsers: state => state.users,
    getLoading: state => state.loading
  },
  actions: {
    initializeUsers() {
      this.users = { users: [], hasNext: false }
    },
    async fetchData(
      limit: number,
      name: string,
      gender: string,
      prefect: string
    ) {
      const prevData = this.users
      const { data } = await apis.getUsers(
        limit,
        prevData.users.length,
        name,
        gender,
        prefect
      )

      this.users = {
        users: [...prevData.users, ...data.users],
        hasNext: data.hasNext
      }
    },
    async fetchAddMembers(limit: number, name: string) {
      const prevData = this.users
      const { data } = await apis.getUsers(
        limit,
        prevData.users.length,
        name,
        '',
        ''
      )

      this.users = {
        users: [...prevData.users, ...data.users],
        hasNext: data.hasNext
      }
    },
    setLoading(value: boolean) {
      this.loading = value
    }
  }
})
