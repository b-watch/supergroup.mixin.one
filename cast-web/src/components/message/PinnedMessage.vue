<template>
  <div
    v-if="show"
    v-ripple
    class="pinned-message"
    @click="handleClick"
  >
    <div
      v-if="thumb"
      class="thumb"
    >
      <v-img
        aspect-ratio="1"
        width="38"
        :src="thumb"
        class="thumb-image"
      />
    </div>
    <div>
      <div class="subtitle-2">
        置顶消息
      </div>
      <div class="text">
        {{ text }}
      </div>
    </div>
  </div>
</template>
<script>
import { mapState } from 'vuex'
import { SUP_MESSAGE_CAT } from '@/constants'
import cutStr from '@/utils/cutStr'

export default {
  name: "PinnedMessage",
  data () {
    return {
    };
  },
  computed: {
    ...mapState('group', {
      information: state => state.information
    }),
    pinnedMessage() {
      return (this.information && this.information.pinned_message) || {}
    },
    show () {
      if (!this.pinnedMessage) { return false }
      const category = this.pinnedMessage.category
      return SUP_MESSAGE_CAT.includes(category)
    },
    text() {
      if (!this.show) { return '' }
      const category = this.pinnedMessage.category
      if (category === 'PLAIN_TEXT') {
        return this.pinnedMessage.text
      } 
      if (category === 'PLAIN_DATA') {
        const name = this.pinnedMessage.attachment.name
        return cutStr(name, 8)
      }
      const textMap = {
        'PLAIN_IMAGE': 'Photo',
        'PLAIN_VIDEO': 'Video',
        'PLAIN_AUDIO': 'Voice Message',
        'PLAIN_LIVE': 'Live'
      }
      return textMap[category] || ''
    },
    thumb() {
      if (!this.show) { return '' }
      const thumbnail = this.pinnedMessage.attachment.thumbnail
      const thumbUrl = this.pinnedMessage.attachment.thumb_url
      if (thumbnail) {
        return 'data:image/jpeg;base64,' + thumbnail
      }
      if (thumbUrl) {
        return thumbUrl
      }
      return ''
    }
  },
  methods: {
    handleClick() {
      this.$root.$emit('viewMessage', this.pinnedMessage)
    }
  }
}
</script>
<style lang="scss" scoped>
.pinned-message {
  border-left: 2px solid;
  border-radius: 2px;
  padding-left: 4px;
  margin: 4px 0;
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
}
</style>