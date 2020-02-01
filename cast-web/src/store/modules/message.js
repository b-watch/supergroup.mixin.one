import Vue from 'vue'
import { SOCKET_STATE } from '@/constants'
import getNameColor from '@/utils/getNameColor'

// id: "98cd796f-6337-46e4-a12e-95e376324b7f"
// speaker_name: "divisey"
// speaker_avatar: "https://mixin-images.zeromesh.net/Y1tgxUK6EyJalixzHoUrzpOLHMiJRTOe-xjTwSsd_GPOJqnEKAzn-dA3ghliJB_m_4C9gjrtXXvntuTIS4EeptQ=s256"
// speaker_id: "5467e9ea-cd04-4b91-b84c-93a0c87cb6a4"
// category: "PLAIN_TEXT"
// data: "dGVzdA=="
// text: "test"
// attachment: {id: "", size: 0, mime_type: "", view_url: ""}
// id: ""
// size: 0
// mime_type: ""
// view_url: ""
// created_at: "2020-01-31T07:27:41.130515Z"


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