import { get, omit } from 'lodash';

import { setSeveralFields, unsetSeveralFieldsWithConditions } from '@/helpers/immutable';
import { textPairsToObject, objectToTextPairs } from '@/helpers/text-pairs';
import { getSecondByUnit, formatSecondByUnit } from '@/helpers/time';

/**
 * Get webhook form field's values (or customizer function)
 * @param {Object} webhook
 * @returns {Object}
 */
function getWebhookFormFields(webhook) {
  const patternsFieldsCustomizer = value => value || [];

  const declareTicketField = webhook.declare_ticket ? omit(webhook.declare_ticket, ['empty_response']) : {};

  return {
    declare_ticket: () => objectToTextPairs(declareTicketField),
    'request.headers': objectToTextPairs,
    'hook.event_patterns': patternsFieldsCustomizer,
    'hook.alarm_patterns': patternsFieldsCustomizer,
    'hook.entity_patterns': patternsFieldsCustomizer,
  };
}

export function webhookToForm(webhook) {
  return {
    emptyResponse: webhook.empty_response || false,
    ...setSeveralFields(webhook, getWebhookFormFields(webhook)),
    retry: webhook.retry && {
      ...webhook.retry,
      delay: formatSecondByUnit(webhook.retry.delay, webhook.retry.unit),
    },
  };
}

/**
 * Tranform webhook declare ticket field to object (editable in the UI)
 * @returns {Function}
 */
function getWebhookDeclareTicketField() {
  return value => ({
    ...textPairsToObject(value),
  });
}

/**
 * Get webhook's auth fields values
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

  const retry = {
    ...form.retry,
    delay: getSecondByUnit(form.retry.delay, form.retry.unit),
  };

  const formData = {
    ...omit(form, ['emptyResponse']),
    retry,
  };

  return setSeveralFields(formData, pathValuesMap);
}

/**
 * Remove empty "patterns" (alarmpattern, entitypattern and eventpattern) fields from webhook
 * @param {Object} webhook
 * @returns {Object}
 */
function removeEmptyPatternsFromWebhook(webhook) {
  const patternsCondition = value => !value || !value.length;

  return unsetSeveralFieldsWithConditions(webhook, {
    'hook.event_patterns': patternsCondition,
    'hook.alarm_patterns': patternsCondition,
    'hook.entity_patterns': patternsCondition,
  });
}

/**
 * Transform webhook "form" object to valid webhook to the API
 * @param {Object} form
 * @returns {Object}
 */
export function formToWebhook(form) {
  const webhook = createWebhookObject(form);

  return removeEmptyPatternsFromWebhook(webhook);
}

