import { TOKEN_PREFIX } from '@/constants'

function setItem(key, value) {
  const prefix = TOKEN_PREFIX + '-storage'
  try {
    let storage = JSON.parse(window.localStorage.getItem(prefix))
    if (storage === null) {
      storage = {}
    }
    if (storage) {
      storage[key] = value
    }
    window.localStorage.setItem(prefix, JSON.stringify(storage))
  } catch (err) {
    return null
  }
  return null
}

function getItem(key) {
  const prefix = TOKEN_PREFIX + '-storage'
  try {
    const storage = JSON.parse(window.localStorage.getItem(prefix))
    if (storage) {
      if (!storage.hasOwnProperty(key)) {
        return null
      }
      return storage[key]
    }
  } catch (err) {
    return null
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