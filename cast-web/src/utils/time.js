import dayjs from 'dayjs'

export default {
  format(t, p = "HH:mm") {
    return dayjs(t).format(p)
  }
}