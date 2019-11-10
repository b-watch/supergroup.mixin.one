const api = require("./net").default;

const Rewards = {
  indexRecipients: async function() {
    return await api.get(`/rewards/recipients`, {});
  },

  createRecipient: async function(identityNumber) {
    return await api.post(
      "/rewards/recipients",
      { identity_number: identityNumber },
      {}
    );
  },

  deleteRecipient: async function(userId) {
    return await api.delete("/rewards/recipients/" + userId, {}, {});
  }
};
export default Rewards;
