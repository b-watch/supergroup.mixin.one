<template>
  <loading :loading="maskLoading" :fullscreen="true">
    <div class="members-page">
      <nav-bar :title="$t('members.title')" :hasTopRight="false" :hasBack="true"></nav-bar>
      <van-cell>
        <van-field
          placeholder="Search"
          left-icon="search"
          @change="searchEnter"
          v-model="searchQuery"
        ></van-field>
      </van-cell>
      <van-list v-model="loading" :finished="finished" finished-text="~ END ~" @load="onLoad">
        <div v-if="admins && admins.length !== 0" class="list-label">{{$t('members.list_label_admins')}}</div>
        <member-item :member="item" v-for="item in admins" role="admin" @member-click="memberClick"></member-item>
        <div v-if="lecturers && lecturers.length !== 0" class="list-label">{{$t('members.list_label_lecturers')}}</div>
        <member-item :member="item" v-for="item in lecturers" role="lecturer" @member-click="memberClick"></member-item>
        <div class="list-label">{{$t('members.list_label_subscribers')}}</div>
        <member-item :member="item" v-for="item in users" @member-click="memberClick"></member-item>
      </van-list>
      <van-action-sheet
        :title="currentMember ? currentMember.full_name : ''"
        v-model="showActionSheet"
        :actions="actions"
        :cancel-text="$t('comm.cancel')"
        @select="onSelectAction"
        @cancel="onCancelAction"
      />
    </div>
  </loading>
</template>

<script>
import NavBar from "@/components/Nav";
import dayjs from "dayjs";
import MemberItem from "@/components/partial/MemberItem";
import Loading from "@/components/Loading";
import { ActionSheet, Toast } from "vant";
import utils from "@/utils";
import storage from '@/utils/localStorage'

export default {
  name: "Members",
  props: {},
  data() {
    return {
      searchQuery: "",
      showActionSheet: false,
      maskLoading: false,
      currentMember: null,
      loading: false,
      finished: false,
      users: [],
      admins: [],
      lecturers: [],
      actions: [
        { name: this.$t("members.kick") },
        { name: this.$t("members.block") }
      ]
    };
  },
  components: {
    NavBar,
    MemberItem,
    Loading,
    "van-action-sheet": ActionSheet
  },
  async mounted() {},
  computed: {
    lastOffset() {
      if (this.users.length) {
        let d = new Date(this.users[this.users.length - 1].subscribed_at);
        d.setSeconds(d.getSeconds() + 1);
        return d.toISOString();
      }
      return 0;
    }
  },
  methods: {
    async onLoad() {
      await this.loadMembers(this.lastOffset, "");
    },
    async loadMembers(offset = 0, query = "", append = true) {
      this.maskLoading = true;
      this.loading = true;
      let resp = await this.GLOBAL.api.account.subscribers(offset, query);
      if (resp.data.users.length < 2) {
        this.finished = true;
      }
      const adminSet = {}
      const admins = resp.data.admins.map(x => {
        x.time = dayjs(x.subscribed_at).format("YYYY/MM/DD");
        x.subscribed = !dayjs(x.subscribed_at).isBefore(dayjs("1900-01-01"));
        adminSet[x.user_id] = 1;
        return x;
      });
      const lecutreSet = {}
      const lecturers = resp.data.lecturers.map(x => {
        x.time = dayjs(x.subscribed_at).format("YYYY/MM/DD");
        x.subscribed = !dayjs(x.subscribed_at).isBefore(dayjs("1900-01-01"));
        lecutreSet[x.user_id] = 1;
        return x;
      });
      const users = resp.data.users.map(x => {
        x.time = dayjs(x.subscribed_at).format("YYYY/MM/DD");
        x.subscribed = !dayjs(x.subscribed_at).isBefore(dayjs("1900-01-01"));
        return x;
      }).filter((x) => {
        return !adminSet.hasOwnProperty(x.user_id) && !lecutreSet.hasOwnProperty(x.user_id)
      });

      this.admins = admins
      this.lecturers = lecturers
      if (append) {
        this.users = this.users.concat(users);
      } else {
        this.users = users;
        this.finished = true;
      }
      this.loading = false;
      this.maskLoading = false;
    },
    memberClick(mem) {
      if (storage.getItem("role") === "admin") {
        this.currentMember = mem;
        this.showActionSheet = true;
      }
    },
    async onSelectAction(item, ix) {
      if (this.currentMember) {
        let mem = this.currentMember;
        this.maskLoading = true;
        if (ix === 0) {
          let result = await this.GLOBAL.api.account.remove(mem.user_id);
          if (result.error) {
            this.maskLoading = false;
            return;
          }
          utils.reloadPage();
        } else if (ix === 1) {
          let result = await this.GLOBAL.api.account.block(mem.user_id);
          if (result.error) {
            this.maskLoading = false;
            return;
          }
          utils.reloadPage();
        } else {
          this.maskLoading = false;
        }
      }
      this.showActionSheet = false;
    },
    onCancelAction(item, ix) {
      this.showActionSheet = false;
    },
    searchEnter() {
      this.loadMembers(0, this.searchQuery, false);
      this.finished = true;
    }
  }
};
</script>

<style scoped>
.members-page {
  padding-top: 60px;
}
.list-label {
  padding: 2px 15px;
  font-size: 12px;
  opacity: 0.6;
}
</style>
