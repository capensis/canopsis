import moment from 'moment';
import { get, omit, cloneDeep } from 'lodash';

import i18n from '@/i18n';
import { PAGINATION_LIMIT } from '@/config';
import {
  WIDGET_TYPES,
  STATS_CALENDAR_COLORS,
  STATS_TYPES,
  STATS_DURATION_UNITS,
  STATS_DISPLAY_MODE,
  STATS_DISPLAY_MODE_PARAMETERS,
  SERVICE_WEATHER_WIDGET_MODAL_TYPES,
  SORT_ORDERS,
} from '@/constants';

import uuid from './uuid';

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
        infoPopups: [],
        periodicRefresh: {
          enabled: false,
          interval: 60,
        },
        sort: {
          order: 'ASC',
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
          order: 'ASC',
        },
      };
      break;

    case WIDGET_TYPES.weather:
      specialParameters = {
        mfilter: {},
        blockTemplate: '',
        modalTemplate: '',
        entityTemplate: '',
        columnSM: 6,
        columnMD: 4,
        columnLG: 3,
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
        duration: `1${STATS_DURATION_UNITS.day}`,
        tstop: moment()
          .startOf('hour')
          .unix(),
        groups: [],
        stats: {},
        statsColors: {},
      };
      break;
    case WIDGET_TYPES.statsCurves:
      specialParameters = {
        mfilter: {},
        dateInterval: {
          periodValue: 1,
          periodUnit: STATS_DURATION_UNITS.day,
          tstart: 'now/d',
          tstop: 'now/d',
        },
        stats: {},
        statsColors: {},
      };
      break;
    case WIDGET_TYPES.statsTable:
      specialParameters = {
        duration: `1${STATS_DURATION_UNITS.day}`,
        tstop: moment()
          .startOf('hour')
          .unix(),
        stats: {},
        mfilter: {},
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
        limit: 10,
        sortOrder: SORT_ORDERS.desc,
        displayMode: {
          mode: STATS_DISPLAY_MODE.criticity,
          parameters: cloneDeep(STATS_DISPLAY_MODE_PARAMETERS),
        },
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
    _id: `${widget._id}_${user.crecord_name}`,
    widget_preferences: {},
    crecord_name: user.crecord_name,
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


export default {
  generateWidgetByType,
  generateViewRow,
  generateView,
  generateUserPreferenceByWidgetAndUser,
  generateUser,
  generateRole,
  generateRight,
  generateRoleRightByChecksum,
  generateCopyOfViewTab,
  generateCopyOfView,

  getViewsTabsWidgetsIdsMappings,
  getViewsWidgetsIdsMappings,
};
