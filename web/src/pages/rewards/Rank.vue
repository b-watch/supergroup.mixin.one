<template>
  <loading :loading="loading" :fullscreen="true">
    <div class="reward-rank-page">
      <nav-bar
        :title="$t('rewards.title')"
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
          你共计打赏 {{currentRank.tip_count}} 次，约 ${{currentRank.tip_usd}} 等值数字货币
        </div>
      </div>
      <van-tabs v-model="active">
        <van-tab title="周排名">
          <rank-member-item
            :member="mem" :rank="ix" v-for="mem, ix in ranks.week"
          ></rank-member-item>
        </van-tab>
        <van-tab title="月排名">
          <rank-member-item
            :member="mem" :rank="ix" v-for="mem, ix in ranks.month"
          ></rank-member-item>
        </van-tab>
        <van-tab title="总排名">
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
  margin: 10px 0 20px 0;
  padding: 0 15px;
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
