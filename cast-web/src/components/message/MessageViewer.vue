<template>
  <v-dialog
    v-model="dialog"
    fullscreen
  >
    <v-card
      v-if="message"
      flat
      tile
      class="viewer-card pa-5"
    >
      <message-item
        v-touch="{ swipe: handleSwipe }"
        :message="message"
      />
    </v-card>
  </v-dialog>
</template>
<script>
import MessageItem from './MessageItem'

export default {
  name: "MessageViewer",
  components: {
    MessageItem
  },
  data () {
    return {
      dialog: false,
      message: null
    };
  },
  methods: {
    show(message) {
      this.message = message
      this.dialog = true
    },
    handleSwipe() {
      this.exit()
    },
    exit() {
      this.dialog = false
      this.message = null
    }
  }
}
</script>
<style lang="scss" scoped>
.viewer-card {
  height: 100%;
  display: flex;
  flex-direction: column;

  .message-wrapper {
    height: 100%;
  }
}
</style>