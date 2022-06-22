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
import CPatternOperatorInformation from '@/components/common/block/c-pattern-operator-information.vue';

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
import CPatternOperatorChip from '@/components/common/chips/c-pattern-operator-chip.vue';

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
import CFabExpandBtn from '@/components/common/buttons/c-fab-expand-btn.vue';
import CFabBtn from '@/components/common/buttons/c-fab-btn.vue';
import CRefreshBtn from '@/components/common/buttons/c-refresh-btn.vue';

/**
 * Fields
 */
import CEnabledField from '@/components/forms/fields/c-enabled-field.vue';
import CDurationField from '@/components/forms/fields/duration/c-duration-field.vue';
import CDisableDuringPeriodsField from '@/components/forms/fields/pbehavior/c-disable-during-periods-field.vue';
import CTriggersField from '@/components/forms/fields/c-triggers-field.vue';
import CActionTypeField from '@/components/forms/fields/c-action-type-field.vue';
import CPatternsField from '@/components/forms/fields/pattern/c-patterns-field.vue';
import CWorkflowField from '@/components/forms/fields/c-workflow-field.vue';
import CChangeStateField from '@/components/forms/fields/c-change-state-field.vue';
import CRequestUrlField from '@/components/forms/fields/c-request-url-field.vue';
import CTextPairsField from '@/components/forms/fields/text-pairs/c-text-pairs-field.vue';
import CTextPairField from '@/components/forms/fields/text-pairs/c-text-pair-field.vue';
import CJsonField from '@/components/forms/fields/c-json-field.vue';
import CRetryField from '@/components/forms/fields/c-retry-field.vue';
import CMixedField from '@/components/forms/fields/c-mixed-field.vue';
import CMixedInputField from '@/components/forms/fields/c-mixed-input-field.vue';
import CInputTypeField from '@/components/forms/fields/c-input-type-field.vue';
import CArrayMixedField from '@/components/forms/fields/c-array-mixed-field.vue';
import CColorPickerField from '@/components/forms/fields/color/c-color-picker-field.vue';
import CEntityTypeField from '@/components/forms/fields/entity/c-entity-type-field.vue';
import CImpactLevelField from '@/components/forms/fields/entity/c-impact-level-field.vue';
import CSearchField from '@/components/forms/fields/c-search-field.vue';
import CAdvancedSearchField from '@/components/forms/fields/c-advanced-search-field.vue';
import CEntityCategoryField from '@/components/forms/fields/entity/c-entity-category-field.vue';
import CStoragesField from '@/components/forms/fields/c-storages-field.vue';
import CStorageField from '@/components/forms/fields/c-storage-field.vue';
import CFileNameMaskField from '@/components/forms/fields/c-file-name-mask-field.vue';
import CPercentsField from '@/components/forms/fields/c-percents-field.vue';
import CColumnsField from '@/components/forms/fields/c-columns-field.vue';
import CColorIndicatorField from '@/components/forms/fields/color/c-color-indicator-field.vue';
import CMiniBarChart from '@/components/common/chart/c-mini-bar-chart.vue';
import CImagesViewer from '@/components/common/images-viewer/c-images-viewer.vue';
import CClickableTooltip from '@/components/common/clickable-tooltip/c-clickable-tooltip.vue';
import CRolePickerField from '@/components/forms/fields/c-role-picker-field.vue';
import CUserPickerField from '@/components/forms/fields/c-user-picker-field.vue';
import CInstructionTypeField from '@/components/forms/fields/c-instruction-type-field.vue';
import CPriorityField from '@/components/forms/fields/c-priority-field.vue';
import CQuickDateIntervalField from '@/components/forms/fields/c-quick-date-interval-field.vue';
import CDateIntervalField from '@/components/forms/fields/date-picker/c-date-interval-field.vue';
import CDateTimeIntervalField from '@/components/forms/fields/date-time-picker/c-date-time-interval-field.vue';
import CQuickDateIntervalTypeField from '@/components/forms/fields/c-quick-date-interval-type-field.vue';
import CEnabledDurationField from '@/components/forms/fields/duration/c-enabled-duration-field.vue';
import CEnabledLimitField from '@/components/forms/fields/c-enabled-limit-field.vue';
import CTimezoneField from '@/components/forms/fields/c-timezone-field.vue';
import CLanguageField from '@/components/forms/fields/c-language-field.vue';
import CSamplingField from '@/components/forms/fields/c-sampling-field.vue';
import CAlarmMetricParametersField from '@/components/forms/fields/alarm/c-alarm-metric-parameters-field.vue';
import CFilterField from '@/components/forms/fields/pattern/c-filter-field.vue';
import CEntityStateField from '@/components/forms/fields/entity/c-entity-state-field.vue';
import CRecordsPerPageField from '@/components/forms/fields/c-records-per-page-field.vue';
import COperatorField from '@/components/forms/fields/c-operator-field.vue';
import CIdField from '@/components/forms/fields/c-id-field.vue';
import CDescriptionField from '@/components/forms/fields/c-description-field.vue';
import CEventFilterTypeField from '@/components/forms/fields/c-event-filter-type-field.vue';
import CDraggableListField from '@/components/forms/fields/c-draggable-list-field.vue';
import CDatePickerField from '@/components/forms/fields/date-picker/c-date-picker-field.vue';
import CEntityStatusField from '@/components/forms/fields/entity/c-entity-status-field.vue';
import CNumberField from '@/components/forms/fields/c-number-field.vue';
import CEntityField from '@/components/forms/fields/entity/c-entity-field.vue';
import CPbehaviorReasonField from '@/components/forms/fields/pbehavior/c-pbehavior-reason-field.vue';
import CPbehaviorTypeField from '@/components/forms/fields/pbehavior/c-pbehavior-type-field.vue';
import CSelectField from '@/components/forms/fields/c-select-field.vue';
import CCollapsePanel from '@/components/common/block/c-collapse-panel.vue';

/**
 * Patterns
 */
import CPatternField from '@/components/forms/fields/pattern/c-pattern-field.vue';
import CPatternAttributeField from '@/components/forms/fields/pattern/c-pattern-attribute-field.vue';
import CPatternInfosAttributeField from '@/components/forms/fields/pattern/c-pattern-infos-attribute-field.vue';
import CPatternExtraInfosAttributeField from '@/components/forms/fields/pattern/c-pattern-extra-infos-attribute-field.vue';
import CPatternOperatorField from '@/components/forms/fields/pattern/c-pattern-operator-field.vue';
import CPatternRuleField from '@/components/forms/fields/pattern/c-pattern-rule-field.vue';
import CPatternRulesField from '@/components/forms/fields/pattern/c-pattern-rules-field.vue';
import CPatternGroupField from '@/components/forms/fields/pattern/c-pattern-group-field.vue';
import CPatternGroupsField from '@/components/forms/fields/pattern/c-pattern-groups-field.vue';
import CPatternsEditorField from '@/components/forms/fields/pattern/c-patterns-editor-field.vue';
import CPatternsAdvancedEditorField from '@/components/forms/fields/pattern/c-patterns-advanced-editor-field.vue';
import CEntityPatternsField from '@/components/forms/fields/entity/c-entity-patterns-field.vue';
import CAlarmPatternsField from '@/components/forms/fields/alarm/c-alarm-patterns-field.vue';
import CPatternPanel from '@/components/forms/fields/pattern/c-pattern-panel.vue';
import CPbehaviorPatternsField from '@/components/forms/fields/pbehavior/c-pbehavior-patterns-field.vue';
import CEventFilterPatternsField from '@/components/forms/fields/event-filter/c-event-filter-patterns-field.vue';

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
Vue.component('c-mixed-input-field', CMixedInputField);
Vue.component('c-input-type-field', CInputTypeField);
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
Vue.component('c-date-picker-field', CDatePickerField);
Vue.component('c-date-interval-field', CDateIntervalField);
Vue.component('c-date-time-interval-field', CDateTimeIntervalField);
Vue.component('c-quick-date-interval-field', CQuickDateIntervalField);
Vue.component('c-quick-date-interval-type-field', CQuickDateIntervalTypeField);
Vue.component('c-enabled-duration-field', CEnabledDurationField);
Vue.component('c-enabled-limit-field', CEnabledLimitField);
Vue.component('c-timezone-field', CTimezoneField);
Vue.component('c-language-field', CLanguageField);
Vue.component('c-filter-field', CFilterField);
Vue.component('c-state-count-changes-chips', CStateCountChangesChips);
Vue.component('c-information-block', CInformationBlock);
Vue.component('c-information-block-row', CInformationBlockRow);
Vue.component('c-responsive-list', CResponsiveList);
Vue.component('c-sampling-field', CSamplingField);
Vue.component('c-alarm-metric-parameters-field', CAlarmMetricParametersField);
Vue.component('c-entity-state-field', CEntityStateField);
Vue.component('c-entity-status-field', CEntityStatusField);
Vue.component('c-records-per-page-field', CRecordsPerPageField);
Vue.component('c-operator-field', COperatorField);
Vue.component('c-id-field', CIdField);
Vue.component('c-description-field', CDescriptionField);
Vue.component('c-event-filter-type-field', CEventFilterTypeField);
Vue.component('c-draggable-list-field', CDraggableListField);
Vue.component('c-number-field', CNumberField);
Vue.component('c-select-field', CSelectField);
Vue.component('c-entity-field', CEntityField);
Vue.component('c-pbehavior-reason-field', CPbehaviorReasonField);
Vue.component('c-pbehavior-type-field', CPbehaviorTypeField);
Vue.component('c-collapse-panel', CCollapsePanel);

Vue.component('c-pattern-attribute-field', CPatternAttributeField);
Vue.component('c-pattern-infos-attribute-field', CPatternInfosAttributeField);
Vue.component('c-pattern-extra-infos-attribute-field', CPatternExtraInfosAttributeField);
Vue.component('c-pattern-operator-field', CPatternOperatorField);
Vue.component('c-pattern-rule-field', CPatternRuleField);
Vue.component('c-pattern-rules-field', CPatternRulesField);
Vue.component('c-pattern-group-field', CPatternGroupField);
Vue.component('c-pattern-groups-field', CPatternGroupsField);
Vue.component('c-patterns-editor-field', CPatternsEditorField);
Vue.component('c-patterns-advanced-editor-field', CPatternsAdvancedEditorField);
Vue.component('c-pattern-field', CPatternField);
Vue.component('c-pattern-operator-information', CPatternOperatorInformation);
Vue.component('c-pattern-operator-chip', CPatternOperatorChip);
Vue.component('c-pattern-panel', CPatternPanel);
Vue.component('c-alarm-patterns-field', CAlarmPatternsField);
Vue.component('c-entity-patterns-field', CEntityPatternsField);
Vue.component('c-pbehavior-patterns-field', CPbehaviorPatternsField);
Vue.component('c-event-filter-patterns-field', CEventFilterPatternsField);

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
    [MODALS.alarmsList]: { fullscreen: true, lazy: true },
    [MODALS.createFilter]: { maxWidth: 1100, lazy: true },
    [MODALS.textEditor]: { maxWidth: 700, lazy: true, persistent: true },
    [MODALS.addInfoPopup]: { maxWidth: 700, lazy: true, persistent: true },
    [MODALS.serviceEntities]: { maxWidth: 920, lazy: true },
    [MODALS.serviceDependencies]: { maxWidth: 1100, lazy: true },
    [MODALS.importExportViews]: { maxWidth: 920, persistent: true },
    [MODALS.createPlaylist]: { maxWidth: 920, lazy: true },
    [MODALS.pbehaviorPlanning]: { fullscreen: true, lazy: true, persistent: true },
    [MODALS.pbehaviorRecurrenceRule]: { fullscreen: true, lazy: true, persistent: true },
    [MODALS.pbehaviorRecurrentChangesConfirmation]: { maxWidth: 400, persistent: true },
    [MODALS.createRemediationInstruction]: { maxWidth: 960 },
    [MODALS.remediationInstructionApproval]: { maxWidth: 960 },
    [MODALS.executeRemediationInstruction]: { maxWidth: 960, persistent: true },
    [MODALS.imageViewer]: { maxWidth: '90%', contentClass: 'v-dialog__image-viewer' },
    [MODALS.imagesViewer]: { maxWidth: '100%', contentClass: 'v-dialog__images-viewer' },
    [MODALS.rate]: { maxWidth: 500 },
    [MODALS.createMetaAlarmRule]: { maxWidth: 920, lazy: true },
    [MODALS.createEventFilter]: { maxWidth: 1280 },
    [MODALS.createEventFilterAction]: { maxWidth: 920 },
    [MODALS.testSuite]: { maxWidth: 920 },
    [MODALS.createPattern]: { maxWidth: 1280 },
    [MODALS.remediationPatterns]: { maxWidth: 1280 },
    [MODALS.pbehaviorPatterns]: { maxWidth: 1280 },
    [MODALS.createIdleRule]: { maxWidth: 1280 },
    [MODALS.createScenario]: { maxWidth: 1280 },
    [MODALS.createKpiFilter]: { maxWidth: 1280 },

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
