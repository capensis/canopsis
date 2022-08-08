import './bootstrap';

/* eslint-disable import/first */
import Vue from 'vue';
import deepFreeze from 'deep-freeze';
import Vuetify from 'vuetify';
import VueMq from 'vue-mq';
import VueFullScreen from 'vue-fullscreen';
import VueClipboard from 'vue-clipboard2';
import VueResizeText from 'vue-resize-text';
import VueAsyncComputed from 'vue-async-computed';
import PortalVue from 'portal-vue';
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
import Filters from '@/filters';

import featuresService from '@/services/features';

import ModalsPlugin from '@/plugins/modals';
import PopupsPlugin from '@/plugins/popups';
import SidebarPlugin from '@/plugins/sidebar';
import ValidatorPlugin from '@/plugins/validator';
import SetSeveralPlugin from '@/plugins/set-several';
import UpdateFieldPlugin from '@/plugins/update-field';
import ToursPlugin from '@/plugins/tours';
import VuetifyReplacerPlugin from '@/plugins/vuetify-replacer';
import DaySpanVuetifyPlugin from '@/plugins/dayspan-vuetify';
import GridPlugin from '@/plugins/grid';
import SocketPlugin from '@/plugins/socket';

import { setSeveralFields } from '@/helpers/immutable';

import CPageHeader from '@/components/common/page/c-page-header.vue';
import CEnabled from '@/components/icons/c-enabled.vue';
import CEllipsis from '@/components/common/table/c-ellipsis.vue';
import CPagination from '@/components/common/pagination/c-pagination.vue';
import CDraggableStepNumber from '@/components/common/drag-drop/c-draggable-step-number.vue';
import CInformationBlock from '@/components/common/block/c-information-block.vue';
import CInformationBlockRow from '@/components/common/block/c-information-block-row.vue';
import CResponsiveList from '@/components/common/responsive-list/c-responsive-list.vue';

/**
 * Overlays
 */
import CAlertOverlay from '@/components/common/overlay/c-alert-overlay.vue';
import CProgressOverlay from '@/components/common/overlay/c-progress-overlay.vue';

/**
 * Chips
 */
import CAlarmChip from '@/components/common/chips/c-alarm-chip.vue';
import CStateCountChangesChips from '@/components/common/chips/c-state-count-changes-chips.vue';
import CTestSuiteChip from '@/components/common/chips/c-test-suite-chip.vue';
import CInstructionJobChip from '@/components/common/chips/c-instruction-job-chip.vue';
import CEngineChip from '@/components/common/chips/c-engine-chip.vue';

/**
 * Table
 */
import AlarmsListTable from '@/components/widgets/alarm/partials/alarms-list-table.vue';
import CAdvancedDataTable from '@/components/common/table/c-advanced-data-table.vue';
import CTreeviewDataTable from '@/components/common/table/c-treeview-data-table.vue';
import CEmptyDataTableColumns from '@/components/common/table/c-empty-data-table-columns.vue';
import CTablePagination from '@/components/common/pagination/c-table-pagination.vue';

/**
 * Buttons
 */
import CExpandBtn from '@/components/common/buttons/c-expand-btn.vue';
import CActionBtn from '@/components/common/buttons/c-action-btn.vue';
import CDownloadBtn from '@/components/common/buttons/c-download-btn.vue';
import CCopyBtn from '@/components/common/buttons/c-copy-btn.vue';
import CActionFabBtn from '@/components/common/buttons/c-action-fab-btn.vue';
import CFabExpandBtn from '@/components/common/buttons/c-fab-expand-btn.vue';
import CFabBtn from '@/components/common/buttons/c-fab-btn.vue';
import CRefreshBtn from '@/components/common/buttons/c-refresh-btn.vue';

/**
 * Fields
 */
import CEnabledField from '@/components/forms/fields/c-enabled-field.vue';
import CDurationField from '@/components/forms/fields/c-duration-field.vue';
import CDisableDuringPeriodsField from '@/components/forms/fields/c-disable-during-periods-field.vue';
import CTriggersField from '@/components/forms/fields/c-triggers-field.vue';
import CActionTypeField from '@/components/forms/fields/c-action-type-field.vue';
import CPatternsField from '@/components/forms/fields/c-patterns-field.vue';
import CWorkflowField from '@/components/forms/fields/c-workflow-field.vue';
import CChangeStateField from '@/components/forms/fields/c-change-state-field.vue';
import CRequestUrlField from '@/components/forms/fields/c-request-url-field.vue';
import CTextPairsField from '@/components/forms/fields/text-pairs/c-text-pairs-field.vue';
import CTextPairField from '@/components/forms/fields/text-pairs/c-text-pair-field.vue';
import CJsonField from '@/components/forms/fields/c-json-field.vue';
import CRetryField from '@/components/forms/fields/c-retry-field.vue';
import CMixedField from '@/components/forms/fields/c-mixed-field.vue';
import CArrayMixedField from '@/components/forms/fields/c-array-mixed-field.vue';
import CColorPickerField from '@/components/forms/fields/c-color-picker-field.vue';
import CEntityTypeField from '@/components/forms/fields/c-entity-type-field.vue';
import CImpactLevelField from '@/components/forms/fields/c-impact-level-field.vue';
import CSearchField from '@/components/forms/fields/c-search-field.vue';
import CAdvancedSearchField from '@/components/forms/fields/c-advanced-search-field.vue';
import CEntityCategoryField from '@/components/forms/fields/c-entity-category-field.vue';
import CStoragesField from '@/components/forms/fields/c-storages-field.vue';
import CStorageField from '@/components/forms/fields/c-storage-field.vue';
import CFileNameMaskField from '@/components/forms/fields/c-file-name-mask-field.vue';
import CPercentsField from '@/components/forms/fields/c-percents-field.vue';
import CColumnsField from '@/components/forms/fields/c-columns-field.vue';
import CColorIndicatorField from '@/components/forms/fields/c-color-indicator-field.vue';
import CMiniBarChart from '@/components/common/chart/c-mini-bar-chart.vue';
import CImagesViewer from '@/components/common/images-viewer/c-images-viewer.vue';
import CClickableTooltip from '@/components/common/clickable-tooltip/c-clickable-tooltip.vue';
import CRolePickerField from '@/components/forms/fields/c-role-picker-field.vue';
import CUserPickerField from '@/components/forms/fields/c-user-picker-field.vue';
import CInstructionTypeField from '@/components/forms/fields/c-instruction-type-field.vue';
import CPriorityField from '@/components/forms/fields/c-priority-field.vue';
import CQuickDateIntervalField from '@/components/forms/fields/c-quick-date-interval-field.vue';
import CQuickDateIntervalTypeField from '@/components/forms/fields/c-quick-date-interval-type-field.vue';
import CEnabledDurationField from '@/components/forms/fields/c-enabled-duration-field.vue';
import CEnabledLimitField from '@/components/forms/fields/c-enabled-limit-field.vue';
import CTimezoneField from '@/components/forms/fields/c-timezone-field.vue';
import CLanguageField from '@/components/forms/fields/c-language-field.vue';
import CSamplingField from '@/components/forms/fields/c-sampling-field.vue';
import CAlarmMetricParametersField from '@/components/forms/fields/c-alarm-metric-parameters-field.vue';
import CFiltersField from '@/components/forms/fields/c-filters-field.vue';
import CStateTypeField from '@/components/forms/fields/c-state-type-field.vue';
import CRecordsPerPageField from '@/components/forms/fields/c-records-per-page-field.vue';
import COperatorField from '@/components/forms/fields/c-operator-field.vue';
import CIconField from '@/components/forms/fields/c-icon-field.vue';

/**
 * Icons
 */
import CHelpIcon from '@/components/common/icons/c-help-icon.vue';
import CNoEventsIcon from '@/components/common/icons/c-no-events-icon.vue';
import BullhornIcon from '@/components/icons/bullhorn.vue';
import AltRouteIcon from '@/components/icons/alt_route.vue';
import SettingsSyncIcon from '@/components/icons/settings_sync.vue';
import EngineeringIcon from '@/components/icons/engineering.vue';
import InsightsIcon from '@/components/icons/insights.vue';
import MiscellaneousServicesIcon from '@/components/icons/miscellaneous_services.vue';
import PublishedWithChangesIcon from '@/components/icons/published_with_changes.vue';
import DensityMediumIcon from '@/components/icons/density_medium.vue';
import DensitySmallIcon from '@/components/icons/density_small.vue';

/**
 * Groups
 */
import CDensityBtnToggle from '@/components/common/groups/c-density-btn-toggle.vue';

import * as modalsComponents from '@/components/modals';
import * as sidebarsComponents from '@/components/sidebars';

/* eslint-enable import/first */

Vue.use(VueAsyncComputed);
Vue.use(VueResizeText);
Vue.use(PortalVue);
Vue.use(Filters);
Vue.use(Vuetify, {
  iconfont: 'md',
  theme: {
    primary: config.COLORS.primary,
    secondary: config.COLORS.secondary,
  },
  icons: {
    bullhorn: {
      component: BullhornIcon,
    },
    alt_route: {
      component: AltRouteIcon,
    },
    settings_sync: {
      component: SettingsSyncIcon,
    },
    engineering: {
      component: EngineeringIcon,
    },
    insights: {
      component: InsightsIcon,
    },
    miscellaneous_services: {
      component: MiscellaneousServicesIcon,
    },
    published_with_changes: {
      component: PublishedWithChangesIcon,
    },
    density_medium: {
      component: DensityMediumIcon,
    },
    density_small: {
      component: DensitySmallIcon,
    },
  },
});

Vue.use(GridPlugin);
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

Vue.component('alarms-list-table', AlarmsListTable);

/* Global custom canopsis components */
Vue.component('c-alarm-chip', CAlarmChip);
Vue.component('c-instruction-job-chip', CInstructionJobChip);
Vue.component('c-test-suite-chip', CTestSuiteChip);
Vue.component('c-engine-chip', CEngineChip);
Vue.component('c-page-header', CPageHeader);
Vue.component('c-advanced-data-table', CAdvancedDataTable);
Vue.component('c-treeview-data-table', CTreeviewDataTable);
Vue.component('c-expand-btn', CExpandBtn);
Vue.component('c-action-btn', CActionBtn);
Vue.component('c-fab-expand-btn', CFabExpandBtn);
Vue.component('c-fab-btn', CFabBtn);
Vue.component('c-refresh-btn', CRefreshBtn);
Vue.component('c-download-btn', CDownloadBtn);
Vue.component('c-copy-btn', CCopyBtn);
Vue.component('c-density-btn-toggle', CDensityBtnToggle);
Vue.component('c-action-fab-btn', CActionFabBtn);
Vue.component('c-empty-data-table-columns', CEmptyDataTableColumns);
Vue.component('c-enabled', CEnabled);
Vue.component('c-ellipsis', CEllipsis);
Vue.component('c-pagination', CPagination);
Vue.component('c-table-pagination', CTablePagination);
Vue.component('c-alert-overlay', CAlertOverlay);
Vue.component('c-progress-overlay', CProgressOverlay);
Vue.component('c-enabled-field', CEnabledField);
Vue.component('c-duration-field', CDurationField);
Vue.component('c-disable-during-periods-field', CDisableDuringPeriodsField);
Vue.component('c-triggers-field', CTriggersField);
Vue.component('c-action-type-field', CActionTypeField);
Vue.component('c-patterns-field', CPatternsField);
Vue.component('c-workflow-field', CWorkflowField);
Vue.component('c-draggable-step-number', CDraggableStepNumber);
Vue.component('c-change-state-field', CChangeStateField);
Vue.component('c-request-url-field', CRequestUrlField);
Vue.component('c-text-pair-field', CTextPairField);
Vue.component('c-text-pairs-field', CTextPairsField);
Vue.component('c-json-field', CJsonField);
Vue.component('c-retry-field', CRetryField);
Vue.component('c-mixed-field', CMixedField);
Vue.component('c-array-mixed-field', CArrayMixedField);
Vue.component('c-color-picker-field', CColorPickerField);
Vue.component('c-entity-type-field', CEntityTypeField);
Vue.component('c-impact-level-field', CImpactLevelField);
Vue.component('c-search-field', CSearchField);
Vue.component('c-advanced-search-field', CAdvancedSearchField);
Vue.component('c-entity-category-field', CEntityCategoryField);
Vue.component('c-storages-field', CStoragesField);
Vue.component('c-storage-field', CStorageField);
Vue.component('c-file-name-mask-field', CFileNameMaskField);
Vue.component('c-percents-field', CPercentsField);
Vue.component('c-color-indicator-field', CColorIndicatorField);
Vue.component('c-columns-field', CColumnsField);
Vue.component('c-mini-bar-chart', CMiniBarChart);
Vue.component('c-images-viewer', CImagesViewer);
Vue.component('c-clickable-tooltip', CClickableTooltip);
Vue.component('c-help-icon', CHelpIcon);
Vue.component('c-no-events-icon', CNoEventsIcon);
Vue.component('c-role-picker-field', CRolePickerField);
Vue.component('c-user-picker-field', CUserPickerField);
Vue.component('c-instruction-type-field', CInstructionTypeField);
Vue.component('c-priority-field', CPriorityField);
Vue.component('c-quick-date-interval-field', CQuickDateIntervalField);
Vue.component('c-quick-date-interval-type-field', CQuickDateIntervalTypeField);
Vue.component('c-enabled-duration-field', CEnabledDurationField);
Vue.component('c-enabled-limit-field', CEnabledLimitField);
Vue.component('c-timezone-field', CTimezoneField);
Vue.component('c-language-field', CLanguageField);
Vue.component('c-filters-field', CFiltersField);
Vue.component('c-state-count-changes-chips', CStateCountChangesChips);
Vue.component('c-information-block', CInformationBlock);
Vue.component('c-information-block-row', CInformationBlockRow);
Vue.component('c-responsive-list', CResponsiveList);
Vue.component('c-sampling-field', CSamplingField);
Vue.component('c-alarm-metric-parameters-field', CAlarmMetricParametersField);
Vue.component('c-state-type-field', CStateTypeField);
Vue.component('c-records-per-page-field', CRecordsPerPageField);
Vue.component('c-operator-field', COperatorField);
Vue.component('c-icon-field', CIconField);

Vue.use(VueMq, {
  breakpoints: config.MEDIA_QUERIES_BREAKPOINTS,
});

VueClipboard.config.autoSetContainer = true;
Vue.use(VueClipboard);

Vue.use(ValidatorPlugin, { i18n });

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
    [MODALS.alarmsList]: { maxWidth: '95%', lazy: true },
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
    [MODALS.createMetaAlarmRule]: { maxWidth: 920, lazy: true },
    [MODALS.createEventFilterRuleAction]: { maxWidth: 920 },
    [MODALS.testSuite]: { maxWidth: 920 },

    ...featuresService.get('components.modals.dialogPropsMap'),
  },
});

Vue.use(PopupsPlugin, { store });
Vue.use(SidebarPlugin, {
  store,

  components: {
    ...sidebarsComponents,
  },
});
Vue.use(SetSeveralPlugin);
Vue.use(UpdateFieldPlugin);
Vue.use(ToursPlugin);
Vue.use(VuetifyReplacerPlugin);
Vue.use(SocketPlugin);

Vue.config.productionTip = false;

Vue.config.errorHandler = (err) => {
  console.error(err);

  store.dispatch('popups/error', { text: err.description || i18n.t('errors.default') });
};

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
