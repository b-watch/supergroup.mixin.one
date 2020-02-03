<template>
  <loading :loading="loading" :fullscreen="true"> </loading>
</template>

<script>
import Loading from "../components/Loading";
import { mapState } from "vuex";
import utils from "@/utils";

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
    try {
      this.loading = true;
      await this.GLOBAL.api.account.me();
      this.GLOBAL.api.plugin.shortcuts().then(resp => {
        const pluginMap = this.getPlugins(resp);
        const pluginId = this.$route.query.plugin || "";
        const pluginPath = this.$route.query.path || "/";
        if (!pluginMap.hasOwnProperty(pluginId)) {
          this.$router.replace("/");
        } else {
          const url =
            pluginMap[pluginId].url +
            "/#" +
            pluginPath +
            "?token=" +
            localStorage.getItem("token");
          window.location.href = url;
        }
      });
    } catch (err) {
      console.log("error", err);
    }
  },
  methods: {
    getPlugins(resp) {
      if (resp.data && resp.data[0] && resp.data[0].items) {
        var ret = {};
        for (let ix = 0; ix < resp.data[0].items.length; ix++) {
          const ele = resp.data[0].items[ix];
          ret[ele.id] = ele;
        }
        return ret;
      }
      return {};
    }
  }
};
</script>

<style lang="scss" scoped>
</style>

