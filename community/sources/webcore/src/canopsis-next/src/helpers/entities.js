import { get, omit, cloneDeep, isObject } from 'lodash';

import i18n from '@/i18n';
import { PAGINATION_LIMIT, DEFAULT_WEATHER_LIMIT, COLORS } from '@/config';
import {
  DEFAULT_SERVICE_DEPENDENCIES_COLUMNS,
  WIDGET_TYPES,
  STATS_CALENDAR_COLORS,
  STATS_TYPES,
  STATS_DURATION_UNITS,
  STATS_QUICK_RANGES,
  STATS_DISPLAY_MODE,
  STATS_DISPLAY_MODE_PARAMETERS,
  SERVICE_WEATHER_WIDGET_MODAL_TYPES,
  SORT_ORDERS,
  ENTITIES_STATES,
  ENTITIES_STATUSES,
  AVAILABLE_COUNTERS,
  DEFAULT_COUNTER_BLOCK_TEMPLATE,
  TIME_UNITS,
  WORKFLOW_TYPES,
  COLOR_INDICATOR_TYPES,
  EXPORT_CSV_SEPARATORS,
  EXPORT_CSV_DATETIME_FORMATS,
  DEFAULT_ALARMS_WIDGET_COLUMNS,
} from '@/constants';

import { widgetToForm } from '@/helpers/forms/widgets/common';
import { alarmListWidgetDefaultParametersToForm } from '@/helpers/forms/widgets/alarm';

import uuid from './uuid';
import uid from './uid';

/**
 * @typedef {Object} Interval
 * @property {number} interval
 * @property {string} unit
 */

/**
 * Convert default columns from constants to columns with prepared by i18n label
 *
 * @param {{ labelKey: string, value: string }[]} [columns = []]
 * @returns {{ label: string, value: string }[]}
 */
export function defaultColumnsToColumns(columns = []) {
  return columns.map(({ labelKey, value }) => ({
    label: i18n.t(labelKey),
    value,
  }));
}

/**
 * Generate id for view tab
 *
 * @returns {string}
 */
export const generateViewTabId = () => uuid('view-tab');

/**
 * Generate id for widget tab
 *
 * @returns {string}
 */
export const generateWidgetId = type => uuid(`widget_${type}`);

/**
 * Generate widget by type
 *
 * @param {string} type
 * @returns {Object}
 */
export function generateWidgetByType(type) {
  const widget = widgetToForm({ type });

  const alarmsListDefaultParameters = alarmListWidgetDefaultParametersToForm();

  let specialParameters = {};

  switch (type) {
    case WIDGET_TYPES.context:
      specialParameters = {
        itemsPerPage: PAGINATION_LIMIT,
        viewFilters: [],
        mainFilter: null,
        mainFilterUpdatedAt: 0,
        widgetColumns: [
          {
            label: i18n.t('common.name'),
            value: 'name',
          },
          {
            label: i18n.t('common.type'),
            value: 'type',
          },
        ],
        serviceDependenciesColumns: [
          {
            label: i18n.t('common.name'),
            value: 'name',
          },
          {
            label: i18n.t('common.type'),
            value: 'type',
          },
        ],
        selectedTypes: [],
        sort: {
          order: SORT_ORDERS.asc,
        },
        exportCsvSeparator: EXPORT_CSV_SEPARATORS.comma,
        exportCsvDatetimeFormat: EXPORT_CSV_DATETIME_FORMATS.datetimeSeconds.value,
        widgetExportColumns: defaultColumnsToColumns(DEFAULT_ALARMS_WIDGET_COLUMNS),
      };
      break;

    case WIDGET_TYPES.serviceWeather:
      specialParameters = {
        viewFilters: [],
        mainFilter: null,
        mainFilterUpdatedAt: 0,
        sort: {
          order: SORT_ORDERS.asc,
        },
        blockTemplate: `<p><strong><span style="font-size: 18px;">{{entity.name}}</span></strong></p>
<hr id="null">
<p>{{ entity.output }}</p>
<p> Dernière mise à jour : {{ timestamp entity.last_update_date }}</p>`,

        modalTemplate: '{{ entities name="entity._id" }}',
        entityTemplate: `<ul>
    <li><strong>Libellé</strong> : {{entity.name}}</li>
</ul>`,

        columnSM: 6,
        columnMD: 4,
        columnLG: 3,
        limit: DEFAULT_WEATHER_LIMIT,
        colorIndicator: COLOR_INDICATOR_TYPES.state,
        serviceDependenciesColumns: defaultColumnsToColumns(DEFAULT_SERVICE_DEPENDENCIES_COLUMNS),
        margin: {
          top: 1,
          right: 1,
          bottom: 1,
          left: 1,
        },
        isCountersEnabled: false,
        heightFactor: 6,
        modalType: SERVICE_WEATHER_WIDGET_MODAL_TYPES.both,
        alarmsList: alarmsListDefaultParameters,
        modalItemsPerPage: PAGINATION_LIMIT,
      };
      break;

    case WIDGET_TYPES.statsHistogram:
      specialParameters = {
        mfilter: {},
        dateInterval: {
          periodValue: 1,
          periodUnit: STATS_DURATION_UNITS.day,
          tstart: STATS_QUICK_RANGES.thisMonthSoFar.start,
          tstop: STATS_QUICK_RANGES.thisMonthSoFar.stop,
        },
        stats: {},
        statsColors: {},
        annotationLine: {},
      };
      break;
    case WIDGET_TYPES.statsCurves:
      specialParameters = {
        mfilter: {},
        dateInterval: {
          periodValue: 1,
          periodUnit: STATS_DURATION_UNITS.day,
          tstart: STATS_QUICK_RANGES.thisMonthSoFar.start,
          tstop: STATS_QUICK_RANGES.thisMonthSoFar.stop,
        },
        stats: {},
        statsColors: {},
        statsPointsStyles: {},
        annotationLine: {},
      };
      break;
    case WIDGET_TYPES.statsTable:
      specialParameters = {
        dateInterval: {
          periodValue: 1,
          periodUnit: STATS_DURATION_UNITS.day,
          tstart: STATS_QUICK_RANGES.thisMonthSoFar.start,
          tstop: STATS_QUICK_RANGES.thisMonthSoFar.stop,
        },
        mfilter: {},
        stats: {},
        sort: {},
      };
      break;
    case WIDGET_TYPES.statsCalendar:
      specialParameters = {
        filters: [],
        alarmsStateFilter: {},
        considerPbehaviors: false,
        criticityLevelsColors: { ...STATS_CALENDAR_COLORS.alarm },
        criticityLevels: {
          minor: 20,
          major: 30,
          critical: 40,
        },
        alarmsList: alarmsListDefaultParameters,
      };
      break;

    case WIDGET_TYPES.statsNumber:
      specialParameters = {
        dateInterval: {
          periodValue: 1,
          periodUnit: STATS_DURATION_UNITS.day,
          tstart: STATS_QUICK_RANGES.thisMonthSoFar.start,
          tstop: STATS_QUICK_RANGES.thisMonthSoFar.stop,
        },
        mfilter: {},
        stat: {
          parameters: {
            recursive: true,
          },
          stat: STATS_TYPES.alarmsCreated,
          title: 'Alarmes créées',
          trend: false,
        },
        limit: 10,
        sortOrder: SORT_ORDERS.desc,
        displayMode: {
          mode: STATS_DISPLAY_MODE.criticity,
          parameters: cloneDeep(STATS_DISPLAY_MODE_PARAMETERS),
        },
      };
      break;

    case WIDGET_TYPES.statsPareto:
      specialParameters = {
        dateInterval: {
          periodValue: 1,
          periodUnit: STATS_DURATION_UNITS.day,
          tstart: 'now/d',
          tstop: 'now/d',
        },
        mfilter: {},
        stat: {
          parameters: {
            recursive: true,
          },
          stat: STATS_TYPES.alarmsCreated,
          title: 'Alarmes créées',
          trend: false,
        },
        statsColors: {},
      };
      break;

    case WIDGET_TYPES.text:
      specialParameters = {
        dateInterval: {
          periodValue: 1,
          periodUnit: STATS_DURATION_UNITS.day,
          tstart: STATS_QUICK_RANGES.thisMonthSoFar.start,
          tstop: STATS_QUICK_RANGES.thisMonthSoFar.stop,
        },
        mfilter: {},
        stats: {},
        template: '',
      };
      break;
    case WIDGET_TYPES.counter:
      specialParameters = {
        viewFilters: [],
        alarmsStateFilter: {
          opened: true,
        },
        blockTemplate: DEFAULT_COUNTER_BLOCK_TEMPLATE,
        columnSM: 6,
        columnMD: 4,
        columnLG: 3,
        margin: {
          top: 1,
          right: 1,
          bottom: 1,
          left: 1,
        },
        heightFactor: 6,
        isCorrelationEnabled: false,
        levels: {
          counter: AVAILABLE_COUNTERS.total,
          colors: {
            ok: COLORS.state.ok,
            minor: COLORS.state.minor,
            major: COLORS.state.major,
            critical: COLORS.state.critical,
          },
          values: {
            minor: 20,
            major: 30,
            critical: 40,
          },
        },
        alarmsList: alarmsListDefaultParameters,
      };
      break;
  }

  widget.parameters = { ...widget.parameters, ...specialParameters };

  return widget;
}

/**
 * Generate view tab
 *
 * @param {string} [title = '']
 * @returns {Object}
 */
export function generateViewTab(title = '') {
  return {
    title,

    _id: generateViewTabId(),
    widgets: [],
  };
}

/**
 * Generate copy of view tab
 *
 * @param {ViewTab} tab
 * @returns {ViewTab}
 */
export function generateCopyOfViewTab(tab) {
  return {
    ...generateViewTab(),
    ...omit(tab, ['_id', 'widgets']),

    widgets: tab.widgets.map(widget => ({
      ...generateWidgetByType(widget.type),
      ...omit(widget, ['_id']),
    })),
  };
}

/**
 * Get mappings for widgets ids from old tab to new tab
 *
 * @param {Object} oldTab
 * @param {Object} newTab
 * @returns {{ oldId: string, newId: string }[]}
 */
export function getViewsTabsWidgetsIdsMappings(oldTab, newTab) {
  return oldTab.widgets.map((widget, widgetIndex) => ({
    oldId: widget._id,
    newId: get(newTab, `widgets.${widgetIndex}._id`, null),
  }));
}

/**
 * Get mappings for widgets from old view to new view
 *
 * @param {View | ViewRequest} oldView
 * @param {View | ViewRequest} newView
 * @returns {{ oldId: string, newId: string }[]}
 */
export function getViewsWidgetsIdsMappings(oldView, newView) {
  return oldView.tabs.reduce((acc, tab, index) =>
    acc.concat(getViewsTabsWidgetsIdsMappings(tab, newView.tabs[index])), []);
}

/**
 * Checks if alarm is resolved
 * @param alarm - alarm entity
 * @returns {boolean}
 */
export function isResolvedAlarm(alarm) {
  return [ENTITIES_STATUSES.off, ENTITIES_STATUSES.cancelled].includes(alarm.v.status.val);
}

/**
 * Checks if alarm have critical state
 *
 * @param alarm - alarm entity
 * @returns {boolean}
 */
export function isWarningAlarmState(alarm) {
  return ENTITIES_STATES.ok !== alarm.v.state.val;
}

/**
 * Function return new title if title is not uniq
 *
 * @param {Object} [entity = {}]
 * @param {Array} [entities = []]
 * @returns {string}
 */
export function getDuplicateEntityTitle(entity = {}, entities = []) {
  const suffixRegexp = '(\\s\\(\\d+\\))?$';
  const clearName = entity.title.replace(new RegExp(suffixRegexp), '');

  const titleRegexp = new RegExp(`^${clearName}${suffixRegexp}`);

  const duplicateEntityCount = entities.reduce((count, { title }) => {
    const isDuplicate = titleRegexp.test(title);

    return isDuplicate ? count + 1 : count;
  }, 0);

  return duplicateEntityCount !== 0 ? `${clearName} (${duplicateEntityCount})` : entity.title;
}

/**
 * Add uniq key field in each entity.
 *
 * @param {Array} entities
 * @return {Array}
 */
export const addKeyInEntities = (entities = []) => entities.map(entity => ({
  ...entity,
  key: uid(),
}));

/**
 * Remove key field from each entity.
 *
 * @param {Array} entities
 * @return {Array}
 */
export const removeKeyFromEntities = (entities = []) => entities.map(entity => omit(entity, ['key']));

/**
 * Get id from entity
 *
 * @param {Object} entity
 * @param {string} idField
 * @return {string}
 */
export const getIdFromEntity = (entity, idField = '_id') =>
  (isObject(entity) ? entity[idField] : entity);

/**
 * Generate an remediation instruction step operation entity
 *
 * @typedef {Object} RemediationInstructionStepOperation
 * @property {string} name
 * @property {string} description
 * @property {Array} jobs
 * @property {DurationForm} time_to_complete
 * @property {string} [key]
 * @return {RemediationInstructionStepOperation}
 */
export const generateRemediationInstructionStepOperation = () => ({
  name: '',
  description: '',
  jobs: [],
  time_to_complete: {
    value: 0,
    unit: TIME_UNITS.minute,
  },
  key: uid(),
});

/**
 * Generate an remediation instruction step entity
 *
 * @typedef {Object} RemediationInstructionStep
 * @property {string} endpoint
 * @property {string} name
 * @property {boolean} stop_on_fail
 * @property {RemediationInstructionStepOperation[]} operations
 * @property {boolean} [saved]
 * @property {string} [key]
 * @return {RemediationInstructionStep}
 */
export const generateRemediationInstructionStep = () => ({
  endpoint: '',
  name: '',
  operations: [generateRemediationInstructionStepOperation()],
  stop_on_fail: WORKFLOW_TYPES.stop,
  saved: false,
  key: uid(),
});
