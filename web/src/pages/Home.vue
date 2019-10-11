<template>
  <loading :loading="loading" :fullscreen="true">
    <div class="home home-page page">
      <van-panel
        :title="welcomeMessage || $t('home.welcome')"
        :desc="
          $t('home.welcome_desc', {
            count: websiteInfo ? websiteInfo.data.users_count : '...'
          })
        "
      >
        <div class="panel-header" slot="header">
          <h1>{{ websiteConf ? websiteConf.data.service_name : "..." }}</h1>
          <div class="announcement">
            <p v-html="announcementText"></p>
            <van-button
              v-if="isAdmin"
              type="default"
              size="mini"
              round
              @click="gotoEditAnnouncement"
              icon="edit"
            >
            </van-button>
          </div>
          <div class="btns">
            <van-button
              type="primary"
              :type="isSubscribed ? 'primary' : 'warning'"
              size="small"
              round
              :plain="!isSubscribed"
              hairline
              @click="togglSubscribe"
              >{{
                isSubscribed
                  ? this.$t("home.op_subscribed")
                  : this.$t("home.op_unsubscribed")
              }}</van-button
            >
            <van-button
              type="default"
              size="small"
              round
              @click="gotoMembers"
              icon="friends-o"
              >{{
                websiteInfo ? websiteInfo.data.users_count : "..."
              }}</van-button
            >
            <van-button
              v-if="isAdmin"
              :type="isProhibited ? 'danger' : 'primary'"
              plain
              hairline
              size="small"
              round
              class="icon-btn"
              @click="toggleProhibit"
              :icon="isProhibited ? 'close' : 'comment-circle-o'"
            ></van-button>
            <van-button
              v-if="isAdmin"
              type="default"
              size="small"
              round
              class="icon-btn"
              @click="gotoMessages"
              icon="comment-o"
            ></van-button>
          </div>
        </div>
      </van-panel>
      <br />
      <template v-for="group in shortcutsGroups">
        <van-panel :title="group.label">
          <cell-table
            :items="group.shortcuts"
            @external="openExternalLink"
          ></cell-table>
        </van-panel>
        <br />
      </template>
      <van-panel :title="$t('home.pane_operations')">
        <cell-table :items="builtinItems"></cell-table>
      </van-panel>
    </div>
  </loading>
</template>

<script>
import NavBar from "../components/Nav";
import CellTable from "../components/CellTable";
import Loading from "../components/Loading";
import { mapState } from "vuex";
import { Dialog } from "vant";
import AssetItem from "@/components/partial/AssetItem";
import utils from "@/utils";

export default {
  data() {
    return {
      loading: false,
      meInfo: null,
      welcomeMessage: "",
      websiteInfo: null,
      websiteConf: null,
      isProhibited: false,
      isSubscribed: false,
      builtinItems: [
        // builtin
        {
          icon: require("../assets/images/luckymoney-circle.png"),
          label: this.$t("home.op_luckycoin"),
          url: "/packets/prepare"
        }
      ],
      invitationItem: {
        icon: require("../assets/images/invitation.png"),
        label: this.$t("invitation.entry"),
        url: "/invitation/details"
      },
      rewardsItem: {
        icon: require("../assets/images/rewards.png"),
        label: this.$t("rewards.entry"),
        url: "/rewards"
      },
      couponsItem: {
        icon: require("../assets/images/coupons.png"),
        label: this.$t("home.op_coupons"),
        url: "/coupons"
      },
      shortcutsGroups: []
    };
  },
  computed: {
    isAdmin() {
      if (this.meInfo) {
        return this.meInfo.data.role === "admin";
      }
      return false;
    },
    isZh() {
      return this.$i18n.locale.indexOf("zh") !== -1;
    },
    announcementText() {
      if (this.websiteInfo) {
        return utils.urlify(this.websiteInfo.data.announcement);
      }
      return "";
    }
  },
  components: {
    NavBar,
    CellTable,
    Loading
  },
  async mounted() {
    try {
      this.loading = true;
      this.websiteInfo = await this.GLOBAL.api.website.amount();
      this.meInfo = await this.GLOBAL.api.account.me();

      this.isProhibited = this.websiteInfo.data.prohibited;
      this.isSubscribed =
        new Date(this.meInfo.data.subscribed_at).getYear() > 1;

      this.websiteConf = await this.GLOBAL.api.website.config();
      if (this.websiteConf.data.home_shortcut_groups) {
        this.shortcutsGroups = this.addToGroups(
          this.websiteConf.data.home_shortcut_groups,
          false
        );
      }
      this.welcomeMessage = this.websiteConf.data.home_welcome_message;
      this.loading = false;

      if (this.meInfo.data.state === "pending") {
        this.$router.push("/pay");
        return;
      }
      // tips visbility
      if (this.websiteConf.data.rewards_enable) {
        this.builtinItems.push(this.rewardsItem);
      }
      // invitation visbility
      if (this.websiteConf.data.invite_to_join) {
        this.builtinItems.push(this.invitationItem);
      }
      // plugins
      this.GLOBAL.api.plugin.shortcuts().then(resp => {
        if (
          resp.data &&
          resp.data[0] &&
          resp.data[0].items &&
          resp.data[0].items.length
        ) {
          let aa = this.buildShortcuts(resp.data[0].items, true);
          this.builtinItems = this.builtinItems.concat(aa);
          console.log(this.builtinItems);
        }
      });
    } catch (err) {
      console.log("error", err);
    }
  },
  methods: {
    openExternalLink(item) {
      console.log("openExternalLink 2");
      this.loading = true;
      window.location.href = item.url;
      setTimeout(() => {
        this.loading = false;
      }, 5000);
    },
    addToGroups(groups, isPlugin) {
      return groups.map(x => {
        x.label = this.isZh ? x.label_zh : x.label_en;
        const items = x.items || x.shortcuts;
        x.shortcuts = this.buildShortcuts(items, isPlugin);
        return x;
      });
    },
    buildShortcuts(items, isPlugin) {
      return items
        .map(z => {
          z.label = this.isZh ? z.label_zh : z.label_en;
          if (isPlugin) {
            // for plugin SPA, use hash mode to pass query
            z.isPlugin = true;
            z.url +=
              "/#/?token=" + encodeURIComponent(localStorage.getItem("token"));
          }
          return z;
        })
        .filter(z => {
          if (isPlugin && z.admin_only === true) {
            return this.isAdmin;
          }
          return true;
        });
    },
    handlePluginRedirect(groupId, itemId) {
      return () => {
        this.GLOBAL.api.plugin
          .redirect(groupId, itemId)
          .then(resp => {
            console.log(resp);
          })
          .catch(err => {
            console.log(err);
          });
      };
    },
    gotoEditAnnouncement() {
      this.$router.push("/announcement/edit");
    },
    gotoMembers() {
      this.$router.push("/members");
    },
    gotoMessages() {
      this.$router.push("/messages");
    },
    async togglSubscribe() {
      if (this.isSubscribed) {
        Dialog.confirm({
          message: this.$t("home.op_unsubscribe_confirm_msg")
        })
          .then(async () => {
            await this.GLOBAL.api.account.unsubscribe();
            this.isSubscribed = false;
          })
          .cancel(() => {});
      } else {
        await this.GLOBAL.api.account.subscribe();
        this.isSubscribed = true;
      }
    },
    async toggleProhibit() {
      if (this.isProhibited) {
        await this.GLOBAL.api.property.create(
          "prohibited-message-property",
          "false"
        );
        this.isProhibited = false;
      } else {
        await this.GLOBAL.api.property.create(
          "prohibited-message-property",
          "true"
        );
        this.isProhibited = true;
      }
      return;
    }
  }
};
</script>

<style lang="scss" scoped>
.home-page {
  // padding-top: 60px;
}
.panel-header {
  padding: 15px;
  h1 {
    font-size: 18px;
    margin: 0 0 10px 0;
  }
  .announcement {
    margin-bottom: 15px;
  }
  .btns {
    margin-top: 10px;
    .van-button {
      margin-right: 6px;
    }
    .van-button.icon-btn {
      min-width: 16px;
      padding: 0;
      min-height: 16px;
      width: 30px;
      height: 30px;
    }
  }
}
</style>

