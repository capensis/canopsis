import { DEFAULT_PERIODIC_REFRESH } from '@/constants';

import { durationWithEnabledToForm } from '@/helpers/date/duration';

/**
 * @typedef {Object} AvailabilityWidgetParameters
 * @property {DurationWithEnabled} periodic_refresh
 */

/**
 * @typedef {AvailabilityWidgetParameters} AvailabilityWidgetParametersForm
 */

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
