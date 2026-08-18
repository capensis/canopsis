import { uniq } from 'lodash';

import { DATETIME_FORMATS, DATETIME_INTERVAL_TYPES, QUICK_RANGES } from '@/constants';

import { isOmitEqual } from '@/helpers/collection';
import { convertDateIntervalToTimestampByTimezone } from '@/helpers/date/date-intervals';

import { isRatioMetric } from './form';

/**
 * Check query is changed with interval.
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
 * Convert widget filter to query.
 *
 * @param {AlarmChart[]} charts
 * @returns {string[]}
 */
export function convertWidgetChartsToPerfDataQuery(charts = []) {
  const allMetrics = charts.reduce((acc, chart) => {
    acc.push(...(chart.parameters?.metrics ?? []).map(({ metric }) => metric));

    return acc;
  }, []);

  return uniq(allMetrics);
}

/**
 * Convert KPI metrics interval to API timestamps.
 *
 * Metrics API validates both `from` and `to` as day starts in Canopsis timezone.
 *
 * @param {{ from: number|string, to: number|string }} interval
 * @param {string} timezone
 * @returns {{ from: number, to: number }}
 */
export const convertMetricsIntervalToTimestamp = (interval, timezone) => ({
  from: convertDateIntervalToTimestampByTimezone(
    interval.from,
    DATETIME_INTERVAL_TYPES.start,
    DATETIME_FORMATS.dateTimePicker,
    undefined,
    timezone,
  ),
  to: convertDateIntervalToTimestampByTimezone(
    interval.to,
    DATETIME_INTERVAL_TYPES.start,
    DATETIME_FORMATS.dateTimePicker,
    undefined,
    timezone,
  ),
});

/**
 * Convert user preference of a chart widget to query.
 *
 * @param {Object} userPreference
 * @returns {{ sampling?: string, interval?: Object, filter?: * }}
 */
export const convertChartUserPreferenceToQuery = ({
  content: { sampling, interval, mainFilter } = {},
}) => {
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
 * Convert chart widget default parameters to query.
 *
 * @param {Object} widget
 * @returns {Object}
 */
export const convertChartWidgetDefaultParametersToQuery = (widget) => {
  const {
    mainFilter,
    default_sampling: defaultSampling,
    default_time_range: defaultTimeRange,
  } = widget.parameters ?? {};
  const range = QUICK_RANGES[defaultTimeRange] ?? QUICK_RANGES.last30Days;

  return {
    lockedFilter: mainFilter,
    sampling: defaultSampling,
    interval: {
      from: range.start,
      to: range.stop,
    },
  };
};

/**
 * Convert bar/line chart widget to query.
 *
 * @param {Object} widget
 * @returns {Object}
 */
export const convertChartWidgetToQuery = (widget) => {
  const { comparison = false, metrics = [] } = widget.parameters ?? {};

  return {
    ...convertChartWidgetDefaultParametersToQuery(widget),
    with_history: comparison,
    parameters: metrics.map(({ metric, aggregate_func: aggregateFunc }) => ({
      metric,
      aggregate_func: aggregateFunc,
    })),
  };
};

/**
 * Convert numbers widget to query.
 *
 * @param {Object} widget
 * @returns {Object}
 */
export const convertNumbersWidgetToQuery = (widget) => {
  const { metrics = [], show_trend: showTrend = false } = widget.parameters ?? {};

  return {
    ...convertChartWidgetDefaultParametersToQuery(widget),
    with_history: showTrend,
    parameters: metrics.map(({ metric, aggregate_func: aggregateFunc }) => ({
      metric,
      aggregate_func: isRatioMetric(metric) ? undefined : aggregateFunc,
    })),
  };
};
