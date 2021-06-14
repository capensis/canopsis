import sha1 from 'sha1';
import { get, omit, cloneDeep, isObject } from 'lodash';

import i18n from '@/i18n';
import { PAGINATION_LIMIT, DEFAULT_WEATHER_LIMIT, COLORS, DEFAULT_CATEGORIES_LIMIT } from '@/config';
import {
  DEFAULT_ALARMS_WIDGET_COLUMNS,
  DEFAULT_ALARMS_WIDGET_GROUP_COLUMNS,
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
  GRID_SIZES, AVAILABLE_COUNTERS,
  DEFAULT_COUNTER_BLOCK_TEMPLATE,
  TIME_UNITS,
  WIDGET_GRID_SIZES_KEYS,
  WIDGET_GRID_COLUMNS_COUNT,
  WORKFLOW_TYPES,
  EXPORT_CSV_SEPARATORS,
} from '@/constants';

import uuid from './uuid';
import uid from './uid';

/**
 * @typedef {Object} Interval
 * @property {number} interval
 * @property {string} unit
 */

/**
 * Generate id for view tab
 *
 * @returns {string}
 */
export const generateViewTabId = () => uuid('view-tab');

/**
 * Generate id for widget by type
 *
 * @param {string} type
 * @returns {string}
 */
export const generateWidgetId = type => uuid(`widget_${type}`);

/**
 * Generate id for action
 *
 * @returns {string}
 */
export const generateActionId = () => uuid('action');

/**
 * Generate id for watcher
 *
 * @returns {string}
 */
export const generateWatcherId = () => uuid('watcher');

/**
 * Generate id for entity
 *
 * @returns {string}
 */
export const generateEntityId = () => uuid('entity');

/**
 * Generate widget by type
 *
 * @param {string} type
 * @returns {Object}
 */
export function generateWidgetByType(type) {
  const widget = {
    type,
    _id: generateWidgetId(type),
    title: '',
    parameters: {},
    gridParameters: Object.values(WIDGET_GRID_SIZES_KEYS).reduce((acc, size) => {
      acc[size] = {
        x: 0,
        y: 0,
        h: 0,
        w: WIDGET_GRID_COLUMNS_COUNT,
        autoHeight: true,
      };

      return acc;
    }, {}),
  };

  const alarmsListDefaultParameters = {
    itemsPerPage: PAGINATION_LIMIT,
    infoPopups: [],
    moreInfoTemplate: '',
    isAckNoteRequired: false,
    isSnoozeNoteRequired: false,
    isMultiAckEnabled: false,
    isHtmlEnabledOnTimeLine: false,
    fastAckOutput: {
      enabled: false,
      value: 'auto ack',
    },
    widgetColumns: DEFAULT_ALARMS_WIDGET_COLUMNS.map(({ labelKey, value }) => ({
      label: i18n.t(labelKey),
      value,
    })),
    widgetGroupColumns: DEFAULT_ALARMS_WIDGET_GROUP_COLUMNS.map(({ labelKey, value }) => ({
      label: i18n.t(labelKey),
      value,
    })),
    linksCategoriesAsList: {
      enabled: false,
      limit: DEFAULT_CATEGORIES_LIMIT,
    },
  };

  let specialParameters = {};

  switch (type) {
    case WIDGET_TYPES.alarmList:
      specialParameters = {
        ...alarmsListDefaultParameters,

        viewFilters: [],
        mainFilter: null,
        mainFilterUpdatedAt: 0,
        infoPopups: [],
        liveReporting: {},
        periodicRefresh: {
          enabled: false,
          interval: 60,
          unit: 's',
        },
        sort: {
          order: SORT_ORDERS.asc,
        },
        alarmsStateFilter: {
          opened: true,
        },
        expandGridRangeSize: [GRID_SIZES.min, GRID_SIZES.max],
        exportCsvSeparator: EXPORT_CSV_SEPARATORS.comma,
      };
      break;

    case WIDGET_TYPES.context:
      specialParameters = {
        itemsPerPage: PAGINATION_LIMIT,
        viewFilters: [],
        mainFilter: null,
        mainFilterUpdatedAt: 0,
        widgetColumns: [
          {
            label: i18n.t('tables.contextList.name'),
            value: 'name',
          },
          {
            label: i18n.t('tables.contextList.type'),
            value: 'type',
          },
        ],
        selectedTypes: [],
        sort: {
          order: SORT_ORDERS.asc,
        },
        exportCsvSeparator: EXPORT_CSV_SEPARATORS.comma,
      };
      break;

    case WIDGET_TYPES.weather:
      specialParameters = {
        mfilter: {},
        sort: {
          order: SORT_ORDERS.asc,
        },
        blockTemplate: `<p><strong><span style="font-size: 18px;">{{entity.display_name}}</span></strong></p>
<hr id="null">
<p>{{ entity.output }}</p>
<p> Dernière mise à jour : {{ timestamp entity.last_update_date }}</p>`,

        modalTemplate: '{{ entities name="entity.entity_id" }}',
        entityTemplate: `<ul>
    <li><strong>Libellé</strong> : {{entity.name}}</li>
</ul>`,

        columnSM: 6,
        columnMD: 4,
        columnLG: 3,
        limit: DEFAULT_WEATHER_LIMIT,
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
 * @returns {Object}
 */
export function generateViewTab() {
  return {
    _id: generateViewTabId(),
    title: '',
    widgets: [],
  };
}

/**
 * Generate view
 *
 * @returns {Object}
 */
export function generateView() {
  const defaultTab = { ...generateViewTab(), title: 'Default' };

  return {
    title: '',
    name: '',
    description: '',
    group_id: null,
    tabs: [defaultTab],
    tags: [],
    enabled: true,
  };
}

/**
 * Generate user preference by widget and user objects
 *
 * @param {Object} widget
 * @param {Object} user
 * @returns {Object}
 */
export function generateUserPreferenceByWidgetAndUser(widget, user) {
  return {
    _id: `${widget._id}_${user._id}`,
    widget_preferences: {},
    crecord_name: user._id,
    widget_id: widget._id,
    widgetXtype: widget.type,
    crecord_type: 'userpreferences',
  };
}

/**
 * Generate user
 *
 * @returns {Object}
 */
export function generateUser() {
  return {
    crecord_write_time: null,
    crecord_type: 'user',
    crecord_creation_time: null,
    crecord_name: null,
    user_contact: null,
    rights: null,
    user_role: null,
    user_groups: null,
    authkey: null,
    role: null,
    external: false,
    defaultview: null,
    id: null,
    _id: null,
    firstname: '',
    lastname: '',
    password: '',
    mail: '',
    enable: true,
    ui_language: 'fr',
  };
}

/**
 * Generate role
 *
 * @returns {Object}
 */
export function generateRole() {
  return {
    crecord_write_time: null,
    enable: true,
    crecord_type: 'role',
    crecord_creation_time: null,
    crecord_name: null,
    rights: null,
    id: null,
  };
}

/**
 * Generate right
 *
 * @returns {Object}
 */
export function generateRight() {
  return {
    crecord_creation_time: null,
    crecord_name: null,
    crecord_type: 'action',
    crecord_write_time: null,
    desc: '',
    enable: true,
    id: null,
    type: '',
    _id: '',
  };
}

/**
 * Generate role right by checksum
 *
 * @param {number} checksum
 * @returns {Object}
 */
export function generateRoleRightByChecksum(checksum) {
  return {
    checksum,
    crecord_type: 'right',
  };
}

/**
 * Generate copy of view tab
 *
 * @param {Object} tab
 * @returns {Object}
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
 * Generate copy of view
 *
 * @param {Object} view
 * @returns {Object}
 */
export function generateCopyOfView(view) {
  return {
    ...generateView(),
    ...omit(view, ['_id', 'tabs']),

    tabs: view.tabs.map(tab => generateCopyOfViewTab(tab)),
  };
}

/**
 * Get mappings for widgets ids from old tab to new tab
 *
 * @param {Object} oldTab
 * @param {Object} newTab
 * @returns {Array.<{ oldId: number, newId: number }>}
 */
export function getViewsTabsWidgetsIdsMappings(oldTab, newTab) {
  return oldTab.widgets.map((acc, widget, widgetIndex) => ({
    oldId: widget._id,
    newId: get(newTab, `widgets.${widgetIndex}._id`, null),
  }));
}

/**
 * Get mappings for widgets from old view to new view
 *
 * @param {Object} oldView
 * @param {Object} newView
 * @returns {Array.<{ oldId: number, newId: number }>}
 */
export function getViewsWidgetsIdsMappings(oldView, newView) {
  return oldView.tabs.reduce((acc, tab, index) =>
    acc.concat(getViewsTabsWidgetsIdsMappings(tab, newView.tabs[index])), []);
}

export function prepareUserByData(data, user = generateUser()) {
  const result = {
    ...omit(user, ['rights']),
    ...omit(data, ['password']),
  };

  if (data.password && data.password !== '') {
    result.shadowpasswd = sha1(data.password);
  }

  if (!result._id && !data._id) {
    result._id = data.crecord_name;
  }

  return result;
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
 * Function return new name if name is not uniq
 * @param {Object} entity
 * @param {Array} entities
 * @returns {string}
 */
export function getDuplicateEntityName(entity, entities) {
  const suffixRegexp = '(\\s\\(\\d+\\))?$';
  const clearName = entity.name.replace(new RegExp(suffixRegexp), '');

  const nameRegexp = new RegExp(`^${clearName}${suffixRegexp}`);

  const duplicateEntityCount = entities.reduce((count, { name }) => {
    const isDuplicate = nameRegexp.test(name);

    return isDuplicate ? count + 1 : count;
  }, 0);

  return duplicateEntityCount !== 0 ? `${clearName} (${duplicateEntityCount})` : entity.name;
}

/**
 * Create default playlist entity
 *
 * @returns {Object}
 */
export function getDefaultPlaylist() {
  return {
    name: '',
    fullscreen: true,
    enabled: true,
    interval: {
      interval: 10,
      unit: TIME_UNITS.second,
    },
    tabs_list: [],
  };
}

/**
 * Add uniq key field in each entity.
 *
 * @param {Array} entities
 * @return {Array}
 */
export const addKeyInEntity = (entities = []) => entities.map(entity => ({
  ...entity,
  key: uid(),
}));

/**
 * Remove key field from each entity.
 *
 * @param {Array} entities
 * @return {Array}
 */
export const removeKeyFromEntity = (entities = []) => entities.map(entity => omit(entity, ['key']));

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
