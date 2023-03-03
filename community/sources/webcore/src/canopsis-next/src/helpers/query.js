import { isUndefined, isEmpty, omit, isArray } from 'lodash';

import { PAGINATION_LIMIT, DEFAULT_WEATHER_LIMIT } from '@/config';
import {
  WIDGET_TYPES,
  QUICK_RANGES,
  SORT_ORDERS,
  ALARMS_OPENED_VALUES,
  DATETIME_FORMATS,
  DEFAULT_ALARMS_WIDGET_GROUP_COLUMNS,
} from '@/constants';

import featuresService from '@/services/features';

import {
  prepareRemediationInstructionsFiltersToQuery,
  getRemediationInstructionsFilters,
} from './filter/remediation-instructions-filter';
import {
  convertStartDateIntervalToTimestamp,
  convertStopDateIntervalToTimestamp,
} from './date/date-intervals';

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
 * Convert widget filter to query
 *
 * @param parameters
 * @returns {{ lockedFilter?: string }}
 */
export function convertWidgetFilterToQuery({ parameters }) {
  if (parameters.mainFilter) {
    return {
      lockedFilter: parameters.mainFilter,
    };
  }

  return {};
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
    widgetColumns = [],
    itemsPerPage,
    sort,
    mainFilter,
    opened = ALARMS_OPENED_VALUES.opened,
  } = widget.parameters;

  const query = {
    opened,
    page: 1,
    limit: itemsPerPage || PAGINATION_LIMIT,
    with_instructions: true,
    with_declare_tickets: true,
    with_links: true,
    multiSortBy: [],
    lockedFilter: mainFilter,
  };

  if (!isEmpty(liveReporting)) {
    query.tstart = liveReporting.tstart;
    query.tstop = liveReporting.tstop;
  } else if (query.opened === ALARMS_OPENED_VALUES.resolved) {
    query.tstart = QUICK_RANGES.last30Days.start;
    query.tstop = QUICK_RANGES.last30Days.stop;
  }

  if (widgetColumns.length) {
    query.active_columns = widgetColumns.map(v => v.value);
  }

  if (sort && sort.column && sort.order) {
    query.multiSortBy.push({
      sortBy: sort.column,
      descending: sort.order === SORT_ORDERS.desc,
    });
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
    mainFilter,
  } = widget.parameters;

  const query = {
    page: 1,
    limit: itemsPerPage || PAGINATION_LIMIT,
    lockedFilter: mainFilter,
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
  const { limit, mainFilter } = widget.parameters;

  return {
    ...convertSortToQuery(widget),
    limit: limit || DEFAULT_WEATHER_LIMIT,
    lockedFilter: mainFilter,
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
    mainFilter,
    isCorrelationEnabled = false,
  } = content;

  const query = {
    correlation: isCorrelationEnabled,
    filter: mainFilter,
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
  const { category, mainFilter, hide_grey: hideGrey = false } = content;

  return { category, filter: mainFilter, hide_grey: hideGrey };
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
  const { category, noEvents, mainFilter } = content;

  return {
    category,
    filter: mainFilter,
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
  const convertersMap = {
    [WIDGET_TYPES.alarmList]: convertAlarmUserPreferenceToQuery,
    [WIDGET_TYPES.context]: convertContextUserPreferenceToQuery,
    [WIDGET_TYPES.serviceWeather]: convertWeatherUserPreferenceToQuery,
    [WIDGET_TYPES.map]: convertMapUserPreferenceToQuery,

    ...featuresService.get('helpers.query.convertUserPreferenceToQuery.convertersMap'),
  };

  const converter = convertersMap[widgetType];

  return converter ? converter(userPreference) : {};
}

/**
 * This function converts widget to query Object
 *
 * @param {Object} widget
 * @returns {{}}
 */
export function convertWidgetToQuery(widget) {
  const convertersMap = {
    [WIDGET_TYPES.alarmList]: convertAlarmWidgetToQuery,
    [WIDGET_TYPES.context]: convertContextWidgetToQuery,
    [WIDGET_TYPES.serviceWeather]: convertWeatherWidgetToQuery,
    [WIDGET_TYPES.statsCalendar]: convertStatsCalendarWidgetToQuery,
    [WIDGET_TYPES.counter]: convertCounterWidgetToQuery,

    ...featuresService.get('helpers.query.convertWidgetToQuery.convertersMap'),
  };

  const converter = convertersMap[widget.type];

  return converter ? converter(widget) : {};
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
 * Convert filter to query filters
 *
 * @param {string | string[]} filter
 * @returns {string[]}
 */
const convertFilterToQuery = filter => (isArray(filter) ? filter : [filter]).filter(Boolean);

/**
 * Convert locked filter and main filter to query filters
 *
 * @param {string | string[]} filter
 * @param {string | string[]} lockedFilter
 * @returns {string[]}
 */
const convertFiltersToQuery = (filter, lockedFilter) => [
  ...convertFilterToQuery(filter),
  ...convertFilterToQuery(lockedFilter),
];

/**
 * Prepare query for alarm details fetching
 *
 * @param {Alarm} alarm
 * @param {Widget} widget
 * @returns {Object}
 */
export const prepareAlarmDetailsQuery = (alarm, widget) => {
  const { sort = {}, widgetGroupColumns = [] } = widget.parameters;
  const columns = widgetGroupColumns.length > 0
    ? widgetGroupColumns
    : DEFAULT_ALARMS_WIDGET_GROUP_COLUMNS;

  const query = {
    _id: alarm._id,
    with_instructions: true,
    with_declare_tickets: true,
    with_links: true,
    opened: widget.parameters.opened,
    steps: {
      reversed: true,
      page: 1,
      limit: PAGINATION_LIMIT,
    },
    children: {
      page: 1,
      limit: PAGINATION_LIMIT,
      multiSortBy: [],
    },
  };

  if (sort.column && sort.order && columns.some(({ value }) => value.endsWith(sort.column))) {
    query.children.multiSortBy.push({ sortBy: sort.column, descending: sort.order === SORT_ORDERS.desc });
  }

  return query;
};

/**
 * Convert multiSortBy query parameter to request
 *
 * @param {Object[]} multiSortBy
 * @returns {string[]}
 */
export const convertMultiSortToRequest = (multiSortBy = []) => multiSortBy
  .map(({ sortBy, descending }) => `${sortBy},${(descending ? SORT_ORDERS.desc : SORT_ORDERS.asc).toLowerCase()}`);

/**
 * Convert alarmDetails query to request
 *
 * @param {Object} query
 * @returns {Object}
 */
export const convertAlarmDetailsQueryToRequest = query => ({
  ...query,

  children: {
    ...omit(query.children, ['multiSortBy']),

    multi_sort: convertMultiSortToRequest(query.children?.multiSortBy),
  },
});

/**
 * Convert widget query to request parameters
 *
 * @param {Object} query
 * @returns {Object}
 */
export const convertWidgetQueryToRequest = (query) => {
  const result = omit(query, [
    'tstart',
    'tstop',
    'sortKey',
    'sortDir',
    'category',
    'multiSortBy',
    'limit',
    'filter',
    'lockedFilter',
  ]);

  const {
    tstart,
    tstop,
    sortKey,
    sortDir,
    category,
    filter,
    lockedFilter,
    multiSortBy = [],
    limit = PAGINATION_LIMIT,
  } = query;

  if (lockedFilter || filter) {
    result.filters = convertFiltersToQuery(filter, lockedFilter);
  }

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
    result.multi_sort = convertMultiSortToRequest(multiSortBy);
  }

  result.limit = limit;

  return result;
};
