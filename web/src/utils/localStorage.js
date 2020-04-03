import { TOKEN_PREFIX } from '@/constants'

function setItem(key, value) {
  const prefix = TOKEN_PREFIX + '-storage'
  const storage = window.localStorage.getItem(prefix)
  if (storage) {
    storage[key] = value
  }
  window.localStorage.setItem(prefix, JSON.stringify(storage))
  return null
}

function getItem(key) {
  const prefix = TOKEN_PREFIX + '-storage'
  const storage = JSON.parse(window.localStorage.getItem(prefix))
  if (storage) {
    if (!storage.hasOwnProperty(key)) {
      return null
    }
    return storage[key]
  }
  return null
}

function clear() {
  const prefix = TOKEN_PREFIX + '-storage'
  window.localStorage.setItem(prefix, JSON.stringify(null))
}

export default {
  setItem,
  getItem,
  clear
}