import { isUndefined } from 'lodash';

import { QUICK_RANGES } from '@/constants';

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
    },
  } = widget;

  const { start, stop } = QUICK_RANGES[defaultTimeRange];

  return {
    displayParameter: defaultDisplayParameter,
    showType: defaultShowType,
    showTrend: false,
    interval: {
      from: start,
      to: stop,
    },
  };
};
