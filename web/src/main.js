import Vue from 'vue'
import App from './App'
import router from './route'
import i18n from './i18n'
import global from './global'

import '@/plugins/vant'
import '@/plugins/vue-qr'
import '@/plugins/infinite-loading'
import '@/plugins/vue-clipboard2'

Vue.config.productionTip = false
Vue.prototype.GLOBAL = global
new Vue({
  components: {App},
  router,
  i18n,
  template: '<App/>'
}).$mount('#app')
