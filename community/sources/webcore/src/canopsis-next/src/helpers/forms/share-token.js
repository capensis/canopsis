import { durationToForm } from '@/helpers/date/duration';
import { TIME_UNITS } from '@/constants';

/**
 * @typedef {Object} ShareToken
 * @property {string} description
 * @property {Duration} duration
 */

/**
 * @typedef {ShareToken} ShareTokenForm
 */

/**
 * Convert share token to form
 *
 * @param {ShareToken} [shareToken = {}]
 * @returns {ShareTokenForm}
 */
export const shareTokenToForm = (shareToken = {}) => ({
  description: '',
  duration: durationToForm(shareToken.duration ?? { value: 1, unit: TIME_UNITS.hour }),
});
