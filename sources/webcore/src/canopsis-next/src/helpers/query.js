import get from 'lodash/get';

import { PAGINATION_LIMIT } from '@/config';
import { WIDGET_TYPES } from '@/constants';

/**
 * WIDGET CONVERTERS
 */

/**
 * This function converts widget.default_sort_column to query Object
 *
 * @param {Object} widget
 * @returns {{}}
 */
export function convertDefaultSortColumnToQuery(widget) {
  const { default_sort_column: defaultSortColumn } = widget;

  if (defaultSortColumn && defaultSortColumn.property) {
    return { sortKey: defaultSortColumn.property, sortDir: defaultSortColumn.direction };
  }

  return { sortKey: null, sortDir: null };
}

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
  const stateFilter = widget.alarms_state_filter;

  if (stateFilter) {
    if (stateFilter.state && (stateFilter.opened === undefined && stateFilter.resolved === undefined)) {
      query.opened = stateFilter.state === 'opened';
      query.resolved = stateFilter.state === 'resolved';
    } else {
      query.opened = Boolean(stateFilter.opened);
      query.resolved = Boolean(stateFilter.resolved);
    }
  }

  if (widget.widget_columns) {
    query.active_columns = widget.widget_columns.map(v => v.value);
  }
  return { ...query, ...convertDefaultSortColumnToQuery(widget) };
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

  return { ...query, ...convertDefaultSortColumnToQuery(widget) };
}

/**
 * USER_PREFERENCE CONVERTERS
 */

/**
 * This function converts userPreference with widgetXtype 'listalarm' to query Object
 *
 * @param {Object} userPreference
 * @returns {{}}
 */
export function convertAlarmUserPreferenceToQuery({ widget_preferences: widgetPreferences }) {
  return {
    limit: get(widgetPreferences, 'itemsPerPage', PAGINATION_LIMIT),
    filter: get(widgetPreferences, 'selected_filter.filter'),
  };
}

/**
 * This function converts userPreference with widgetXtype 'crudcontext' to query Object
 *
 * @param {Object} userPreference
 * @returns {{}}
 */
export function convertContextUserPreferenceToQuery({ widget_preferences: widgetPreferences }) {
  return {
    limit: get(widgetPreferences, 'itemsPerPage', PAGINATION_LIMIT),
    selectedTypes: get(widgetPreferences, 'selectedTypes', []),
  };
}

/**
 * MAIN CONVERTERS
 */

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
