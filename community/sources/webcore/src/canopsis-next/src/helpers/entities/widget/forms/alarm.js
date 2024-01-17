import { cloneDeep, isBoolean, isNull, omit } from 'lodash';

import {
  ALARM_DENSE_TYPES,
  ALARM_FIELDS,
  ALARM_FIELDS_TO_LABELS_KEYS,
  ALARM_UNSORTABLE_FIELDS,
  ALARMS_RESIZING_CELLS_CONTENTS_BEHAVIORS,
  COLOR_INDICATOR_TYPES,
  DEFAULT_ALARMS_WIDGET_COLUMNS,
  DEFAULT_ALARMS_WIDGET_GROUP_COLUMNS,
  DEFAULT_LINKS_INLINE_COUNT,
  DEFAULT_PERIODIC_REFRESH,
  DEFAULT_SERVICE_DEPENDENCIES_COLUMNS,
  ENTITY_FIELDS_TO_LABELS_KEYS,
  EXPORT_CSV_DATETIME_FORMATS,
  EXPORT_CSV_SEPARATORS,
  GRID_SIZES,
  SORT_ORDERS,
  TIME_UNITS,
  WIDGET_TYPES,
} from '@/constants';
import { PAGINATION_LIMIT } from '@/config';

import { setSeveralFields } from '@/helpers/immutable';
import { convertDateToStringWithFormatForToday } from '@/helpers/date/date';
import { convertDurationToString, durationWithEnabledToForm, isValidUnit } from '@/helpers/date/duration';
import { addKeyInEntities, removeKeyFromEntities } from '@/helpers/array';
import { kioskParametersToForm } from '@/helpers/entities/shared/kiosk/form';
import { convertAlarmWidgetParametersToActiveColumns } from '@/helpers/entities/alarm/query';

import ALARM_EXPORT_PDF_TEMPLATE from '@/assets/templates/alarm-export-pdf.html';

import { formToWidgetTemplateValue, widgetTemplateValueToForm } from '../template/form';
import { formToWidgetColumns, widgetColumnsToForm } from '../column/form';
import { getWidgetColumnLabel, getWidgetColumnSortable } from '../list';

import { barChartWidgetParametersToForm, formToBarChartWidgetParameters } from './bar-chart';
import { formToLineChartWidgetParameters, lineChartWidgetParametersToForm } from './line-chart';
import { formToNumbersWidgetParameters, numbersWidgetParametersToForm } from './numbers-chart';

/**
 * @typedef {'BarChart', 'LineChart', 'Numbers'} AlarmChartType
 */

/**
 * @typedef { 'wrap' | 'truncate' } AlarmsResizingBehaviors
 */

/**
 * @typedef {Object} AlarmsListDataTableColumn
 * @property {string} value
 * @property {string} text
 * @property {boolean} [isHtml]
 * @property {boolean} [colorIndicator]
 */

/**
 * @typedef {Object} WidgetFastActionOutput
 * @property {boolean} enabled
 * @property {string} value
 */

/**
 * @typedef {Object} WidgetLinksCategoriesAsList
 * @property {boolean} enabled
 * @property {number} limit
 */

/**
 * @typedef {Object} WidgetLiveReporting
 * @property {string} [tstart]
 * @property {string} [tstop]
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
 * @typedef {Object} AlarmChartBarChartWidgetParameters
 * @property {MetricPreset[]} metrics
 * @property {boolean} stacked
 * @property {string} chart_title
 * @property {string} default_time_range
 * @property {Sampling} default_sampling
 * @property {boolean} comparison
 */

/**
 * @typedef {Object} AlarmChartLineChartWidgetParameters
 * @property {MetricPreset[]} metrics
 * @property {string} chart_title
 * @property {string} default_time_range
 * @property {Sampling} default_sampling
 * @property {boolean} comparison
 */

/**
 * @typedef {Object} AlarmChartNumbersWidgetParameters
 * @property {MetricPreset[]} metrics
 * @property {string} chart_title
 * @property {string} default_time_range
 * @property {Sampling} default_sampling
 * @property {string} show_trend
 * @property {number} [font_size]
 */

/**
 * @typedef {Object} AlarmChart
 * @property {string} title
 * @property {AlarmChartType} type
 * @property {
 *   AlarmChartBarChartWidgetParameters
 *   | AlarmChartLineChartWidgetParameters
 *   | AlarmChartNumbersWidgetParameters
 * } parameters
 */

/**

 * @typedef {Object} AlarmListBaseParameters
 * @property {number} itemsPerPage
 * @property {WidgetSort} sort
 * @property {string} moreInfoTemplate
 * @property {string} moreInfoTemplateTemplate
 * @property {WidgetInfoPopup[]} infoPopups
 * @property {string} widgetColumnsTemplate
 * @property {WidgetColumn[]} widgetColumns
 * @property {string} exportPdfTemplate
 * @property {string} exportPdfTemplateTemplate
 * @property {boolean} showRootCauseByStateClick
 */

/**
 * @typedef {Object} AlarmListWidgetDefaultParameters
 * @property {WidgetFastActionOutput} fastAckOutput
 * @property {WidgetFastActionOutput} fastCancelOutput
 * @property {number} inlineLinksCount
 * @property {number} itemsPerPage
 * @property {WidgetInfoPopup[]} infoPopups
 * @property {string} moreInfoTemplate
 * @property {string} moreInfoTemplateTemplate
 * @property {string} widgetColumnsTemplate
 * @property {string} widgetGroupColumnsTemplate
 * @property {string} widgetExportColumnsTemplate
 * @property {string} serviceDependenciesColumnsTemplate
 * @property {string} exportPdfTemplate
 * @property {string} exportPdfTemplateTemplate
 * @property {WidgetColumn[]} widgetColumns
 * @property {WidgetColumn[]} widgetGroupColumns
 * @property {WidgetColumn[]} widgetExportColumns
 * @property {WidgetColumn[]} serviceDependenciesColumns
 * @property {boolean} isAckNoteRequired
 * @property {boolean} isSnoozeNoteRequired
 * @property {boolean} isRemoveAlarmsFromMetaAlarmCommentRequired
 * @property {boolean} isUncancelAlarmsCommentRequired
 * @property {boolean} isMultiAckEnabled
 * @property {boolean} isMultiDeclareTicketEnabled
 * @property {boolean} isHtmlEnabledOnTimeLine
 * @property {boolean} isActionsAllowWithOkState
 * @property {boolean} sticky_header
 * @property {boolean} dense
 * @property {boolean} dense
 * @property {boolean} showRootCauseByStateClick
 */

/**
 * @typedef {AlarmListWidgetDefaultParameters} AlarmListWidgetParameters
 * @property {DurationWithEnabled} periodic_refresh
 * @property {boolean} liveWatching
 * @property {string | null} mainFilter
 * @property {WidgetLiveReporting} liveReporting
 * @property {WidgetSort} sort
 * @property {boolean | null} opened
 * @property {number[]} expandGridRangeSize
 * @property {WidgetCsvSeparator} exportCsvSeparator
 * @property {string} exportCsvDatetimeFormat
 * @property {boolean} clearFilterDisabled
 * @property {WidgetKioskParameters} kiosk
 * @property {AlarmChart[]} charts
 * @property {WidgetColumnsParameters} [columns]
 * @property {string[]} [usedAlarmProperties]
 */

/**
 * @typedef {AlarmListBaseParameters} AlarmListBaseParametersForm
 * @property {string | Symbol} widgetColumnsTemplate
 * @property {WidgetColumnForm[]} widgetColumns
 */

/**
 * @typedef {AlarmChart & ObjectKey} AlarmChartForm
 */

/**
 * @typedef {AlarmListWidgetDefaultParameters} AlarmListWidgetDefaultParametersForm
 * @property {string | Symbol} widgetColumnsTemplate
 * @property {string | Symbol} widgetGroupColumnsTemplate
 * @property {string | Symbol} widgetExportColumnsTemplate
 * @property {string | Symbol} serviceDependenciesColumnsTemplate
 * @property {WidgetColumnForm[]} widgetColumns
 * @property {WidgetColumnForm[]} widgetGroupColumns
 * @property {WidgetColumnForm[]} widgetExportColumns
 * @property {WidgetColumnForm[]} serviceDependenciesColumns
 */

/**
 * @typedef {WidgetColumnsParameters} WidgetColumnsParametersForm
 */

/**
 * @typedef {AlarmListWidgetDefaultParametersForm & AlarmListWidgetParameters} AlarmListWidgetParametersForm
 * @property {AlarmChartForm[]} charts
 * @property {WidgetColumnsParametersForm} columns
 */

/**
 * Convert opened field widget
 *
 * @param  {boolean | null} [opened]
 * @returns {boolean | null}
 */
export const openedToForm = (opened) => {
  if (isBoolean(opened) || isNull(opened)) {
    return opened;
  }

  return true;
};

/**
 * Convert columns parameters field widget
 *
 * @param  {WidgetColumnsParameters} [columns]
 * @returns {WidgetColumnsParametersForm}
 */
export const columnsParametersToForm = (columns = {}) => ({
  draggable: columns.draggable ?? false,
  resizable: columns.resizable ?? false,
  cells_content_behavior: columns.cells_content_behavior ?? ALARMS_RESIZING_CELLS_CONTENTS_BEHAVIORS.wrap,
});

/**
 * Convert alarm list infoPopups parameters to form
 *
 * @param {WidgetInfoPopup[]} [infoPopups = []]
 * @return {WidgetInfoPopup[]}
 */
export const infoPopupsToForm = (infoPopups = []) => infoPopups.map(infoPopup => ({ ...infoPopup }));

/**
 * Convert widget sort parameters to form
 *
 * @param {WidgetSort} [sort = {}]
 * @return {WidgetSort}
 */
export const widgetSortToForm = (sort = {}) => ({
  order: sort.order ?? SORT_ORDERS.asc,
  column: sort.column ?? '',
});

/**
 * Convert alarm list base parameters (we are using it inside another widgets with alarmList) to form
 *
 * @param {AlarmListBaseParameters} [alarmListParameters = {}]
 * @return {AlarmListBaseParametersForm}
 */
export const alarmListBaseParametersToForm = (alarmListParameters = {}) => ({
  sort: widgetSortToForm(alarmListParameters.sort),
  itemsPerPage: alarmListParameters.itemsPerPage ?? PAGINATION_LIMIT,
  moreInfoTemplate: alarmListParameters.moreInfoTemplate ?? '',
  moreInfoTemplateTemplate: widgetTemplateValueToForm(alarmListParameters.moreInfoTemplateTemplate),
  infoPopups: infoPopupsToForm(alarmListParameters.infoPopups),
  widgetColumnsTemplate: widgetTemplateValueToForm(alarmListParameters.widgetColumnsTemplate),
  widgetColumns: widgetColumnsToForm(alarmListParameters.widgetColumns ?? DEFAULT_ALARMS_WIDGET_COLUMNS),
  exportPdfTemplate: alarmListParameters.exportPdfTemplate ?? ALARM_EXPORT_PDF_TEMPLATE,
  exportPdfTemplateTemplate: widgetTemplateValueToForm(alarmListParameters.exportPdfTemplateTemplate),
  showRootCauseByStateClick: alarmListParameters.showRootCauseByStateClick ?? true,
});

/**
 * Convert widget periodic refresh to form duration
 *
 * @param {DurationWithEnabled} periodicRefresh
 * @returns {DurationWithEnabled}
 */
export const periodicRefreshToDurationForm = (periodicRefresh = DEFAULT_PERIODIC_REFRESH) => {
  /**
   * @link https://git.canopsis.net/canopsis/canopsis-pro/-/issues/4390
   */
  const unit = isValidUnit(periodicRefresh.unit)
    ? periodicRefresh.unit
    : TIME_UNITS.second;

  return durationWithEnabledToForm({ ...periodicRefresh, unit });
};

/**
 * Convert alarm list chart to form
 *
 * @param {AlarmChart} [chart = {}]
 * @returns {AlarmChart}
 */
export const alarmListChartToForm = (chart = {}) => {
  const convertersMap = {
    [WIDGET_TYPES.barChart]: barChartWidgetParametersToForm,
    [WIDGET_TYPES.lineChart]: lineChartWidgetParametersToForm,
    [WIDGET_TYPES.numbers]: numbersWidgetParametersToForm,
  };

  const type = chart.type ?? WIDGET_TYPES.barChart;
  const converter = convertersMap[type];

  return {
    type,
    title: chart.title ?? '',
    parameters: omit(converter ? converter(chart.parameters) : chart.parameters, ['periodic_refresh']),
  };
};

/**
 * Convert alarm list widget parameters to form
 *
 * @param {AlarmListWidgetDefaultParameters} [parameters = {}]
 * @return {AlarmListWidgetDefaultParametersForm}
 */
export const alarmListWidgetDefaultParametersToForm = (parameters = {}) => ({
  itemsPerPage: parameters.itemsPerPage ?? PAGINATION_LIMIT,
  infoPopups: infoPopupsToForm(parameters.infoPopups),
  moreInfoTemplate: parameters.moreInfoTemplate ?? '',
  moreInfoTemplateTemplate: widgetTemplateValueToForm(parameters.moreInfoTemplateTemplate),
  isAckNoteRequired: !!parameters.isAckNoteRequired,
  isSnoozeNoteRequired: !!parameters.isSnoozeNoteRequired,
  isRemoveAlarmsFromMetaAlarmCommentRequired: parameters.isRemoveAlarmsFromMetaAlarmCommentRequired ?? true,
  isUncancelAlarmsCommentRequired: parameters.isUncancelAlarmsCommentRequired ?? true,
  isMultiAckEnabled: !!parameters.isMultiAckEnabled,
  isMultiDeclareTicketEnabled: !!parameters.isMultiDeclareTicketEnabled,
  isHtmlEnabledOnTimeLine: parameters.isHtmlEnabledOnTimeLine ?? true,
  isActionsAllowWithOkState: !!parameters.isActionsAllowWithOkState,
  sticky_header: !!parameters.sticky_header,
  dense: parameters.dense ?? ALARM_DENSE_TYPES.large,
  fastAckOutput: parameters.fastAckOutput
    ? { ...parameters.fastAckOutput }
    : {
      enabled: false,
      value: 'auto ack',
    },
  fastCancelOutput: parameters.fastCancelOutput
    ? { ...parameters.fastCancelOutput }
    : {
      enabled: false,
      value: 'auto cancel',
    },
  widgetColumnsTemplate: widgetTemplateValueToForm(parameters.widgetColumnsTemplate),
  widgetGroupColumnsTemplate: widgetTemplateValueToForm(parameters.widgetGroupColumnsTemplate),
  serviceDependenciesColumnsTemplate: widgetTemplateValueToForm(parameters.serviceDependenciesColumnsTemplate),
  widgetExportColumnsTemplate: widgetTemplateValueToForm(parameters.widgetExportColumnsTemplate),
  widgetColumns:
    widgetColumnsToForm(parameters.widgetColumns ?? DEFAULT_ALARMS_WIDGET_COLUMNS),
  widgetGroupColumns:
    widgetColumnsToForm(parameters.widgetGroupColumns ?? DEFAULT_ALARMS_WIDGET_GROUP_COLUMNS),
  serviceDependenciesColumns:
    widgetColumnsToForm(parameters.serviceDependenciesColumns ?? DEFAULT_SERVICE_DEPENDENCIES_COLUMNS),
  widgetExportColumns:
    widgetColumnsToForm(parameters.widgetExportColumns ?? DEFAULT_ALARMS_WIDGET_COLUMNS),
  inlineLinksCount: parameters.inlineLinksCount ?? DEFAULT_LINKS_INLINE_COUNT,
  exportPdfTemplate: parameters.exportPdfTemplate ?? ALARM_EXPORT_PDF_TEMPLATE,
  exportPdfTemplateTemplate: widgetTemplateValueToForm(parameters.exportPdfTemplateTemplate),
  showRootCauseByStateClick: parameters.showRootCauseByStateClick ?? true,
});

/**
 * Convert alarm list widget parameters to form
 *
 * @param {AlarmListWidgetParameters} [parameters = {}]
 * @return {AlarmListWidgetParametersForm}
 */
export const alarmListWidgetParametersToForm = (parameters = {}) => ({
  ...parameters,
  ...alarmListWidgetDefaultParametersToForm(parameters),

  periodic_refresh: periodicRefreshToDurationForm(parameters.periodic_refresh),
  liveWatching: parameters.liveWatching ?? false,
  mainFilter: parameters.mainFilter ?? null,
  clearFilterDisabled: parameters.clearFilterDisabled ?? false,
  liveReporting: parameters.liveReporting
    ? cloneDeep(parameters.liveReporting)
    : {},
  sort: widgetSortToForm(parameters.sort),
  opened: openedToForm(parameters.opened),
  expandGridRangeSize: parameters.expandGridRangeSize
    ? [...parameters.expandGridRangeSize]
    : [GRID_SIZES.min, GRID_SIZES.max],
  exportCsvSeparator: parameters.exportCsvSeparator ?? EXPORT_CSV_SEPARATORS.comma,
  exportCsvDatetimeFormat: parameters.exportCsvDatetimeFormat ?? EXPORT_CSV_DATETIME_FORMATS.datetimeSeconds.value,
  kiosk: kioskParametersToForm(parameters.kiosk),
  columns: columnsParametersToForm(parameters.columns),
  charts: addKeyInEntities(parameters.charts),
  usedAlarmProperties: [],
});

/**
 * Convert form to alarm list base parameters (we are using it inside another widgets with alarmList)
 *
 * @param {AlarmListBaseParametersForm} [form = {}]
 * @return {AlarmListBaseParameters}
 */
export const formToAlarmListBaseParameters = (form = {}) => ({
  ...form,

  moreInfoTemplateTemplate: formToWidgetTemplateValue(form.moreInfoTemplateTemplate),
  exportPdfTemplateTemplate: formToWidgetTemplateValue(form.exportPdfTemplateTemplate),
  widgetColumnsTemplate: formToWidgetTemplateValue(form.widgetColumnsTemplate),
  widgetColumns: formToWidgetColumns(form.widgetColumns),
});

/**
 * Convert form to alarm list chart
 *
 * @param {AlarmChartForm} chart
 * @returns {AlarmChart}
 */
export const formToAlarmListChart = ({ type, title, parameters }) => {
  const convertersMap = {
    [WIDGET_TYPES.barChart]: formToBarChartWidgetParameters,
    [WIDGET_TYPES.lineChart]: formToLineChartWidgetParameters,
    [WIDGET_TYPES.numbers]: formToNumbersWidgetParameters,
  };

  const converter = convertersMap[type];

  return {
    type,
    title,
    parameters: omit(converter ? converter(parameters) : parameters, ['periodic_refresh']),
  };
};

/**
 * Convert form parameters to alarm list widget parameters
 *
 * @param {AlarmListWidgetParametersForm} form
 * @return {AlarmListWidgetParameters}
 */
export const formToAlarmListWidgetParameters = (form) => {
  const parameters = {
    ...form,

    moreInfoTemplateTemplate: formToWidgetTemplateValue(form.moreInfoTemplateTemplate),
    exportPdfTemplateTemplate: formToWidgetTemplateValue(form.exportPdfTemplateTemplate),
    widgetColumnsTemplate: formToWidgetTemplateValue(form.widgetColumnsTemplate),
    widgetGroupColumnsTemplate: formToWidgetTemplateValue(form.widgetGroupColumnsTemplate),
    serviceDependenciesColumnsTemplate: formToWidgetTemplateValue(form.serviceDependenciesColumnsTemplate),
    widgetExportColumnsTemplate: formToWidgetTemplateValue(form.widgetExportColumnsTemplate),
    widgetColumns: formToWidgetColumns(form.widgetColumns),
    widgetGroupColumns: formToWidgetColumns(form.widgetGroupColumns),
    widgetExportColumns: formToWidgetColumns(form.widgetExportColumns),
    serviceDependenciesColumns: formToWidgetColumns(form.serviceDependenciesColumns),
    charts: removeKeyFromEntities(form.charts),
  };

  parameters.usedAlarmProperties = convertAlarmWidgetParametersToActiveColumns(parameters);

  return parameters;
};

/**
 * Prepared alarms list widget for displaying
 *
 * @param {Object} widget
 * @returns {Object}
 */
export const prepareAlarmListWidget = (widget = {}) => setSeveralFields(widget, {
  'parameters.widgetColumns': (columns = []) => (
    columns.map(column => ({
      ...column,

      sortable: getWidgetColumnSortable(column, ALARM_UNSORTABLE_FIELDS),
      text: getWidgetColumnLabel(column, ALARM_FIELDS_TO_LABELS_KEYS),
    }))
  ),

  'parameters.widgetGroupColumns': (columns = DEFAULT_ALARMS_WIDGET_GROUP_COLUMNS) => (
    columns.map(column => ({
      ...column,

      sortable: getWidgetColumnSortable(column, ALARM_UNSORTABLE_FIELDS),
      text: getWidgetColumnLabel(column, ALARM_FIELDS_TO_LABELS_KEYS),
    }))
  ),

  'parameters.serviceDependenciesColumns': (columns = DEFAULT_SERVICE_DEPENDENCIES_COLUMNS) => (
    columns.map(column => ({
      ...column,

      sortable: false,
      text: getWidgetColumnLabel(column, ENTITY_FIELDS_TO_LABELS_KEYS),
    }))
  ),

  'parameters.widgetExportColumns': (columns = []) => (
    columns.map(column => ({
      ...column,

      text: getWidgetColumnLabel(column, ALARM_FIELDS_TO_LABELS_KEYS),
    }))
  ),
});

/**
 * Get filter for alarms list widget column value
 *
 * @param {string} value
 * @returns {Function | undefined}
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
 * Get component getter for alarms list widget column
 *
 * @param {string} value
 * @param {boolean} [onlyIcon]
 * @param {number} [inlineLinksCount]
 * @param {boolean} [showRootCauseByStateClick]
 * @returns {Function}
 */
export const getAlarmsListWidgetColumnComponentGetter = (
  { value, onlyIcon, inlineLinksCount },
  { showRootCauseByStateClick } = {},
) => {
  switch (value) {
    case ALARM_FIELDS.state:
      return (context) => {
        const component = {
          bind: {
            is: 'alarm-column-value-state',
            alarm: context.alarm,
            small: context.small,
          },
        };

        if (showRootCauseByStateClick) {
          component.bind.class = 'cursor-pointer';
          component.on = {
            click: () => context.$emit('click:state', context.alarm.entity),
          };
        }

        return component;
      };

    case ALARM_FIELDS.status:
      return context => ({
        bind: {
          is: 'alarm-column-value-status',
          alarm: context.alarm,
          small: context.small,
        },
      });

    case ALARM_FIELDS.impactState:
      return context => ({
        bind: {
          is: 'color-indicator-wrapper',
          type: COLOR_INDICATOR_TYPES.impactState,
          entity: context.alarm.entity,
          alarm: context.alarm,
        },
      });

    case ALARM_FIELDS.links:
      return context => ({
        bind: {
          onlyIcon,

          is: 'c-alarm-links-chips',
          alarm: context.alarm,
          small: context.small,
          inlineCount: inlineLinksCount,
        },
        on: {
          activate: context.$listeners.activate,
        },
      });

    case ALARM_FIELDS.extraDetails:
      return context => ({
        bind: {
          is: 'alarm-column-value-extra-details',
          alarm: context.alarm,
          small: context.small,
        },
      });

    case ALARM_FIELDS.tags:
      return context => ({
        bind: {
          is: 'c-alarm-tags-chips',
          alarm: context.alarm,
          selectedTag: context.selectedTag,
          small: context.small,
        },
        on: {
          select: context.$listeners['select:tag'],
        },
      });
  }

  if (value.startsWith('links.')) {
    const category = value.replace('links.', '');

    return context => ({
      bind: {
        category,
        onlyIcon,

        is: 'c-alarm-links-chips',
        alarm: context.alarm,
        small: context.small,
        inlineCount: inlineLinksCount,
      },
    });
  }

  return context => ({
    bind: {
      is: 'c-ellipsis',
      class: 'alarm-column-cell__text',
      title: context.value,
      text: context.value,
    },
  });
};
