import Vue from 'vue'
import VuePlyr from 'vue-plyr'

Vue.use(VuePlyr, {
  plyr: {
    fullscreen: { enabled: false },
    loadSprite: false,
    resetOnEnd: true
  },
  emit: ['ended']
})