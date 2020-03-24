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
import VueAsyncComputed from 'vue-async-computed';
import sanitizeHTML from 'sanitize-html';

import 'vue-tour/dist/vue-tour.css';
import 'vuetify/dist/vuetify.min.css';
import 'dayspan-vuetify/dist/lib/dayspan-vuetify.min.css';

import * as config from '@/config';
import * as constants from '@/constants';
import App from '@/app.vue';
import router from '@/router';
import store from '@/store';
import i18n from '@/i18n';
import filters from '@/filters';

import featuresService from '@/services/features';

import ModalsPlugin from '@/plugins/modals';
import PopupsPlugin from '@/plugins/popups';
import SetSeveralPlugin from '@/plugins/set-several';
import UpdateFieldPlugin from '@/plugins/update-field';
import ToursPlugin from '@/plugins/tours';
import VuetifyReplacerPlugin from '@/plugins/vuetify-replacer';

import DsCalendarEvent from '@/components/other/stats/calendar/day-span/partial/calendar-event.vue';
import DsCalendarEventTime from '@/components/other/stats/calendar/day-span/partial/calendar-event-time.vue';

import AlarmChips from '@/components/other/alarm/alarm-chips.vue';

import WebhookIcon from '@/components/icons/webhook.vue';

import * as modalsComponents from '@/components/modals';
/* eslint-enable import/first */

Vue.use(VueAsyncComputed);
Vue.use(VueResizeText);
Vue.use(filters);
Vue.use(Vuetify, {
  iconfont: 'md',
  theme: {
    primary: config.COLORS.primary,
    secondary: config.COLORS.secondary,
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

Vue.component('alarm-chips', AlarmChips);

Vue.use(VueMq, {
  breakpoints: config.MEDIA_QUERIES_BREAKPOINTS,
});

VueClipboard.config.autoSetContainer = true;
Vue.use(VueClipboard);

Vue.use(VeeValidate, {
  i18n,
  inject: false,
  silentTranslationWarn: false,
  dictionary: {
    en: enValidationMessages,
    fr: frValidationMessages,
  },
});

const { MODALS } = constants;

Vue.use(ModalsPlugin, {
  store,

  components: {
    ...modalsComponents,
    ...featuresService.get('components.modals.components'),
  },

  dialogPropsMap: {
    [MODALS.createPbehavior]: { maxWidth: 920, lazy: true },
    [MODALS.pbehaviorList]: { maxWidth: 1280, lazy: true },
    [MODALS.createWidget]: { maxWidth: 500, lazy: true },
    [MODALS.alarmsList]: { fullscreen: true, lazy: true },
    [MODALS.createFilter]: { maxWidth: 920, lazy: true },
    [MODALS.textEditor]: { maxWidth: 700, lazy: true, persistent: true },
    [MODALS.addInfoPopup]: { maxWidth: 700, lazy: true, persistent: true },
    [MODALS.watcher]: { maxWidth: 920, lazy: true },

    ...featuresService.get('components.modals.dialogPropsMap'),
  },
});

Vue.use(PopupsPlugin, { store });
Vue.use(SetSeveralPlugin);
Vue.use(UpdateFieldPlugin);
Vue.use(ToursPlugin);
Vue.use(VuetifyReplacerPlugin);

Vue.config.productionTip = false;

/**
 * TODO: Update it to Vue.config.errorHandler after updating to 2.6.0+ Vue version
 */
window.addEventListener('unhandledrejection', (err) => {
  store.dispatch('popups/error', { text: err.description || i18n.t('errors.default') });
});

if (process.env.NODE_ENV === 'development') {
  Vue.config.devtools = true;
  Vue.config.performance = true;
}

Vue.prototype.$constants = deepFreeze(constants);
Vue.prototype.$config = deepFreeze(config);
Vue.prototype.$sanitize = sanitizeHTML;

new Vue({
  router,
  store,
  i18n,
  render: h => h(App),
}).$mount('#app');

