import Vue from 'vue'
import { COLORS } from '@/constants'
import timeUtils from '@/utils/time'
import Touch from '@/utils/touch'

import "viewerjs/dist/viewer.css";

Vue.use(Touch)

Vue.prototype.$colors = COLORS

Vue.prototype.$timeUtil = timeUtils