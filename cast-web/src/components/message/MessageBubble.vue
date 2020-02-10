<template>
  <v-layout
    ref="message"
    :class="[messageClass, 'message', {'focus': isFocus}]"
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
      <quote-message :message="message" />
      <message-item
        :message="message"
        v-bind="$attrs"
      />
      <message-time :message="message" />
    </div>
  </v-layout>
</template>

<script>
import { mapGetters } from 'vuex'
import MessageTime from './MessageTime'
import MessageItem from './MessageItem'
import QuoteMessage from './QuoteMessage'

const MESSAGE_CLASS = {
  PLAIN_TEXT: 'message-text',
  PLAIN_IMAGE: 'message-image',
  PLAIN_VIDEO: 'message-video',
  PLAIN_AUDIO: 'message-audio',
  PLAIN_LIVE: 'message-live'
}

export default {
  name: 'MessageBubble',
  components: {
    MessageItem,
    MessageTime,
    QuoteMessage
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
  data() {
    return {
      isFocus: false
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
    }
  },
  mounted() {
    this.$root.$on('focusMessage', (message) => {
      if (this.message.id === message.id) {
        this.isFocus = true
        setTimeout(() => {
          this.isFocus = false
        }, 1000)
        const el = this.$refs.message
        const content = document.querySelector('#chatContent')
        content.scrollTop = el.offsetTop
      }
    })
  }
};
</script>

<style lang="scss" scoped>
.message {
  margin: 8px 0;
  position: relative;

  &::after {
    content: '';
    z-index: -1;
    position: absolute;
    top: -4px;
    bottom: -4px;
    right: -12px;
    left: -12px;
    transition: all 0.4s ease;
  }

  &.focus {
    z-index: 1;

    &::after {
      background-color: #ECEFF1;
    }
  }

  .bubble {
    margin-left: 15px;
    padding: 4px;
    background: white;
    position: relative;
    box-shadow: 0 0 10px rgba(0,0,0,0.04);
    border-radius: 7px;
    font-size: 14px;
    overflow: hidden;

    .speaker {
      margin: 4px;
      font-size: 12px;
      font-weight: bold;
      line-height: 14px;
    }
  }
}

.message-image, 
.message-video, 
.message-live {
  .bubble {
    max-width: 240px;
  }
}
</style>