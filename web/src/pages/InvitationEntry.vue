<template>
  <div class="container">
    <h2>{{this.$t("invitation.welcome") + configData.service_name}}</h2>
    <h3>{{configData.home_welcome_message}}</h3>
    <div class="notice">
      {{this.$t('invitation.notice.title')}}
      <br />
      {{this.$t('invitation.notice.line1')}}
      <br />
      {{this.line2}}
    </div>
    <van-field v-model="code" :placeholder="$t('invitation.input_placeholder')" autosize />
    <van-button class="button" type="info" @click="apply">{{this.$t("invitation.submit")}}</van-button>
  </div>
</template>

<script>
export default {
  name: "InvitationEntry",

  data() {
    return {
      code: "",
      configData: {
        service_name: "",
        home_welcome_message: "",
        pay_to_join: false,
        auto_estimate_base: "",
        auto_estimate_currency: "",
        accept_asset_list: []
      },
      currencyData: []
    };
  },

  mounted() {
    this.GLOBAL.api.website
      .config()
      .then(conf => {
        if (conf.data) {
          this.configData = conf.data;
        }
      })
      .then(() => {
        this.GLOBAL.api.payment.currency().then(res => {
          if (res.data) {
            this.currencyData = res.data;
          }
        });
      });
  },

  computed: {
    line2() {
      if (this.configData.pay_to_join) {
        return this.$t("invitation.notice.line2_2")
          .replace("{amount}", `${this.amount}`)
          .replace("{symbol}", this.symbol);
      } else {
        return this.$t("invitation.notice.line2_1");
      }
    },
    amount() {
      if (
        this.configData &&
        this.currencyData &&
        this.currencyData.length > 0
      ) {
        const target = this.currencyData.find(c => {
          return c.symbol == this.symbol;
        });
        if (target) {
          switch (this.configData.auto_estimate_currency) {
            case "cny":
              return (
                Number(this.configData.auto_estimate_base) /
                Number(target.price_cny)
              );
            case "usd":
              return (
                Number(this.configData.auto_estimate_base) /
                Number(target.price_usd)
              );
            default:
              return 0;
          }
        }
        return 0;
      }
    },
    symbol() {
      if (
        this.configData.accept_asset_list &&
        this.configData.accept_asset_list.length > 0
      ) {
        return this.configData.accept_asset_list[0].symbol;
      }
      return "";
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
.container {
  margin-left: 2rem;
  margin-right: 2rem;

  h2 {
    margin-top: 0;
    padding-top: 10rem;
    text-align: center;
  }

  h3 {
    text-align: center;
    font-size: 14px;
    color: rbga(0, 0, 0);
    opacity: 0.6;
    margin-bottom: 3rem;
    font-weight: normal;
  }

  .notice {
    font-size: 14px;
    color: #b32424;
    text-align: left;
    background-color: #fffae0;
    border: 1px solid #e2b4a5;
    border-radius: 3px;
    padding: 1rem;
    margin-bottom: 1.5rem;
  }

  .van-cell.van-field {
    padding: 0.5rem;
    border-radius: 4px;
    margin-bottom: 1.5rem;
    border: 1px solid #dddddd;
  }
  .button {
    width: 100%;
    border-radius: 4px;
    height: 3rem;
  }
}
</style>

