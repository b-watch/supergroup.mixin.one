<template>
  <v-scroll-y-transition>
    <v-system-bar
      v-if="!isBroadcasting"
      fixed
      window
      dark
      color="error"
      class="bottom-bar"
      @click="handleSymtemBarClick"
    >
      <span>大群广播模式未开启</span>
    </v-system-bar>
    <v-system-bar
      v-show="show"
      v-else
      fixed
      window
      dark
      :color="stateColor"
      class="bottom-bar"
      @click="handleSymtemBarClick"
    >
      <template v-if="hasNewMessage">
        <v-icon>mdi-message</v-icon>
        <span>有新信息</span>
        <v-spacer />
      </template>
      <template v-else>
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
      </template>
    </v-system-bar>
  </v-scroll-y-transition>
</template>
<script>
import { mapState, mapMutations } from 'vuex'
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
  computed: {
    ...mapState('app', {
      systemBar: state => state.systemBar
    }),
    ...mapState('group', {
      information: state => state.information
    }),
    ...mapState('message', {
      state: state => state.state,
      hasNewMessage: state => state.hasNewMessage
    }),
    isBroadcasting() {
      return this.information && this.information.broadcast === 'on'
    },
    show () {
      return this.systemBar
    },
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
    },
    newMsgText() {
      return this.hasNewMessage ? '有新信息' : ''
    },
    newMsgIcon() {
      return this.hasNewMessage ? 'mdi-message' : ''
    },
    text() {
      return newMsgText || this.stateText
    },
    icon() {
      return newMsgIcon || this.stateIcon
    }
  },
  watch: {
    state(val) {
      if (val === SOCKET_STATE.CONNECTED) {
        setTimeout(() => {
          this.setSystemBar(false)
        }, 1000)
      } else if (val === SOCKET_STATE.DISCONNECT) {
        this.setSystemBar(true)
      }
    },
    hasNewMessage(value) {
      this.setSystemBar(value)
    }
  },
  methods: {
    ...mapMutations('app', ['setSystemBar']),
    handleReconnect() {
      this.$socket.reconnect()
    },
    handleCheckMsg() {
      this.$root.$emit('CHECK_NEW_MESSAGE')
    },
    handleSymtemBarClick() {
      if (this.hasNewMessage) {
        this.handleCheckMsg()
      }
    }
  }
}
</script>
<style lang="scss" scoped>
.bottom-bar {
  bottom: 0!important;
  top: auto;
  z-index: 6;
}
</style>