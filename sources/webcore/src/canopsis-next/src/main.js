import Vue from 'vue';
import Vuetify from 'vuetify';
import VeeValidate from 'vee-validate';
import enValidationMessages from 'vee-validate/dist/locale/en';
import frValidationMessages from 'vee-validate/dist/locale/fr';
import VueMq from 'vue-mq';

import 'vuetify/dist/vuetify.min.css';

import { MEDIA_QUERIES_BREAKPOINTS } from '@/config';
import App from '@/app.vue';
import router from '@/router';
import store from '@/store';
import i18n from '@/i18n';
import filters from '@/filters';

Vue.use(filters);
Vue.use(Vuetify);

Vue.use(VueMq, {
  breakpoints: MEDIA_QUERIES_BREAKPOINTS,
});

Vue.use(VeeValidate, {
  i18n,
  inject: false,
  dictionary: {
    en: enValidationMessages,
    fr: frValidationMessages,
  },
});

Vue.config.productionTip = false;

new Vue({
  router,
  store,
  i18n,
  render: h => h(App),
}).$mount('#app');

