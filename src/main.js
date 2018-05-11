import Vue from 'vue';
import Vuetify from 'vuetify';
import VeeValidate from 'vee-validate';

import 'vuetify/dist/vuetify.min.css';

import App from './app.vue';
import router from './router';
import store from './store';
import i18n from './i18n';

Vue.use(Vuetify);
Vue.use(VeeValidate);

Vue.config.productionTip = false;

new Vue({
  router,
  store,
  i18n,
  render: h => h(App),
}).$mount('#app');
