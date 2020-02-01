import Vue from 'vue'
import { COLORS } from '@/constants'
import timeUtils from '@/utils/time'

Vue.prototype.$colors = COLORS

Vue.prototype.$timeUtil = timeUtils