import { DEFAULT_CLOCK_WIDGET_PARAMETERS } from '@/constants';

import { addKeyInEntities, removeKeyFromEntities } from '@/helpers/array';
import { getLocalTimezone } from '@/helpers/date/date';

/**
 * @typedef {Object} ClockWidgetClock
 * @property {string} timezone
 * @property {string} label
 */

/**
 * @typedef {ClockWidgetClock & ObjectKey} ClockWidgetClockForm
 */

/**
 * @typedef {Object} ClockWidgetParameters
 * @property {'digital' | 'seven_segments' | 'analog'} display
 * @property {'s' | 'm' | 'l'} size
 * @property {'automatic' | 'fixed'} spacing
 * @property {boolean} is_12_hour_format
 * @property {boolean} show_seconds
 * @property {boolean} show_date
 * @property {boolean} show_utc_offset
 * @property {boolean} show_timezone
 * @property {ClockWidgetClock[]} clocks
 */

/**
 * @typedef {ClockWidgetParameters} ClockWidgetParametersForm
 * @property {ClockWidgetClockForm[]} clocks
 */

/**
 * Create the initial clock shown by a new widget.
 *
 * @return {ClockWidgetClock}
 */
export const сlockWidgetClockToForm = () => ({
  timezone: getLocalTimezone(),
  label: '',
});

/**
 * Convert Clock widget parameters to form values.
 *
 * @param {ClockWidgetParameters} [parameters={}]
 * @return {ClockWidgetParametersForm}
 */
export const clockWidgetParametersToForm = (parameters = {}) => ({
  ...DEFAULT_CLOCK_WIDGET_PARAMETERS,
  ...parameters,
  clocks: addKeyInEntities(
    parameters.clocks?.length ? parameters.clocks : [сlockWidgetClockToForm()],
  ),
});

/**
 * Convert Clock widget form values to API parameters.
 *
 * @param {ClockWidgetParametersForm} form
 * @return {ClockWidgetParameters}
 */
export const formToClockWidgetParameters = form => ({
  ...form,
  clocks: removeKeyFromEntities(form.clocks).map(clock => ({
    ...clock,
    label: clock.label?.trim() ?? '',
  })),
});
