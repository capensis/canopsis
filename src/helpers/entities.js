import { WIDGET_TYPES } from '@/constants';

import uuid from './uuid';

export function generateWidgetByType(type) {
  const id = uuid(`widget_${type}`);
  const widget = {
    id,
    widgetId: id,
    title: null,
    preference_id: uuid(),
    xtype: type,
    tagName: null,
    mixins: [],
    default_sort_column: {
      direction: 'ASC',
    },
    columns: [],
    popup: [],
  };

  if (type === WIDGET_TYPES.alarmList) {
    widget.alarms_state_filter = {};
    widget.widget_columns = [];
  }

  return widget;
}

export default {
  generateWidgetByType,
};
