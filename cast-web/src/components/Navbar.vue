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
      return this.config.service_name || '大群广播'
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