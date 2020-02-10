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
    <div
      class="body-2"
      @click="handleDownload"
    >
      下载 Mixin Messenger
    </div>
    <v-spacer />
    <v-btn
      small
      text
      @click.stop="handleHide"
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
    ...mapState('group', {
      information: state => state.information
    }),
    isBroadcasting() {
      return this.information && this.information.broadcast === 'on'
    },
    show() {
      return !this.systemBar && this.isBroadcasting && this.pageFooter
    }
  },
  methods: {
    ...mapMutations('app', ['setPageFooter']),
    handleHide() {
      this.setPageFooter(false)
    },
    handleDownload() {
      this.$downloadApp()
    }
  }
}
</script>
<style lang="scss" scoped>
</style>