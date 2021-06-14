import { TIME_UNITS } from '@/constants';

import { toSeconds } from '@/helpers/date/duration';

/**
 * @typedef {Object} SnoozeAction
 * @property {string} duration
 * @property {string} output
 */

/**
 * @typedef {Object} SnoozeActionForm
 * @property {DurationForm} duration
 * @property {string} output
 */

/**
 * Convert snooze object to form snooze
 *
 * @param snooze
 * @returns {SnoozeActionForm}
 */
export const snoozeToForm = (snooze = {}) => ({
  duration: {
    unit: snooze.duration ? snooze.duration.unit : TIME_UNITS.minute,
    value: snooze.duration ? snooze.duration.seconds : 1,
  },
  output: snooze.output || '',
});

/**
 * Convert form snooze object to API snooze
 *
 * @param {Object} form
 * @returns {SnoozeAction}
 */
export const formToSnooze = (form = {}) => ({
  duration: toSeconds(parseInt(form.duration.value, 10), form.duration.unit),
  output: form.output || '',
});
