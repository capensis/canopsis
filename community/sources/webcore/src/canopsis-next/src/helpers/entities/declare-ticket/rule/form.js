import { isEmpty } from 'lodash';

import {
  formToRequest,
  formToRequestAuthToken,
  requestAuthTokenToForm,
  requestToForm,
} from '@/helpers/entities/shared/request/form';
import { filterPatternsToForm, formFilterToPatterns } from '@/helpers/entities/filter/form';
import { objectToTextPairs, textPairsToObject } from '@/helpers/text-pairs';
import { removeKeyFromEntities } from '@/helpers/array';
import { flattenErrorMap } from '@/helpers/entities/shared/form';
import { uid } from '@/helpers/uid';

/**
 * @typedef {Object} DeclareTicketRuleCheckTicketStatus
 * @property {Request} request
 * @property {RequestAuthToken} auth_token
 * @property {Object} status_mapping
 * @property {string} [ticket_status]
 * @property {string} [ticket_status_tpl]
 * @property {boolean} reuse_headers_and_auth
 */

/**
 * @typedef {Object} DeclareTicketRuleWebhookDeclareTicket
 * @property {DeclareTicketRuleCheckTicketStatus} [check_ticket_status]
 * @property {string} empty_response
 * @property {string} is_regexp
 * @property {string} ticket_id
 * @property {string} [ticket_id_tpl]
 * @property {string} [ticket_url]
 * @property {string} [ticket_url_tpl]
 * @property {string} [ticket_url_title]
 */

/**
 * @typedef {Object} DeclareTicketRuleWebhook
 * @property {Request} request
 * @property {RequestAuthToken} auth_token
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
 * @typedef {Object} DeclareTicketRuleWebhookTicketTemplateForm
 * @property {boolean} template
 * @property {string} value
 */

/**
 * @typedef {Object} DeclareTicketRuleCheckTicketStatusForm
 * @property {RequestForm} request
 * @property {RequestAuthTokenForm} auth_token
 * @property {TextPairObject[]} status_mapping
 * @property {DeclareTicketRuleWebhookTicketTemplateForm} ticket_status
 * @property {boolean} reuse_headers_and_auth
 */

/**
 * @typedef {DeclareTicketRuleWebhookDeclareTicket} DeclareTicketRuleWebhookDeclareTicketForm
 * @property {boolean} enabled
 * @property {TextPairObject[]} mapping
 * @property {DeclareTicketRuleWebhookTicketTemplateForm} ticket_id
 * @property {DeclareTicketRuleWebhookTicketTemplateForm} ticket_url
 */

/**
 * @typedef {DeclareTicketRuleWebhook} DeclareTicketRuleWebhookForm
 * @property {DeclareTicketRuleWebhookDeclareTicketForm} declare_ticket
 * @property {RequestForm} request
 * @property {RequestAuthTokenForm} auth_token
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
 * Convert declare ticket rule check ticket status object to form compatible object
 *
 * @param {DeclareTicketRuleCheckTicketStatus} declareTicketStatus
 * @returns {DeclareTicketRuleCheckTicketStatusForm}
 */
export const declareTicketRuleCheckTicketStatusToForm = (declareTicketStatus = {}) => {
  const {
    request = {},
    auth_token: authToken = {},
    status_mapping: statusMapping = [],
    ticket_status: ticketStatus = '',
    ticket_status_tpl: ticketStatusTpl = '',
    reuse_headers_and_auth: reuseHeadersAndAuth = false,
  } = declareTicketStatus;

  return {
    enabled: !isEmpty(declareTicketStatus),
    request: requestToForm(request, authToken),
    auth_token: requestAuthTokenToForm(authToken),
    status_mapping: objectToTextPairs(statusMapping),
    ticket_status: {
      template: !!ticketStatusTpl,
      value: ticketStatusTpl || ticketStatus,
    },
    reuse_headers_and_auth: reuseHeadersAndAuth,
  };
};

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
    ticket_url_title: ticketUrlTitle = '',
    check_ticket_status: checkTicketStatus,
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
    ticket_url_title: ticketUrlTitle,
    mapping: objectToTextPairs(fields),
    check_ticket_status: declareTicketRuleCheckTicketStatusToForm(checkTicketStatus),
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
  request: requestToForm(webhook.request, webhook.auth_token),
  auth_token: requestAuthTokenToForm(webhook.auth_token),
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
 * Convert declare ticket rule check ticket status form to API compatible object
 *
 * @param {DeclareTicketRuleCheckTicketStatusForm} form
 * @returns {DeclareTicketRuleCheckTicketStatus}
 */
export const formToDeclareTicketRuleCheckTicketStatus = (form) => {
  const {
    enabled,
    ticket_status: ticketStatus,
  } = form;

  if (!enabled) {
    return null;
  }

  const result = {
    request: formToRequest(form.request),
    auth_token: formToRequestAuthToken(form.auth_token),
    status_mapping: textPairsToObject(form.status_mapping),
    reuse_headers_and_auth: form.reuse_headers_and_auth || false,
  };

  if (ticketStatus.template) {
    result.ticket_status = '';
    result.ticket_status_tpl = ticketStatus.value;
  } else {
    result.ticket_status = ticketStatus.value;
    result.ticket_status_tpl = '';
  }

  return result;
};

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

    check_ticket_status: formToDeclareTicketRuleCheckTicketStatus(form.check_ticket_status),
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
  auth_token: formToRequestAuthToken(webhook.auth_token, webhook.request.auth?.type),
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
