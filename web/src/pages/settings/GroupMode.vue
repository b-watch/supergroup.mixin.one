<template>
  <loading :loading="loading" :fullscreen="true">
    <div class="group-mode-page">
      <nav-bar
        :title="$t('group_mode.title')"
        :hasTopRight="false"
        :hasBack="true"
      ></nav-bar>

      <van-cell-group
        :title="isAdmin ? $t('group_mode.mode_title') : $t('group_mode.mode_title_alt') "
      >
        <van-cell v-bind:key="mode.name" v-for="mode in modes"
          :title="mode.label" :label="mode.desc" :icon="mode.icon"
          :class="!isAdmin && !mode.selected ? 'faded': ''"
          @click="selectMode(mode)">
          <van-image
            round
            width="32"
            height="32"
            slot="icon"
            :src="mode.icon"
            style="margin: 6px 10px 0 0; border-radius: 99em;"
          />
          <van-icon
            v-if="mode.selected"
            slot="right-icon"
            name="success"
            style="line-height: inherit;"
          />
        </van-cell>
      </van-cell-group>

      <van-cell-group v-if="isAdmin" :title="$t('group_mode.broadcast_title')">
        <van-cell :title="$t('group_mode.broadcast_label')">
          <van-switch v-model="isBroadcast" />
        </van-cell>
        <van-cell :title="broadcastUrl">
          <van-button type="info" plain size="small"
            >{{$t('group_mode.btn_copy_broadcast_url')}}</van-button>
        </van-cell>
      </van-cell-group>

      <van-row v-if="isAdmin" style="padding: 20px">
        <van-col span="24">
          <van-button
            style="width: 100%"
            type="info"
            @click="applyMode()"
            :disabled="!isChanged"
            >{{ $t("group_mode.switch_btn_label") }}</van-button>
        </van-col>
      </van-row>

    </div>

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
  name: "GroupMode",
  props: {
    msg: String
  },
  data() {
    return {
      loading: false,
      isAdmin: false,
      isBroadcast: false,
      isChanged: false,
      currentModeName: 'free',
      modes: [
        {
          name: 'free',
          label: this.$t('group_mode.mode_free_title'),
          icon: require('@/assets/images/mode_free.png'),
          desc: this.$t('group_mode.mode_free_text'),
          selected: false
        },
        {
          name: 'mute',
          label: this.$t('group_mode.mode_mute_title'),
          icon: require('@/assets/images/mode_mute.png'),
          desc: this.$t('group_mode.mode_mute_text'),
          selected: false
        },
        {
          name: 'lecture',
          label: this.$t('group_mode.mode_lecture_title'),
          icon: require('@/assets/images/mode_lecture.png'),
          desc: this.$t('group_mode.mode_lecture_text'),
          selected: false
        }
      ]
    };
  },
  components: {
    NavBar,
    RowSelect,
    Loading,
  },
  async mounted() {
    this.loading = true;
    this.GLOBAL.api.website.config().then((resp) => {
      this.broadcastUrl = resp.data.broadcast_host
    })
    let amountInfo = await this.GLOBAL.api.website.amount();
    if (amountInfo && amountInfo.data) {
      this.currentModeName = amountInfo.data.mode
      for (let ix = 0; ix < this.modes.length; ix++) {
        this.modes[ix];
        this.modes[ix].selected = amountInfo.data.mode === this.modes[ix].name
      }
      this.isBroadcast = amountInfo.data.broadcast === 'on'
    }

    this.isAdmin = window.localStorage.getItem("role") === "admin";
    this.loading = false;
  },
  computed: {
    selectedMode() {
      for (let ix = 0; ix < this.modes.length; ix++) {
        if (this.modes[ix].selected) {
          return this.modes[ix];
        }
      }
      return this.modes[0];
    }
  },
  methods: {
    async applyMode() {
      this.loading = true
      await this.GLOBAL.api.property.create(
        "group-mode-property",
        this.selectedMode.name
      );
      await this.GLOBAL.api.property.create(
        "broadcast-property",
        this.isBroadcast ? 'on' : 'off'
      );
      this.loading = false
      this.$router.back()
    },
    selectMode(mode) {
      if (this.isAdmin) {
        this.modes = this.modes.map(x => {
          x.selected = x.name === mode.name;
          return x;
        });
        this.isChanged = mode.name !== this.currentModeName
      }
    }
  }
};
</script>

<!-- Add "scoped" attribute to limit CSS to this component only -->
<style scoped>
.group-mode-page {
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
.faded {
  display: none;
}
</style>
