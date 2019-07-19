<template>
  <div>
    <h2>{{this.$t("invitation.welcome")}}</h2>
    <div class="action">
      <van-field v-model="code" :placeholder="placeholder" autosize />
      <van-button class="button" type="info" size="small" @click="apply">{{this.$t("invitation.verify")}}</van-button>
    </div>
  </div>
</template>

<script>
export default {
  name: "InvitationEntry",

  data() {
    return {
      placeholder: this.$t("invitation.code"),
      meInfo: null,
      code: ""
    };
  },

  computed: {
    fullName() {
      if (this.meInfo) {
        return this.meInfo.data.full_name;
      } else {
        return "";
      }
    }
  },

  methods: {
    apply() {
      this.GLOBAL.api.invitation
        .apply(this.code)
        .then(response => {
          if (response && response.data) {
            this.$toast(this.$t("invitation.code_available"));
            this.$router.push("/pay");
          } else {
            this.$toast(this.$t("invitation.code_unavailable"));
          }
        })
        .catch(error => {
          this.$toast(this.$t("invitation.code_unavailable"));
        });
    }
  }
};
</script>

<style lang="scss" scoped>
h2 {
  padding-top: 10rem;
  text-align: center;
}

.action {
  margin-left: 2rem;
  margin-right: 2rem;
}

.button {
  margin-top: 2rem;
  width: 100%;
  border-radius: 4px;
}
</style>

