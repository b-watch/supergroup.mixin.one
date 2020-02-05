<template>
  <v-dialog
    v-model="dialog"
    fullscreen
  >
    <v-overlay opacity="0.9">
      <div
        class="overlay"
        @click="handleClick"
      >
        <span class="arrow">
          <v-icon>mdi-navigation</v-icon>
        </span>
        点击右上角在浏览器中打开
      </div>
    </v-overlay>
  </v-dialog>
</template>
<script>
export default {
  name: "OpenInBrowser",
  data () {
    return {
      dialog: false
    };
  },
  mounted() {
    this.$root.$on('openInBrowser', () => {
      this.dialog = true
    })
  },
  methods: {
    handleClick() {
      this.dialog = false
    }
  }
}
</script>
<style lang="scss" scoped>
.overlay {
  position: fixed;
  height: 100%;
  width: 100%;
  top: 0;
  left: 0;
  display: flex;
  align-items: center;
  justify-content: center;

  .arrow {
    position: absolute;
    top: 10px;
    right: 15px;
  }
}
</style>