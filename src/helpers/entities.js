import { PAGINATION_LIMIT } from '@/config';
import { WIDGET_TYPES } from '@/constants';

import uuid from './uuid';

export function generateWidgetByType(type) {
  const widget = {
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

    default:
  }

  widget.parameters = { ...widget.parameters, ...specialParameters };

  return widget;
}

export default {
  generateWidgetByType,
};
