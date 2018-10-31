import moment from 'moment';
import deepFreeze from 'deep-freeze';
import Vue from 'vue';
import Vuetify from 'vuetify';
import VeeValidate from 'vee-validate';
import enValidationMessages from 'vee-validate/dist/locale/en';
import frValidationMessages from 'vee-validate/dist/locale/fr';
import VueMq from 'vue-mq';
import VueFullScreen from 'vue-fullscreen';
import DaySpanVuetify from 'dayspan-vuetify';

import 'vuetify/dist/vuetify.min.css';
import 'dayspan-vuetify/dist/lib/dayspan-vuetify.min.css';

import * as config from '@/config';
import * as constants from '@/constants';
import App from '@/app.vue';
import router from '@/router';
import store from '@/store';
import i18n from '@/i18n';
import filters from '@/filters';

import DsCalendarEvent from '@/components/other/stats/day-span/partial/calendar-event.vue';
import DsCalendarEventTime from '@/components/other/stats/day-span/partial/calendar-event-time.vue';

import VCheckboxFunctional from '@/components/forms/v-checkbox-functional.vue';

Vue.use(filters);
Vue.use(Vuetify, {
  theme: {
    primary: '#2fab63',
    secondary: '#2b3e4f',
  },
});
Vue.use(VueFullScreen);
Vue.use(DaySpanVuetify, {
  methods: {
    getPrefix: () => '',
  },
  data: {
    defaults: {
      dsWeeksView: {
        weekdays: moment.weekdaysShort(true),
      },
      dsCalendarEventTime: {
        placeholderStyle: false,
        disabled: false,
        popoverProps: {
          nudgeWidth: 200,
          closeOnContentClick: false,
          offsetOverflow: true,
          offsetX: true,
          maxWidth: 500,
          openOnHover: true,
        },
      },
    },
  },
});

Vue.component('dsCalendarEvent', DsCalendarEvent);
Vue.component('dsCalendarEventTime', DsCalendarEventTime);
Vue.component('vCheckboxFunctional', VCheckboxFunctional);

Vue.use(VueMq, {
  breakpoints: config.MEDIA_QUERIES_BREAKPOINTS,
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

if (process.env.NODE_ENV === 'development') {
  Vue.config.devtools = true;
  Vue.config.performance = true;
}

Vue.prototype.$constants = deepFreeze(constants);
Vue.prototype.$config = deepFreeze(config);

new Vue({
  router,
  store,
  i18n,
  render: h => h(App),
}).$mount('#app');

