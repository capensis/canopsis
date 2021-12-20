import { TIME_UNITS } from '@/constants';

import { toSeconds } from '@/helpers/date/duration';

/**
 * @typedef {Object} SnoozeAction
 * @property {number} duration
 * @property {string} output
 */

/**
 * @typedef {Object} SnoozeActionForm
 * @property {Duration} duration
 * @property {string} output
 */

/**
 * Convert snooze object to form snooze
 *
 * @param {SnoozeAction} snooze
 * @returns {SnoozeActionForm}
 */
export const snoozeToForm = (snooze = {}) => ({
  duration: {
    unit: snooze.duration?.unit ?? TIME_UNITS.minute,
    value: snooze.duration?.seconds ?? 1,
  },
  output: snooze.output ?? '',
});

/**
 * Convert form snooze object to API snooze
 *
 * @param {SnoozeActionForm} form
 * @returns {SnoozeAction}
 */
export const formToSnooze = form => ({
  duration: toSeconds(parseInt(form.duration.value, 10), form.duration.unit),
  output: form.output,
});
