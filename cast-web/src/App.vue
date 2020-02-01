<template>
  <v-app class="app">
    <navbar />
    <v-content>
      <chat-area />
    </v-content>
    <system-bar />
  </v-app>
</template>

<script>
import ChatArea from './components/ChatArea';
import Navbar from './components/Navbar'
import SystemBar from './components/SystemBar'
import { mapActions } from 'vuex'

export default {
  name: 'App',

  components: {
    ChatArea,
    SystemBar,
    Navbar
  },

  mounted() {
    this.loadConfig()
    this.loadInformation()
  },

  methods: {
    ...mapActions('group', ['loadConfig', 'loadInformation'])
  }
};
</script>

<style lang="scss" scoped>
.app {
  background-color: #f3f3f3;
}
</style>