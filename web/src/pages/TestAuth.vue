<template>
  <div></div>
</template>

<script>
import { Toast } from "vant";

export default {
  async mounted() {
    const code = this.$route.query.code;
    try {
      let resp = await this.GLOBAL.api.account.authenticate(code);
      console.log("authenticate", resp);
      if (resp.data.authentication_token) {
        if (resp.data.state == "unverified") {
          this.$router.push("/invitation/entry");
        } else {
          this.$router.push("/");
        }
      }
    } catch (err) {
      Toast("OAuth Failed:" + err.toString());
      this.$router.push("/");
    }
  }
};
</script>
