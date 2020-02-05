<template>
  <div
    v-if="show"
    class="quote-message"
    :style="[{ 'border-color': nameColor }]"
    @click="handleQuoteClick"
  >
    <div
      v-if="thumb"
      class="thumb"
    >
      <v-img
        aspect-ratio="1"
        width="28"
        :src="thumb"
        class="thumb-image"
      />
    </div>
    <div>
      <div
        :style="[{ color: nameColor }]"
        class="speaker"
      >
        {{ quoteMessage.speaker_name }}
      </div>
      <div class="text">
        {{ text }}
      </div>
    </div>
  </div>
</template>
<script>
import { mapGetters } from 'vuex'
import { SUP_MESSAGE_CAT } from '@/constants'
import cutStr from '@/utils/cutStr'

export default {
  name: "QuoteMessage",
  props: {
    message: {
      type: Object,
      default: () => {}
    }
  },
  data () {
    return {
    };
  },
  computed: {
    ...mapGetters('message', ['getMessageById', 'getNameColor']),
    quoteMessageId() {
      return this.message.quote_message_id
    },
    quoteMessage() {
      if (!this.quoteMessageId) { return null }
      return this.getMessageById(this.quoteMessageId)
    },
    nameColor() {
      if (!this.quoteMessage) { return '' }
      return this.getNameColor(this.quoteMessage)
    },
    show () {
      if (!this.quoteMessage) { return false }
      const category = this.quoteMessage.category
      return SUP_MESSAGE_CAT.includes(category)
    },
    text() {
      if (!this.show) { return '' }
      const category = this.quoteMessage.category
      if (category === 'PLAIN_TEXT') {
        return cutStr(this.quoteMessage.text, 8)
      } 
      if (category === 'PLAIN_DATA') {
        const name = this.quoteMessage.attachment.name
        return cutStr(name, 8)
      }
      const textMap = {
        'PLAIN_IMAGE': 'Photo',
        'PLAIN_VIDEO': 'Video',
        'PLAIN_AUDIO': 'Voice Message'
      }
      return textMap[category] || ''
    },
    thumb() {
      if (!this.show) { return '' }
      const thumbnail = this.quoteMessage.attachment.thumbnail
      if (thumbnail) {
        return 'data:image/jpeg;base64,' + thumbnail
      }
      return ''
    }
  },
  methods: {
    handleQuoteClick() {
      this.$root.$emit('focusMessage', this.quoteMessage)
    }
  }
}
</script>
<style lang="scss" scoped>
.quote-message {
  border-left: 2px solid;
  padding-left: 4px;
  font-size: 12px;
  display: flex;

  .thumb {
    margin-right: 4px;
    display: flex;
    align-items: center;

    .thumb-image {
      border-radius: 2px;
    }
  }

  .speaker {
    line-height: 12px;
  }
}
</style>