<template>
  <div
    class="video"
    :class="[{ 'floating': floating }]"
  >
    <div>
      <div class="label">
        LIVE
      </div>
      <v-img
        min-width="200"
        :src="thumb"
        class="thumb"
      />
    </div>
    <div
      v-touch="{ enable: floating, limitClient: true }"
      class="live-wrapper"
    >
      <div class="label">
        LIVE
      </div>
      <div
        v-if="pip && !pipNativeSupport"
        class="pip"
      >
        <v-btn
          icon
          small
          color="#fff"
          @click="handleToggleFloating"
        >
          <v-icon
            v-if="!floating"
            size="18"
          >
            mdi-picture-in-picture-top-right
          </v-icon>
          <v-icon
            v-else
            size="18"
          >
            mdi-close-box
          </v-icon>
        </v-btn>
      </div>
      
      <vue-plyr
        ref="plyr"
        :options="options"
        class="live"
      >
        <video>
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
  name: "VideoMessage",
  mixins: [plyrMixin],
  props: {
    pip: {
      type: Boolean,
      default: false
    },
    message: {
      type: Object,
      default: () => {}
    }
  },
  data () {
    return {
      floating: false,
      pipNativeSupport: false,
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
    this.pipNativeSupport = 'pictureInPictureEnabled' in document
  },
  methods: {
    handleToggleFloating() {
      this.floating = !this.floating
    }
  }
}
</script>
<style lang="scss" scoped>
.video {
  position: relative;
  z-index: 1;

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

  .thumb {
    border-radius: 2px;
    overflow: hidden;
    visibility: hidden;
    position: absolute;
  }

  .live-wrapper {
    position: relative;
  }

  &.floating {
    .thumb {
      position: relative;
      visibility: visible;
    }

    .live-wrapper {
      position: fixed;
      max-width: 280px;
      z-index: 11;
      top: 10px;
      right: 10px;
    }
  }
}
</style>
