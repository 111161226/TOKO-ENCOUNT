import dayjs from 'dayjs'

export const formatDate = (date: string) => {
  return dayjs(new Date(date)).format('YYYY/MM/DD HH:mm')
}
