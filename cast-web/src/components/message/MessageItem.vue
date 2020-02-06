<template>
  <component
    :is="messageComponent"
    :message="message"
    v-bind="$attrs"
  />
</template>

<script>
import { mapGetters } from 'vuex'
import AudioMessage from './AudioMessage.vue'
import ImageMessage from './ImageMessage.vue'
import TextMessage from './TextMessage.vue'
import VideoMessage from './VideoMessage.vue'
import UnSupportMessage from './UnSupportMessage.vue'
import FileMessage from './FileMessage.vue'
import LiveMessage from './LiveMessage.vue'

const MESSAGE_COMPONENT = {
  PLAIN_TEXT: 'text-message',
  PLAIN_IMAGE: 'image-message',
  PLAIN_VIDEO: 'video-message',
  PLAIN_AUDIO: 'audio-message',
  PLAIN_DATA: 'file-message',
  PLAIN_LIVE: 'live-message'
}

export default {
  name: 'MessageItem',
  components: {
    AudioMessage, 
    ImageMessage,
    TextMessage, 
    VideoMessage,
    UnSupportMessage,
    FileMessage,
    LiveMessage
  },
  props: {
    message: {
      type: Object,
      default: () => {}
    }
  },
  computed: {
    messageComponent() {
      if (Object.keys(MESSAGE_COMPONENT).includes(this.message.category)) {
        return MESSAGE_COMPONENT[this.message.category]
      }
      return 'un-support-message'
    }
  }
};
</script>