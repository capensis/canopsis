import sha1 from 'sha1';
import { get, omit, cloneDeep } from 'lodash';

import i18n from '@/i18n';
import { PAGINATION_LIMIT, DEFAULT_WEATHER_LIMIT } from '@/config';
import {
  WIDGET_TYPES,
  STATS_CALENDAR_COLORS,
  STATS_TYPES,
  STATS_DURATION_UNITS,
  STATS_QUICK_RANGES,
  STATS_DISPLAY_MODE,
  STATS_DISPLAY_MODE_PARAMETERS,
  SERVICE_WEATHER_WIDGET_MODAL_TYPES,
  SORT_ORDERS,
  ACTION_TYPES,
  DURATION_UNITS,
  ENTITIES_STATES,
} from '@/constants';

import uuid from './uuid';
import { pbehaviorToForm } from './forms/pbehavior';

/**
 * Generate widget by type
 *
 * @param {string} type
 * @returns {Object}
 */
export function generateWidgetByType(type) {
  const widget = {
    type,
    _id: uuid(`widget_${type}`),
    title: '',
    parameters: {},
    size: {
      sm: 3,
      md: 3,
      lg: 3,
    },
  };

  const alarmsListDefaultParameters = {
    itemsPerPage: PAGINATION_LIMIT,
    infoPopups: [],
    moreInfoTemplate: '',
    isAckNoteRequired: false,
    isMultiAckEnabled: false,
    isHtmlEnabledOnTimeLine: false,
    fastAckOutput: {
      enabled: false,
      value: 'auto ack',
    },
    widgetColumns: [
      {
        label: i18n.t('tables.alarmGeneral.connector'),
        value: 'v.connector',
      },
      {
        label: i18n.t('tables.alarmGeneral.connectorName'),
        value: 'v.connector_name',
      },
      {
        label: i18n.t('tables.alarmGeneral.component'),
        value: 'v.component',
      },
      {
        label: i18n.t('tables.alarmGeneral.resource'),
        value: 'v.resource',
      },
      {
        label: i18n.t('tables.alarmGeneral.output'),
        value: 'v.output',
      },
      {
        label: i18n.t('tables.alarmGeneral.extraDetails'),
        value: 'extra_details',
      },
      {
        label: i18n.t('tables.alarmGeneral.state'),
        value: 'v.state.val',
      },
      {
        label: i18n.t('tables.alarmGeneral.status'),
        value: 'v.status.val',
      },
    ],
  };

  let specialParameters = {};

  switch (type) {
    case WIDGET_TYPES.alarmList:
      specialParameters = {
        ...alarmsListDefaultParameters,

        viewFilters: [],
        mainFilter: null,
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
      };
      break;

    case WIDGET_TYPES.context:
      specialParameters = {
        itemsPerPage: PAGINATION_LIMIT,
        viewFilters: [],
        mainFilter: null,
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
      };
      break;

    case WIDGET_TYPES.weather:
      specialParameters = {
        mfilter: {},
        sort: {
          order: SORT_ORDERS.asc,
        },
        blockTemplate: '',
        modalTemplate: '',
        entityTemplate: '',
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
        heightFactor: 1,
        modalType: SERVICE_WEATHER_WIDGET_MODAL_TYPES.moreInfo,
        alarmsList: alarmsListDefaultParameters,
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
  }

  widget.parameters = { ...widget.parameters, ...specialParameters };

  return widget;
}

/**
 * Generate view row
 *
 * @returns {Object}
 */
export function generateViewRow() {
  return {
    _id: uuid('view-row'),
    title: '',
    widgets: [],
  };
}

/**
 * Generate view tab
 *
 * @returns {Object}
 */
export function generateViewTab() {
  return {
    _id: uuid('view-tab'),
    title: '',
    rows: [],
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

    rows: tab.rows.map(row => ({
      ...generateViewRow(),

      title: row.title,
      widgets: row.widgets.map(widget => ({
        ...generateWidgetByType(widget.type),
        ...omit(widget, ['_id']),
      })),
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

    tabs: view.tabs.map(tab => ({
      ...generateCopyOfViewTab(tab),

      ...omit(tab, ['_id', 'rows']),
    })),
  };
}

/**
 * Generate an 'action' entity
 * @returns {Object}
 */
export function generateAction() {
  const defaultHook = {
    event_patterns: [],
    alarm_patterns: [],
    entity_patterns: [],
    triggers: [],
  };

  // Get basic action parameters
  const generalParameters = {
    _id: uuid('action'),
    type: ACTION_TYPES.snooze,
    hook: defaultHook,
  };

  // Default 'snooze' action parameters
  const snoozeParameters = {
    message: '',
    duration: {
      duration: 1,
      durationType: DURATION_UNITS.minute.value,
    },
  };

  // Default 'pbehavior' action parameters
  const pbehaviorParameters = {
    general: { ...pbehaviorToForm() },
    comments: [],
    exdate: [],
  };

  // Default 'changestate' action parameters
  const changeStateParameters = {
    state: ENTITIES_STATES.minor,
    output: '',
  };

  return {
    generalParameters,
    snoozeParameters,
    pbehaviorParameters,
    changeStateParameters,
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
  return oldTab.rows.reduce((acc, row, rowIndex) => {
    const widgetsIds = row.widgets.map((widget, widgetIndex) => ({
      oldId: widget._id,
      newId: get(newTab, `rows.${rowIndex}.widgets.${widgetIndex}._id`, null),
    }));

    return acc.concat(widgetsIds);
  }, []);
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
  const result = { ...user, ...omit(data, ['password']) };

  if (data.password && data.password !== '') {
    result.shadowpasswd = sha1(data.password);
  }

  return result;
}
