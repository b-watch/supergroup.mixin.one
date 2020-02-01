<template>
  <vue-plyr
    ref="plyr"
    :options="options"
  >
    <video>
      <source
        :src="message.attachment.view_url"
        :type="message.attachment.mime_type"
      >
    </video>
  </vue-plyr>
</template>
<script>
export default {
  name: "VideoMessage",
  props: {
    message: {
      type: Object,
      default: () => {}
    }
  },
  data () {
    return {
      options: {
        controls: ['play-large', 'play', 'progress', 'current-time', 'mute', 'fullscreen']
      }
    };
  },
  computed: {
    player() {
      return this.$refs.plyr.player
    }
  },
  mounted() {
    this.player.on("play", this.onPlay);
    this.player.on("ended", this.onEnd);
    this.player.on("timeupdate", this.onTimeUpdate);
  },
  methods: {
    onPlay() {},
    onEnd() {},
    onTimeUpdate() {}
  }
}
</script>
<style lang="scss" scoped>
</style>