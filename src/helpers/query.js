import get from 'lodash/get';
import isUndefined from 'lodash/isUndefined';
import isEmpty from 'lodash/isEmpty';

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
export function convertDefaultSortColumnToQuery({ parameters }) {
  const { sortColumn, sortOrder } = parameters;

  if (sortColumn && sortOrder) {
    return { sortKey: sortColumn, sortDir: sortOrder };
  }

  return { sortKey: null, sortDir: null };
}

/**
 * This function converts widget with type 'AlarmsList' to query Object
 *
 * @param {Object} widget
 * @returns {{}}
 */
export function convertAlarmWidgetToQuery(widget) {
  const {
    displayOpenAlarms,
    displayResolvedAlarms,
    columnTranslations,
    itemsPerPage,
    mainFilter,
  } = widget.parameters;

  const query = {
    page: 1,
    limit: itemsPerPage || PAGINATION_LIMIT,
  };

  if (!isEmpty(mainFilter)) {
    query.filter = mainFilter.filter;
  }

  if (!isUndefined(displayOpenAlarms)) {
    query.opened = displayOpenAlarms;
  }

  if (!isUndefined(displayResolvedAlarms)) {
    query.resolved = displayResolvedAlarms;
  }

  if (columnTranslations) {
    query.active_columns = columnTranslations.map(v => v.value);
  }

  return { ...query, ...convertDefaultSortColumnToQuery(widget) };
}

/**
 * This function converts widget with type 'Context' to query Object
 *
 * @param {Object} widget
 * @returns {{}}
 */
export function convertContextWidgetToQuery(widget) {
  const query = {
    page: 1,
    limit: get(widget, 'parameters.itemsPerPage', PAGINATION_LIMIT),
    selectedTypes: get(widget, 'parameters.selectedTypes', []),
  };

  return { ...query, ...convertDefaultSortColumnToQuery(widget) };
}

/**
 * USER_PREFERENCE CONVERTERS
 */

/**
 * This function converts userPreference with widgetXtype 'AlarmsList' to query Object
 *
 * @param {Object} userPreference
 * @returns {{}}
 */
export function convertAlarmUserPreferenceToQuery({ widget_preferences: widgetPreferences }) {
  const query = {
    limit: get(widgetPreferences, 'itemsPerPage', PAGINATION_LIMIT),
  };

  if (!isEmpty(widgetPreferences.mainFilter)) {
    query.filter = widgetPreferences.mainFilter.filter;
  }

  return query;
}

/**
 * This function converts userPreference with widgetXtype 'Context' to query Object
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
  switch (widget.type) {
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
