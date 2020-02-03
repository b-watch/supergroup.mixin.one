
const api = require('./net').default

const Website = {
  amount: async function () {
    return await api.get('/amount', {})
  },
  config: async function () {
    let resp = await api.get('/config', {})
    if (resp.data) {
      window.localStorage.setItem('cfg_client_id', resp.data.mixin_client_id);
      window.localStorage.setItem('cfg_host', resp.data.host);
      window.localStorage.setItem('cfg_invite_to_join', resp.data.invite_to_join);
    }
    return resp
  }
};

export default Website;
