<template>
  <div
    v-if="message"
    class="video"
  >
    <div
      v-touch="{ enable: message, limitClient: true }"
      class="live-wrapper"
    >
      <div class="label">
        LIVE
      </div>
      <div class="pip">
        <v-btn
          icon
          small
          color="#fff"
          @click="handleToggleFloating"
        >
          <v-icon size="18">
            mdi-close-box
          </v-icon>
        </v-btn>
      </div>
      
      <vue-plyr
        ref="plyr"
        :options="options"
        class="live"
      >
        <video
          ref="media"
          controls
          crossorigin
          playsinline
          :poster="thumb"
        >
          <source
            :src="message.attachment.view_url"
            :type="message.attachment.mime_type"
          >
        </video>
      </vue-plyr>
    </div>
  </div>
</template>
<script>
import plyrMixin from '@/mixins/plyr'
export default {
  name: "LiveFloatingMessage",
  mixins: [plyrMixin],
  data () {
    return {
      message: null,
      options: {
        controls: ['play', 'play-large', 'pip', 'fullscreen']
      }
    };
  },
  computed: {
    thumb() {
      return this.message.attachment.thumb_url
    }
  },
  mounted() {
    this.$root.$on('floatLive', (message) => {
      this.message = message
    })
  },
  methods: {
    handleToggleFloating() {
      this.message = null
    }
  }
}
</script>
<style lang="scss" scoped>
</style>
<style lang="scss" scoped>
.video {
  .label {
    position: absolute;
    color: #fff;
    z-index: 1;
    top: 4px;
    left: 4px;
    font-size: 12px;
    font-weight: bold;
    background: rgba(0, 0, 0, 0.6);
    padding: 0px 4px;
    border-radius: 2px;
  }

  .pip {
    position: absolute;
    top: 0px;
    right: 4px;
    z-index: 1;
    color: #fff;
    border-radius: 2px;
    background: rgba(0, 0, 0, 0.6);
  }

  .live-wrapper {
    position: fixed;
    max-width: 280px;
    max-height: 280px;
    z-index: 1;
    top: 10px;
    right: 10px;
  }
}
</style>