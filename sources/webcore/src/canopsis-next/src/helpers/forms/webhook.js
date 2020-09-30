import { get, omit } from 'lodash';

import { setSeveralFields, unsetSeveralFieldsWithConditions } from '@/helpers/immutable';
import { textPairsToObject, objectToTextPairs } from '@/helpers/text-pairs';
import { getConditionsForRemovingEmptyPatterns } from '@/helpers/forms/shared/patterns';

/**
 * Get webhook form field's values (or customizer function)
 *
 * @param {Object} webhook
 * @returns {Object}
 */
function getWebhookFormFields(webhook) {
  const patternsFieldsCustomizer = value => value || [];

  const declareTicketField = webhook.declare_ticket ? omit(webhook.declare_ticket, ['empty_response']) : {};
  const pathValuesMap = {
    declare_ticket: () => objectToTextPairs(declareTicketField),
    'request.headers': objectToTextPairs,
    'hook.event_patterns': patternsFieldsCustomizer,
    'hook.alarm_patterns': patternsFieldsCustomizer,
    'hook.entity_patterns': patternsFieldsCustomizer,
  };

  if (webhook.combine_meta_alarm_request) {
    pathValuesMap['combine_meta_alarm_request.headers'] = objectToTextPairs;
  }

  return pathValuesMap;
}

export function webhookToForm(webhook) {
  return {
    emptyResponse: webhook.empty_response || false,
    enabled: webhook.enabled || true,
    ...setSeveralFields(webhook, getWebhookFormFields(webhook)),
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
 * Get fields with customizers for request subform
 *
 * @param {Object} form
 * @param {string} [pathPrefix = 'request']
 * @returns {Object}
 */
function getRequestFormFields(form, pathPrefix = 'request') {
  const authPath = `${pathPrefix}.auth`;
  const auth = get(form, authPath);

  const pathValuesMap = {
    [`${pathPrefix}.headers`]: textPairsToObject,
  };

  if (auth) {
    pathValuesMap[authPath] = {
      username: auth.username,
      password: auth.password,
    };
  }

  return pathValuesMap;
}

/**
 * Create a webhook object that is valid to the API
 *
 * @param {Object} [form = {}]
 * @returns {Object}
 */
function createWebhookObjectFromForm(form = {}) {
  const pathValuesMap = {
    ...getRequestFormFields(form),
    ...(form.combine_meta_alarm_request && getRequestFormFields(form, 'combine_meta_alarm_request')),

    empty_response: form.emptyResponse || false,
  };

  if (form.declare_ticket) {
    pathValuesMap.declare_ticket = getWebhookDeclareTicketField(form);
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

  return unsetSeveralFieldsWithConditions(createWebhookObjectFromForm(form), {
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

