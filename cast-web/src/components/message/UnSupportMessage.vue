<template>
  <div
    class="px-1 caption"
    @click="handleClick"
  >
    <a>[ 请下载 Mixin Messenger 查看这条消息 ]<span class="hidden">AA:AA</span></a>
  </div>
</template>
<script>
export default {
  name: "UnSupportMessage",
  props: {
    message: {
      type: Object,
      default: () => {}
    }
  },
  methods: {
    handleClick() {
      this.$downloadApp()
    }
  }
}
</script>
<style lang="scss" scoped>
.hidden {
  visibility: hidden;
}
</style>