import { isEmpty, isString } from 'lodash';

import { FILTER_DEFAULT_VALUES } from '@/constants';

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
 * Get main filter for the widget by widget and userPreference parameters
 *
 * @param {Object} widget
 * @param {Object} userPreference
 * @returns {Object|null}
 */
export function getMainFilter(widget, userPreference) {
  const {
    mainFilter: userMainFilter,
    mainFilterUpdatedAt: userMainFilterUpdatedAt = 0,
  } = userPreference.widget_preferences;

  const {
    mainFilter: widgetMainFilter,
    mainFilterUpdatedAt: widgetMainFilterUpdatedAt = 0,
  } = widget.parameters;

  let mainFilter = userMainFilter;

  if (isEmpty(mainFilter) && !isEmpty(widgetMainFilter) && widgetMainFilterUpdatedAt >= userMainFilterUpdatedAt) {
    mainFilter = widgetMainFilter;
  }

  return mainFilter || null;
}

/**
 * Remove spaces from filter string
 *
 * @param {string} filter
 * @return {string}
 */
export function removeSpacesFromStringFilter(filter) {
  if (!isString(filter)) {
    throw new Error('Incorrect using of the function. Filter should be a string.');
  }

  return JSON.stringify(JSON.parse(filter));
}
