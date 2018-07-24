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
    widget.alarms_state_filter = null;
    widget.hide_resources = false;
    widget.widget_columns = [];
    widget.columns = [
      'connector_name',
      'component',
      'resource',
      'state',
      'status',
      'last_update_date',
      'extra_details',
    ];
  }

  return widget;
}

export default {
  generateWidgetByType,
};
