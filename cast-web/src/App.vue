<template>
  <v-app class=".app">
    <v-app-bar
      app
      color="#00926F"
      dark
      :height="isPanelExpand ? '256px' : '74px'"
    >
      <v-flex class="mt-1" style="height: 100%">
        <v-layout>
          <v-flex>
            <h1 class="title">{{groupName}}</h1>
            <div class="caption">{{usersCount}} listeners</div>
          </v-flex>
        </v-layout>

        <v-flex v-if="isPanelExpand" class="pt-2 pb-2" >
          <div class="announcement my-2" >{{announcement}}</div>
          <v-divider
            class="my-2"
          ></v-divider>
          <v-flex class="py-2">
            <v-btn outlined block>Join Group</v-btn>
          </v-flex>
        </v-flex>

        <div
          class="text-center"
          style="margin-top: -10px;"
          @click="togglePanel"
        >
          <v-icon v-if="isPanelExpand">mdi-chevron-up</v-icon>
          <v-icon v-else>mdi-chevron-down</v-icon>
        </div>
      </v-flex>
    </v-app-bar>

    <v-content>
      <ChatArea :messages="messages" />
    </v-content>

    <div v-if="disconnected" class="mask flex">
      <div class="hint font-weight-bold">Can’t establish a connection to the server</div>
      <div class="hint">Try to connect in {{(waitingTime - waitingTick)/1000}} sec</div>
    </div>
  </v-app>
</template>

<script>
import ChatArea from './components/ChatArea';
import ws from '@/ws';
import api from '@/api';
import { WS_BASE_URL } from '@/constants';

const wsUrl = WS_BASE_URL + "/messages";

const supportedCategories = {
  'PLAIN_TEXT': 1,
  'PLAIN_IMAGE': 1,
  'PLAIN_VIDEO': 1,
  'PLAIN_AUDIO': 1
}

export default {
  name: 'App',

  components: {
    ChatArea,
  },

  data: () => ({
    //
    groupName: "",
    usersCount: 0,
    announcement: "",
    isPanelExpand: false,
    disconnected: false,
    waitingTick: 0,
    waitingTime: 1000,
    last2ndWaitingTime: 0,
    last1stWaitingTime: 1000,
    messages: []
  }),

  mounted () {
    api.website.amount().then((resp) => {
      this.usersCount = resp.data.users_count
      this.announcement = resp.data.announcement
    })
    api.website.config().then((resp) => {
      this.groupName = resp.data.service_name
    })
    this.connect()
  },

  methods: {
    onClose (evt) {
      console.log('disconnected')
      this.disconnected = true
      this.reconnect()
    },
    onOpen (evt) {
      console.log('connected')
    },
    onMessage (msg) {
      if (supportedCategories.hasOwnProperty(msg.category)) {
        this.messages.push(msg)
        setTimeout(()=> {
          let html = document.documentElement;
          let height = Math.max(html.clientHeight, html.scrollHeight, html.offsetHeight);
          if ((window.innerHeight + window.scrollY) >= document.body.offsetHeight) {
            window.scrollTo(0, height)
          }
        }, 1000)
      }
    },
    onError (err) {
      console.log('error', err)
      this.disconnected = true
    },
    reconnect () {
      // reconnect
      this.waitingTick = 0
      this.waitingTime = this.genNextWaitingTime()
      console.log('wait for', this.waitingTime/1000, 'sec')
      this.retry()
    },
    togglePanel () {
      this.isPanelExpand = !this.isPanelExpand
    },
    connect () {
      this.disconnected = false
      ws.connect(wsUrl, {
        onopen: this.onOpen,
        onmessage: this.onMessage,
        onclose: this.onClose,
        onerror: this.onError,
      })
    },
    retry () {
      setTimeout(() => {
        if (this.waitingTick >= this.waitingTime) {
          this.connect()
        } else {
          this.waitingTick += 1000
          this.retry()
        }
      }, 1000)
    },
    genNextWaitingTime () {
      this.last2ndWaitingTime = this.last1stWaitingTime
      this.last1stWaitingTime = this.waitingTime
      return this.last2ndWaitingTime + this.last1stWaitingTime
    }
  }
};
</script>

<style lang="scss" scoped>
.app {
  background-color: #f3f3f3;
}
.announcement {
  font-size: 14px;
  opacity: 0.5;
  line-height: 1.2em;
  /* max-height = line-height (1.2) * lines max number (5) */
  height: 6em;
  overflow: hidden;
}
.mask {
  position: fixed;
  top: 0;
  bottom: 0;
  left: 0;
  right: 0;
  background: rgba(0, 0, 0, 0.4);
  z-index: 10000;
  display: flex;
  justify-content: center;
  align-items: center;
  flex-direction: column;
  .hint {
    font-size: 16px;
    color: white;
    text-align: center;
  }
}

</style>