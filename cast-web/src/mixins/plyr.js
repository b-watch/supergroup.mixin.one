import { mapState } from 'vuex'
import Hls from 'hls.js'

export default {
  props: {
    idx: {
      type: Number
    }
  },
  data() {
    return {
      hls: null,
      loading: false,
      player: null
    }
  },
  computed: {
    ...mapState('message', {
      messages: state => state.messages
    }),
    isHls() {
      const source = this.message.attachment.view_url
      const mime = this.message.attachment.mime_type
      return source.endsWith('.m3u8') || 
      source.endsWith('.m3u') || 
      mime === 'application/vnd.apple.mpegurl' || 
      mime === 'audio/mpegurl'
    }
  },
  watch: {
    message() {
      this.$nextTick(() => {
        this.init()
      })
    }
  },
  mounted() {
    this.init()
  },
  methods: {
    init() {
      const plyr = this.$refs.plyr
      this.player = plyr && plyr.player
      if (!this.message) { return }
      const element = this.$refs.media
      const source = this.message.attachment.view_url
      if (this.isHls && Hls.isSupported()) {
        const hls = new Hls({ autoStartLoad: false });
        this.hls = hls
        hls.attachMedia(element);
        hls.on(Hls.Events.MEDIA_ATTACHED, () => {
          hls.loadSource(source);
        });
      }
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
      this.player.on("pause", this.onPause)
    },
    onPlay() {
      if (this.hls) { 
        this.hls.startLoad(-1)
      }
      this.$root.$emit("INST_MUTE_ALL", { message: this.message, idx: this.idx });
    },
    onPause() {
      if (this.hls) { 
        this.hls.stopLoad()
      }
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
  },
  beforeDestroy() {
    if (this.hls) {
      this.hls.destroy()
    }
  }
}