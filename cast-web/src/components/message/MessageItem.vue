<template>
  <v-layout
    ref="message"
    :class="[messageClass, 'message']"
  >
    <v-avatar
      height="32"
      width="32"
      min-width="32"
    >
      <img
        v-show="!isSameSpeaker"
        :src="message.speaker_avatar"
        :alt="message.speaker_name"
      >
    </v-avatar>
    <div class="bubble">
      <div
        v-show="!isSameSpeaker"
        class="speaker"
        :style="[{ color: getNameColor(message) }]"
      >
        {{ message.speaker_name }}
      </div>
      <component
        :is="messageComponent"
        :message="message"
      />
      <message-time :message="message" />
      <div
        v-show="!isSameSpeaker"
        class="caret"
      />
    </div>
  </v-layout>
</template>

<script>
import { mapGetters } from 'vuex'
import AudioMessage from './AudioMessage.vue'
import ImageMessage from './ImageMessage.vue'
import TextMessage from './TextMessage.vue'
import VideoMessage from './VideoMessage.vue'
import MessageTime from './MessageTime.vue'

const MESSAGE_CLASS = {
  PLAIN_TEXT: 'message-text',
  PLAIN_IMAGE: 'message-image',
  PLAIN_VIDEO: 'message-video',
  PLAIN_AUDIO: 'message-audio'
}

const MESSAGE_COMPONENT = {
  PLAIN_TEXT: 'text-message',
  PLAIN_IMAGE: 'image-message',
  PLAIN_VIDEO: 'video-message',
  PLAIN_AUDIO: 'audio-message'
}

export const wrapper = Symbol('messageWrapper')

export default {
  name: 'Message',
  components: {
    AudioMessage, 
    ImageMessage,
    TextMessage, 
    VideoMessage,
    MessageTime
  },
  props: {
    message: {
      type: Object,
      default: () => {}
    },
    prev: {
      type: Object,
      default: () => {}
    }
  },
  computed: {
    ...mapGetters('message', ['getNameColor']),
    createdAt() {
      return this.message.created_at
    },
    isSameSpeaker() {
      return this.prev && this.prev.speaker_name === this.message.speaker_name
    },
    messageClass() {
      return MESSAGE_CLASS[this.message.category]
    },
    messageComponent() {
      return MESSAGE_COMPONENT[this.message.category]
    }
  }
};
</script>

<style lang="scss" scoped>
.message {
  margin: 8px 0;

  .bubble {
    margin-left: 15px;
    padding: 4px;
    background: white;
    position: relative;
    box-shadow: 0 0 10px rgba(0,0,0,0.04);
    border-radius: 7px;
    font-size: 14px;

    .speaker {
      margin: 4px;
      font-size: 12px;
      font-weight: bold;
      line-height: 14px;
    }

    .caret {
      position: absolute;
      left: -20px;
      top: 10px;
      width: 0;
      height: 0;
      border-width: 8px 12px 8px 12px;
      border-color: transparent white transparent transparent;
      border-style: solid;
    }
  }
}
</style>