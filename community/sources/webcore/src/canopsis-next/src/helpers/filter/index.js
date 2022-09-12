import { isEmpty } from 'lodash';

/**
 * Get main filter and condition for the widget by widget and userPreference parameters
 *
 * @param {Object} widget
 * @param {Object} userPreference
 * @returns {string}
 */
export const getMainFilter = (widget, userPreference) => {
  const {
    mainFilter: userMainFilter,
    mainFilterUpdatedAt: userMainFilterUpdatedAt = 0,
  } = userPreference.content;

  const {
    mainFilter: widgetMainFilter,
    mainFilterUpdatedAt: widgetMainFilterUpdatedAt = 0,
  } = widget.parameters;

  if (!isEmpty(widgetMainFilter) && widgetMainFilterUpdatedAt >= userMainFilterUpdatedAt) {
    return widgetMainFilter;
  }

  return userMainFilter;
};
