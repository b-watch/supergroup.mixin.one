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
      <div class="close-btn">
        <v-btn
          icon
          @click="exit"
        >
          <v-icon>mdi-close</v-icon>
        </v-btn>
      </div>
      <message-item
        v-touch="{ swipe: handleSwipe, enable: touchless, scale: true }"
        :message="message"
        class="viewer"
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
      message: null,
      touchless: false
    };
  },
  mounted() {
    this.$root.$on('viewMessage', (message, { touchless = false } = {}) => {
      this.show(message, touchless)
    })
  },
  methods: {
    show(message, touchless) {
      this.message = message
      this.dialog = true
      this.touchless = touchless
    },
    handleSwipe() {
      this.exit()
    },
    exit() {
      this.dialog = false
      this.message = null
      this.touchless = false 
    }
  }
}
</script>
<style lang="scss" scoped>
.viewer-card {
  height: 100%;
  position: relative;
  display: flex;
  flex-direction: column;
  justify-content: center;

  .close-btn {
    position: absolute;
    top: 10px;
    right: 10px;
  }

  .message-wrapper {
    height: 100%;
  }
}
</style>