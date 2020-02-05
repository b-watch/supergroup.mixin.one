import Vue from 'vue'
import { SOCKET_STATE } from '@/constants'
import getNameColor from '@/utils/getNameColor'

const MAX_MESSAGE_AMOUNT = 200

const state = {
  messages: [],
  state: SOCKET_STATE.DISCONNECT,
  nameColorMap: {},
  hasNewMessage: false
}

const getters = {
  getNameColor(state) {
    return message => {
      const color = getNameColor(message.speaker_name || '', state.nameColorMap)
      Vue.set(state.nameColorMap, message.speaker_name, color)
      return color
    }
  },
  getMessageById(state) {
    return (id) => {
      return state.messages.find(m => m.id === id)
    }
  }
}

const mutations = {
  SET_HAS_NEW_MESSAGE(state, value) {
    state.hasNewMessage = value
  },
  ADD_MESSAGE(state, message) {
    state.messages = [...state.messages, message]
    if (state.messages.length > MAX_MESSAGE_AMOUNT) {
      state.messages.shift()
    }
  },
  CONNECT_FAILED(state) {
    state.state = SOCKET_STATE.DISCONNECT
  },
  CONNECTED(state) {
    state.state = SOCKET_STATE.CONNECTED
  },
  CONNECTING(state) {
    state.state = SOCKET_STATE.CONNECTING
  }
}

const actions = {}

export default {
  namespaced: true,
  state,
  getters,
  mutations,
  actions
}