import { DEFAULT_WEATHER_LIMIT } from '@/config';

import { convertSortToQuery } from '@/helpers/entities/shared/query';

/**
 * Converts a ServiceWeather widget to a query object.
 *
 * @param {Object} widget
 * @returns {Object}
 */
export const convertWeatherWidgetToQuery = (widget) => {
  const { limit, mainFilter } = widget.parameters;

  return {
    ...convertSortToQuery(widget),
    itemsPerPage: limit || DEFAULT_WEATHER_LIMIT,
    lockedFilter: mainFilter,
    with_tag_colors: true,
  };
};

/**
 * Converts a ServiceWeather user preference to a query object.
 *
 * @param {Object} userPreference
 * @returns {{ category: string, filter: *, hide_grey: boolean }}
 */
export const convertWeatherUserPreferenceToQuery = ({ content }) => {
  const { category, mainFilter, hide_grey: hideGrey = false } = content;

  return { category, filter: mainFilter, hide_grey: hideGrey };
};
