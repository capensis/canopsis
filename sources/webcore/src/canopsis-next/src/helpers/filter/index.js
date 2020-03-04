import { get, isEmpty, isString } from 'lodash';

import { FILTER_DEFAULT_VALUES } from '@/constants';

/**
 * Prepare user main filter to query filter
 *
 * @param {Object|Array} [filterObject = {}]
 * @param {string} [condition = FILTER_DEFAULT_VALUES.condition]
 * @returns {Object|undefined}
 */
export function prepareMainFilterToQueryFilter(
  filterObject = {},
  condition = FILTER_DEFAULT_VALUES.condition,
) {
  let filter;

  if (Array.isArray(filterObject) && filterObject.length) {
    filter = {
      [condition]: filterObject.map((item) => {
        if (isString(item.filter)) {
          return JSON.parse(item.filter);
        }

        return item.filter || {};
      }),
    };
  } else if (!Array.isArray(filterObject)) {
    filter = get(filterObject, 'filter');
  }

  return filter;
}

/**
 * Get main filter for the widget by widget and userPreference parameters
 *
 * @param {Object} widget
 * @param {Object} userPreference
 * @returns {Object|null}
 */
export function getWidgetMainFilter(widget, userPreference) {
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

