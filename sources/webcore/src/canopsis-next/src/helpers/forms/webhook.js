import { get, isUndefined, omit } from 'lodash';

import { setSeveralFields, unsetSeveralFieldsWithConditions } from '@/helpers/immutable';
import { textPairsToObject, objectToTextPairs } from '@/helpers/text-pairs';
import { getConditionsForRemovingEmptyPatterns } from '@/helpers/forms/shared/patterns';

export function webhookToForm(webhook = {}) {
  const declareTicketField = webhook.declare_ticket ? omit(webhook.declare_ticket, ['empty_response']) : {};

  return {
    _id: webhook._id,
    retry: webhook.retry || {},
    hook: {
      triggers: get(webhook, 'hook.triggers', []),
      event_patterns: get(webhook, 'hook.event_patterns', []),
      alarm_patterns: get(webhook, 'hook.alarm_patterns', []),
      entity_patterns: get(webhook, 'hook.entity_patterns', []),
    },
    request: {
      method: get(webhook, 'request.method', ''),
      url: get(webhook, 'request.url', ''),
      headers: webhook.request && webhook.request.headers
        ? objectToTextPairs(webhook.request.headers)
        : [],
      payload: webhook.request && webhook.request.payload
        ? webhook.request.payload
        : '{}',
    },
    declare_ticket: objectToTextPairs(declareTicketField),
    disable_during_periods: webhook.disable_during_periods || [],

    emptyResponse: !!webhook.empty_response,
    enabled: !isUndefined(webhook.enabled) ? webhook.enabled : true,
  };
}

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
  });
}

