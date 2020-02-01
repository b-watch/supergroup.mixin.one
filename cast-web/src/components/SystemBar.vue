<template>
  <v-scroll-y-transition>
    <v-system-bar
      v-show="show"
      fixed
      window
      dark
      :color="stateColor"
      class="bottom-bar"
    >
      <v-icon
        v-if="stateIcon"
        @click="show = false"
      >
        {{ stateIcon }}
      </v-icon>
      <v-progress-circular
        v-show="connecting"
        indeterminate
        size="16"
        width="1"
        class="mr-2"
        color="rgba(255, 255, 255, 0.7)"
      />
      <span>{{ stateText }}</span>
      <v-spacer />
      <v-btn
        v-if="disconnected"
        text
        small
        @click="handleReconnect"
      >
        重新连接
      </v-btn>
    </v-system-bar>
  </v-scroll-y-transition>
</template>
<script>
import { mapState } from 'vuex'
import { SOCKET_STATE } from '@/constants'

const STATE_META = {
  [SOCKET_STATE.DISCONNECT]: {
    text: '与消息服务器的连接已断开',
    color: 'error',
    icon: 'mdi-close'
  },
  [SOCKET_STATE.CONNECTED]: {
    text: '已连接到消息服务器',
    color: 'success',
    icon: 'mdi-check'
  },
  [SOCKET_STATE.CONNECTING]: {
    text: '正在连接消息服务器',
    color: 'warning',
  }
}

export default {
  name: "SystemBar",
  data () {
    return {
      show: true
    };
  },
  computed: {
    ...mapState('message', {
      state: state => state.state
    }),
    connecting() {
      return this.state === SOCKET_STATE.CONNECTING
    },
    disconnected() {
      return this.state === SOCKET_STATE.DISCONNECT
    },
    stateText() {
      return STATE_META[this.state].text
    },
    stateColor() {
      return STATE_META[this.state].color
    },
    stateIcon() {
      return STATE_META[this.state].icon
    }
  },
  watch: {
    state(val) {
      if (val === SOCKET_STATE.CONNECTED) {
        setTimeout(() => {
          this.show = false
        }, 1000)
      } else if (val === SOCKET_STATE.DISCONNECT) {
        this.show = true
      }
    }
  },
  methods: {
    handleReconnect() {
      this.$socket.reconnect()
    }
  }
}
</script>
<style lang="scss" scoped>
.bottom-bar {
  bottom: 0!important;
  top: auto;
}
</style>