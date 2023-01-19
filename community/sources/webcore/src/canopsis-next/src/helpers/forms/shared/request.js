import { isNumber, pick } from 'lodash';

import { objectToTextPairs, textPairsToObject } from '@/helpers/text-pairs';
import { durationToForm, isValidDuration } from '@/helpers/date/duration';

/**
 * @typedef {Object} RequestAuth
 * @property {string} username
 * @property {string} password
 */

/**
 * @typedef {Object} Request
 * @property {string} method
 * @property {string} url
 * @property {RequestAuth} auth
 * @property {Object} headers
 * @property {boolean} skip_verify
 * @property {Duration} [timeout]
 * @property {number} [retry_count]
 * @property {?Duration} [retry_delay]
 * @property {boolean} [empty_response]
 * @property {string} payload
 */

/**
 * @typedef {RequestAuth} RequestAuthForm
 * @property {boolean} enabled
 */

/**
 * @typedef {Request} RequestForm
 * @property {TextPairObject[]} headers
 * @property {RequestAuthForm} auth
 */

/**
 * @typedef {Object} RetryParameters
 * @property {number} retry_count
 * @property {Duration} retry_delay
 */

/**
 * Convert request field to form object
 *
 * @param {Request} request
 * @returns {RequestForm}
 */
export const requestToForm = (request = {}) => ({
  method: request.method ?? '',
  url: request.url ?? '',
  skip_verify: !!request.skip_verify,
  empty_response: !!request.empty_response,
  timeout: request.timeout
    ? durationToForm(request.timeout)
    : { value: undefined, unit: undefined },
  retry_count: request.retry_count,
  retry_delay: request.retry_delay
    ? durationToForm(request.retry_delay)
    : { value: undefined, unit: undefined },
  auth: request.auth
    ? { enabled: true, ...request.auth }
    : { enabled: false, username: '', password: '' },
  headers: request.headers ? objectToTextPairs(request.headers) : [],
  payload: request.payload ?? '',
});

/**
 * Convert retry parameters to form object
 *
 * @param {RetryParameters} parameters
 * @returns {RetryDuration}
 */
export const retryToForm = (parameters = {}) => (
  parameters.retry_delay
    ? { count: parameters.retry_count, ...durationToForm(parameters.retry_delay) }
    : { count: '', unit: '', value: '' }
);

/**
 * Convert form object to request field
 *
 * @param {RequestForm} form
 * @returns {Request}
 */
export const formToRequest = form => ({
  ...form,

  retry_delay: isValidDuration(form.retry_delay)
    ? form.retry_delay
    : undefined,
  timeout: isNumber(form.timeout.value) ? form.timeout : null,
  auth: form.auth.enabled ? pick(form.auth, ['username', 'password']) : null,
  headers: textPairsToObject(form.headers),
});

/**
 * Convert form object to retry parameters
 *
 * @param {RetryDuration} parameters
 * @returns {RetryParameters}
 */
export const formToRetry = ({ value, unit, count }) => (
  value ? { retry_count: count, retry_delay: { value, unit } } : {}
);
