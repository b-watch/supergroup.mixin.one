const api = require("./net").default;

let Property = {
  async create(key, val) {
    return await api.post("/properties", { key: key, value: val }, {});
  }
};

export default Property;
