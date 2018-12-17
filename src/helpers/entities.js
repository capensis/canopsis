import moment from 'moment';
import i18n from '@/i18n';
import { PAGINATION_LIMIT } from '@/config';
import { WIDGET_TYPES, STATS_CALENDAR_COLORS, STATS_DURATION_UNITS, SERVICE_WEATHER_WIDGET_MODAL_TYPES } from '@/constants';

import uuid from './uuid';

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

  let specialParameters = {};

  switch (type) {
    case WIDGET_TYPES.alarmList:
      specialParameters = {
        itemsPerPage: PAGINATION_LIMIT,
        moreInfoTemplate: '',
        alarmsStateFilter: {},
        widgetColumns: [
          {
            label: i18n.t('tables.alarmGeneral.connector'),
            value: 'alarm.connector',
          },
          {
            label: i18n.t('tables.alarmGeneral.connectorName'),
            value: 'alarm.connector_name',
          },
          {
            label: i18n.t('tables.alarmGeneral.component'),
            value: 'alarm.component',
          },
          {
            label: i18n.t('tables.alarmGeneral.resource'),
            value: 'alarm.resource',
          },
          {
            label: i18n.t('tables.alarmGeneral.output'),
            value: 'alarm.output',
          },
          {
            label: i18n.t('tables.alarmGeneral.extraDetails'),
            value: 'extra_details',
          },
          {
            label: i18n.t('tables.alarmGeneral.state'),
            value: 'alarm.state.val',
          },
          {
            label: i18n.t('tables.alarmGeneral.status'),
            value: 'alarm.status.val',
          },
        ],
        viewFilters: [],
        infoPopups: [],
        periodicRefresh: {
          enabled: false,
          interval: 60,
        },
        sort: {
          order: 'ASC',
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
        duration: `1${STATS_DURATION_UNITS.day}`,
        tstop: moment()
          .startOf('hour')
          .unix(),
        periods: 2,
        stats: {},
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
        alarmsList: {
          itemsPerPage: PAGINATION_LIMIT,
          infoPopups: [],
          moreInfoTemplate: '',
          widgetColumns: [
            {
              label: i18n.t('tables.alarmGeneral.connector'),
              value: 'alarm.connector',
            },
            {
              label: i18n.t('tables.alarmGeneral.connectorName'),
              value: 'alarm.connector_name',
            },
            {
              label: i18n.t('tables.alarmGeneral.component'),
              value: 'alarm.component',
            },
            {
              label: i18n.t('tables.alarmGeneral.resource'),
              value: 'alarm.resource',
            },
            {
              label: i18n.t('tables.alarmGeneral.output'),
              value: 'alarm.output',
            },
            {
              label: i18n.t('tables.alarmGeneral.extraDetails'),
              value: 'extra_details',
            },
            {
              label: i18n.t('tables.alarmGeneral.state'),
              value: 'alarm.state.val',
            },
            {
              label: i18n.t('tables.alarmGeneral.status'),
              value: 'alarm.status.val',
            },
          ],
        },
      };
      break;

    case WIDGET_TYPES.statsNumber:
      specialParameters = {
        duration: `1${STATS_DURATION_UNITS.day}`,
        tstop: moment()
          .startOf('hour')
          .unix(),
        mfilter: {},
        stat: {},
        yesNoMode: false,
        criticityLevels: {
          minor: 20,
          major: 30,
          critical: 40,
        },
        statColors: {
          ok: '#66BB6A',
          minor: '#FFEE58',
          major: '#FFA726',
          critical: '#FF7043',
        },
      };
      break;
  }

  widget.parameters = { ...widget.parameters, ...specialParameters };

  return widget;
}

export function generateViewRow() {
  return {
    _id: uuid('view-row'),
    title: '',
    widgets: [],
  };
}

export function generateViewTab() {
  return {
    _id: uuid('view-tab'),
    title: '',
    rows: [],
  };
}

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

export function generateRoleRightByChecksum(checksum) {
  return {
    checksum,
    crecord_type: 'right',
  };
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
};
