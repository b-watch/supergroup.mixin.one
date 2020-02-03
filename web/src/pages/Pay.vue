<template>
  <loading :loading="loading" :fullscreen="true">
    <div class="pay-page" style="padding-top: 60px;">
      <nav-bar :title="$t('pay.title')" :hasTopRight="false"></nav-bar>
      <van-panel
        :title="$t('pay.welcome')"
        :desc="$t('pay.welcome_desc')"
      ></van-panel>
      <br />
      <van-panel :title="$t('pay.method_crypto')">
        <row-select
          :index="0"
          :title="$t('pay.select_assets')"
          :columns="assets"
          placeholder="Tap to Select"
          @change="onChangeAsset"
        >
          <span slot="text">{{
            selectedAsset ? selectedAsset.text : "Tap to Select"
          }}</span>
        </row-select>
        <van-cell
          :title="
            $t('pay.price_label', {
              price: currentCryptoPrice,
              unit: selectedAsset ? selectedAsset.text : '...'
            })
          "
        >
          <span v-if="currentEstimatedPrice"
            >≈{{ currencySymbol
            }}{{ currentEstimatedPrice.toLocaleString() }}</span
          >
        </van-cell>
        <div slot="footer">
          <van-cell>
            <van-button
              style="width: 100%"
              type="info"
              :disabled="selectedAsset === null || loading"
              @click="payCrypto"
              >{{ $t("pay.pay_crypto") }}</van-button
            >
          </van-cell>
          <!-- <van-cell>
          <van-button style="width: 100%" type="warning" :disabled="selectedAsset === null" @click="payCrypto">{{$t('pay.pay_foxone')}}</van-button>
          </van-cell>-->
        </div>
      </van-panel>
      <br />
      <van-panel v-if="acceptWechatPayment" :title="$t('pay.method_wechat')">
        <van-cell
          :title="
            $t('pay.price_label', {
              price: wechatPaymentAmount,
              unit: $t('currency.' + autoEstimateCurrency)
            })
          "
        ></van-cell>
        <div slot="footer">
          <van-cell>
            <van-button
              style="width: 100%"
              type="primary"
              @click="payWechatMobile"
              >{{ $t("pay.pay_wechat") }}</van-button
            >
          </van-cell>
        </div>
      </van-panel>
      <br />
      <van-panel v-if="acceptCouponPayment" :title="$t('pay.method_coupon')">
        <van-cell>
          <van-field
            :placeholder="$t('pay.coupon_placeholder')"
            v-model="couponCode"
          ></van-field>
        </van-cell>
        <div slot="footer">
          <van-cell>
            <van-button
              style="width: 100%"
              type="info"
              plain
              @click="payCoupon"
              :disabled="loading"
              >{{ $t("pay.pay_coupon") }}</van-button
            >
          </van-cell>
        </div>
      </van-panel>
    </div>
  </loading>
</template>

<script>
import NavBar from "@/components/Nav";
import RowSelect from "@/components/RowSelect";
import Loading from "../components/Loading";
import Row from "@/components/Nav";
import uuid from "uuid";
import { Toast } from "vant";
import { CLIENT_ID, WEB_ROOT } from "@/constants";

export default {
  name: "Pay",
  props: {
    msg: String
  },
  data() {
    return {
      loading: true,
      config: null,
      meInfo: null,
      selectedAsset: null,
      autoEstimate: false,
      autoEstimateCurrency: "usd",
      acceptWechatPayment: false,
      acceptCouponPayment: false,
      wechatPaymentAmount: "100",
      cryptoEsitmatedUsdMap: {},
      currencyRates: [],
      currentCryptoPrice: 0,
      currentEstimatedPrice: 0,
      couponCode: "",
      assets: []
    };
  },
  components: {
    NavBar,
    RowSelect,
    Loading
  },
  async mounted() {
    this.loading = true;
    let config = await this.GLOBAL.api.website.config();
    console.log(config);
    this.assets = config.data.accept_asset_list.map(x => {
      x.text = x.symbol;
      return x;
    });
    this.selectedAsset = this.assets[0];
    this.autoEstimate = config.data.auto_estimate;
    this.autoEstimateCurrency = config.data.auto_estimate_currency;
    this.autoEstimateBase = config.data.auto_estimate_base;
    this.acceptWechatPayment = config.data.accept_wechat_payment;
    this.wechatPaymentAmount = config.data.wechat_payment_amount;
    this.acceptCouponPayment = config.data.accept_coupon_payment;

    if (this.autoEstimate) {
      let currencyInfo = await this.GLOBAL.api.payment.currency();
      for (let ix = 0; ix < currencyInfo.data.length; ix++) {
        const ele = currencyInfo.data[ix];
        this.currencyRates[ele.symbol] = ele;
      }
    }

    this.meInfo = await this.GLOBAL.api.account.me();
    setTimeout(() => {
      this.updatePrice();
      this.loading = false;
    }, 2000);
  },
  computed: {
    currencySymbol() {
      if (this.autoEstimate) {
        if (this.autoEstimateCurrency === "cny") return "¥";
        if (this.autoEstimateCurrency === "usd") return "$";
      }
      return "";
    }
  },
  methods: {
    async payCrypto() {
      this.loading = true;
      // get order info
      let orderInfo = await this.GLOBAL.api.payment.create({
        method: "crypto",
        asset_id: this.selectedAsset.asset_id,
        user_id: this.meInfo.data.user_id
      });
      let orderId = orderInfo.data.order.order_id;
      let amount = orderInfo.data.order.amount;
      let assetId = orderInfo.data.order.asset_id;
      setTimeout(async () => {
        await this.waitForPayment();
      }, 1000);
      window.location.href = `mixin://pay?recipient=${CLIENT_ID}&asset=${assetId}&amount=${amount}&trace=${orderId}&memo=PAY_TO_JOIN`;
      console.log(
        `mixin://pay?recipient=${CLIENT_ID}&asset=${assetId}&amount=${amount}&trace=${orderId}&memo=PAY_TO_JOIN`
      );
    },
    async onChangeAsset(ix) {
      this.loading = true;
      this.selectedAsset = this.assets[ix];
      await this.updatePrice();
      this.loading = false;
    },
    async updatePrice() {
      if (this.selectedAsset.amount === "auto") {
        let priceCny = this.currencyRates[this.selectedAsset.symbol].price_cny;
        let priceUsd = this.currencyRates[this.selectedAsset.symbol].price_usd;
        if (this.autoEstimateCurrency === "usd") {
          this.currentCryptoPrice = (this.autoEstimateBase / priceUsd).toFixed(
            8
          );
          this.currentEstimatedPrice = this.autoEstimateBase;
        } else {
          this.currentCryptoPrice = (this.autoEstimateBase / priceCny).toFixed(
            8
          );
          this.currentEstimatedPrice = this.autoEstimateBase;
        }
      } else {
        this.currentCryptoPrice = parseFloat(this.selectedAsset.amount).toFixed(
          8
        );
        this.currentEstimatedPrice = 0;
      }
    },
    async waitForPayment() {
      let meInfo = await this.GLOBAL.api.account.me();
      if (meInfo.error) {
        setTimeout(async () => {
          await this.waitForPayment();
        }, 1500);
        return;
      }
      if (meInfo.data.state === "paid") {
        Toast(this.$t("pay.success_toast"));
        this.$router.push("/");
        this.loading = false;
        return;
      }
      setTimeout(async () => {
        await this.waitForPayment();
      }, 1500);
    },
    payWechatMobile() {
      this.$router.push(
        `/pay/wxqr/?qr_url=${encodeURIComponent(
          WEB_ROOT + "/wechat/request/" + this.meInfo.data.user_id
        )}`
      );
    },
    async payCoupon() {
      this.loading = true;
      try {
        let resp = await this.GLOBAL.api.coupon.occupy(this.couponCode);
        if (resp && resp.data) {
          Toast(this.$t("pay.correct_coupon_code_toast"));
          this.$router.push("/");
          this.loading = false;
        } else {
          Toast(this.$t("pay.incorrect_coupon_code_toast"));
          this.loading = false;
        }
      } catch (err) {
        Toast(this.$t("pay.incorrect_coupon_code_toast"));
        this.loading = false;
      }
    }
  }
};
</script>

<style scoped>
.pay-page {
  padding-top: 60px;
}
</style>
