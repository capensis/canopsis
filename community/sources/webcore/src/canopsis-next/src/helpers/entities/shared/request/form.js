import { isNumber, pick } from 'lodash';

import {
  REQUEST_AUTH_TYPES,
  REQUEST_AUTH_TOKEN_TYPES,
  REQUEST_AUTH_TOKEN_TYPES_TO_HEADERS,
  HEADERS_TO_REQUEST_AUTH_TOKEN_TYPES,
} from '@/constants';

import { objectToTextPairs, textPairsToObject } from '@/helpers/text-pairs';
import { durationToForm, isValidDuration } from '@/helpers/date/duration';

/**
 * @typedef { 'none' | 'credentials' | 'token' } RequestAuthType
 */

/**
 * @typedef {Object} RequestAuth
 * @property {string} username
 * @property {string} password
 */

/**
 * @typedef {Object} RequestAuthToken
 * @property {string} rule
 * @property {string} [header]
 * @property {string} [query_string]
 */

/**
 * @typedef {Object} Request
 * @property {string} method
 * @property {string} url
 * @property {Object} headers
 * @property {boolean} skip_verify
 * @property {RequestAuth} [auth]
 * @property {Duration} [timeout]
 * @property {number} [retry_count]
 * @property {?Duration} [retry_delay]
 * @property {string} payload
 */

/**
 * @typedef {RequestAuth} RequestAuthForm
 * @property {RequestAuthType} type
 * @property {string} username
 * @property {string} password
 */

/**
 * @typedef {Object} RequestAuthTokenForm
 * @property {string} type
 * @property {string} parameter
 * @property {string} rule
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
 * Convert request auth token entity to request auth token form
 *
 * @param {RequestAuthToken} requestAuthToken
 * @returns {RequestAuthTokenForm}
 */
export const requestAuthTokenToForm = (requestAuthToken = {}) => {
  const form = {
    type: requestAuthToken.type ?? '',
    parameter: requestAuthToken.parameter ?? '',
    rule: requestAuthToken.rule ?? '',
  };

  if (requestAuthToken.query_string) {
    form.type = REQUEST_AUTH_TOKEN_TYPES.url;
    form.parameter = requestAuthToken.query_string;
  } else if (requestAuthToken.payload_field) {
    form.type = REQUEST_AUTH_TOKEN_TYPES.payload;
    form.parameter = requestAuthToken.payload_field;
  } else if (requestAuthToken.header) {
    const type = HEADERS_TO_REQUEST_AUTH_TOKEN_TYPES[requestAuthToken.header];

    if (type) {
      form.type = type;
    } else {
      form.type = REQUEST_AUTH_TOKEN_TYPES.headerCustomParameter;
      form.parameter = requestAuthToken.header;
    }
  }

  return form;
};

/**
 * Convert request field to form object
 *
 * @param {Request} request
 * @param {RequestAuthToken} [authToken = {}]
 * @returns {RequestForm}
 */
export const requestToForm = (request = {}, authToken = {}) => {
  const form = {
    method: request.method ?? '',
    url: request.url ?? '',
    skip_verify: !!request.skip_verify,
    timeout: request.timeout
      ? durationToForm(request.timeout)
      : { value: undefined, unit: undefined },
    retry_count: request.retry_count,
    retry_delay: request.retry_delay
      ? durationToForm(request.retry_delay)
      : { value: undefined, unit: undefined },
    auth: {
      type: REQUEST_AUTH_TYPES.none,
      username: '',
      password: '',
    },
    headers: request.headers ? objectToTextPairs(request.headers) : [],
    payload: request.payload ?? '',
  };

  if (request.auth?.username) {
    form.auth.type = REQUEST_AUTH_TYPES.credentials;
    form.auth.username = request.auth.username;
    form.auth.password = request.auth.password;
  } else if (authToken.rule) {
    form.auth.type = REQUEST_AUTH_TYPES.token;
  }

  return form;
};

/**
 * Convert request auth token form to request auth token entity
 *
 * @param {RequestAuthTokenForm} form
 * @param {RequestAuthType} requestAuthType
 * @returns {RequestAuthToken}
 */
export const formToRequestAuthToken = (form = {}, requestAuthType = REQUEST_AUTH_TYPES.none) => {
  if (requestAuthType !== REQUEST_AUTH_TYPES.token) {
    return null;
  }

  const result = {
    rule: form.rule,
  };

  const header = REQUEST_AUTH_TOKEN_TYPES_TO_HEADERS[form.type];

  if (header) {
    result.header = header;
  } else if (form.type === REQUEST_AUTH_TOKEN_TYPES.headerCustomParameter) {
    result.header = form.parameter;
  } else if (form.type === REQUEST_AUTH_TOKEN_TYPES.url) {
    result.query_string = form.parameter;
  } else if (form.type === REQUEST_AUTH_TOKEN_TYPES.payload) {
    result.payload_field = form.parameter;
  }

  return result;
};

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
  auth: form.auth.type === REQUEST_AUTH_TYPES.credentials && form.auth?.username
    ? pick(form.auth, ['username', 'password'])
    : null,
  headers: textPairsToObject(form.headers),
});

/**
 * Convert error structure to form structure
 *
 * @param {Object[]} headersErrors
 * @param {Object[]} headers
 * @return {FlattenErrors}
 */
export const requestHeadersTemplateVariablesErrorsToForm = (
  headersErrors,
  headers,
) => headersErrors.reduce((acc, { is_valid: isValid, err }, index) => {
  const header = headers[index];

  if (!isValid) {
    acc[header.key] = {
      value: err.message,
    };
  }

  return acc;
}, {});

/**
 * Convert template variables errors structure to form structure
 *
 * @param {Object} errorsObject
 * @param {Object} form
 * @return {FlattenErrors}
 */
export const requestTemplateVariablesErrorsToForm = (errorsObject, form) => {
  const { url, payload, headers } = errorsObject;

  const requestErrors = {};

  if (url && !url.is_valid) {
    requestErrors.url = url.err.message;
  }

  if (payload && !payload.is_valid) {
    requestErrors.payload = `${payload.err.line}|${payload.err.message}`;
  }

  if (headers?.some(({ is_valid: isValid }) => !isValid)) {
    requestErrors.headers = requestHeadersTemplateVariablesErrorsToForm(headers, form.headers);
  }

  return requestErrors;
};
