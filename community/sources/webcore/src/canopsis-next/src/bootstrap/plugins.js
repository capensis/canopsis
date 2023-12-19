import VueMq from 'vue-mq';
import VueFullScreen from 'vue-fullscreen';
import VueClipboard from 'vue-clipboard2';
import PortalVue from 'portal-vue';
import frDaySpanVuetifyMessages from 'dayspan-vuetify/src/locales/fr';

import 'dayspan-vuetify/dist/lib/dayspan-vuetify.min.css';

import { MODALS } from '@/constants';
import { DEFAULT_THEME_COLORS, MEDIA_QUERIES_BREAKPOINTS } from '@/config';

import Vuetify from '@/plugins/vuetify';
import ValidatorPlugin from '@/plugins/validator';
import ModalsPlugin from '@/plugins/modals';
import PopupsPlugin from '@/plugins/popups';
import SidebarPlugin from '@/plugins/sidebar';
import SetSeveralPlugin from '@/plugins/set-several';
import UpdateFieldPlugin from '@/plugins/update-field';
import SocketPlugin from '@/plugins/socket';
import DaySpanVuetifyPlugin from '@/plugins/dayspan-vuetify';

import featuresService from '@/services/features';

import store from '@/store';
import i18n from '@/i18n';
import Filters from '@/filters';

import { setSeveralFields } from '@/helpers/immutable';
import { themePropertiesToCSSVariables } from '@/helpers/entities/theme/entity';

import * as modalsComponents from '@/components/modals';
import * as sidebarsComponents from '@/components/sidebars';

/**
 * @param {import('vue').VueConstructor | import('vue').Vue} Vue
 * @returns {*}
 */
export const bootstrapApplicationPlugins = (Vue) => {
  Vue.use(PortalVue);
  Vue.use(Filters);
  Vue.use(Vuetify, {
    theme: themePropertiesToCSSVariables(DEFAULT_THEME_COLORS),
  });

  Vue.use(VueFullScreen);
  Vue.use(DaySpanVuetifyPlugin, {
    data: {
      locales: {
        fr: setSeveralFields(frDaySpanVuetifyMessages, {
          'defaults.dsScheduleFrequencyDayOfWeek.weekdays': ['Lundi', 'Mardi', 'Mercredi', 'Jeudi', 'Vendredi', 'Samedi', 'Dimanche'],
          'defaults.dsDayPicker.weekdays': ['Lundi', 'Mardi', 'Mercredi', 'Jeudi', 'Vendredi', 'Samedi', 'Dimanche'],
          'defaults.dsWeeksView.weekdays': ['lun.', 'mar.', 'mer.', 'jeu.', 'ven.', 'sam.', 'dim.'],
        }),
      },
      defaults: {
        dsCalendarEventTime: {
          placeholderStyle: false,
          disabled: false,
          popoverProps: {
            openOnHover: true,
            transition: 'fade-transition',
            attach: '.modals-wrapper',
            offsetOverflow: true,
          },
        },
        dsCalendarEvent: {
          popoverProps: {
            openOnHover: true,
            transition: 'fade-transition',
            attach: '.modals-wrapper',
            offsetOverflow: true,
          },
        },
        dsCalendarEventPlaceholder: {
          popoverProps: {
            openOnHover: true,
            transition: 'fade-transition',
            attach: '.modals-wrapper',
            offsetOverflow: true,
          },
        },
        dsCalendarEventTimePlaceholder: {
          popoverProps: {
            openOnHover: true,
            transition: 'fade-transition',
            attach: '.modals-wrapper',
            offsetOverflow: true,
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

  Vue.use(VueMq, {
    breakpoints: MEDIA_QUERIES_BREAKPOINTS,
  });

  VueClipboard.config.autoSetContainer = true;
  Vue.use(VueClipboard);

  Vue.use(ValidatorPlugin, { i18n });

  Vue.use(ModalsPlugin, {
    store,

    components: {
      ...modalsComponents,
      ...featuresService.get('components.modals.components'),
    },

    dialogPropsMap: {
      [MODALS.pbehaviorList]: { maxWidth: 1280, lazy: true },
      [MODALS.createWidget]: { maxWidth: 500, lazy: true },
      [MODALS.createWidgetTemplate]: { maxWidth: 920, lazy: true },
      [MODALS.alarmsList]: { maxWidth: '95%', lazy: true },
      [MODALS.entitiesList]: { maxWidth: '95%', lazy: true },
      [MODALS.createFilter]: { maxWidth: 1100, lazy: true },
      [MODALS.textEditor]: { maxWidth: 700, lazy: true, persistent: true },
      [MODALS.addInfoPopup]: { maxWidth: 700, lazy: true, persistent: true },
      [MODALS.serviceEntities]: { maxWidth: 920, lazy: true },
      [MODALS.serviceDependencies]: { maxWidth: 1100, lazy: true },
      [MODALS.importExportViews]: { maxWidth: 920, persistent: true },
      [MODALS.createPlaylist]: { maxWidth: 920, lazy: true },
      [MODALS.pbehaviorPlanning]: { maxWidth: '95%', lazy: true, persistent: true },
      [MODALS.pbehaviorsCalendar]: { maxWidth: '95%', lazy: true, persistent: true },
      [MODALS.pbehaviorRecurrenceRule]: { maxWidth: '95%', lazy: true, persistent: true },
      [MODALS.pbehaviorRecurrentChangesConfirmation]: { maxWidth: 400, persistent: true },
      [MODALS.createRemediationInstruction]: { maxWidth: 960 },
      [MODALS.remediationInstructionApproval]: { maxWidth: 960 },
      [MODALS.executeRemediationInstruction]: { maxWidth: 960, persistent: true },
      [MODALS.imageViewer]: { maxWidth: '90%', contentClass: 'v-dialog__image-viewer' },
      [MODALS.imagesViewer]: { maxWidth: '100%', contentClass: 'v-dialog__images-viewer' },
      [MODALS.rate]: { maxWidth: 500 },
      [MODALS.createMetaAlarmRule]: { maxWidth: 1280, lazy: true },
      [MODALS.createEventFilter]: { maxWidth: 1280 },
      [MODALS.testSuite]: { maxWidth: 920 },
      [MODALS.createPattern]: { maxWidth: 1280 },
      [MODALS.remediationPatterns]: { maxWidth: 1280 },
      [MODALS.pbehaviorPatterns]: { maxWidth: 1280 },
      [MODALS.createIdleRule]: { maxWidth: 1280 },
      [MODALS.createScenario]: { maxWidth: 1280 },
      [MODALS.createKpiFilter]: { maxWidth: 1280 },
      [MODALS.createDynamicInfo]: { maxWidth: 1280 },
      [MODALS.createAlarmStatusRule]: { maxWidth: 1280 },
      [MODALS.createService]: { maxWidth: 1280 },
      [MODALS.createMap]: { maxWidth: 500, lazy: true },
      [MODALS.createMermaidMap]: { maxWidth: 1600 },
      [MODALS.createTreeOfDependenciesMap]: { maxWidth: 1334 },
      [MODALS.createGeoMap]: { maxWidth: 1280 },
      [MODALS.createFlowchartMap]: { maxWidth: 1600 },
      [MODALS.entityDependenciesList]: { maxWidth: 1600 },
      [MODALS.createDeclareTicketRule]: { maxWidth: 1280 },
      [MODALS.createDeclareTicketEvent]: { maxWidth: 1280 },
      [MODALS.executeDeclareTickets]: { maxWidth: 920 },
      [MODALS.createLinkRule]: { maxWidth: 920 },
      [MODALS.createAlarmChart]: { maxWidth: 500 },
      [MODALS.createTag]: { maxWidth: 920 },

      ...featuresService.get('components.modals.dialogPropsMap'),
    },
  });

  Vue.use(PopupsPlugin, { store });
  Vue.use(SidebarPlugin, {
    store,

    components: {
      ...sidebarsComponents,
      ...featuresService.get('components.sidebars.components'),
    },
  });
  Vue.use(SetSeveralPlugin);
  Vue.use(UpdateFieldPlugin);
  Vue.use(SocketPlugin);
};
