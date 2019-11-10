const api = require('./net').default

const Invitation = {
  index: async function (isHistory) {
    return await api.get(`/invitations?history=${isHistory}`, {})
  },

  create: async function () {
    return await api.post('/invitations', {}, {})
  },

  apply: async function (code) {
    return await api.put('/invitations/' + code, {}, {})
  },

  checkRule: async function () {
    return await api.get(`/invite_rule`, {})
  },
}
export default Invitation;
