import { objectToTextPairs, textPairsToObject } from '@/helpers/text-pairs';
import { durationToForm } from '@/helpers/date/duration';

/**
 * @typedef {Duration} RetryDuration
 * @property {number} count
 */

/**
 * @typedef {Object} RequestParameter
 * @property {string} method
 * @property {string} url
 * @property {{ username: string, password: string }} auth
 * @property {Object} headers
 * @property {string} payload
 * @property {boolean} skip_verify
 */

/**
 * @typedef {RequestParameter} RequestFormParameter
 * @property {TextPairObject[]} headers
 */

/**
 * @typedef {Object} RetryParameters
 * @property {number} retry_count
 * @property {Duration} retry_delay
 */

/**
 * Convert request field to form object
 *
 * @param {RequestParameter} request
 * @returns {RequestFormParameter}
 */
export const requestToForm = (request = {}) => ({
  method: request.method ?? '',
  url: request.url ?? '',
  auth: request.auth,
  headers: request.headers ? objectToTextPairs(request.headers) : [],
  payload: request.payload ?? '{}',
  skip_verify: !!request.skip_verify,
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
 * @param {RequestFormParameter} form
 * @returns {RequestParameter}
 */
export const formToRequest = form => ({
  ...form,

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
