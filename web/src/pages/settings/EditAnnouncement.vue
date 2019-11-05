<template>
  <loading :loading="maskLoading" :fullscreen="true">
    <div class="announcement-page">
      <nav-bar
        :title="$t('announcement.title')"
        :hasTopRight="false"
        :hasBack="true"
      ></nav-bar>
      <van-cell-group>
        <van-cell>
          <van-field
            v-model="text"
            rows="3"
            autosize
            type="textarea"
            :label="$t('announcement.label')"
            :placeholder="
              $t('announcement.placeholder', { max: maxMemoLength })
            "
            :maxlength="maxMemoLength"
          ></van-field>
        </van-cell>
      </van-cell-group>
      <van-row style="padding: 20px">
        <van-col span="24">
          <van-button
            style="width: 100%"
            type="info"
            :disabled="!validated"
            @click="send"
            >{{ $t("announcement.send") }}</van-button
          >
        </van-col>
      </van-row>
    </div>
  </loading>
</template>

<script>
import NavBar from "@/components/Nav";
import dayjs from "dayjs";
import Loading from "@/components/Loading";
import { Toast } from "vant";
import utils from "@/utils";

export default {
  name: "Announcement",
  props: {},
  data() {
    return {
      showActionSheet: false,
      maskLoading: false,
      currentMessage: null,
      loading: false,
      finished: false,
      text: ""
    };
  },
  components: {
    NavBar,
    Loading
  },
  async mounted() {
    let websiteInfo = await this.GLOBAL.api.website.amount();
    this.text = websiteInfo.data.announcement;
  },
  computed: {
    maxMemoLength() {
      return 512;
    },
    validated() {
      return this.text.trim().length !== 0;
    }
  },
  methods: {
    async send() {
      await this.GLOBAL.api.property.create(
        "announcement-message-property",
        this.text
      );
      this.$router.replace("/");
    }
  }
};
</script>

<style scoped>
.announcement-page {
  padding-top: 60px;
}
</style>
