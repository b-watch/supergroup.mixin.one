<template>
  <v-app-bar
    class="app-bar"
    app
    dark
    :color="$colors.NAV_COLOR"
    :height="isPanelExpand ? '256px' : '74px'"
  >
    <v-flex
      class="mt-1"
      style="height: 100%"
    >
      <v-layout>
        <v-flex>
          <h1 class="title">
            {{ groupName }}
          </h1>
          <div class="caption">
            {{ usersCount }} listeners
          </div>
        </v-flex>
      </v-layout>

      <v-flex
        v-if="isPanelExpand"
        class="pt-2 pb-2"
      >
        <div class="announcement my-2">
          {{ announcement }}
        </div>
        <v-divider
          class="my-2"
        />
        <v-flex class="py-2">
          <v-btn
            outlined
            block
            @click="handleJoin"
          >
            Join Group
          </v-btn>
        </v-flex>
      </v-flex>
      <div
        class="text-center"
        style="margin-top: -10px;"
        @click="togglePanel"
      >
        <v-icon v-if="isPanelExpand">
          mdi-chevron-up
        </v-icon>
        <v-icon v-else>
          mdi-chevron-down
        </v-icon>
      </div>
    </v-flex>
  </v-app-bar>
</template>
<script>
import { mapState } from 'vuex'
export default {
  name: "Navbar",

  data () {
    return {
      isPanelExpand: false
    };
  },
  computed: {
    ...mapState('group', {
      config: state => state.config,
      information: state => state.information
    }),
    groupName() {
      return this.config.service_name
    },
    usersCount() {
      return this.information.users_count
    },
    announcement() {
      return this.information.announcement
    }
  },
  watch: {
    groupName() {
      document.title = this.groupName
    }
  },
  mounted() {
    document.title = this.groupName
  },
  methods: {
    togglePanel() {
      this.isPanelExpand = !this.isPanelExpand
    },
    handleJoin() {
      this.$socket._onmessage({ data: '{"id":"9535dac1-8495-49c6-94dd-02206fcc46f3","speaker_name":"divisey","speaker_avatar":"https://mixin-images.zeromesh.net/Y1tgxUK6EyJalixzHoUrzpOLHMiJRTOe-xjTwSsd_GPOJqnEKAzn-dA3ghliJB_m_4C9gjrtXXvntuTIS4EeptQ=s256","speaker_id":"5467e9ea-cd04-4b91-b84c-93a0c87cb6a4","category":"PLAIN_TEXT","data":"dGVzdA==","text":"test","attachment":{"id":"","size":0,"mime_type":"","view_url":""},"created_at":"2020-02-02T12:39:59.913099Z"}' })
      // window.location.href = 'https://mixin.one/codes/'
    }
  }
}
</script>
<style lang="scss" scoped>
.announcement {
  font-size: 14px;
  opacity: 0.7;
  line-height: 1.2em;
  height: 6em;
  overflow: hidden;
}
</style>