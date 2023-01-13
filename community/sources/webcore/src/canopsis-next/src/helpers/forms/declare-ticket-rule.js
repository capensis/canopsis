import { pick } from 'lodash';

import { enabledToForm } from '@/helpers/forms/shared/common';
import { filterPatternsToForm, formFilterToPatterns } from '@/helpers/forms/filter';
import { formToRequest, requestToForm } from '@/helpers/forms/shared/request';
import { objectToTextPairs, textPairsToObject } from '@/helpers/text-pairs';

/**
 * @typedef {Object} DeclareTicketRuleWebhookDeclareTicket
 * @property {string} empty_response
 * @property {string} is_regexp
 * @property {string} ticket_id
 * @property {string} ticket_url
 */

/**
 * @typedef {Object} DeclareTicketRuleWebhook
 * @property {Request} request
 * @property {DeclareTicketRuleWebhookDeclareTicket} declare_ticket
 * @property {boolean} stop_on_fail
 */

/**
 * @typedef {DeclareTicketRuleWebhook[]} DeclareTicketRuleWebhooks
 */

/**
 * @typedef {FilterPatterns} DeclareTicketRule
 * @property {boolean} enabled
 * @property {boolean} emit_trigger
 * @property {string} name
 * @property {string} system_name
 * @property {DeclareTicketRuleWebhooks} webhooks
 */

/**
 * @typedef {Object} DeclareTicketRuleWebhookDeclareTicketForm
 * @property {string} empty_response
 * @property {string} is_regexp
 * @property {TextPairObject[]} mapping
 */

/**
 * @typedef {DeclareTicketRuleWebhook} DeclareTicketRuleWebhookForm
 * @property {DeclareTicketRuleWebhookDeclareTicketForm} declare_ticket
 * @property {RequestForm} request
 */

/**
 * @typedef {DeclareTicketRuleWebhookForm[]} DeclareTicketRuleWebhooksForm
 */

/**
 * @typedef {DeclareTicketRule} DeclareTicketRuleForm
 * @property {DeclareTicketRuleWebhooksForm} webhooks
 * @property {FilterPatternsForm} patterns
 */

export const declareTicketRuleWebhookDeclareTicketToForm = (declareTicket = {}) => {
  const { empty_response: emptyResponse, is_regexp: isRegexp, ...fields } = declareTicket;

  return {
    empty_response: declareTicket.empty_response ?? true,
    is_regexp: declareTicket.is_regexp ?? false,
    mapping: objectToTextPairs(fields),
  };
};

/**
 * Convert declare ticket rule webhook object to form compatible object
 *
 * @param {DeclareTicketRuleWebhook} webhook
 * @returns {DeclareTicketRuleWebhookForm}
 */
export const declareTicketRuleWebhookToForm = (webhook = {}) => ({
  declare_ticket: declareTicketRuleWebhookDeclareTicketToForm(webhook.declare_ticket),
  request: requestToForm(webhook.request),
  stop_on_fail: webhook.stop_on_fail ?? true,
});

/**
 * Convert declare ticket rule webhooks object to form compatible object
 *
 * @param {DeclareTicketRuleWebhooks} webhooks
 * @returns {DeclareTicketRuleWebhooksForm}
 */
export const declareTicketRuleWebhooksToForm = (webhooks = []) => webhooks.map(declareTicketRuleWebhookToForm);

/**
 * Convert declare ticket rule object to form compatible object
 *
 * @param {DeclareTicketRule} [declareTicketRule = {}]
 * @return {DeclareTicketRuleForm}
 */
export const declareTicketRuleToForm = (declareTicketRule = {}) => ({
  enabled: enabledToForm(declareTicketRule.enabled),
  emit_trigger: enabledToForm(declareTicketRule.emit_trigger),
  name: declareTicketRule.name ?? '',
  system_name: declareTicketRule.system_name ?? '',
  webhooks: declareTicketRuleWebhooksToForm(declareTicketRule.webhooks),
  patterns: filterPatternsToForm(declareTicketRule),
});

/**
 * Convert declare ticket rule webhook form to API compatible object
 *
 * @param {DeclareTicketRuleWebhookDeclareTicketForm} form
 * @returns {DeclareTicketRuleWebhookDeclareTicket}
 */
export const formToDeclareTicketRuleWebhookDeclareTicket = form => ({
  ...pick(form, ['empty_response', 'is_regexp']),
  ...textPairsToObject(form.mapping),
});

/**
 * Convert declare ticket rule webhook form to API compatible object
 *
 * @param {DeclareTicketRuleWebhookForm} webhook
 * @returns {DeclareTicketRuleWebhook}
 */
export const formToDeclareTicketRuleWebhook = webhook => ({
  ...webhook,
  declare_ticket: formToDeclareTicketRuleWebhookDeclareTicket(webhook.declare_ticket),
  request: formToRequest(webhook.request),
});

/**
 * Convert declare ticket rule webhooks form to API compatible object
 *
 * @param {DeclareTicketRuleWebhooksForm} webhooks
 * @returns {DeclareTicketRuleWebhooks}
 */
export const formToDeclareTicketRuleWebhooks = (webhooks = []) => webhooks.map(formToDeclareTicketRuleWebhook);

/**
 * Convert form object to declare ticket API compatible object
 *
 * @param {DeclareTicketRuleForm} form
 * @return {DeclareTicketRule}
 */
export const formToDeclareTicketRule = (form) => {
  const { patterns, webhooks, ...declareTicketRule } = form;

  return {
    ...declareTicketRule,
    webhooks: formToDeclareTicketRuleWebhooks(webhooks),
    ...formFilterToPatterns(patterns),
  };
};
