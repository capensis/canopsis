import { isUndefined, isEmpty } from 'lodash';

import { PAGINATION_LIMIT, DEFAULT_WEATHER_LIMIT } from '@/config';
import {
  WIDGET_TYPES,
  QUICK_RANGES,
  ALARMS_LIST_WIDGET_ACTIVE_COLUMNS_MAP,
  SORT_ORDERS,
  ALARMS_OPENED_VALUES,
} from '@/constants';

import { prepareMainFilterToQueryFilter, getMainFilterAndCondition } from './filter';
import {
  prepareRemediationInstructionsFiltersToQuery,
  getRemediationInstructionsFilters,
} from './filter/remediation-instructions-filter';

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
    query.typesFilter = {
      $or: selectedTypes.map(type => ({ type })),
    };
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
  const {
    filters,
    considerPbehaviors,
  } = widget.parameters;

  const query = {
    considerPbehaviors,
    filters: filters || [],
  };

  return { ...query, ...convertAlarmStateFilterToQuery(widget) };
}

/**
 *
 * @param widget
 * @returns {{filters: *}}
 */
export function convertCounterWidgetToQuery(widget) {
  const { viewFilters = [], isCorrelationEnabled = false } = widget.parameters;

  return {
    ...convertAlarmStateFilterToQuery(widget),

    correlation: isCorrelationEnabled,
    filters: viewFilters.map(({ filter }) => filter),
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

  const WIDGET_FILTER_KEYS_MAP = {
    [WIDGET_TYPES.alarmList]: 'filter',
    [WIDGET_TYPES.context]: 'mainFilter',
    [WIDGET_TYPES.serviceWeather]: 'filter',
  };

  const filterKey = WIDGET_FILTER_KEYS_MAP[widget.type];

  if (filterKey) {
    const { mainFilter: activeMainFilter, condition } = getMainFilterAndCondition(widget, userPreference);

    if (activeMainFilter) {
      query[filterKey] = prepareMainFilterToQueryFilter(activeMainFilter, condition);
    }
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
