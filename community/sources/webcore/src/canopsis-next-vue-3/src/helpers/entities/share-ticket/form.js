import { TIME_UNITS } from '@/constants';

import { durationToForm } from '@/helpers/date/duration';

/**
 * Convert share token to form
 *
 * @param {ShareToken} [shareToken = {}]
 * @returns {Object}
 */
export const shareTokenToForm = (shareToken = {}) => ({
  description: shareToken.description ?? '',
  duration_enabled: !!shareToken.duration,
  duration: durationToForm(shareToken.duration ?? { value: 1, unit: TIME_UNITS.hour }),
});

/**
 * Convert share token form to request payload
 *
 * @param {Object} form
 * @returns {Object}
 */
export const formToShareToken = (form) => {
  const { duration_enabled: durationEnabled, duration, ...shareToken } = form;

  if (durationEnabled) {
    shareToken.duration = duration;
  }

  return shareToken;
};
