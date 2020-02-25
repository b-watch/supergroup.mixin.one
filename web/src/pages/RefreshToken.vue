<template>
  <loading :loading="loading" :fullscreen="true"> </loading>
</template>

<script>
import Loading from "../components/Loading";
import { mapState } from "vuex";
import utils from "@/utils";
import { CLIENT_ID, OAUTH_CALLBACK_URL } from "@/constants";

export default {
  data() {
    return {
      loading: false
    };
  },
  computed: {
    isZh() {
      return this.$i18n.locale.indexOf("zh") !== -1;
    }
  },
  components: {
    Loading
  },
  async mounted() {
    this.redirectToOAuth()
  },
  methods: {
    redirectToOAuth() {
      let url = `https://mixin-www.zeromesh.net/oauth/authorize?client_id=${CLIENT_ID}&scope=PROFILE:READ+ASSETS:READ+MESSAGES:REPRESENT&response_type=code&return_to=${encodeURIComponent(
        OAUTH_CALLBACK_URL
      )}`;
      window.location.href = url;
    }
  }
};
</script>

<style lang="scss" scoped>
</style>

