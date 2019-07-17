const api = require('./net').default

const Invitation = {
  index: async function () {
    return await api.get('/invitations', {})
  },

  create: async function (params) {
    return await api.post('/invitations', params, {})
  },

  apply: async function (code) {
    return await api.put('/invitations/' + code, {}, {})
  }
}
export default Invitation;
