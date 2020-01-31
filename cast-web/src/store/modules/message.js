import { SOCKET_STATE } from '@/constants'

const state = {
  messages: [],
  state: SOCKET_STATE.DISCONNECT
}

const getters = {}

const mutations = {
  ADD_MESSAGE(state, message) {
    state.message = [...state.message, message]
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