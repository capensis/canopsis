import { PAGINATION_LIMIT } from '@/config';
import { WIDGET_TYPES } from '@/constants';

export function convertAlarmWidgetToQuery(widget) {
  const query = {
    page: 1,
  };

  const { default_sort_column: defaultSortColumn } = widget;

  if (defaultSortColumn && defaultSortColumn.property) {
    query.sortBy = defaultSortColumn.property.startsWith('v.') ?
      defaultSortColumn.property : `v.${defaultSortColumn.property}`;

    query.descending = defaultSortColumn.direction === 'DESC';
  } else {
    query.sortBy = null;
    query.descending = null;
  }

  return query;
}

export function convertContextWidgetToQuery(widget) {
  const query = {
    page: 1,
  };

  const { default_sort_column: defaultSortColumn } = widget;

  if (defaultSortColumn && defaultSortColumn.property) {
    query.sortBy = defaultSortColumn.property;
    query.descending = defaultSortColumn.direction === 'DESC';
  } else {
    query.sortBy = null;
    query.descending = null;
  }

  return query;
}

export function convertAlarmUserPreferenceToQuery(userPreference) {
  return {
    rowsPerPage: userPreference.widget_preferences.itemsPerPage || PAGINATION_LIMIT,
  };
}

export function convertContextUserPreferenceToQuery(userPreference) {
  return {
    rowsPerPage: userPreference.widget_preferences.itemsPerPage || PAGINATION_LIMIT,
    selectedTypes: userPreference.widget_preferences.selectedTypes || [],
  };
}

export function convertUserPreferenceToQuery(userPreference) {
  switch (userPreference.widgetXtype) {
    case WIDGET_TYPES.alarmList:
      return convertAlarmUserPreferenceToQuery(userPreference);
    case WIDGET_TYPES.context:
      return convertContextUserPreferenceToQuery(userPreference);
    default:
      return {};
  }
}

export function convertWidgetToQuery(widget) {
  switch (widget.xtype) {
    case WIDGET_TYPES.alarmList:
      return convertAlarmWidgetToQuery(widget);
    case WIDGET_TYPES.context:
      return convertContextWidgetToQuery(widget);
    default:
      return {};
  }
}

export default {
  convertContextWidgetToQuery,
  convertContextUserPreferenceToQuery,
  convertUserPreferenceToQuery,
  convertWidgetToQuery,
};
