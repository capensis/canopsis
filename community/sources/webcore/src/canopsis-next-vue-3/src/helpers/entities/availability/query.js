import { isUndefined } from 'lodash';

import {
  AVAILABILITY_TREND_TYPE_BY_QUICK_RANGE,
  AVAILABILITY_TREND_TYPES,
  QUICK_RANGES,
} from '@/constants';

import { findQuickRangeByInterval } from '@/helpers/date/date-intervals';
import { getInfosWidgetColumn } from '@/helpers/entities/widget/column/form';

/**
 * Convert userPreference with widget availability type to query object.
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
 * Remove useless part of infos value.
 *
 * We should remove .value, because it's useless for this endpoint.
 *
 * @param {string} value
 * @returns {string}
 */
export const prepareAvailabilityWidgetColumnValue = value => getInfosWidgetColumn(value) ? value.replace(/.value$/, '') : value;

/**
 * Convert availability widget default parameters to query object.
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
    active_columns: widgetColumns.map(column => prepareAvailabilityWidgetColumnValue(column.value)),
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
 * Find quick availability range by interval.
 *
 * @param {Object} interval
 * @returns {string}
 */
export const getAvailabilitiesTrendByInterval = (interval) => {
  const selectedQuickRange = findQuickRangeByInterval(interval);
  const selectedQuickRangeType = selectedQuickRange?.value;

  return (
    AVAILABILITY_TREND_TYPE_BY_QUICK_RANGE[selectedQuickRangeType] ||
    AVAILABILITY_TREND_TYPES.custom
  );
};
