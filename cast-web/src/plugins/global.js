import Vue from 'vue'
import { COLORS } from '@/constants'
import timeUtils from '@/utils/time'
import touch from '@/directives/touch'
import mutate from 'vuetify/lib/directives/mutate'
import intersect from 'vuetify/lib/directives/intersect'

import "viewerjs/dist/viewer.css";

Vue.use(touch)
Vue.directive('mutate', mutate)
Vue.directive('intersect', intersect)

Vue.prototype.$colors = COLORS

Vue.prototype.$timeUtil = timeUtils

Vue.prototype.$downloadApp = function() {
  window.alert('尽情期待')
}