<template>
  <div
    class="time"
    :class="[{ 'time-text': textStyle }]"
  >
    {{ $timeUtil.format(message.created_at) }}
  </div>
</template>
<script>
import { SUP_MESSAGE_CAT } from '@/constants'

export default {
  name: "MessageTime",
  props: {
    message: {
      type: Object,
      default: () => {}
    }
  },
  computed: {
    textStyle() {
      return ['PLAIN_TEXT', 'PLAIN_AUDIO', 'PLAIN_DATA'].includes(this.message.category) || this.unsupportMessage
    },
    unsupportMessage() {
      return !SUP_MESSAGE_CAT.includes(this.message.category)
    }
  }
}
</script>
<style lang="scss" scoped>
.time {
  font-size: 12px;
  line-height: 14px;
  position: absolute;
  color: #fff;
  background: rgba(0, 0, 0, 0.6);
  bottom: 8px;
  right: 8px;
  padding: 2px 4px;
  border-radius: 10px;

  &.time-text {
    background: none;
    color: rgba(0, 0, 0, 0.4);
    bottom: 4px;
    right: 4px;
  }
}
</style>