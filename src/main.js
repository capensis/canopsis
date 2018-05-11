import Vue from 'vue';
import Vuetify from 'vuetify';
import VeeValidate from 'vee-validate';

import 'vuetify/dist/vuetify.min.css';

import App from './app.vue';
import router from './router';
import store from './store';
import i18n from './i18n';

Vue.use(Vuetify);
<<<<<<< HEAD
Vue.use(VeeValidate, {
  i18n,
  inject: false,
  dictionary: {
    en: enValidationMessages,
    fr: frValidationMessages,
  },
});
=======
Vue.use(VeeValidate);
>>>>>>> parent of c7eef04... Add base modal structure and finish some modals

Vue.config.productionTip = false;

new Vue({
  router,
  store,
  i18n,
  render: h => h(App),
}).$mount('#app');
