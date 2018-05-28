import { MOBILE_BREAKPOINT, TABLET_BREAKPOINT, LAPTOP_BREAKPOINT } from '@/config';

import Vue from 'vue';
import Vuetify from 'vuetify';
import VueMoment from 'vue-moment';
import VeeValidate from 'vee-validate';
import enValidationMessages from 'vee-validate/dist/locale/en';
import frValidationMessages from 'vee-validate/dist/locale/fr';
import VueMq from 'vue-mq';

import 'vuetify/dist/vuetify.min.css';

import App from './app.vue';
import router from './router';
import store from './store';
import i18n from './i18n';


Vue.use(Vuetify);
Vue.use(VueMoment);

Vue.use(VueMq, {
  breakpoints: {
    mobile: MOBILE_BREAKPOINT,
    tablet: TABLET_BREAKPOINT,
    laptop: LAPTOP_BREAKPOINT,
  },
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

