import { TIME_UNITS } from '@/constants';

import { durationToForm } from '@/helpers/date/duration';

/**
 * Prepare the public API token creation form.
 *
 * @returns {Object}
 */
export const apiTokenToForm = () => ({
  name: '',
  duration_enabled: false,
  duration: durationToForm({ value: 1, unit: TIME_UNITS.hour }),
});

/**
 * Convert the public API token form to the create request payload.
 *
 * @param {Object} form
 * @returns {Object}
 */
export const formToApiToken = (form) => {
  const { duration_enabled: durationEnabled, duration, ...apiToken } = form;

  if (durationEnabled) {
    apiToken.duration = duration;
  }

  return apiToken;
};
