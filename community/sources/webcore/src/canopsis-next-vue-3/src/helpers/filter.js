import { get } from 'lodash';

import { convertDateToString, convertDateToTimezoneDateString } from '@/helpers/date/date';
import { convertDurationToMaxUnitDurationString, convertDurationToString } from '@/helpers/date/duration';
import { stringifyJsonFilter } from '@/helpers/json';
import { convertNumberToFixedString, convertNumberToRoundedPercentString } from '@/helpers/string';

/**
 * Named column display filters, matching the old Vue 2 global filters.
 *
 * @type {Object.<string, Function>}
 */
export const FILTERS = {
  get,
  date: convertDateToString,
  timezone: convertDateToTimezoneDateString,
  duration: convertDurationToString,
  json: stringifyJsonFilter,
  percentage: convertNumberToRoundedPercentString,
  maxDurationByUnit: convertDurationToMaxUnitDurationString,
  fixed: convertNumberToFixedString,
};
