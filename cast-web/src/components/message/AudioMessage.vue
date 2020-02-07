<template>
  <div class="audio">
    <div
      class="audio-processing"
      :style="[{ width: `${precent}%` }]"
    />
    <vue-plyr
      ref="plyr"
      :options="options"
    >
      <audio ref="media">
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
import plyrMixin from '@/mixins/plyr'

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
  mixins: [plyrMixin],
  props: {
    message: {
      type: Object,
      default: () => {}
    }
  },
  data() {
    return {
      currentTime: 0,
      duration: 0,
      options: {
        controls
      }
    }
  },
  computed: {
    precent() {
      if (!this.duration) return 0;
      return ((this.currentTime / this.duration) * 100).toFixed(2);
    },
  },
  mounted() {
    this.player.on("timeupdate", this.onTimeUpdate);
  },
  methods: {
    onTimeUpdate() {
      this.currentTime = Number(this.player.currentTime).toFixed();
      this.duration = Number(this.player.duration).toFixed();
    }
  }
}
</script>
<style lang="scss" scoped>
.audio-processing {
  position: absolute;
  height: 100%;
  top: 0;
  left: 0;
  background: #9bcec836;
  transition: width 0.5s ease;
}
</style>