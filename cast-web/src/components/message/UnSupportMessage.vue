<template>
  <div class="px-1 caption">
    <a>[ 请下载学习 App 查看这条消息 ]<span class="hidden">AA:AA</span></a>
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
  data () {
    return {
    };
  },
}
</script>
<style lang="scss" scoped>
.hidden {
  visibility: hidden;
}
</style>