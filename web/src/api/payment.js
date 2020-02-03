const api = require('./net').default

const Payment = {
  create: async function (params) {
    let resp = await api.post('/payment/create', params, {})
    return resp
  },

  currency: async function () {
    let resp = await api.get(`/payment/currency`, {})
    return resp
  },

  create_wx_pay: async function (params) {
    let resp = await api.post('/wechat/pay/create', params, {})
    return resp
  },

  check_wx_pay: async function (order_id) {
    let resp = await api.get(`/wechat/pay/${order_id}`, {})
    return resp
  }
};

export default Payment
