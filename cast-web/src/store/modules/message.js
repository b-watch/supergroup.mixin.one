import Vue from 'vue'
import { SOCKET_STATE } from '@/constants'
import getNameColor from '@/utils/getNameColor'

const state = {
  messages: [],
  state: SOCKET_STATE.DISCONNECT,
  nameColorMap: {}
}

const getters = {
  getNameColor(state) {
    return message => {
      const color = getNameColor(message.speaker_name || '', state.nameColorMap)
      Vue.set(state.nameColorMap, message.speaker_name, color)
      return color
    }
  }
}

const mutations = {
  ADD_MESSAGE(state, message) {
    state.messages = [...state.messages, message]
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