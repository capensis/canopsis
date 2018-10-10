import moment from 'moment';

import { PAGINATION_LIMIT } from '@/config';
import { WIDGET_TYPES } from '@/constants';

import uuid from './uuid';

export function generateWidgetByType(type) {
  const widget = {
    type,
    _id: uuid(`widget_${type}`),
    title: '',
    minColumns: 6,
    maxColumns: 12,
    parameters: {},
  };

  let specialParameters = {};

  switch (type) {
    case WIDGET_TYPES.alarmList:
      specialParameters = {
        itemsPerPage: PAGINATION_LIMIT,
        moreInfoTemplate: '',
        alarmsStateFilter: {},
        widgetColumns: [],
        viewFilters: [],
        infoPopups: [],
        periodicRefresh: {
          enabled: false,
        },
        sort: {
          order: 'ASC',
        },
      };
      break;

    case WIDGET_TYPES.context:
      specialParameters = {
        itemsPerPage: PAGINATION_LIMIT,
        widgetColumns: [],
        selectedTypes: [],
        sort: {
          order: 'ASC',
        },
      };
      break;

    case WIDGET_TYPES.weather:
      specialParameters = {
        blockTemplate: '',
        modalTemplate: '',
        entityTemplate: '',
        columnSM: 0,
        columnMD: 0,
        columnLG: 0,
        columnHG: 0,
      };
      break;
    case WIDGET_TYPES.statsHistogram:
      specialParameters = {
        mfilter: {},
        duration: '1m',
        tstop: moment().unix(),
        groups: [],
        stats: {},
      };
      break;
    case WIDGET_TYPES.statsTable:
      specialParameters = {
        duration: '1m',
        tstop: moment().unix(),
        stats: {},
        mfilter: {},
      };
      break;

    case WIDGET_TYPES.statsNumber:
      specialParameters = {
        duration: '1d',
        tstop: moment().unix(),
        mfilter: {},
        stat: {
          stat: 'alarms_created',
          trend: true,
        },
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

export function generateRow() {
  return {
    _id: uuid('row'),
    title: '',
    widgets: [],
  };
}

export function generateView() {
  return {
    title: '',
    name: '',
    description: '',
    group_id: null,
    rows: [],
    tags: [],
    enabled: true,
  };
}

export default {
  generateWidgetByType,
  generateRow,
  generateView,
};
