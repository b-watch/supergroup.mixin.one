<template>
  <loading :loading="loading" :fullscreen="true">
    <div class="reward-rank-page">
      <nav-bar
        :title="$t('reward_rank.title')"
        :hasTopRight="false"
        :hasBack="true"
      ></nav-bar>
      <div v-if="currentRank" class="my-rank mb-10">
        <van-image
          round
          width="48px"
          height="48px"
          :src="currentRank.avatar_url"
        />
        <div class="info">
          {{ $t('reward_rank.current_info', {count: currentRank.tip_count, usd: currentRank.tip_usd}) }}
        </div>
      </div>
      <van-row style="padding: 20px">
        <van-col span="24">
          <van-button
            style="width: 100%"
            type="info"
            @click="send"
            >{{ $t("reward_rank.send") }}</van-button
          >
        </van-col>
      </van-row>
      <van-tabs v-model="active">
        <van-tab :title="$t('reward_rank.week')">
          <rank-member-item
            :member="mem" :rank="ix" v-for="mem, ix in ranks.week"
          ></rank-member-item>
        </van-tab>
        <van-tab :title="$t('reward_rank.month')">
          <rank-member-item
            :member="mem" :rank="ix" v-for="mem, ix in ranks.month"
          ></rank-member-item>
        </van-tab>
        <van-tab :title="$t('reward_rank.all')">
          <rank-member-item
            :member="mem" :rank="ix" v-for="mem, ix in ranks.all"
          ></rank-member-item>
        </van-tab>
      </van-tabs>
    </div>
  </loading>
</template>

<script>
import NavBar from "@/components/Nav";
import Loading from "@/components/Loading";
import RankItem from "@/components/partial/RankItem";
import utils from "@/utils";
import { Toast, Dialog } from "vant";
import { CLIENT_ID } from "@/constants";

export default {
  name: "RewardRank",
  props: {
    msg: String
  },
  data() {
    return {
      loading: false,
      currentRank: null,
      ranks: [],
      active: 0,
    };
  },
  components: {
    NavBar,
    Loading,
    "rank-member-item": RankItem
  },
  async mounted() {
    this.loading = true;
    let rankInfo = await this.GLOBAL.api.rewards.ranks();
    if (rankInfo && rankInfo.data) {
      this.currentRank = rankInfo.data.current_rank
      this.ranks = rankInfo.data.ranks
    }
    this.isAdmin = window.localStorage.getItem("role") === "admin";
    this.loading = false;
  },
  computed: {
  },
  methods: {
    send () {
      this.$router.push('/rewards/send')
    },
  }
};
</script>

<!-- Add "scoped" attribute to limit CSS to this component only -->
<style scoped lang="scss">
.reward-rank-page {
  padding-top: 60px;
}
.my-rank {
  display: flex;
  margin: 10px 0 10px 0;
  padding: 0 20px;
  align-items: center;
  .info {
    margin-left: 10px;
  }
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
