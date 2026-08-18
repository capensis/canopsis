import { omit } from 'lodash';

import { ACTION_TYPES, ALARM_STATES, TIME_UNITS } from '@/constants';

import { durationToForm } from '@/helpers/date/duration';

/**
 * @typedef {
 *   "snooze" |
 *   "unsnooze" |
 *   "pbehavior" |
 *   "pbehaviorremove" |
 *   "changestate" |
 *   "webhook" |
 *   "ack" |
 *   "ackremove" |
 *   "assocticket" |
 *   "cancel"
 * } ActionType
 */

/**
 * @typedef {Object} ActionForwardAuthorParameters
 * @property {boolean} [forward_author]
 * @property {string} [author]
 */

/**
 * @typedef {ActionForwardAuthorParameters} ActionDefaultParameters
 * @property {string} output
 */

/**
 * @typedef {ActionDefaultParameters} ActionSnoozeParameters
 * @property {Duration} duration
 */

/**
 * @typedef {ActionDefaultParameters} ActionChangeStateParameters
 * @property {number} state
 */

/**
 * @typedef {ActionDefaultParameters} AssociateTicketParameters
 * @property {string} ticket
 * @property {string} ticket_url
 * @property {string} ticket_url_title
 * @property {string} ticket_system_name
 * @property {Object} ticket_data
 * @property {boolean} ticket_resources
 */

/**
 * @typedef {ActionDefaultParameters & AssociateTicketParameters} ActionAssocTicketParameters
 */

/**
 * @typedef {AssociateTicketParameters} ActionAssocTicketParametersForm
 * @property {TextPairObject[]} mapping
 */

/**
 * @typedef {Object} ActionWebhookParameters
 * @property {Request} request
 * @property {RequestAuthToken} auth_token
 * @property {?DeclareTicketRuleWebhookDeclareTicket} [declare_ticket]
 * @property {boolean} [forward_author]
 * @property {boolean} [empty_response]
 * @property {boolean} skip_for_child
 * @property {boolean} skip_for_instruction
 * @property {boolean} [stop_on_fail]
 * @property {boolean} [stop_on_success]
 * @property {string} [author]
 * @property {string} [ticket_system_name]
 */

/**
 * @typedef {ActionWebhookParameters} ActionWebhookFormParameters
 * @property {RequestForm} request
 * @property {DeclareTicketRuleWebhookDeclareTicketForm} declare_ticket
 * @property {boolean} [multiple_urls]
 */

/**
 * @typedef {
 *   Pbehavior |
 *   ActionDefaultParameters |
 *   ActionSnoozeParameters |
 *   ActionChangeStateParameters |
 *   ActionWebhookParameters |
 *   ActionAssocTicketParameters
 * } ActionParameters
 */

/**
 * @typedef {Action} ActionForm
 * @property {FilterPatternsForm} patterns
 * @property {Object.<ActionType, ActionFormParameters>} parameters
 */

/**
 * @typedef {
 *   PbehaviorForm |
 *   ActionDefaultParameters |
 *   ActionSnoozeParameters |
 *   ActionChangeStateParameters |
 *   ActionWebhookFormParameters |
 *   ActionAssocTicketParametersForm
 * } ActionFormParameters
 */

/**
 * Creates default action parameters form.
 *
 * @param {Object} [parameters={}] Action parameters.
 * @returns {Object}
 */
export const defaultActionParametersToForm = (parameters = {}) => ({
  output: parameters.output ?? '',
  forward_author: parameters.forward_author ?? true,
  author: parameters.author ?? '',
});

/**
 * Creates snooze action parameters form.
 *
 * @param {Object} [parameters={}] Action parameters.
 * @returns {Object}
 */
export const snoozeActionParametersToForm = (parameters = {}) => ({
  ...defaultActionParametersToForm(parameters),
  duration: parameters.duration
    ? durationToForm(parameters.duration)
    : { value: 1, unit: TIME_UNITS.hour },
});

/**
 * Creates change state action parameters form.
 *
 * @param {Object} [parameters={}] Action parameters.
 * @returns {Object}
 */
export const changeStateActionParametersToForm = (parameters = {}) => ({
  ...defaultActionParametersToForm(parameters),
  state: parameters.state ?? ALARM_STATES.critical,
});

/**
 * Creates associate ticket action parameters form.
 *
 * @param {Object} [parameters={}] Action parameters.
 * @returns {Object}
 */
export const associateTicketActionParametersToForm = (parameters = {}) => ({
  ...defaultActionParametersToForm(parameters),
  ticket: parameters.ticket ?? '',
  ticket_url: parameters.ticket_url ?? '',
  ticket_url_title: parameters.ticket_url_title ?? '',
  ticket_system_name: parameters.ticket_system_name ?? '',
  ticket_resources: parameters.ticket_resources ?? [],
});

/**
 * Check if action type is associate ticket
 *
 * @param {ActionType} type
 * @returns {boolean}
 */
export const isAssociateTicketActionType = type => type === ACTION_TYPES.assocticket;

/**
 * Check action type is pbehavior
 *
 * @param {ActionType} type
 * @return {boolean}
 */
export const isPbehaviorActionType = type => type === ACTION_TYPES.pbehavior;

/**
 * Converts action parameters to form for all supported idle-rule actions.
 *
 * @param {Object} action Action entity.
 * @returns {Object}
 */
export const actionParametersToForm = (action = {}) => ({
  [ACTION_TYPES.snooze]: snoozeActionParametersToForm(
    action.type === ACTION_TYPES.snooze ? action.parameters : undefined,
  ),
  [ACTION_TYPES.ack]: defaultActionParametersToForm(
    action.type === ACTION_TYPES.ack ? action.parameters : undefined,
  ),
  [ACTION_TYPES.ackremove]: defaultActionParametersToForm(
    action.type === ACTION_TYPES.ackremove ? action.parameters : undefined,
  ),
  [ACTION_TYPES.cancel]: defaultActionParametersToForm(
    action.type === ACTION_TYPES.cancel ? action.parameters : undefined,
  ),
  [ACTION_TYPES.assocticket]: associateTicketActionParametersToForm(
    action.type === ACTION_TYPES.assocticket ? action.parameters : undefined,
  ),
  [ACTION_TYPES.changeState]: changeStateActionParametersToForm(
    action.type === ACTION_TYPES.changeState ? action.parameters : undefined,
  ),
  [ACTION_TYPES.pbehavior]: defaultActionParametersToForm(
    action.type === ACTION_TYPES.pbehavior ? action.parameters : undefined,
  ),
});

const actionParametersToEntity = (parameters = {}) => {
  const preparedParameters = omit(parameters, ['mapping']);

  if (!preparedParameters.forward_author) {
    return preparedParameters;
  }

  return omit(preparedParameters, ['author']);
};

/**
 * Converts action form to API entity.
 *
 * @param {Object} form Action form.
 * @returns {Object}
 */
export const formToAction = (form = {}) => {
  const parameters = form.parameters?.[form.type] ?? {};

  return {
    type: form.type,
    parameters: actionParametersToEntity(parameters),
  };
};
