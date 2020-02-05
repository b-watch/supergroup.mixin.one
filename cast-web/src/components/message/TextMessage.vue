<template>
  <div
    ref="text"
    class="text"
  >
    <div
      ref="text1"
      v-html="text"
    />
  </div>
</template>
<script>
import MessageTime from './MessageTime.vue'
import urlify from '@/utils/urlify'

export default {
  name: "TextMessage",
  props: {
    message: {
      type: Object,
      default: () => {}
    }
  },
  data () {
    return {
      sameline: false
    };
  },
  computed: {
    text() {
      return this.urlify(this.message.text) + '<span style="visibility: hidden">AA:AA</span>'
    }
  },
  methods: {
    urlify
  }
}
</script>
<style lang="scss" scoped>
.text {
  margin: 4px;
}
</style>