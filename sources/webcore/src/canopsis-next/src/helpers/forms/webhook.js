import { get, isUndefined, omit } from 'lodash';

import { setSeveralFields, unsetSeveralFieldsWithConditions } from '@/helpers/immutable';
import { textPairsToObject, objectToTextPairs } from '@/helpers/text-pairs';
import { getConditionsForRemovingEmptyPatterns } from '@/helpers/forms/shared/patterns';

/**
 * @typedef {Object} WebhookRequest
 * @property {string} method
 * @property {string} url
 * @property {Array|Object} headers
 * @property {string} payload
 */

/**
 * @typedef {Object} WebhookHook
 * @property {Array} triggers
 * @property {Array} event_patterns
 * @property {Array} alarm_patterns
 * @property {Array} entity_patterns
 */

/**
 * @typedef {Object} Webhook
 * @property {Array} _id
 * @property {Object} retry
 * @property {Object} hook
 * @property {Object} request
 * @property {Array|Object} declare_ticket
 * @property {Array} disable_during_periods
 * @property {boolean} [emptyResponse]
 * @property {boolean} [empty_response]
 * @property {boolean} enabled
 */

/**
 * Convert webhook request field object to webhook request form object
 *
 * @param {WebhookRequest} request
 * @returns {WebhookRequest}
 */
const webhookRequestToForm = (request = {}) => ({
  method: request.method || '',
  url: request.url || '',
  headers: request.headers
    ? objectToTextPairs(request.headers)
    : [],
  payload: request.payload || '{}',
});

/**
 * Convert webhook hook field object to webhook hook form object
 *
 * @param {WebhookHook} hook
 * @returns {WebhookHook}
 */
const webhookHookToForm = (hook = {}) => ({
  triggers: hook.triggers || [],
  event_patterns: hook.event_patterns || [],
  alarm_patterns: hook.alarm_patterns || [],
  entity_patterns: hook.entity_patterns || [],
});

/**
 * Convert webhook object to webhook form object
 *
 * @param {Webhook} webhook
 * @returns {Webhook}
 */
export const webhookToForm = (webhook = {}) => {
  const declareTicketField = webhook.declare_ticket ? omit(webhook.declare_ticket, ['empty_response']) : {};

  return {
    _id: webhook._id,
    retry: webhook.retry || {},
    hook: webhookHookToForm(webhook.hook),
    request: webhookRequestToForm(webhook.request),
    declare_ticket: objectToTextPairs(declareTicketField),
    disable_during_periods: webhook.disable_during_periods || [],

    emptyResponse: !!webhook.empty_response,
    enabled: !isUndefined(webhook.enabled) ? webhook.enabled : true,
  };
};

/**
 * Tranform webhook declare ticket field to object (editable in the UI)
 *
 * @returns {Function}
 */
function getWebhookDeclareTicketField() {
  return value => ({
    ...textPairsToObject(value),
  });
}

/**
 * Get webhook's auth fields values
 *
 * @param {Object} form
 * @returns {Object}
 */
function getWebhookAuthField(form) {
  return {
    username: form.request.auth.username,
    password: form.request.auth.password,
  };
}

/**
 * Create a webhook object that is valid to the API
 *
 * @param {Object} form
 * @returns {Object}
 */
function createWebhookObject(form) {
  const hasAuth = get(form, 'request.auth');

  const pathValuesMap = {
    'request.headers': textPairsToObject,
    empty_response: form.emptyResponse || false,
  };

  if (form.declare_ticket) {
    pathValuesMap.declare_ticket = getWebhookDeclareTicketField(form);
  }

  if (hasAuth) {
    pathValuesMap['request.auth'] = getWebhookAuthField(form);
  }

  return setSeveralFields(omit(form, ['emptyResponse']), pathValuesMap);
}

/**
 * Transform webhook "form" object to valid webhook to the API
 *
 * @param {Object} form
 * @returns {Object}
 */
export function formToWebhook(form) {
  const hasValue = v => !v;

  return unsetSeveralFieldsWithConditions(createWebhookObject(form), {
    ...getConditionsForRemovingEmptyPatterns([
      'hook.alarm_patterns',
      'hook.entity_patterns',
      'hook.event_patterns',
    ]),

    'retry.count': hasValue,
    'retry.unit': hasValue,
    'retry.delay': hasValue,
    _id: hasValue,
  });
}

