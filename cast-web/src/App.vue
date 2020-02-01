<template>
  <v-app class="app">
    <navbar />
    <v-content>
      <chat-area />
    </v-content>
  </v-app>
</template>

<script>
import ChatArea from './components/ChatArea';
import Navbar from './components/Navbar'
import { mapActions } from 'vuex'

export default {
  name: 'App',

  components: {
    ChatArea,
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