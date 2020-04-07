const api = require("./net").default;
const storage = require("@/utils/localStorage").default;

const Account = {
  me: async function() {
    return await api.get("/me", {});
  },

  assets: async function(scene) {
    return await api.get("/me/assets?scene=" + scene, {});
  },

  subscribe: async function() {
    return await api.post("/subscribe", {}, {});
  },

  unsubscribe: async function() {
    return await api.post("/unsubscribe", {}, {});
  },

  subscribers: async function(t = 0, q = "") {
    return await api.get("/subscribers?offset=" + t + "&q=" + q, {});
  },

  remove: async function(id) {
    return await api.post("/users/" + id + "/remove", {}, {});
  },

  block: async function(id) {
    return await api.post("/users/" + id + "/block", {}, {});
  },

  authenticate: async function(authorizationCode) {
    var params = {
      code: authorizationCode
    };
    let resp = await api.post("/auth", params, {});
    if (resp.data) {
      storage.setItem("token", resp.data.authentication_token);
      storage.setItem("user_id", resp.data.user_id);
      storage.setItem("role", resp.data.role);
    }
    return resp;
  },

  create_wx_pay: async function(params) {
    let resp = await api.post("/wechat/pay/create", params, {});
    return resp;
  },

  check_wx_pay: async function(order_id) {
    let resp = await api.get(`/wechat/pay/${order_id}`, {}, {});
    return resp;
  },

  userId: function() {
    return storage.getItem("user_id");
  },

  role: function() {
    return storage.getItem("role");
  },

  token: function() {
    return storage.getItem("token");
  },

  clear: function(callback) {
    storage.clear();
    if (typeof callback === "function") {
      callback();
    }
  }
};

export default Account;
