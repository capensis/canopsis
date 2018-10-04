import moment from 'moment';
import Vue from 'vue';
import Vuetify from 'vuetify';
import VeeValidate from 'vee-validate';
import enValidationMessages from 'vee-validate/dist/locale/en';
import frValidationMessages from 'vee-validate/dist/locale/fr';
import VueMq from 'vue-mq';
import DaySpanVuetify from 'dayspan-vuetify';

import 'vuetify/dist/vuetify.min.css';
import 'dayspan-vuetify/dist/lib/dayspan-vuetify.min.css';

import { MEDIA_QUERIES_BREAKPOINTS } from '@/config';
import App from '@/app.vue';
import router from '@/router';
import store from '@/store';
import i18n from '@/i18n';
import filters from '@/filters';

import DsDaysView from '@/components/other/stats/day-span/partial/days-view.vue';

Vue.use(filters);
Vue.use(Vuetify);
Vue.use(DaySpanVuetify, {
  methods: {
    getDefaultEventColor: () => '#1976d2',
    getPrefix: (calendarEvent, sameDay) =>
      (sameDay.length === 1 ? sameDay[0].start.format('HH[h]') : `(${sameDay.length})`),
  },
  data: {
    defaults: {
      dsWeeksView: {
        weekdays: moment.weekdaysShort(true),
      },
    },
  },
});
Vue.component('dsDaysView', DsDaysView);

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

