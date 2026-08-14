import { cloneDeep, map } from 'lodash';

import ALARM_EXPORT_PDF_TEMPLATE from '@/assets/templates/alarm-export-pdf.html?raw';
import { EXPAND_DEFAULT_MAX_LETTERS, PAGINATION_LIMIT } from '@/config';
import {
  ALARM_FIELDS,
  ALARM_FIELDS_TO_LABELS_KEYS,
  ALARM_UNSORTABLE_FIELDS,
  DEFAULT_ALARMS_QUICK_ACTIONS,
  DEFAULT_ALARMS_WIDGET_GROUP_COLUMNS,
  DEFAULT_LINKS_INLINE_COUNT,
  DEFAULT_PERIODIC_REFRESH,
  DEFAULT_SERVICE_DEPENDENCIES_COLUMNS,
  DENSE_TYPES,
  ENTITY_FIELDS_TO_LABELS_KEYS,
  COLOR_INDICATOR_TYPES,
  DEFAULT_ALARM_MORE_INFO_TEMPLATE,
  DEFAULT_ALARMS_WIDGET_COLUMNS,
  EXPORT_CSV_DATETIME_FORMATS,
  CSV_SEPARATORS,
  GRID_SIZES,
  RESIZING_CELLS_CONTENTS_BEHAVIORS,
  TIME_UNITS,
} from '@/constants';
import { addKeyInEntities, addKeyInEntity, removeKeyFromEntities } from '@/helpers/array';
import { convertDateToStringWithFormatForToday } from '@/helpers/date/date';
import {
  convertDurationToString,
  durationWithEnabledToForm,
  isValidUnit,
} from '@/helpers/date/duration';
import { convertAlarmWidgetParametersToActiveColumns } from '@/helpers/entities/alarm/query';
import { hasStateSetting } from '@/helpers/entities/entity/entity';
import { kioskParametersToForm } from '@/helpers/entities/shared/kiosk/form';
import { availabilityFieldToForm } from '@/helpers/entities/widget/forms/availability';
import { getWidgetColumnLabel, getWidgetColumnSortable } from '@/helpers/entities/widget/list';
import { setSeveralFields } from '@/helpers/immutable';
import { uid } from '@/helpers/uid';

import { formToWidgetColumns, widgetColumnsToForm } from '../column/form';
import { formToWidgetQuickActions, widgetQuickActionsToForm } from '../quick-action/form';
import { formToWidgetSortColumns, widgetSortColumnsToForm } from '../sort-column/form';
import { formToWidgetTemplateValue, widgetTemplateValueToForm } from '../template/form';

import { barChartWidgetParametersToForm, formToBarChartWidgetParameters } from './bar-chart';
import { lineChartWidgetParametersToForm, formToLineChartWidgetParameters } from './line-chart';
import { numbersWidgetParametersToForm, formToNumbersWidgetParameters } from './numbers-chart';

/**
 * @typedef { 'wrap' | 'truncate' } AlarmsResizingBehaviors
 */

/**
 * @typedef {Object} WidgetColumnsParameters
 * @property {boolean} draggable
 * @property {boolean} resizable
 * @property {AlarmsResizingBehaviors} cells_content_behavior
 */

/**
 * @typedef {Object} WidgetInfoPopup
 * @property {string} column
 * @property {string} template
 */

/**
 * @typedef {Object} AlarmListBaseParameters
 * @property {number} itemsPerPage
 * @property {WidgetSortColumn} sort
 * @property {string} sortTemplate
 * @property {string} moreInfoTemplate
 * @property {string} moreInfoTemplateTemplate
 * @property {WidgetInfoPopup[]} infoPopups
 * @property {string} widgetColumnsTemplate
 * @property {WidgetColumn[]} widgetColumns
 * @property {string} exportPdfTemplate
 * @property {string} exportPdfTemplateTemplate
 * @property {boolean} showRootCauseByStateClick
 * @property {ColorIndicator} rootCauseColorIndicator
 * @property {WidgetColumnsParameters} columns
 */

/**
 * @typedef {AlarmListBaseParameters} AlarmListBaseParametersForm
 * @property {string | Symbol} widgetColumnsTemplate
 * @property {WidgetColumnForm[]} widgetColumns
 */

/**
 * Convert opened field widget.
 *
 * @param {boolean | null} [opened]
 * @returns {boolean | null}
 */
export const openedToForm = opened =>
  typeof opened === 'boolean' || opened === null ? opened : true;

export const alarmListChartToForm = (chart = {}) => {
  const chartParametersToForm = {
    BarChart: barChartWidgetParametersToForm,
    LineChart: lineChartWidgetParametersToForm,
    Numbers: numbersWidgetParametersToForm,
  }[chart.type ?? 'BarChart'];

  return {
    ...cloneDeep(chart),
    type: chart.type ?? 'BarChart',
    title: chart.title ?? '',
    parameters: chartParametersToForm(chart.parameters),
  };
};

export const formToAlarmListChart = ({ type, title, parameters }) => {
  const formToChartParameters = {
    BarChart: formToBarChartWidgetParameters,
    LineChart: formToLineChartWidgetParameters,
    Numbers: formToNumbersWidgetParameters,
  }[type];

  return { type, title, parameters: formToChartParameters(parameters) };
};

/**
 * Convert columns parameters field widget.
 *
 * @param {WidgetColumnsParameters} [columns={}]
 * @returns {WidgetColumnsParameters}
 */
export const columnsParametersToForm = (columns = {}) => ({
  draggable: columns.draggable ?? false,
  resizable: columns.resizable ?? false,
  cells_content_behavior: columns.cells_content_behavior ?? RESIZING_CELLS_CONTENTS_BEHAVIORS.wrap,
});

/**
 * Convert alarm list infoPopups parameters to form.
 *
 * @param {WidgetInfoPopup[]} [infoPopups=[]]
 * @return {WidgetInfoPopup[]}
 */
export const infoPopupsToForm = (infoPopups = []) =>
  infoPopups.map(infoPopup => ({ ...infoPopup }));

/**
 * Convert alarm list base parameters to form.
 *
 * @param {AlarmListBaseParameters} [alarmListParameters={}]
 * @return {AlarmListBaseParametersForm}
 */
export const alarmListBaseParametersToForm = (alarmListParameters = {}) => ({
  sort: widgetSortColumnsToForm(alarmListParameters.sort),
  sortTemplate: widgetTemplateValueToForm(alarmListParameters.sortTemplate),
  itemsPerPage: alarmListParameters.itemsPerPage ?? PAGINATION_LIMIT,
  moreInfoTemplate: alarmListParameters.moreInfoTemplate ?? DEFAULT_ALARM_MORE_INFO_TEMPLATE,
  moreInfoTemplateTemplate: widgetTemplateValueToForm(alarmListParameters.moreInfoTemplateTemplate),
  infoPopups: infoPopupsToForm(alarmListParameters.infoPopups),
  widgetColumnsTemplate: widgetTemplateValueToForm(alarmListParameters.widgetColumnsTemplate),
  widgetColumns: widgetColumnsToForm(
    alarmListParameters.widgetColumns ?? DEFAULT_ALARMS_WIDGET_COLUMNS,
  ),
  columns: columnsParametersToForm(alarmListParameters.columns),
  exportPdfTemplate: alarmListParameters.exportPdfTemplate ?? ALARM_EXPORT_PDF_TEMPLATE,
  exportPdfTemplateTemplate: widgetTemplateValueToForm(
    alarmListParameters.exportPdfTemplateTemplate,
  ),
  showRootCauseByStateClick: alarmListParameters.showRootCauseByStateClick ?? true,
  rootCauseColorIndicator:
    alarmListParameters.rootCauseColorIndicator ?? COLOR_INDICATOR_TYPES.state,
});

/**
 * Convert form to alarm list base parameters.
 *
 * @param {AlarmListBaseParametersForm} [form={}]
 * @return {AlarmListBaseParameters}
 */
export const formToAlarmListBaseParameters = (form = {}) => ({
  ...form,
  moreInfoTemplateTemplate: formToWidgetTemplateValue(form.moreInfoTemplateTemplate),
  exportPdfTemplateTemplate: formToWidgetTemplateValue(form.exportPdfTemplateTemplate),
  widgetColumnsTemplate: formToWidgetTemplateValue(form.widgetColumnsTemplate),
  widgetColumns: formToWidgetColumns(form.widgetColumns),
  sort: formToWidgetSortColumns(form.sort),
  sortTemplate: formToWidgetTemplateValue(form.sortTemplate),
});

/**
 * Convert widget periodic refresh to the duration form.
 *
 * @param {DurationWithEnabled} [periodicRefresh]
 * @returns {DurationWithEnabled}
 */
export const periodicRefreshToDurationForm = (periodicRefresh = DEFAULT_PERIODIC_REFRESH) => ({
  ...durationWithEnabledToForm(periodicRefresh),
  unit: isValidUnit(periodicRefresh.unit) ? periodicRefresh.unit : TIME_UNITS.second,
});

/**
 * Convert one fast pbehavior setting to form values.
 *
 * @param {Object} [fastPbehavior={}]
 * @returns {Object}
 */
export const alarmListFastPbehaviorParametersItemToForm = (fastPbehavior = {}) =>
  addKeyInEntity({
    name_prefix: fastPbehavior.name_prefix ?? '',
    type: fastPbehavior.type ?? '',
    reason: fastPbehavior.reason ?? '',
  });

/**
 * Convert fast pbehavior settings to form values.
 *
 * @param {Object[]} [fastPbehavior=[]]
 * @returns {Object[]}
 */
export const alarmListFastPbehaviorParametersToForm = (fastPbehavior = []) =>
  fastPbehavior.map(alarmListFastPbehaviorParametersItemToForm);

/**
 * Convert alarm list widget parameters to form values.
 *
 * @param {Object} [parameters={}]
 * @returns {Object}
 */
export const alarmListWidgetParametersToForm = (parameters = {}) => ({
  ...cloneDeep(parameters),
  itemsPerPage: parameters.itemsPerPage ?? PAGINATION_LIMIT,
  infoPopups: infoPopupsToForm(parameters.infoPopups),
  moreInfoTemplate: parameters.moreInfoTemplate ?? DEFAULT_ALARM_MORE_INFO_TEMPLATE,
  moreInfoTemplateTemplate: widgetTemplateValueToForm(parameters.moreInfoTemplateTemplate),
  isAckNoteRequired: !!parameters.isAckNoteRequired,
  isSnoozeNoteRequired: !!parameters.isSnoozeNoteRequired,
  isRemoveAlarmsFromMetaAlarmCommentRequired:
    parameters.isRemoveAlarmsFromMetaAlarmCommentRequired ?? true,
  isUncancelAlarmsCommentRequired: parameters.isUncancelAlarmsCommentRequired ?? true,
  isMultiAckEnabled: !!parameters.isMultiAckEnabled,
  isMultiDeclareTicketEnabled: !!parameters.isMultiDeclareTicketEnabled,
  isHtmlEnabledOnTimeLine: parameters.isHtmlEnabledOnTimeLine ?? true,
  openExpandPanelByDoubleClick: parameters.openExpandPanelByDoubleClick ?? false,
  pausePeriodicRefreshOnExpandPanel: parameters.pausePeriodicRefreshOnExpandPanel ?? false,
  isActionsAllowWithOkState: !!parameters.isActionsAllowWithOkState,
  keepSelectedAfterAction: !!parameters.keepSelectedAfterAction,
  isVirtualScrollEnabled: !!parameters.isVirtualScrollEnabled,
  isCorrelationEnabled: !!parameters.isCorrelationEnabled,
  fast_pbehaviors: alarmListFastPbehaviorParametersToForm(parameters.fast_pbehaviors),
  sticky_header: !!parameters.sticky_header,
  sticky_horizontal_scroll: !!parameters.sticky_horizontal_scroll,
  dense: parameters.dense ?? DENSE_TYPES.large,
  fastAckOutput: parameters.fastAckOutput
    ? { ...parameters.fastAckOutput }
    : { enabled: false, value: 'auto ack' },
  fastCancelOutput: parameters.fastCancelOutput
    ? { ...parameters.fastCancelOutput }
    : { enabled: false, value: 'auto cancel' },
  widgetColumnsTemplate: widgetTemplateValueToForm(parameters.widgetColumnsTemplate),
  widgetGroupColumnsTemplate: widgetTemplateValueToForm(parameters.widgetGroupColumnsTemplate),
  serviceDependenciesColumnsTemplate: widgetTemplateValueToForm(
    parameters.serviceDependenciesColumnsTemplate,
  ),
  widgetExportColumnsTemplate: widgetTemplateValueToForm(parameters.widgetExportColumnsTemplate),
  widgetColumns: widgetColumnsToForm(parameters.widgetColumns ?? DEFAULT_ALARMS_WIDGET_COLUMNS),
  widgetGroupColumns: widgetColumnsToForm(
    parameters.widgetGroupColumns ?? DEFAULT_ALARMS_WIDGET_GROUP_COLUMNS,
  ),
  serviceDependenciesColumns: widgetColumnsToForm(
    parameters.serviceDependenciesColumns ?? DEFAULT_SERVICE_DEPENDENCIES_COLUMNS,
  ),
  widgetExportColumns: widgetColumnsToForm(
    parameters.widgetExportColumns ?? DEFAULT_ALARMS_WIDGET_COLUMNS,
  ),
  inlineLinksCount: parameters.inlineLinksCount ?? DEFAULT_LINKS_INLINE_COUNT,
  exportPdfTemplate: parameters.exportPdfTemplate ?? ALARM_EXPORT_PDF_TEMPLATE,
  exportPdfTemplateTemplate: widgetTemplateValueToForm(parameters.exportPdfTemplateTemplate),
  showRootCauseByStateClick: parameters.showRootCauseByStateClick ?? true,
  rootCauseColorIndicator: parameters.rootCauseColorIndicator ?? COLOR_INDICATOR_TYPES.state,
  availability: availabilityFieldToForm(parameters.availability),
  quickActionsTemplate: widgetTemplateValueToForm(parameters.quickActionsTemplate),
  quickMassActionsTemplate: widgetTemplateValueToForm(parameters.quickMassActionsTemplate),
  quickActions: widgetQuickActionsToForm(parameters.quickActions ?? DEFAULT_ALARMS_QUICK_ACTIONS),
  quickMassActions: widgetQuickActionsToForm(
    parameters.quickMassActions ?? DEFAULT_ALARMS_QUICK_ACTIONS,
  ),
  hideMassActions: parameters.hideMassActions ?? false,
  comment_templates: (parameters.comment_templates ?? []).map(template => ({
    template,
    key: uid(),
  })),
  sort: widgetSortColumnsToForm(parameters.sort),
  sortTemplate: widgetTemplateValueToForm(parameters.sortTemplate),
  periodic_refresh: periodicRefreshToDurationForm(parameters.periodic_refresh),
  liveWatching: parameters.liveWatching ?? false,
  mainFilter: parameters.mainFilter ?? null,
  clearFilterEnabled: parameters.clearFilterEnabled ?? true,
  liveReporting: parameters.liveReporting ? cloneDeep(parameters.liveReporting) : {},
  opened: openedToForm(parameters.opened),
  expandGridRangeSize: parameters.expandGridRangeSize
    ? [...parameters.expandGridRangeSize]
    : [GRID_SIZES.min, GRID_SIZES.max],
  exportCsvSeparator: parameters.exportCsvSeparator ?? CSV_SEPARATORS.comma,
  exportCsvDatetimeFormat:
    parameters.exportCsvDatetimeFormat ?? EXPORT_CSV_DATETIME_FORMATS.datetimeSeconds.value,
  kiosk: kioskParametersToForm(parameters.kiosk),
  columns: columnsParametersToForm(parameters.columns),
  charts: addKeyInEntities(parameters.charts),
  usedAlarmProperties: [],
});

/**
 * Convert alarm list widget form values to widget parameters.
 *
 * @param {Object} form
 * @returns {Object}
 */
export const formToAlarmListWidgetParameters = (form = {}) => {
  const parameters = {
    ...cloneDeep(form),
    moreInfoTemplateTemplate: formToWidgetTemplateValue(form.moreInfoTemplateTemplate),
    exportPdfTemplateTemplate: formToWidgetTemplateValue(form.exportPdfTemplateTemplate),
    widgetColumnsTemplate: formToWidgetTemplateValue(form.widgetColumnsTemplate),
    widgetGroupColumnsTemplate: formToWidgetTemplateValue(form.widgetGroupColumnsTemplate),
    serviceDependenciesColumnsTemplate: formToWidgetTemplateValue(
      form.serviceDependenciesColumnsTemplate,
    ),
    widgetExportColumnsTemplate: formToWidgetTemplateValue(form.widgetExportColumnsTemplate),
    widgetColumns: formToWidgetColumns(form.widgetColumns),
    widgetGroupColumns: formToWidgetColumns(form.widgetGroupColumns),
    widgetExportColumns: formToWidgetColumns(form.widgetExportColumns),
    serviceDependenciesColumns: formToWidgetColumns(form.serviceDependenciesColumns),
    charts: removeKeyFromEntities(form.charts),
    quickActionsTemplate: formToWidgetTemplateValue(form.quickActionsTemplate),
    quickMassActionsTemplate: formToWidgetTemplateValue(form.quickMassActionsTemplate),
    quickActions: formToWidgetQuickActions(form.quickActions),
    quickMassActions: formToWidgetQuickActions(form.quickMassActions),
    comment_templates: map(form.comment_templates, 'template'),
    fast_pbehaviors: removeKeyFromEntities(form.fast_pbehaviors),
    sort: formToWidgetSortColumns(form.sort),
    sortTemplate: formToWidgetTemplateValue(form.sortTemplate),
  };

  parameters.usedAlarmProperties = convertAlarmWidgetParametersToActiveColumns(parameters);

  return parameters;
};

/**
 * Prepared alarms list widget for displaying.
 *
 * @param {Object} widget
 * @returns {Object}
 */
export const prepareAlarmListWidget = (widget = {}) =>
  setSeveralFields(widget, {
    'parameters.widgetColumns': (columns = []) =>
      columns.map(column => ({
        ...column,

        sortable: getWidgetColumnSortable(column, ALARM_UNSORTABLE_FIELDS),
        text: getWidgetColumnLabel(column, ALARM_FIELDS_TO_LABELS_KEYS),
      })),

    'parameters.widgetGroupColumns': (columns = DEFAULT_ALARMS_WIDGET_GROUP_COLUMNS) =>
      columns.map(column => ({
        ...column,

        sortable: getWidgetColumnSortable(column, ALARM_UNSORTABLE_FIELDS),
        text: getWidgetColumnLabel(column, ALARM_FIELDS_TO_LABELS_KEYS),
      })),

    'parameters.serviceDependenciesColumns': (columns = DEFAULT_SERVICE_DEPENDENCIES_COLUMNS) =>
      columns.map(column => ({
        ...column,

        sortable: false,
        text: getWidgetColumnLabel(column, ENTITY_FIELDS_TO_LABELS_KEYS),
      })),

    'parameters.widgetExportColumns': (columns = []) =>
      columns.map(column => ({
        ...column,

        text: getWidgetColumnLabel(column, ALARM_FIELDS_TO_LABELS_KEYS),
      })),
  });

/**
 * Get a display formatter for an AlarmList column value.
 *
 * @param {string} value
 * @returns {Function|undefined}
 */
export const getAlarmsListWidgetColumnValueFilter = (value) => {
  switch (value) {
    case ALARM_FIELDS.lastUpdateDate:
    case ALARM_FIELDS.creationDate:
    case ALARM_FIELDS.lastEventDate:
    case ALARM_FIELDS.activationDate:
    case ALARM_FIELDS.ackAt:
    case ALARM_FIELDS.stateAt:
    case ALARM_FIELDS.statusAt:
    case ALARM_FIELDS.resolved:
    case ALARM_FIELDS.timestamp:
    case ALARM_FIELDS.entityLastPbehaviorDate:
      return convertDateToStringWithFormatForToday;
    case ALARM_FIELDS.duration:
    case ALARM_FIELDS.currentStateDuration:
    case ALARM_FIELDS.activeDuration:
    case ALARM_FIELDS.snoozeDuration:
    case ALARM_FIELDS.pbhInactiveDuration:
      return convertDurationToString;
    default:
      return undefined;
  }
};

/**
 * Get a component descriptor for a special AlarmList column value.
 *
 * @param {Object} column
 * @param {Object} [options]
 * @returns {Function}
 */
export const getAlarmsListWidgetColumnComponentGetter = (
  { value, onlyIcon, inlineLinksCount },
  { showRootCauseByStateClick = true, isCellContentTruncated } = {},
) => {
  if (value === ALARM_FIELDS.state) {
    return (context) => {
      const state = context.alarm.v?.state;
      const component = {
        bind: {
          is: 'c-alarm-state-chip',
          value: state?.val,
          type: state?._t,
          badgeValue: context.alarm.v?.events_count,
          small: context.small,
        },
      };

      if (showRootCauseByStateClick && hasStateSetting(context.alarm.entity)) {
        component.bind.class = 'cursor-pointer';
        component.bind.appendIconName = '$vuetify.icons.root_cause';
        component.on = {
          click: () => context.$emit('click:state', { ...context.alarm.entity, state: state?.val }),
        };
      }

      return component;
    };
  }

  if (value === ALARM_FIELDS.maxState) {
    return context => ({
      bind: {
        is: 'c-alarm-state-chip',
        value: context.alarm.v?.max_state,
        small: context.small,
      },
    });
  }

  if (value === ALARM_FIELDS.initialState) {
    return context => ({
      bind: {
        is: 'c-alarm-state-chip',
        value: context.alarm.v?.initial_state,
        small: context.small,
      },
    });
  }

  if (value === ALARM_FIELDS.status) {
    return context => ({
      bind: { is: 'alarm-column-value-status', alarm: context.alarm, small: context.small },
    });
  }

  if (value === ALARM_FIELDS.extraDetails) {
    return context => ({
      bind: { is: 'alarm-column-value-extra-details', alarm: context.alarm, small: context.small },
    });
  }

  if (value === ALARM_FIELDS.tags) {
    return context => ({
      bind: {
        is: 'c-alarm-tags-chips',
        alarm: context.alarm,
        selectedTags: context.selectedTags,
        small: context.small,
      },
      on: {
        select: context.$listeners['select:tag'],
      },
    });
  }

  if (value === ALARM_FIELDS.impactState) {
    return context => ({
      bind: {
        is: 'color-indicator-wrapper',
        type: COLOR_INDICATOR_TYPES.impactState,
        entity: context.alarm.entity,
        alarm: context.alarm,
      },
    });
  }

  if (value === ALARM_FIELDS.links || value.startsWith('links.')) {
    return context => ({
      bind: {
        is: 'c-alarm-links-chips',
        alarm: context.alarm,
        category: value.startsWith('links.') ? value.replace('links.', '') : undefined,
        onlyIcon,
        inlineCount: inlineLinksCount,
        small: context.small,
      },
      on: {
        activate: context.$listeners.activate,
      },
    });
  }

  const maxLetters = isCellContentTruncated ? Infinity : EXPAND_DEFAULT_MAX_LETTERS;

  return context => ({
    bind: {
      is: 'c-ellipsis',
      class: 'alarm-column-cell__text',
      title: context.value,
      text: context.value,
      maxLetters,
    },
  });
};

/**
 * Convert comment form to create comment event payload
 *
 * @param {Object} form - Comment form data
 * @returns {Object} Comment event payload with either comment or struct_comment
 */
export const createCommentFormToCreateCommentEvent = (form) => {
  if (!form.template) {
    return {
      comment: form.comment,
    };
  }

  return {
    struct_comment: (form.template?.fields ?? []).map(({ name: field }) => ({
      field,
      message: form[field],
    })).filter(({ message }) => message),
  };
};
