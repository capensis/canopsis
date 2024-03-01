import { AVAILABILITY_SHOW_TYPE, DEFAULT_PERIODIC_REFRESH, QUICK_RANGES } from '@/constants';

import { durationWithEnabledToForm } from '@/helpers/date/duration';

/**
 * @typedef {Object} AvailabilityWidgetParameters
 * @property {DurationWithEnabled} periodic_refresh
 */

/**
 * @typedef {Object} AvailabilityField
 * @property {boolean} enabled
 * @property {number} default_show_type
 * @property {string} default_time_range
 */

/**
 * @typedef {AvailabilityWidgetParameters} AvailabilityWidgetParametersForm
 */

/**
 * Convert availability widget field to form
 *
 * @param {AvailabilityField} [availability = {}]
 * @return {AvailabilityField}
 */
export const availabilityFieldToForm = (availability = {}) => ({
  enabled: availability.enabled ?? false,
  default_show_type: availability.default_show_type ?? AVAILABILITY_SHOW_TYPE.percent,
  default_time_range: availability.default_time_range ?? QUICK_RANGES.today.value,
});

/**
 * Convert form to availability widget parameters to form
 *
 * @param {AvailabilityWidgetParameters} [parameters = {}]
 * @returns {AvailabilityWidgetParametersForm}
 */
export const availabilityWidgetParametersToForm = parameters => ({
  periodic_refresh: durationWithEnabledToForm(parameters.periodic_refresh ?? DEFAULT_PERIODIC_REFRESH),
});

/**
 * Convert form to statistics widget parameters to form
 *
 * @param {AvailabilityWidgetParametersForm} form
 * @returns {StatisticsWidgetParameters}
 */
export const formToAvailabilityWidgetParameters = form => form;
