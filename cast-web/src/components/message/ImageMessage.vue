<template>
  <v-img
    :id="'message' + message.id"
    :lazy-src="'data:image/jpeg;base64,' + message.attachment.thumbnail"
    :src="message.attachment.view_url"
    :max-height="isViewer ? '100%' : '280'"
    contain
    min-width="200"
    class="image"
    @click="handleViewImage"
  />
</template>
<script>
import Viewer from "viewerjs";

export default {
  name: "ImageMessage",
  props: {
    message: {
      type: Object,
      default: () => {}
    },
    isViewer: {
      type: Boolean,
      default: false
    }
  },
  data () {
    return {
    };
  },
  methods: {
    handleViewImage() {
      this.$root.$emit('viewMessage', this.message, { touchless: true })
    }
  }
}
</script>
<style lang="scss" scoped>
.image {
  border-radius: 4px;
}
</style>