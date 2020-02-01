<template>
  <v-img
    :lazy-src="'data:image/jpeg;base64,' + message.attachment.thumbnail"
    :src="message.attachment.view_url"
    aspect-ratio="1"
    max-width="500"
    max-height="200"
    min-width="200"
    class="image"
  />
</template>
<script>
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
  }
}
</script>
<style lang="scss" scoped>
.image {
  border-radius: 4px;
}
</style>