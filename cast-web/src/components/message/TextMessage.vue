<template>
  <div
    ref="text"
    class="text"
  >
    <div
      ref="text1"
      v-html="urlify(message.text)"
    />
    <div
      ref="text2"
      style="visibility: hidden"
      v-html="urlify(message.text + 'AA:AA')"
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
  mounted () {
    const text = this.$refs.text
    const text1 = this.$refs.text1
    const text2 = this.$refs.text2
    const { width: w1, height: h1 } = text1.getBoundingClientRect()
    const { width: w2, height: h2 } = text2.getBoundingClientRect()
    text.style.width = Math.max(w1, w2) + 'px'
    text.style.height = Math.max(h1, h2) + 'px'
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