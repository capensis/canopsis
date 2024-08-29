import { omit } from 'lodash';

import { WIDGET_TYPES } from '@/constants';

import featuresService from '@/services/features';

import { mapIds } from '@/helpers/array';

import {
  prepareRemediationInstructionsFiltersToQuery,
  getRemediationInstructionsFilters,
} from '../remediation/instruction-filter/query';
import { convertAlarmUserPreferenceToQuery, convertAlarmWidgetToQuery } from '../alarm/query';
import {
  convertChartUserPreferenceToQuery,
  convertChartWidgetToQuery,
  convertNumbersWidgetToQuery,
  convertPieChartWidgetToQuery,
  convertStatisticsUserPreferenceToQuery,
  convertStatisticsWidgetParametersToQuery,
} from '../metric/query';
import { convertContextUserPreferenceToQuery, convertContextWidgetToQuery } from '../entity/query';
import { convertWeatherUserPreferenceToQuery, convertWeatherWidgetToQuery } from '../service-weather/query';
import { convertMapUserPreferenceToQuery, convertMapWidgetToQuery } from '../map/query';
import { convertCounterWidgetToQuery } from '../counter/query';
import { convertStatsCalendarWidgetToQuery } from '../stats/query';
import {
  convertAvailabilityUserPreferenceToQuery,
  convertAvailabilityWidgetParametersToQuery,
} from '../availability/query';

/**
 * This function converts userPreference to query Object
 *
 * @param {Object} userPreference
 * @param {WidgetType} widgetType
 * @returns {Object}
 */
export function convertUserPreferenceToQuery(userPreference, widgetType) {
  const convertersMap = {
    [WIDGET_TYPES.alarmList]: convertAlarmUserPreferenceToQuery,
    [WIDGET_TYPES.context]: convertContextUserPreferenceToQuery,
    [WIDGET_TYPES.serviceWeather]: convertWeatherUserPreferenceToQuery,
    [WIDGET_TYPES.map]: convertMapUserPreferenceToQuery,
    [WIDGET_TYPES.barChart]: convertChartUserPreferenceToQuery,
    [WIDGET_TYPES.lineChart]: convertChartUserPreferenceToQuery,
    [WIDGET_TYPES.pieChart]: convertChartUserPreferenceToQuery,
    [WIDGET_TYPES.numbers]: convertChartUserPreferenceToQuery,
    [WIDGET_TYPES.userStatistics]: convertStatisticsUserPreferenceToQuery,
    [WIDGET_TYPES.alarmStatistics]: convertStatisticsUserPreferenceToQuery,
    [WIDGET_TYPES.availability]: convertAvailabilityUserPreferenceToQuery,

    ...featuresService.get('helpers.query.convertUserPreferenceToQuery.convertersMap'),
  };

  const converter = convertersMap[widgetType];

  return converter ? converter(userPreference) : {};
}

/**
 * This function converts widget to query Object
 *
 * @param {Object} widget
 * @returns {{}}
 */
export function convertWidgetToQuery(widget) {
  const convertersMap = {
    [WIDGET_TYPES.alarmList]: convertAlarmWidgetToQuery,
    [WIDGET_TYPES.context]: convertContextWidgetToQuery,
    [WIDGET_TYPES.serviceWeather]: convertWeatherWidgetToQuery,
    [WIDGET_TYPES.map]: convertMapWidgetToQuery,
    [WIDGET_TYPES.statsCalendar]: convertStatsCalendarWidgetToQuery,
    [WIDGET_TYPES.counter]: convertCounterWidgetToQuery,
    [WIDGET_TYPES.barChart]: convertChartWidgetToQuery,
    [WIDGET_TYPES.lineChart]: convertChartWidgetToQuery,
    [WIDGET_TYPES.pieChart]: convertPieChartWidgetToQuery,
    [WIDGET_TYPES.numbers]: convertNumbersWidgetToQuery,
    [WIDGET_TYPES.userStatistics]: convertStatisticsWidgetParametersToQuery,
    [WIDGET_TYPES.alarmStatistics]: convertStatisticsWidgetParametersToQuery,
    [WIDGET_TYPES.availability]: convertAvailabilityWidgetParametersToQuery,

    ...featuresService.get('helpers.query.convertWidgetToQuery.convertersMap'),
  };

  const converter = convertersMap[widget.type];

  return converter ? converter(widget) : {};
}

/**
 * Prepare query by widget and userPreference objects
 *
 * @param {Object} widget
 * @param {Object} userPreference
 * @returns {Object}
 */
export function prepareWidgetQuery(widget, userPreference) {
  const { filters: widgetFilters = [] } = widget;
  const { filters: userPreferenceFilters } = userPreference;
  const widgetQuery = convertWidgetToQuery(widget);
  const userPreferenceQuery = convertUserPreferenceToQuery(userPreference, widget.type);

  let query = {
    ...widgetQuery,
    ...userPreferenceQuery,
  };

  if (query.filter) {
    const allFiltersIds = mapIds([...widgetFilters, ...userPreferenceFilters]);

    if (!allFiltersIds.includes(query.filter)) {
      query = omit(query, ['filter']);
    }
  }

  const remediationInstructionsFilters = getRemediationInstructionsFilters(widget, userPreference);

  if (remediationInstructionsFilters.length) {
    query = {
      ...query,
      ...prepareRemediationInstructionsFiltersToQuery(remediationInstructionsFilters),
    };
  }

  return query;
}
