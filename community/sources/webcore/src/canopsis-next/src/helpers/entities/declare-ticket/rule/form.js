import flatten from 'flat';

import { DECLARE_TICKET_EXECUTION_STATUSES } from '@/constants';

import {
  formToRequest,
  requestTemplateVariablesErrorsToForm,
  requestToForm,
} from '@/helpers/entities/shared/request/form';
import { filterPatternsToForm, formFilterToPatterns } from '@/helpers/entities/filter/form';
import { objectToTextPairs, textPairsToObject } from '@/helpers/text-pairs';
import { removeKeyFromEntities } from '@/helpers/array';
import { flattenErrorMap } from '@/helpers/entities/shared/form';
import { uid } from '@/helpers/uid';

/**
 * @typedef {Object} DeclareTicketRuleWebhookDeclareTicket
 * @property {string} empty_response
 * @property {string} is_regexp
 * @property {string} ticket_id
 * @property {string} [ticket_id_tpl]
 * @property {string} [ticket_url]
 * @property {string} [ticket_url_tpl]
 */

/**
 * @typedef {Object} DeclareTicketRuleWebhook
 * @property {Request} request
 * @property {?DeclareTicketRuleWebhookDeclareTicket} declare_ticket
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
 * @typedef {Object} DeclareTicketRuleWebhookTickerUrlForm
 * @property {boolean} template
 * @property {string} value
 */

/**
 * @typedef {Object} DeclareTicketRuleWebhookTickerIdForm
 * @property {boolean} template
 * @property {string} value
 */

/**
 * @typedef {DeclareTicketRuleWebhookDeclareTicket} DeclareTicketRuleWebhookDeclareTicketForm
 * @property {boolean} enabled
 * @property {TextPairObject[]} mapping
 * @property {DeclareTicketRuleWebhookTickerIdForm} ticket_id
 * @property {DeclareTicketRuleWebhookTickerUrlForm} ticket_url
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

/**
 * Check declare ticket execution status is waiting
 *
 * @param {number} status
 * @returns {boolean}
 */
export const isDeclareTicketExecutionWaiting = ({ status }) => status === DECLARE_TICKET_EXECUTION_STATUSES.waiting;

/**
 * Check declare ticket execution status is running
 *
 * @param {number} status
 * @returns {boolean}
 */
export const isDeclareTicketExecutionRunning = ({ status }) => status === DECLARE_TICKET_EXECUTION_STATUSES.running;

/**
 * Check declare ticket execution status is succeeded
 *
 * @param {number} status
 * @returns {boolean}
 */
export const isDeclareTicketExecutionSucceeded = ({ status }) => status === DECLARE_TICKET_EXECUTION_STATUSES.succeeded;

/**
 * Check declare ticket execution status is failed
 *
 * @param {number} status
 * @returns {boolean}
 */
export const isDeclareTicketExecutionFailed = ({ status }) => status === DECLARE_TICKET_EXECUTION_STATUSES.failed;

/**
 * Convert declare ticket object to form compatible object
 *
 * @param {DeclareTicketRuleWebhookDeclareTicket} declareTicket
 * @returns {DeclareTicketRuleWebhookDeclareTicketForm}
 */
export const declareTicketRuleWebhookDeclareTicketToForm = (declareTicket) => {
  const {
    empty_response: emptyResponse,
    is_regexp: isRegexp,
    ticket_id: ticketId = '',
    ticket_id_tpl: ticketIdTpl = '',
    ticket_url: ticketUrl = '',
    ticket_url_tpl: ticketUrlTpl = '',
    ...fields
  } = declareTicket ?? {};

  return {
    enabled: !!declareTicket,
    empty_response: emptyResponse ?? false,
    is_regexp: isRegexp ?? false,
    ticket_id: {
      template: !!ticketIdTpl,
      value: ticketIdTpl || ticketId,
    },
    ticket_url: {
      template: !!ticketUrlTpl,
      value: ticketUrlTpl || ticketUrl,
    },
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
  key: uid(),
  declare_ticket: declareTicketRuleWebhookDeclareTicketToForm(webhook.declare_ticket),
  request: requestToForm(webhook.request),
  stop_on_fail: webhook.stop_on_fail ?? false,
});

/**
 * Convert declare ticket rule webhooks object to form compatible object
 *
 * @param {DeclareTicketRuleWebhooks} webhooks
 * @returns {DeclareTicketRuleWebhooksForm}
 */
export const declareTicketRuleWebhooksToForm = (webhooks = [undefined]) => webhooks.map(declareTicketRuleWebhookToForm);

/**
 * Convert declare ticket rule object to form compatible object
 *
 * @param {DeclareTicketRule} [declareTicketRule = {}]
 * @return {DeclareTicketRuleForm}
 */
export const declareTicketRuleToForm = (declareTicketRule = {}) => ({
  enabled: declareTicketRule.enabled ?? true,
  emit_trigger: declareTicketRule.emit_trigger ?? true,
  name: declareTicketRule.name ?? '',
  system_name: declareTicketRule.system_name ?? '',
  webhooks: declareTicketRuleWebhooksToForm(declareTicketRule.webhooks),
  patterns: filterPatternsToForm(declareTicketRule),
});

/**
 * Convert declare ticket rule webhook form to API compatible object
 *
 * @param {DeclareTicketRuleWebhookDeclareTicketForm} form
 * @returns {DeclareTicketRuleWebhookDeclareTicket | null}
 */
export const formToDeclareTicketRuleWebhookDeclareTicket = (form) => {
  const {
    enabled,
    mapping,
    ticket_url: ticketUrl,
    ticket_id: ticketId,
    ...rest
  } = form;

  if (!enabled) {
    return null;
  }

  const declareTicket = {
    ...rest,
    ...textPairsToObject(mapping),
  };

  if (ticketUrl.template) {
    declareTicket.ticket_url = '';
    declareTicket.ticket_url_tpl = ticketUrl.value;
  } else {
    declareTicket.ticket_url = ticketUrl.value;
    declareTicket.ticket_url_tpl = '';
  }

  if (ticketId.template) {
    declareTicket.ticket_id = '';
    declareTicket.ticket_id_tpl = ticketId.value;
  } else {
    declareTicket.ticket_id = ticketId.value;
    declareTicket.ticket_id_tpl = '';
  }

  return declareTicket;
};

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
export const formToDeclareTicketRuleWebhooks = (webhooks = []) => removeKeyFromEntities(
  webhooks.map(formToDeclareTicketRuleWebhook),
);

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

/**
 * Convert error structure to form structure
 *
 * @param {FlattenErrors} errors
 * @param {DeclareTicketRuleForm} form
 * @return {FlattenErrors}
 */
export const declareTicketRuleErrorsToForm = (errors, form) => {
  const prepareWebhooksErrors = (errorsObject) => {
    const { webhooks, ...errorMessages } = errorsObject;

    if (webhooks) {
      errorMessages.webhooks = webhooks.reduce((acc, messages, index) => {
        const webhook = form.webhooks[index];
        acc[webhook.key] = messages;

        return acc;
      }, {});
    }

    return errorMessages;
  };

  return flattenErrorMap(errors, prepareWebhooksErrors);
};

/**
 * Convert template variables errors structure to form structure
 *
 * @param {Object} errorsObject
 * @param {DeclareTicketRuleForm} form
 * @return {FlattenErrors}
 */
export const declareTicketRuleTemplateVariablesErrorsToForm = (errorsObject, form) => {
  const { webhooks } = errorsObject;

  return flatten({
    webhooks: webhooks.reduce((acc, { request }, index) => {
      const webhook = form.webhooks[index];

      acc[webhook.key] = {
        request: requestTemplateVariablesErrorsToForm(request, webhook.request),
      };

      return acc;
    }, {}),
  });
};
