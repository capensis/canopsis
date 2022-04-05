import { isEmpty } from 'lodash';

import { FILTER_DEFAULT_VALUES, FILTER_MONGO_OPERATORS } from '@/constants';

/**
 * Prepare user main filter to query filter
 *
 * @param {Object|Array} [mainFilter = null]
 * @param {string} [condition = FILTER_DEFAULT_VALUES.condition]
 * @returns {Object|string|undefined}
 */
export function prepareMainFilterToQueryFilter(mainFilter = null, condition = FILTER_DEFAULT_VALUES.condition) {
  if (!mainFilter) {
    return undefined;
  }

  if (Array.isArray(mainFilter)) {
    if (!mainFilter.length) {
      return undefined;
    }

    return {
      condition,
      filters: mainFilter,
    };
  }

  return { filter: mainFilter };
}

/**
 * Get main filter and condition for the widget by widget and userPreference parameters
 *
 * @param {Object} widget
 * @param {Object} userPreference
 * @returns {{ mainFilter: Object, condition: string }}
 */
export function getMainFilterAndCondition(widget, userPreference) {
  const {
    mainFilter: userMainFilter,
    mainFilterCondition: userMainFilterCondition = FILTER_MONGO_OPERATORS.and,
    mainFilterUpdatedAt: userMainFilterUpdatedAt = 0,
  } = userPreference.content;

  const {
    mainFilter: widgetMainFilter,
    mainFilterCondition: widgetMainFilterCondition = FILTER_MONGO_OPERATORS.and,
    mainFilterUpdatedAt: widgetMainFilterUpdatedAt = 0,
  } = widget.parameters;

  let mainFilter = userMainFilter;
  let condition = userMainFilterCondition;

  if (!isEmpty(widgetMainFilter) && widgetMainFilterUpdatedAt >= userMainFilterUpdatedAt) {
    mainFilter = widgetMainFilter;
    condition = widgetMainFilterCondition;
  }

  return { mainFilter, condition };
}
