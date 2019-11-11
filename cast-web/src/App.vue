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
  </v-app>
</template>

<script>
import ChatArea from './components/ChatArea';
import ws from '@/ws';
import api from '@/api';
import { WS_BASE_URL } from '@/constants';

const wsUrl = WS_BASE_URL + "/messages";

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
    messages: [
    ]
  }),

  async mounted () {
    api.website.amount().then((resp) => {
      this.usersCount = resp.data.users_count
      this.announcement = resp.data.announcement
    })
    api.website.config().then((resp) => {
      this.groupName = resp.data.service_name
    })
    ws.init(wsUrl, {
      onopen: this.onOpen,
      onmessage: this.onMessage,
      onerror: this.onError,
    })
  },

  methods: {
    onOpen (evt) {
      console.log('connected')
    },
    onMessage (msg) {
      this.messages.push(msg)
      setTimeout(()=> {
        let html = document.documentElement;
        let height = Math.max(html.clientHeight, html.scrollHeight, html.offsetHeight);
        if ((window.innerHeight + window.scrollY) >= document.body.offsetHeight) {
          window.scrollTo(0, height)
        }
      }, 1000)
    },
    onError (evt) {
      console.log(evt.data)
    },
    togglePanel () {
      this.isPanelExpand = !this.isPanelExpand
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

</style>