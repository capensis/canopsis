import { TIME_UNITS } from '@/constants';

import { durationToForm } from '@/helpers/date/duration';

/**
 * @typedef {Object} ShareToken
 * @property {string} description
 * @property {Duration} duration
 */

/**
 * @typedef {ShareToken} ShareTokenForm
 * @property {boolean} duration_enabled
 */

/**
 * Convert share token to form
 *
 * @param {ShareToken} [shareToken = {}]
 * @returns {ShareTokenForm}
 */
export const shareTokenToForm = (shareToken = {}) => ({
  description: '',
  duration_enabled: !!shareToken.duration,
  duration: durationToForm(shareToken.duration ?? { value: 1, unit: TIME_UNITS.hour }),
});

/**
 * Convert share token to form
 *
 * @param {ShareTokenForm} [form = {}]
 * @returns {ShareToken}
 */
export const formToShareToken = (form) => {
  const { duration_enabled: durationEnabled, duration, ...shareToken } = form;

  if (durationEnabled) {
    shareToken.duration = duration;
  }

  return shareToken;
};
