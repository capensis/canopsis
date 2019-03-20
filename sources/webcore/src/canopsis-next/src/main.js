import './bootstrap';

/* eslint-disable import/first */
import Vue from 'vue';
import moment from 'moment';
import deepFreeze from 'deep-freeze';
import Vuetify from 'vuetify';
import VeeValidate from 'vee-validate';
import enValidationMessages from 'vee-validate/dist/locale/en';
import frValidationMessages from 'vee-validate/dist/locale/fr';
import VueMq from 'vue-mq';
import VueFullScreen from 'vue-fullscreen';
import DaySpanVuetify from 'dayspan-vuetify';
import VueClipboard from 'vue-clipboard2';
import VueResizeText from 'vue-resize-text';

import 'vuetify/dist/vuetify.min.css';
import 'dayspan-vuetify/dist/lib/dayspan-vuetify.min.css';

import * as config from '@/config';
import * as constants from '@/constants';
import App from '@/app.vue';
import router from '@/router';
import store from '@/store';
import i18n from '@/i18n';
import filters from '@/filters';

import DsCalendarEvent from '@/components/other/stats/calendar/day-span/partial/calendar-event.vue';
import DsCalendarEventTime from '@/components/other/stats/calendar/day-span/partial/calendar-event-time.vue';

import VCheckboxFunctional from '@/components/forms/fields/v-checkbox-functional.vue';
import VExpansionPanelContent from '@/components/tables/v-expansion-panel-content.vue';

import WebhookIcon from '@/components/icons/webhook.vue';
/* eslint-enable import/first */

Vue.use(VueResizeText);
Vue.use(filters);
Vue.use(Vuetify, {
  iconfont: 'md',
  theme: {
    primary: '#2fab63',
    secondary: '#2b3e4f',
  },
  icons: {
    webhook: {
      component: WebhookIcon,
    },
  },
});

Vue.use(VueFullScreen);
Vue.use(DaySpanVuetify, {
  methods: {
    getPrefix: () => '',
    getStyleColor(details, calendarEvent, past, cancelled) {
      let { color } = details;

      if (!past && !cancelled) {
        color = this.blend(color, this.inactiveBlendAmount, this.inactiveBlendTarget);
      }

      return color;
    },
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

Vue.component('v-checkbox-functional', VCheckboxFunctional);
Vue.component('v-expansion-panel-content', VExpansionPanelContent);

Vue.use(VueMq, {
  breakpoints: config.MEDIA_QUERIES_BREAKPOINTS,
});

Vue.use(VueClipboard);

Vue.use(VeeValidate, {
  i18n,
  inject: false,
  dictionary: {
    en: enValidationMessages,
    fr: frValidationMessages,
  },
});

Vue.config.productionTip = false;

/**
 * TODO: Update it to Vue.config.errorHandler after updating to 2.6.0+ Vue version
 */
window.addEventListener('unhandledrejection', () => {
  store.dispatch('popup/add', { type: 'error', text: i18n.t('errors.default') });
});

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

