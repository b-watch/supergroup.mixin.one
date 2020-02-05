const state = {
  systemBar: true,
  pageFooter: true
}

const mutations = {
  setSystemBar(state, value) {
    state.systemBar = value 
  },
  setPageFooter(state, value) {
    state.pageFooter = value
  }
}

export default {
  namespaced: true,
  state,
  mutations,
}