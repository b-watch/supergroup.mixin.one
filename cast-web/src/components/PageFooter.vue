<template>
  <v-footer
    v-if="show"
    fixed
    app
    dark
    color="green lighten-1"
  >
    <v-icon class="mr-2">
      mdi-download
    </v-icon>
    <div class="body-2">
      下载 学习 APP
    </div>
    <v-spacer />
    <v-btn
      small
      text
      @click="handleHide"
    >
      不再显示
    </v-btn>
  </v-footer>
</template>
<script>
import { mapState, mapMutations } from 'vuex'
export default {
  name: "PageFooter",
  computed: {
    ...mapState('app', {
      systemBar: state => state.systemBar,
      pageFooter: state => state.pageFooter
    }),
    show() {
      return !this.systemBar && this.pageFooter
    }
  },
  methods: {
    ...mapMutations('app', ['setPageFooter']),
    handleHide() {
      this.setPageFooter(false)
    }
  }
}
</script>
<style lang="scss" scoped>
</style>