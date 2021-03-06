import { omit, isUndefined, isEmpty } from 'lodash';

import { PAGINATION_LIMIT, DEFAULT_WEATHER_LIMIT } from '@/config';
import { WIDGET_TYPES, STATS_QUICK_RANGES } from '@/constants';

import { prepareMainFilterToQueryFilter, getMainFilter } from './filter';
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
 *  This function converts widget.parameters.alarmsStateFilter to query Object
 *
 * @param {Object} parameters
 * @returns {{}}
 */
export function convertAlarmStateFilterToQuery({ parameters }) {
  const { alarmsStateFilter } = parameters;
  const query = {};

  if (alarmsStateFilter) {
    if (!isUndefined(alarmsStateFilter.opened)) {
      query.opened = alarmsStateFilter.opened;
    }

    if (!isUndefined(alarmsStateFilter.resolved)) {
      query.resolved = alarmsStateFilter.resolved;
    }
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
    alarmsStateFilter = {},
    liveReporting = {},
    widgetColumns,
    itemsPerPage,
  } = widget.parameters;

  const query = {
    page: 1,
    opened: alarmsStateFilter.opened || false,
    resolved: alarmsStateFilter.resolved || false,
    limit: itemsPerPage || PAGINATION_LIMIT,
  };

  if (!isEmpty(liveReporting)) {
    query.tstart = liveReporting.tstart;
    query.tstop = liveReporting.tstop;
  } else if (query.resolved) {
    query.tstart = STATS_QUICK_RANGES.last30Days.start;
    query.tstop = STATS_QUICK_RANGES.last30Days.stop;
  }

  if (widgetColumns) {
    query.active_columns = widgetColumns.map(v => v.value);
  }

  return { ...query, ...convertSortToQuery(widget) };
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
 * This function converts widget with type stats field to query Object
 *
 * @param {Object} widget
 * @returns {{}}
 */
export function convertWidgetStatsParameterToQuery(widget) {
  const statsList = Object.keys(widget.parameters.stats).reduce((acc, stat) => {
    acc[stat] = {
      ...omit(widget.parameters.stats[stat], ['position']),
      stat: widget.parameters.stats[stat].stat.value,
    };
    return acc;
  }, {});

  return {
    ...widget.parameters,

    stats: statsList,
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
 * This function converts widget with type 'Stats number' to query Object
 *
 * @param {Object} widget
 * @returns {{}}
 */
export function convertStatsNumberWidgetToQuery(widget) {
  const { stat } = widget.parameters;
  const query = {
    ...omit(widget.parameters, [
      'statColors',
      'criticityLevels',
      'yesNoMode',
      'statName',
    ]),

    trend: true,
  };

  if (stat) {
    query.stats = {
      [stat.title]: {
        parameters: stat.parameters,
        stat: stat.stat.value,
        trend: true,
      },
    };
  }

  query.trend = true;

  return query;
}

/**
 * This function converts widget with type 'Stats Pareto diagram' to query Object
 *
 * @param {Object} widget
 * @returns {{}}
 */
export function convertStatsParetoWidgetToQuery(widget) {
  const { stat } = widget.parameters;
  const query = { ...widget.parameters };

  if (stat) {
    query.stats = {
      [stat.title]: {
        ...omit(stat, ['title']),
        stat: stat.stat.value,
        aggregate: ['sum'],
      },
    };
  }

  return query;
}

/**
 *
 * @param widget
 * @returns {{filters: *}}
 */
export function convertCounterWidgetToQuery(widget) {
  const { viewFilters, isCorrelationEnabled = false } = widget.parameters;

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
export function convertAlarmUserPreferenceToQuery({ widget_preferences: widgetPreferences }) {
  const {
    itemsPerPage,
    isCorrelationEnabled = false,
  } = widgetPreferences;

  const query = {
    correlation: isCorrelationEnabled,
  };

  if (itemsPerPage) {
    query.limit = itemsPerPage;
  }

  return query;
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
    case WIDGET_TYPES.weather:
      return convertWeatherWidgetToQuery(widget);
    case WIDGET_TYPES.statsCurves:
    case WIDGET_TYPES.statsHistogram:
    case WIDGET_TYPES.statsTable:
    case WIDGET_TYPES.text:
      return convertWidgetStatsParameterToQuery(widget);
    case WIDGET_TYPES.statsNumber:
      return convertStatsNumberWidgetToQuery(widget);
    case WIDGET_TYPES.statsPareto:
      return convertStatsParetoWidgetToQuery(widget);
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
  const userPreferenceQuery = convertUserPreferenceToQuery(userPreference);
  let query = {
    ...widgetQuery,
    ...userPreferenceQuery,
  };

  const WIDGET_FILTER_KEYS_MAP = {
    [WIDGET_TYPES.alarmList]: 'filter',
    [WIDGET_TYPES.context]: 'mainFilter',
  };

  const filterKey = WIDGET_FILTER_KEYS_MAP[widget.type];

  if (filterKey) {
    const activeMainFilter = getMainFilter(widget, userPreference);

    if (activeMainFilter) {
      query[filterKey] = prepareMainFilterToQueryFilter(activeMainFilter);
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
