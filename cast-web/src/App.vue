<template>
  <v-app class="app">
    <div>
      <navbar />
    </div>
    <v-content id="chatContent">
      <chat-area />
    </v-content>
    <system-bar />
    <open-in-browser />
    <page-footer />
  </v-app>
</template>

<script>
import ChatArea from './components/ChatArea';
import Navbar from './components/Navbar'
import SystemBar from './components/SystemBar'
import OpenInBrowser from './components/OpenInBrowser'
import PageFooter from './components/PageFooter'
import { mapActions } from 'vuex'

export default {
  name: 'App',

  components: {
    ChatArea,
    SystemBar,
    Navbar,
    OpenInBrowser,
    PageFooter
  },

  mounted() {
    this.loadConfig()
    this.loadInformation()
    if(window.location.pathname === '/recent') this.loadRecentMessage()
  },

  methods: {
    ...mapActions('group', ['loadConfig', 'loadInformation']),
    ...mapActions('message', ['loadRecentMessage'])
  }
};
</script>

<style lang="scss" scoped>
.app {
  background-color: #f3f3f3;
}
</style>