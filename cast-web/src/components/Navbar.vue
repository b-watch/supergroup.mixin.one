<template>
  <v-app-bar
    class="app-bar align-start"
    dark
    :color="$colors.NAV_COLOR"
    height="auto"
  >
    <div class="flex-fill content">
      <div class="flex-fill d-flex align-center">
        <div class="flex-fill">
          <h1 class="subtitle-1 font-weight-bold">
            {{ groupName }}
          </h1>
          <div class="caption">
            {{ usersCount }} listeners
          </div>
        </div>
        <v-btn
          icon
          @click="togglePanel"
        >
          <v-icon :class="['icon', { 'expand': isPanelExpand }]">
            mdi-chevron-up
          </v-icon>
        </v-btn>
      </div>
      <div
        v-if="isPanelExpand"
        class="py-2"
      >
        <div class="announcement">
          {{ announcement }}
        </div>
        <v-divider
          class="my-2"
        />
        <v-btn
          outlined
          block
          @click="handleJoin"
        >
          Join Group
        </v-btn>
      </div>
      <pinned-message v-else />
    </div>
  </v-app-bar>
</template>
<script>
import PinnedMessage from './message/PinnedMessage'

import { mapState } from 'vuex'
export default {
  name: "Navbar",
  components: {
    PinnedMessage
  },
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
      this.$downloadApp()
    }
  }
}
</script>
<style lang="scss" scoped>
.icon {
  transform: rotate(180deg);

  &.expand {
    transform: rotate(0deg);
  }
}

.announcement {
  font-size: 14px;
  opacity: 0.7;
  overflow: hidden;
}
</style>