import { mapState } from 'vuex'

export default {
  props: {
    idx: {
      type: Number
    }
  },
  computed: {
    ...mapState('message', {
      messages: state => state.messages
    }),
    player() {
      return this.$refs.plyr.player
    }
  },
  mounted() {
    this.$root.$on('INST_PLAY', ({ message }) => {
      if (message.id === this.message.id) {
        if (this.player) {
          this.player.play()
        }
      }
    })
    this.$root.$on("INST_MUTE_ALL", ({ message }) => {
      if (message.id !== this.message.id) {
        if (this.player) {
          this.player.stop()
        }
      }
    })
    this.player.on("play", this.onPlay)
    this.player.on("ended", this.onEnd)
  },
  methods: {
    onPlay() {
      this.$root.$emit("INST_MUTE_ALL", { message: this.message, idx: this.idx });
    },
    onEnd() {
      if (this.idx < this.messages.length - 1) {
        const nextMsg = this.getNextPlayableMessage(this.idx + 1)
        if (nextMsg !== null) {
          this.$root.$emit("INST_PLAY", { message: nextMsg })
        }
      }
    },
    getNextPlayableMessage(idx) {
      for (let i = idx; i < this.messages.length; i++) {
        const msg = this.messages[i];
        if (msg.category === "PLAIN_AUDIO" || msg.category === "PLAIN_VIDEO") {
          return msg;
        }
      }
      return null;
    },
  }
}