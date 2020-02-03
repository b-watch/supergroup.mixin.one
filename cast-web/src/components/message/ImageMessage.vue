<template>
  <v-img
    :id="'message' + message.id"
    :lazy-src="'data:image/jpeg;base64,' + message.attachment.thumbnail"
    :src="message.attachment.view_url"
    contain
    min-width="200"
    max-height="100%"
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
    }
  },
  data () {
    return {
    };
  },
  methods: {
    handleViewImage() {
    }
  }
}
</script>
<style lang="scss" scoped>
.image {
  border-radius: 4px;
}
</style>