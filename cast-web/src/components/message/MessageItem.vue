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
import MessageTime from './MessageTime.vue'

const MESSAGE_COMPONENT = {
  PLAIN_TEXT: 'text-message',
  PLAIN_IMAGE: 'image-message',
  PLAIN_VIDEO: 'video-message',
  PLAIN_AUDIO: 'audio-message'
}

export const wrapper = Symbol('messageWrapper')

export default {
  name: 'MessageItem',
  components: {
    AudioMessage, 
    ImageMessage,
    TextMessage, 
    VideoMessage
  },
  props: {
    message: {
      type: Object,
      default: () => {}
    }
  },
  computed: {
    messageComponent() {
      return MESSAGE_COMPONENT[this.message.category]
    }
  }
};
</script>