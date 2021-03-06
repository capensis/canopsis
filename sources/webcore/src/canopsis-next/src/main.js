import './bootstrap';

/* eslint-disable import/first */
import Vue from 'vue';
import moment from 'moment-timezone';
import deepFreeze from 'deep-freeze';
import Vuetify from 'vuetify';
import VeeValidate, { Validator } from 'vee-validate';
import enValidationMessages from 'vee-validate/dist/locale/en';
import frValidationMessages from 'vee-validate/dist/locale/fr';
import VueMq from 'vue-mq';
import VueFullScreen from 'vue-fullscreen';
import VueClipboard from 'vue-clipboard2';
import VueResizeText from 'vue-resize-text';
import VueAsyncComputed from 'vue-async-computed';
import PortalVue from 'portal-vue';
import sanitizeHTML from 'sanitize-html';
import frDaySpanVuetifyMessages from 'dayspan-vuetify/src/locales/fr';

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
import DaySpanVuetifyPlugin from '@/plugins/dayspan-vuetify';
import GridPlugin from '@/plugins/grid';

import { isValidUrl } from '@/helpers/validators/is-valid-url';
import { isValidJson } from '@/helpers/validators/is-valid-json';

import AlarmsListTable from '@/components/widgets/alarm/partials/alarms-list-table.vue';
import AlarmChips from '@/components/widgets/alarm/alarm-chips.vue';
import CAdvancedDataTable from '@/components/common/table/c-advanced-data-table.vue';
import CThePageHeader from '@/components/common/page/c-the-page-header.vue';
import CExpandBtn from '@/components/common/buttons/c-expand-btn.vue';
import CActionBtn from '@/components/common/buttons/c-action-btn.vue';
import CFabExpandBtn from '@/components/common/buttons/c-fab-expand-btn.vue';
import CFabBtn from '@/components/common/buttons/c-fab-btn.vue';
import CRefreshBtn from '@/components/common/buttons/c-refresh-btn.vue';
import CEmptyDataTableColumns from '@/components/common/table/c-empty-data-table-columns.vue';
import CEnabled from '@/components/icons/c-enabled.vue';
import CEllipsis from '@/components/common/table/c-ellipsis.vue';
import CPagination from '@/components/common/pagination/c-pagination.vue';
import CTablePagination from '@/components/common/pagination/c-table-pagination.vue';
import CAlertOverlay from '@/components/common/overlay/c-alert-overlay.vue';
import CProgressOverlay from '@/components/common/overlay/c-progress-overlay.vue';

import WebhookIcon from '@/components/icons/webhook.vue';
import BullhornIcon from '@/components/icons/bullhorn.vue';
import AltRouteIcon from '@/components/icons/alt_route.vue';
import SettingsSyncIcon from '@/components/icons/settings_sync.vue';

import * as modalsComponents from '@/components/modals';

/* eslint-enable import/first */

Vue.use(VueAsyncComputed);
Vue.use(VueResizeText);
Vue.use(PortalVue);
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
    bullhorn: {
      component: BullhornIcon,
    },
    alt_route: {
      component: AltRouteIcon,
    },
    settings_sync: {
      component: SettingsSyncIcon,
    },
  },
});

Vue.use(GridPlugin);
Vue.use(VueFullScreen);
Vue.use(DaySpanVuetifyPlugin, {
  data: {
    locales: {
      fr: frDaySpanVuetifyMessages,
    },
    defaults: {
      dsWeeksView: {
        // dayspan-vuetify doesn't not supported first day in weekend, because return weekdays without locale sort.
        weekdays: moment.weekdaysShort(),
      },
      dsCalendarEventTime: {
        placeholderStyle: false,
        disabled: false,
        popoverProps: {
          nudgeWidth: 200,
          closeOnContentClick: false,
          transition: 'fade-transition',
          offsetOverflow: true,
          offsetX: true,
          maxWidth: 500,
          openOnHover: true,
        },
      },
      dsCalendarEvent: {
        popoverProps: {
          offsetY: true,
          openOnHover: true,
          transition: 'fade-transition',
        },
      },
      dsCalendarEventPlaceholder: {
        popoverProps: {
          offsetY: true,
          openOnHover: true,
          transition: 'fade-transition',
        },
      },
      dsCalendarEventTimePlaceholder: {
        popoverProps: {
          openOnHover: true,
          transition: 'fade-transition',
        },
      },
    },
  },
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
});

Vue.component('alarm-chips', AlarmChips);
Vue.component('alarms-list-table', AlarmsListTable);

/* Global custom canopsis components */
Vue.component('c-the-page-header', CThePageHeader);
Vue.component('c-advanced-data-table', CAdvancedDataTable);
Vue.component('c-expand-btn', CExpandBtn);
Vue.component('c-action-btn', CActionBtn);
Vue.component('c-fab-expand-btn', CFabExpandBtn);
Vue.component('c-fab-btn', CFabBtn);
Vue.component('c-refresh-btn', CRefreshBtn);
Vue.component('c-empty-data-table-columns', CEmptyDataTableColumns);
Vue.component('c-enabled', CEnabled);
Vue.component('c-ellipsis', CEllipsis);
Vue.component('c-pagination', CPagination);
Vue.component('c-table-pagination', CTablePagination);
Vue.component('c-alert-overlay', CAlertOverlay);
Vue.component('c-progress-overlay', CProgressOverlay);

Vue.use(VueMq, {
  breakpoints: config.MEDIA_QUERIES_BREAKPOINTS,
});

VueClipboard.config.autoSetContainer = true;
Vue.use(VueClipboard);

Validator.extend('json', {
  getMessage: () => i18n.$t('errors.JSONNotValid'),
  validate: isValidJson,
});

Validator.extend('url', {
  validate: isValidUrl,
});

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
    [MODALS.pbehaviorList]: { maxWidth: 1280, lazy: true },
    [MODALS.createWidget]: { maxWidth: 500, lazy: true },
    [MODALS.alarmsList]: { fullscreen: true, lazy: true },
    [MODALS.createFilter]: { maxWidth: 920, lazy: true },
    [MODALS.textEditor]: { maxWidth: 700, lazy: true, persistent: true },
    [MODALS.addInfoPopup]: { maxWidth: 700, lazy: true, persistent: true },
    [MODALS.watcher]: { maxWidth: 920, lazy: true },
    [MODALS.importExportViews]: { maxWidth: 920, persistent: true },
    [MODALS.createPlaylist]: { maxWidth: 920, lazy: true },
    [MODALS.pbehaviorPlanning]: { fullscreen: true, lazy: true, persistent: true },
    [MODALS.pbehaviorRecurrentChangesConfirmation]: { maxWidth: 400, persistent: true },
    [MODALS.createRemediationInstruction]: { maxWidth: 960 },
    [MODALS.executeRemediationInstruction]: { maxWidth: 960, persistent: true },
    [MODALS.imageViewer]: { maxWidth: '90%', contentClass: 'v-dialog__image-viewer' },
    [MODALS.rate]: { maxWidth: 400 },
    [MODALS.createMetaAlarmRule]: { maxWidth: 920, lazy: true },

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

