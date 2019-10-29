import { get, omit } from 'lodash';

import { setSeveralFieldsInObject, unsetSeveralFieldInObjectWithConditions } from '@/helpers/immutable';
import { textPairsToObject, objectToTextPairs } from '@/helpers/text-pairs';

export function webhookToForm(webhook) {
  return {
    emptyResponse: webhook.empty_response || false,
    ...setSeveralFieldsInObject(webhook, getWebhookFormFields(webhook)),
  };
}

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

export function formToWebhook(form) {
  const webhook = createWebhookObject(form);

  return removeEmptyPatternsFromWebhook(webhook);
}

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
    pathValuesMap['request.auth'] = getWebhookAuthField();
  }

  return setSeveralFieldsInObject(omit(form, ['emptyResponse']), pathValuesMap);
}

function getWebhookDeclareTicketField() {
  return value => ({
    ...textPairsToObject(value),
  });
}

function getWebhookAuthField() {
  return auth => ({
    username: auth.username || null,
    password: auth.password || null,
  });
}

function removeEmptyPatternsFromWebhook(webhook) {
  const patternsCondition = value => !value || !value.length;

  return unsetSeveralFieldInObjectWithConditions(webhook, {
    'hook.event_patterns': patternsCondition,
    'hook.alarm_patterns': patternsCondition,
    'hook.entity_patterns': patternsCondition,
  });
}

