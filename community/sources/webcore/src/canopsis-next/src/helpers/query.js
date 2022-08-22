import { isUndefined, isEmpty, omit } from 'lodash';

import { PAGINATION_LIMIT, DEFAULT_WEATHER_LIMIT } from '@/config';
import {
  WIDGET_TYPES,
  QUICK_RANGES,
  ALARMS_LIST_WIDGET_ACTIVE_COLUMNS_MAP,
  SORT_ORDERS,
  ALARMS_OPENED_VALUES,
  DATETIME_FORMATS,
} from '@/constants';

import { getMainFilter } from './filter';
import {
  prepareRemediationInstructionsFiltersToQuery,
  getRemediationInstructionsFilters,
} from './filter/remediation-instructions-filter';
import { convertStartDateIntervalToTimestamp, convertStopDateIntervalToTimestamp } from './date/date-intervals';

/**
 * WIDGET CONVERTERS
 */

/**
 * This function converts widget.parameters.sort to query Object
 *
 * @param {Object} widget
 * @returns {{}}
 */
export function convertSortToQuery({ parameters }) {
  const { sort } = parameters;

  if (sort && sort.column && sort.order) {
    return { sortKey: sort.column, sortDir: sort.order };
  }

  return { sortKey: null, sortDir: null };
}

/**
 *  This function converts widget.parameters.opened to query Object
 *
 * @param {Object} parameters
 * @returns {{ opened: boolean }}
 */
export function convertAlarmStateFilterToQuery({ parameters }) {
  const { opened = ALARMS_OPENED_VALUES.opened } = parameters;
  const query = {};

  if (!isUndefined(opened)) {
    query.opened = opened;
  }

  return query;
}

/**
 * This function converts widget with type 'AlarmsList' to query Object
 *
 * @param {Object} widget
 * @returns {{}}
 */
export function convertAlarmWidgetToQuery(widget) {
  const {
    liveReporting = {},
    widgetColumns,
    itemsPerPage,
    sort,
    opened = ALARMS_OPENED_VALUES.opened,
  } = widget.parameters;

  const query = {
    opened,
    page: 1,
    limit: itemsPerPage || PAGINATION_LIMIT,
    with_instructions: true,
    with_links: true,
    multiSortBy: [],
  };

  if (!isEmpty(liveReporting)) {
    query.tstart = liveReporting.tstart;
    query.tstop = liveReporting.tstop;
  } else if (query.opened === ALARMS_OPENED_VALUES.resolved) {
    query.tstart = QUICK_RANGES.last30Days.start;
    query.tstop = QUICK_RANGES.last30Days.stop;
  }

  if (widgetColumns) {
    query.active_columns = widgetColumns.map(v => (ALARMS_LIST_WIDGET_ACTIVE_COLUMNS_MAP[v.value] || v.value));
  }

  if (sort && sort.column && sort.order) {
    query.multiSortBy.push({ sortBy: sort.column, descending: sort.order === SORT_ORDERS.desc });
  }

  return query;
}

/**
 * This function converts widget with type 'Context' to query Object
 *
 * @param {Object} widget
 * @returns {{}}
 */
export function convertContextWidgetToQuery(widget) {
  const {
    itemsPerPage,
    selectedTypes,
    widgetColumns,
  } = widget.parameters;

  const query = {
    page: 1,
    limit: itemsPerPage || PAGINATION_LIMIT,
  };

  if (widgetColumns) {
    query.active_columns = widgetColumns.map(v => v.value);
  }

  if (!isEmpty(selectedTypes)) {
    query.type = selectedTypes;
  }

  return { ...query, ...convertSortToQuery(widget) };
}

/**
 * This function converts widget with type 'ServiceWeather' to query Object
 *
 * @param {Object} widget
 * @returns {{}}
 */
export function convertWeatherWidgetToQuery(widget) {
  const { limit } = widget.parameters;

  return {
    ...convertSortToQuery(widget),
    limit: limit || DEFAULT_WEATHER_LIMIT,
  };
}

/**
 * This function converts widget with type 'StatsCalendar' to query Object
 *
 * @param {Object} widget
 * @returns {{}}
 */
export function convertStatsCalendarWidgetToQuery(widget) {
  const { filters = [], parameters: { considerPbehaviors = false } } = widget;

  return {
    ...convertAlarmStateFilterToQuery(widget),

    considerPbehaviors,
    filters,
    time_field: 't',
  };
}

/**
 * This function converts widget with type 'counter' widget to query Object
 *
 * @param widget
 * @returns {{filters: *}}
 */
export function convertCounterWidgetToQuery(widget) {
  const { filters = [], parameters: { isCorrelationEnabled = false } } = widget;

  return {
    ...convertAlarmStateFilterToQuery(widget),

    correlation: isCorrelationEnabled,
    filters: filters.map(({ _id: id }) => id),
  };
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
export function convertAlarmUserPreferenceToQuery({ content }) {
  const {
    itemsPerPage,
    category,
    isCorrelationEnabled = false,
  } = content;

  const query = {
    correlation: isCorrelationEnabled,
    category,
  };

  if (itemsPerPage) {
    query.limit = itemsPerPage;
  }

  return query;
}

/**
 * This function converts userPreference with widget type 'ServiceWeather' to query Object
 *
 * @param {Object} userPreference
 * @returns {{ category: string }}
 */
export function convertWeatherUserPreferenceToQuery({ content }) {
  const { category } = content;

  return { category };
}

/**
 * This function converts userPreference with widget type 'Map' to query Object
 *
 * @param {Object} userPreference
 * @returns {{ category: string }}
 */
export const convertMapUserPreferenceToQuery = ({ content: { category } }) => ({ category });

/**
 * This function converts userPreference with widgetXtype 'Context' to query Object
 *
 * @param {Object} userPreference
 * @returns {{ category: string }}
 */
export function convertContextUserPreferenceToQuery({ content }) {
  const { category, noEvents } = content;

  return {
    category,
    no_events: noEvents,
  };
}

/**
 * MAIN CONVERTERS
 */

/**
 * This function converts userPreference to query Object
 *
 * @param {Object} userPreference
 * @param {WidgetType} widgetType
 * @returns {Object}
 */
export function convertUserPreferenceToQuery(userPreference, widgetType) {
  switch (widgetType) {
    case WIDGET_TYPES.alarmList:
      return convertAlarmUserPreferenceToQuery(userPreference);
    case WIDGET_TYPES.context:
      return convertContextUserPreferenceToQuery(userPreference);
    case WIDGET_TYPES.serviceWeather:
      return convertWeatherUserPreferenceToQuery(userPreference);
    case WIDGET_TYPES.map:
      return convertMapUserPreferenceToQuery(userPreference);
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
    case WIDGET_TYPES.serviceWeather:
      return convertWeatherWidgetToQuery(widget);
    case WIDGET_TYPES.statsCalendar:
      return convertStatsCalendarWidgetToQuery(widget);
    case WIDGET_TYPES.counter:
      return convertCounterWidgetToQuery(widget);
    default:
      return {};
  }
}

/**
 * Prepare query by widget and userPreference objects
 *
 * @param {Object} widget
 * @param {Object} userPreference
 * @returns {Object}
 */
export function prepareQuery(widget, userPreference) {
  const widgetQuery = convertWidgetToQuery(widget);
  const userPreferenceQuery = convertUserPreferenceToQuery(userPreference, widget.type);

  let query = {
    ...widgetQuery,
    ...userPreferenceQuery,
  };

  const filter = getMainFilter(widget, userPreference);

  if (filter) {
    query.filter = filter;
  }

  const remediationInstructionsFilters = getRemediationInstructionsFilters(widget, userPreference);

  if (remediationInstructionsFilters.length) {
    query = {
      ...query,
      ...prepareRemediationInstructionsFiltersToQuery(remediationInstructionsFilters),
    };
  }

  return query;
}

/**
 * Prepare query for alarm details fetching
 *
 * @param {Alarm} alarm
 * @param {Widget} widget
 * @returns {Object}
 */
export const prepareAlarmDetailsQuery = (alarm, widget) => ({
  _id: alarm._id,
  with_instructions: true,
  opened: widget.parameters.opened,
  steps: {
    page: 1,
    limit: PAGINATION_LIMIT,
  },
  children: {
    page: 1,
    limit: PAGINATION_LIMIT,
    sort_by: '',
    sort: '',
    multi_sort: [],
  },
});

/**
 * Convert alarms list query to request parameters
 *
 * @param {Object} query
 * @returns {Object}
 */
export const convertAlarmsListQueryToRequest = (query) => {
  const result = omit(query, [
    'tstart',
    'tstop',
    'sortKey',
    'sortDir',
    'category',
    'multiSortBy',
    'limit',
  ]);

  const {
    tstart,
    tstop,
    sortKey,
    sortDir,
    category,
    multiSortBy = [],
    limit = PAGINATION_LIMIT,
  } = query;

  if (tstart) {
    result.tstart = convertStartDateIntervalToTimestamp(tstart, DATETIME_FORMATS.dateTimePicker);
  }

  if (tstop) {
    result.tstop = convertStopDateIntervalToTimestamp(tstop, DATETIME_FORMATS.dateTimePicker);
  }

  if (sortKey) {
    result.sort_by = sortKey;
    result.sort = sortDir.toLowerCase();
  }

  if (category) {
    result.category = category;
  }

  if (multiSortBy.length) {
    result.multi_sort = multiSortBy
      .map(({ sortBy, descending }) => `${sortBy},${(descending ? SORT_ORDERS.desc : SORT_ORDERS.asc).toLowerCase()}`);
  }

  result.limit = limit;

  return result;
};
