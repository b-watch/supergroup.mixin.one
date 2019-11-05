<template>
  <loading :loading="loading" :fullscreen="true">
    <div class="group-mode-page">
      <nav-bar
        :title="$t('group_mode.title')"
        :hasTopRight="false"
        :hasBack="true"
      ></nav-bar>

      <van-cell-group
        :title="$t('group_mode.mode_title')"
      >
        <van-cell v-bind:key="mode.name" v-for="mode in modes"
          :title="mode.label" :label="mode.desc" :icon="mode.icon"
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

      <van-row style="padding: 20px">
        <van-col span="24">
          <van-button
            style="width: 100%"
            type="info"
            @click="applyMode()"
            >{{ $t("group_mode.switch") }}</van-button>
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
      modes: [
        {
          name: 'free',
          label: 'Chit-Chat',
          icon: require('@/assets/images/mode_free.png'),
          desc: 'Everyone is chatting in this mode',
          selected: false
        },
        {
          name: 'mute',
          label: 'Mute',
          icon: require('@/assets/images/mode_mute.png'),
          desc: 'Speakers speak, others send rewards and lucky coins',
          selected: false
        },
        {
          name: 'lecture',
          label: 'Lecture',
          icon: require('@/assets/images/mode_lecture.png'),
          desc: 'Only speakers speak',
          selected: false
        }
      ]
    };
  },
  components: {
    NavBar,
    RowSelect,
    Loading
  },
  async mounted() {
    this.loading = true;
    let amountInfo = await this.GLOBAL.api.website.amount();
    if (amountInfo && amountInfo.data) {
      for (let ix = 0; ix < this.modes.length; ix++) {
        this.modes[ix];
        this.modes[ix].selected = amountInfo.data.mode === this.modes[ix].name
      }
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
      console.log('switch to', this.selectedMode)
      await this.GLOBAL.api.property.create(
        "group-mode-property",
        this.selectedMode.name
      );
    },
    selectMode(mode) {
      this.modes = this.modes.map(x => {
        x.selected = x.name === mode.name;
        return x;
      });
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
</style>
