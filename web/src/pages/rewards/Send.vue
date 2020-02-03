<template>
  <loading :loading="loading" :fullscreen="true">
    <div class="rewards-page">
      <nav-bar
        :title="$t('rewards.title')"
        :hasTopRight="false"
        :hasBack="true"
      ></nav-bar>
      <van-cell-group
        v-if="recipients.length !== 0"
        :title="$t('rewards.recipient_section_title')"
      >
        <van-swipe-cell v-bind:key="user.user_id" v-for="user in recipients">
          <van-cell :title="user.full_name" @click="selectUser(user)">
            <van-image
              round
              width="28"
              height="28"
              slot="icon"
              :src="user.avatar_url"
              style="margin-right: 10px; border-radius: 99em;"
            />
            <van-icon
              v-if="user.selected"
              slot="right-icon"
              name="success"
              style="line-height: inherit;"
            />
          </van-cell>
          <template v-if="isAdmin" slot="right">
            <van-button
              square
              type="danger"
              text="Remove"
              @click="removeUser(user)"
            />
          </template>
        </van-swipe-cell>
      </van-cell-group>
      <template v-else>
        <div v-if="!loading" style="text-align: center; margin-bottom: 20px;">
          <van-icon name="warning-o" size="64"></van-icon>
          <div class>{{ $t("rewards.no_recipient") }}</div>
        </div>
      </template>
      <van-cell-group v-if="isAdmin">
        <van-cell
          :title="$t('rewards.add_label')"
          icon="plus"
          @click="showAddDialog = true"
        ></van-cell>
      </van-cell-group>

      <van-cell-group
        v-if="selectedUser"
        :title="$t('rewards.rewards_section_title')"
      >
        <row-select
          :index="0"
          :title="$t('rewards.select_assets')"
          :columns="assets"
          placeholder="Tap to Select"
          @change="onChangeAsset"
        >
          <span slot="text">{{
            selectedAsset ? selectedAsset.text : "Tap to Select"
          }}</span>
        </row-select>
        <van-cell>
          <van-field
            type="number"
            v-model="form.amount"
            :label="$t('rewards.amount')"
            :placeholder="$t('rewards.placeholder_amount', { min: minAmount })"
          >
            <span slot="right-icon">{{
              selectedAsset ? selectedAsset.symbol : ""
            }}</span>
          </van-field>
        </van-cell>
        <van-cell title=" " :value="esitmatedValue"></van-cell>
      </van-cell-group>
      <van-row v-if="selectedUser" style="padding: 20px">
        <van-col span="24">
          <van-button
            style="width: 100%"
            type="info"
            :disabled="!validated"
            @click="pay"
            >{{ $t("rewards.pay") }}</van-button
          >
        </van-col>
      </van-row>
      <van-row style="padding: 20px">
        <van-col span="24">
          <van-button
            style="width: 100%"
            type="info"
            plain
            @click="rank"
            >{{ $t("rewards.rank") }}</van-button
          >
        </van-col>
      </van-row>
    </div>
    <van-dialog
      v-model="showAddDialog"
      title="Add Recipientss"
      show-cancel-button
      @confirm="addUser"
    >
      <div style="padding: 20px 0;">
        <van-cell-group>
          <van-cell>
            <van-field
              v-model="addUserId"
              placeholder="identity number or user id"
            ></van-field>
          </van-cell>
        </van-cell-group>
      </div>
    </van-dialog>
  </loading>
</template>

<script>
import NavBar from "@/components/Nav";
import RowSelect from "@/components/RowSelect";
import Row from "@/components/Nav";
import Loading from "@/components/Loading";
import uuid from "uuid";
import utils from "@/utils";
import { Toast, Dialog } from "vant";
import { CLIENT_ID } from "@/constants";

export default {
  name: "SendReward",
  props: {
    msg: String
  },
  data() {
    return {
      loading: false,
      showAddDialog: false,
      isAdmin: false,
      rewardsMinAmountBase: "0.01",
      addUserId: "",
      recipients: [],
      assets: [],
      selectedAsset: null,
      form: {
        amount: ""
      }
    };
  },
  components: {
    NavBar,
    RowSelect,
    Loading
  },
  async mounted() {
    this.loading = true;
    let confInfo = await this.GLOBAL.api.website.config();
    if (confInfo && confInfo.data) {
      this.rewardsMinAmountBase = confInfo.data.rewards_min_amount_base;
    }
    let prepareInfo = await this.GLOBAL.api.account.assets("rewards");
    if (prepareInfo) {
      this.assets = prepareInfo.data.assets.map(x => {
        x.text = `${x.symbol} (${x.balance})`;
        return x;
      });
      if (this.assets.length) {
        this.selectedAsset = this.assets[0];
        this.form.memo = this.$t("rewards.default_memo", {
          symbol: this.selectedAsset.symbol
        });
      }
    }
    let recipientsInfo = await this.GLOBAL.api.rewards.indexRecipients();
    if (recipientsInfo && recipientsInfo.data) {
      this.recipients = recipientsInfo.data.map(x => {
        x.selected = false;
        return x;
      });
      if (this.recipients.length) {
        this.recipients[0].selected = true;
      }
    }
    this.isAdmin = window.localStorage.getItem("role") === "admin";
    this.loading = false;
  },
  computed: {
    validated() {
      if (
        this.form.amount &&
        this.selectedAsset &&
        parseFloat(this.form.amount) >= parseFloat(this.minAmount)
      ) {
        return true;
      }
      return false;
    },
    esitmatedValue() {
      let val = 0;
      if (this.form.amount && this.selectedAsset) {
        val = this.form.amount * this.selectedAsset.price_usd;
      }
      return "≈$" + val.toLocaleString();
    },
    minAmount() {
      const base = this.rewardsMinAmountBase; // 1 usd
      if (this.selectedAsset) {
        return (base / this.selectedAsset.price_usd).toFixed(4);
      }
      return 0;
    },
    selectedUser() {
      for (let ix = 0; ix < this.recipients.length; ix++) {
        if (this.recipients[ix].selected) {
          return this.recipients[ix];
        }
      }
      return this.recipients[0];
    }
  },
  methods: {
    async pay() {
      let memo = JSON.stringify({
        a: "rewards",
        p1: this.selectedUser.user_id,
        p2: this.selectedAsset.symbol
      });
      window.location.href = `mixin://pay?recipient=${CLIENT_ID}&asset=${
        this.selectedAsset.asset_id
      }&amount=${this.form.amount}&trace=${uuid.v4()}&memo=${encodeURIComponent(
        memo
      )}`;
      Dialog.confirm({
        title: this.$t("rewards.dialog_confrim_title"),
        confirmButtonText: this.$t("rewards.dialog_confrim_ok"),
        cancelButtonText: this.$t("rewards.dialog_confrim_cancel"),
        message: this.$t("rewards.dialog_confrim_desc", {
          name: this.selectedUser.full_name
        })
      })
        .then(() => {
          window.close();
        })
        .catch(() => {});
    },
    rank () {
      this.$router.push('/rewards/rank')
    },
    onChangeAsset(ix) {
      this.selectedAsset = this.assets[ix];
      this.form.memo = this.$t("rewards.default_memo", {
        symbol: this.selectedAsset.symbol
      });
    },
    async addUser() {
      try {
        const resp = await this.GLOBAL.api.rewards.createRecipient(
          this.addUserId
        );
        utils.reloadPage();
      } catch (err) {
        Toast(`User ${this.addUserId} not found.`);
      }
    },
    async removeUser(user) {
      try {
        const resp = await this.GLOBAL.api.rewards.deleteRecipient(
          user.user_id
        );
        utils.reloadPage();
      } catch (err) {
        Toast(`User ${user.full_name} not found.`);
      }
    },
    selectUser(user) {
      this.recipients = this.recipients.map(x => {
        x.selected = x.user_id === user.user_id;
        return x;
      });
    }
  }
};
</script>

<!-- Add "scoped" attribute to limit CSS to this component only -->
<style scoped>
.rewards-page {
  padding-top: 60px;
}
h3 {
  margin: 40px 0 0;
}
ul {
  list-style-type: none;
  padding: 0;
}
li {
  display: inline-block;
  margin: 0 10px;
}
a {
  color: #42b983;
}
</style>
