import { PAGINATION_LIMIT } from '@/config';
import { WIDGET_TYPES } from '@/constants';

/**
 * This function converts widget with type 'listalarm' to query Object
 *
 * @param {Object} widget
 * @returns {{}}
 */
export function convertAlarmWidgetToQuery(widget) {
  const query = {
    page: 1,
  };

  const { default_sort_column: defaultSortColumn } = widget;

  if (defaultSortColumn && defaultSortColumn.property) {
    query.sortKey = defaultSortColumn.property.startsWith('v.') ?
      defaultSortColumn.property : `v.${defaultSortColumn.property}`;

    query.sortDir = defaultSortColumn.direction;
  } else {
    query.sortKey = null;
    query.sortDir = null;
  }

  return query;
}

/**
 * This function converts widget with type 'crudcontext' to query Object
 *
 * @param {Object} widget
 * @returns {{}}
 */
export function convertContextWidgetToQuery(widget) {
  const query = {
    page: 1,
  };

  const { default_sort_column: defaultSortColumn } = widget;

  if (defaultSortColumn && defaultSortColumn.property) {
    query.sortKey = defaultSortColumn.property;
    query.sortDir = defaultSortColumn.direction;
  } else {
    query.sortKey = null;
    query.sortDir = null;
  }

  return query;
}

/**
 * This function converts userPreference with widgetXtype 'listalarm' to query Object
 *
 * @param {Object} userPreference
 * @returns {{}}
 */
export function convertAlarmUserPreferenceToQuery(userPreference) {
  return {
    limit: userPreference.widget_preferences.itemsPerPage || PAGINATION_LIMIT,
  };
}

/**
 * This function converts userPreference with widgetXtype 'crudcontext' to query Object
 *
 * @param {Object} userPreference
 * @returns {{}}
 */
export function convertContextUserPreferenceToQuery(userPreference) {
  return {
    limit: userPreference.widget_preferences.itemsPerPage || PAGINATION_LIMIT,
    selectedTypes: userPreference.widget_preferences.selectedTypes || [],
  };
}

/**
 * This function converts userPreference to query Object
 *
 * @param {Object} userPreference
 * @returns {{}}
 */
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

/**
 * This function converts widget to query Object
 *
 * @param {Object} widget
 * @returns {{}}
 */
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
