<template>
  <div class="audio">
    <div class="processing" />
    <vue-plyr
      ref="plyr"
      :options="options"
    >
      <audio>
        <source
          :src="message.attachment.view_url"
          :type="message.attachment.mime_type"
        >
      </audio>
    </vue-plyr>
  </div>
</template>
<script>
import Plyr from "plyr";

const controls = `
<div class="plyr__controls">
    <button type="button" class="plyr__control" aria-label="Play, {title}" data-plyr="play">
      <i class="plyr-mdi-icon mdi mdi-play icon--not-pressed" role="presentation"></i>
      <i class="plyr-mdi-icon mdi mdi-pause icon--pressed" role="presentation"></i>
    </button>
    <div class="flex-fill text-left">
      <div class="plyr__time plyr__time--current" aria-label="Current time">00:00</div>
    </div>
</div>
`;

export default {
  name: "AudioMessage",
  props: {
    message: {
      type: Object,
      default: () => {}
    }
  },
  data() {
    return {
      options: {
        controls
      }
    }
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