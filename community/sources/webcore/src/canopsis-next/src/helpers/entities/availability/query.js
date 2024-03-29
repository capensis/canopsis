import { isUndefined } from 'lodash';

import { AVAILABILITY_TREND_TYPE_BY_QUICK_RANGE, AVAILABILITY_TREND_TYPES, QUICK_RANGES } from '@/constants';

import { findQuickRangeByInterval } from '@/helpers/date/date-intervals';

/**
 * This function converts userPreference with widget availability type to query Object
 *
 * @param {Object} userPreference
 * @returns {Object}
 */
export const convertAvailabilityUserPreferenceToQuery = (userPreference) => {
  const {
    content: {
      mainFilter,
      show_trend: showTrend,
      show_type: showType,
      display_parameter: displayParameter,
    },
  } = userPreference;

  const query = {
    filter: mainFilter,
  };

  if (!isUndefined(showTrend)) {
    query.showTrend = showTrend;
  }

  if (!isUndefined(showType)) {
    query.showType = showType;
  }

  if (!isUndefined(displayParameter)) {
    query.displayParameter = displayParameter;
  }

  return query;
};

/**
 * This function converts availability widgets default parameters to query Object
 *
 * @param {Widget} widget
 * @returns {Object}
 */
export const convertAvailabilityWidgetParametersToQuery = (widget) => {
  const {
    parameters: {
      default_display_parameter: defaultDisplayParameter,
      default_show_type: defaultShowType,
      default_time_range: defaultTimeRange,
      widget_columns: widgetColumns,
      mainFilter,
    },
  } = widget;

  const { start, stop } = QUICK_RANGES[defaultTimeRange];

  return {
    active_columns: widgetColumns.map(v => v.value),
    lockedFilter: mainFilter,
    displayParameter: defaultDisplayParameter,
    showType: defaultShowType,
    showTrend: false,
    interval: {
      from: start,
      to: stop,
    },
  };
};

/**
 * Find quick availability range by interval
 *
 * @param {Object} interval
 * @return {string}
 */
export const getAvailabilitiesTrendByInterval = (interval) => {
  const selectedQuickRange = findQuickRangeByInterval(interval);
  const selectedQuickRangeType = selectedQuickRange?.value;

  return AVAILABILITY_TREND_TYPE_BY_QUICK_RANGE[selectedQuickRangeType]
    || AVAILABILITY_TREND_TYPES.custom;
};
