import api from '@/api'

const state = {
  config: '',
  information: ''
}

const mutations = {
  SET_CONFIG(state, config) {
    state.config = config
  },
  SET_INFORMATION(state, information) {
    state.information = information
  }
}

const actions = {
  async loadConfig({ commit }) {
    const res = await api.website.config()
    commit('SET_CONFIG', res.data)
  },
  async loadInformation({ commit }) {
    const res = await api.website.amount()
    commit('SET_INFORMATION', res.data)
  }
}

export default {
  namespaced: true,
  state,
  mutations,
  actions
}