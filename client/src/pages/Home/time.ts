import { ChatData } from '@/lib/apis'
import dayjs from 'dayjs'
import relativeTime from 'dayjs/plugin/relativeTime'

dayjs.extend(relativeTime)

export const getRelativeTime = (date: string) => {
  const past = dayjs(date)
  const now = dayjs()

  if (past.isSame(now, 'day')) {
    return past.format('HH:mm')
  } else {
    return past.fromNow()
  }
}

/** comparing function for sorting ChatData[] */
export const compareFunc = (a: ChatData, b: ChatData) => {
  const aTime = dayjs(a.latestMessage.createdAt)
  const bTime = dayjs(b.latestMessage.createdAt)

  if (aTime.isBefore(bTime)) {
    return 1
  } else if (aTime.isAfter(bTime)) {
    return -1
  } else {
    return 0
  }
}
