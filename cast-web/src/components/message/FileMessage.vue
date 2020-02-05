<template>
  <div
    class="file-message"
    @click="handleDownload"
  >
    <v-icon>mdi-download</v-icon>
    <div class="ml-1">
      <div>
        {{ name }}
        <span class="hidden">AA:AA</span>
      </div>
      <div class="caption text--secondary">
        {{ size }}
      </div>
    </div>
  </div>
</template>
<script>
import byteSize from '@/utils/byteSize'
import download from '@/utils/download'
import { isWeixin } from '@/utils/env'
import cutStr from '@/utils/cutStr'

export default {
  name: "FileMessage",
  props: {
    message: {
      type: Object,
      default: () => {}
    }
  },
  data () {
    return {
      textLength: 16
    };
  },
  computed: {
    name() {
      const name = this.message.attachment.name
      return cutStr(name, 8)
    },
    size() {
      return byteSize(this.message.attachment.size)
    }
  },
  methods: {
    handleDownload() {
      if (isWeixin()) {
        this.$root.$emit('openInBrowser')
        return
      }
      const name = this.message.attachment.name
      const link = this.message.attachment.view_url
      download(link, name)
    }
  }
}
</script>
<style lang="scss" scoped>
.file-message  {
  display: flex;
  align-items: center;
  padding: 0 4px
}

.hidden {
  visibility: hidden;
}
</style>