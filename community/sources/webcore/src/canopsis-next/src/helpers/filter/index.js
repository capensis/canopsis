import { isEmpty, isString } from 'lodash';

import { FILTER_DEFAULT_VALUES, FILTER_MONGO_OPERATORS } from '@/constants';

/**
 * @typedef {Object} FilterWithTitle
 * @property {string} title
 * @property {Object} filter
 */

/**
 * @typedef {FilterWithTitle | Object} Filter
 */

/**
 * Prepare user main filter to query filter
 *
 * @param {Object|Array} [filterObject = {}]
 * @param {string} [condition = FILTER_DEFAULT_VALUES.condition]
 * @returns {Object|string|undefined}
 */
export function prepareMainFilterToQueryFilter(
  filterObject = {},
  condition = FILTER_DEFAULT_VALUES.condition,
) {
  let result;

  if (Array.isArray(filterObject) && filterObject.length) {
    result = {
      [condition]: filterObject.map((item) => {
        if (isString(item.filter)) {
          return JSON.parse(item.filter);
        }

        return item.filter || {};
      }),
    };
  } else if (!Array.isArray(filterObject)) {
    const { filter } = filterObject;

    result = isString(filter) ? JSON.parse(filter) : filter;
  }

  return result;
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
