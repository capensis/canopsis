import { DEFAULT_WEATHER_LIMIT } from '@/config';

import { convertSortToQuery } from '@/helpers/entities/shared/query';

/**
 * This function converts widget with type 'ServiceWeather' to query Object
 *
 * @param {Object} widget
 * @returns {{}}
 */
export function convertWeatherWidgetToQuery(widget) {
  const { limit, mainFilter } = widget.parameters;

  return {
    ...convertSortToQuery(widget),
    limit: limit || DEFAULT_WEATHER_LIMIT,
    lockedFilter: mainFilter,
    define_state: true,
  };
}

/**
 * This function converts userPreference with widget type 'ServiceWeather' to query Object
 *
 * @param {Object} userPreference
 * @returns {{ category: string }}
 */
export function convertWeatherUserPreferenceToQuery({ content }) {
  const { category, mainFilter, hide_grey: hideGrey = false } = content;

  return { category, filter: mainFilter, hide_grey: hideGrey };
}
