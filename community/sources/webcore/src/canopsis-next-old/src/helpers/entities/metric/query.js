import { uniq } from 'lodash';

import { QUICK_RANGES } from '@/constants';

import { isOmitEqual } from '@/helpers/collection';
import { isRatioMetric, isCustomCriteria } from '@/helpers/entities/metric/form';

/**
 * Check query is changed with interval
 *
 * @param {Object} query
 * @param {Object} oldQuery
 * @param {number} minDate
 * @returns {boolean}
 */
export const isMetricsQueryChanged = (query, oldQuery, minDate) => {
  const isFromChanged = query.interval?.from !== oldQuery.interval?.from;
  const isFromEqualMinDate = query.interval?.from === minDate;
  const isToChanged = query.interval?.to !== oldQuery.interval?.to;
  const isQueryWithoutIntervalChanged = !isOmitEqual(query, oldQuery, ['interval']);

  return isQueryWithoutIntervalChanged || (isFromChanged && !isFromEqualMinDate) || isToChanged;
};

/**
 * Convert widget filter to query
 *
 * @param {AlarmChart[]} charts
 * @returns {string[]}
 */
export function convertWidgetChartsToPerfDataQuery(charts) {
  const allMetrics = charts.reduce((acc, chart) => {
    acc.push(...chart.parameters.metrics.map(({ metric }) => metric));

    return acc;
  }, []);

  return uniq(allMetrics);
}

/**
 * This function converts userPreference with widget type 'Map' to query Object
 *
 * @param {Object} userPreference
 * @returns {{ sampling: string, interval: Object }}
 */
export const convertChartUserPreferenceToQuery = ({ content: { sampling, interval, mainFilter } }) => {
  const query = {
    filter: mainFilter,
  };

  if (sampling) {
    query.sampling = sampling;
  }

  if (interval) {
    query.interval = interval;
  }

  return query;
};

/**
 * This function converts chart widgets default parameters to query Object
 *
 * @param {Widget} widget
 * @returns {{ lockedFilter: string | null, sampling: string, interval: Object }}
 */
export function convertChartWidgetDefaultParametersToQuery(widget) {
  const {
    mainFilter,
    default_sampling: defaultSampling,
    default_time_range: defaultTimeRange,
  } = widget.parameters;

  return {
    lockedFilter: mainFilter,
    sampling: defaultSampling,
    interval: {
      from: QUICK_RANGES[defaultTimeRange].start,
      to: QUICK_RANGES[defaultTimeRange].stop,
    },
  };
}

/**
 * This function converts bar chart widget to query Object
 *
 * @param {Widget} widget
 * @returns {Object}
 */
export function convertChartWidgetToQuery(widget) {
  const { comparison = false, metrics = [] } = widget.parameters;

  return {
    ...convertChartWidgetDefaultParametersToQuery(widget),

    with_history: comparison,
    parameters: metrics.map(({ metric, aggregate_func: aggregateFunc }) => ({ metric, aggregate_func: aggregateFunc })),
  };
}

/**
 * This function converts pie chart widget to query Object
 *
 * @param {Widget} widget
 * @returns {Object}
 */
export function convertPieChartWidgetToQuery(widget) {
  const {
    metrics = [],
    aggregate_func: widgetAggregateFunc,
  } = widget.parameters;

  return {
    ...convertChartWidgetDefaultParametersToQuery(widget),

    parameters: metrics.map(({ metric, aggregate_func: metricAggregateFunc }) => ({
      metric,
      aggregate_func: metricAggregateFunc || widgetAggregateFunc,
    })),
  };
}

/**
 * This function converts numbers widget to query Object
 *
 * @param {Widget} widget
 * @returns {Object}
 */
export function convertNumbersWidgetToQuery(widget) {
  const { metrics = [], show_trend: showTrend = false } = widget.parameters;

  return {
    ...convertChartWidgetDefaultParametersToQuery(widget),

    with_history: showTrend,
    parameters: metrics.map(({ metric, aggregate_func: aggregateFunc }) => ({
      metric,
      aggregate_func: isRatioMetric(metric) ? undefined : aggregateFunc,
    })),
  };
}

/**
 * This function converts statistics widgets default parameters to query Object
 *
 * @param {Widget} widget
 * @returns {Object}
 */
export function convertStatisticsWidgetParametersToQuery(widget) {
  const {
    mainFilter,
    mainParameter = {},
    widgetColumns = [],
    default_time_range: defaultTimeRange,
  } = widget.parameters;

  const query = {
    interval: {
      from: QUICK_RANGES[defaultTimeRange].start,
      to: QUICK_RANGES[defaultTimeRange].stop,
    },
    parameters: widgetColumns.map(({ metric, criteria }) => {
      const parameter = {
        metric,
      };

      if (criteria) {
        parameter.criteria = criteria;
      }

      return parameter;
    }),
    entity_patterns: mainParameter.patterns.map(({ title, entity_pattern: pattern }) => ({ title, pattern })),
    lockedFilter: mainFilter,
  };

  if (!isCustomCriteria(mainParameter.criteria)) {
    query.criteria = mainParameter.criteria;
  }

  return query;
}

/**
 * This function converts userPreference with widget statistics type to query Object
 *
 * @param {Object} userPreference
 * @returns {Object}
 */
export const convertStatisticsUserPreferenceToQuery = ({ content: { interval, mainFilter } }) => {
  const query = {
    filter: mainFilter,
  };

  if (interval) {
    query.interval = interval;
  }

  return query;
};
