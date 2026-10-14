import { isEmpty, omit } from 'lodash';

import {
  ACTION_TYPES,
  ALARM_STATES,
  MIN_DURATION_VALUE,
  PATTERNS_FIELDS,
  TIME_UNITS,
  WORKFLOW_TYPES,
} from '@/constants';

import { durationToForm } from '@/helpers/date/duration';
import {
  declareTicketRuleWebhookDeclareTicketToForm,
  formToDeclareTicketRuleWebhookDeclareTicket,
} from '@/helpers/entities/declare-ticket/rule/form';
import { filterPatternsToForm, formFilterToPatterns } from '@/helpers/entities/filter/form';
import {
  formToPbehavior,
  pbehaviorToForm,
  pbehaviorToRequest,
} from '@/helpers/entities/pbehavior/form';
import {
  formToRequest,
  formToRequestAuthToken,
  requestAuthTokenToForm,
  requestToForm,
} from '@/helpers/entities/shared/request/form';
import { objectToTextPairs, textPairsToObject } from '@/helpers/text-pairs';
import { uid } from '@/helpers/uid';

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
 *   "cancel" |
 *   "pause" |
 *   "resume"
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
 * @typedef {Object} ActionPauseParameters
 * @property {string} output
 * @property {string} type
 * @property {boolean} inherited
 */

/**
 * @typedef {Object} ActionResumeParameters
 * @property {string} output
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
 * @typedef {PbehaviorForm & ActionForwardAuthorParameters} ActionPbehaviorFormParameters
 * @property {boolean} start_on_trigger
 * @property {Duration} duration
 */

/**
 * @typedef {
 *   Pbehavior |
 *   ActionDefaultParameters |
 *   ActionSnoozeParameters |
 *   ActionChangeStateParameters |
 *   ActionWebhookParameters |
 *   ActionAssocTicketParameters |
 *   ActionPauseParameters |
 *   ActionResumeParameters
 * } ActionParameters
 */

/**
 * @typedef {Action} ActionForm
 * @property {FilterPatternsForm} patterns
 * @property {Object.<ActionType, ActionFormParameters>} parameters
 */

/**
 * @typedef {
 *   ActionPbehaviorFormParameters |
 *   ActionDefaultParameters |
 *   ActionSnoozeParameters |
 *   ActionChangeStateParameters |
 *   ActionWebhookFormParameters |
 *   ActionAssocTicketParametersForm |
 *   ActionPauseParameters |
 *   ActionResumeParameters
 * } ActionFormParameters
 */

/**
 * Check action type is pbehavior.
 *
 * @param {ActionType} type
 * @return {boolean}
 */
export const isPbehaviorActionType = type => type === ACTION_TYPES.pbehavior;

/**
 * Check action type is pbehaviorremove.
 *
 * @param {ActionType} type
 * @return {boolean}
 */
export const isPbehaviorRemoveActionType = type => type === ACTION_TYPES.pbehaviorRemove;

/**
 * Check action type is webhook.
 *
 * @param {ActionType} type
 * @return {boolean}
 */
export const isWebhookActionType = type => type === ACTION_TYPES.webhook;

/**
 * Check if action type is associate ticket.
 *
 * @param {ActionType} type
 * @returns {boolean}
 */
export const isAssociateTicketActionType = type => type === ACTION_TYPES.assocticket;

/**
 * Creates default forward-author parameters form.
 *
 * @param {Object} [parameters={}]
 * @returns {Object}
 */
const defaultActionForwardAuthorToForm = (parameters = {}) => ({
  forward_author: parameters.forward_author ?? true,
  author: parameters.author ?? '',
});

/**
 * Creates default action parameters form.
 *
 * @param {Object} [parameters={}] Action parameters.
 * @returns {Object}
 */
export const defaultActionParametersToForm = (parameters = {}) => ({
  ...defaultActionForwardAuthorToForm(parameters),
  output: parameters.output ?? '',
});

/**
 * Creates pause action parameters form.
 *
 * @param {Object} [parameters={}]
 * @returns {Object}
 */
export const pauseActionParametersToForm = (parameters = {}) => ({
  ...defaultActionForwardAuthorToForm(parameters),
  output: parameters.output ?? '',
  type: parameters.type ?? '',
  inherited: parameters.inherited ?? false,
});

/**
 * Creates resume action parameters form.
 *
 * @param {Object} [parameters={}]
 * @returns {Object}
 */
const resumeActionParametersToForm = (parameters = {}) => ({
  ...defaultActionForwardAuthorToForm(parameters),
  output: parameters.output ?? '',
});

/**
 * Creates webhook action parameters form.
 *
 * @param {Object} [parameters={}]
 * @returns {Object}
 */
const webhookActionParametersToForm = (parameters = {}) => {
  const emptyParameters = isEmpty(parameters);

  return {
    ...defaultActionForwardAuthorToForm(parameters),
    declare_ticket: declareTicketRuleWebhookDeclareTicketToForm(parameters.declare_ticket),
    request: requestToForm(parameters.request, parameters.auth_token),
    auth_token: requestAuthTokenToForm(parameters.auth_token),
    multiple_urls: parameters.multiple_urls ?? false,
    skip_for_child: parameters.skip_for_child ?? false,
    skip_for_instruction: parameters.skip_for_instruction ?? false,
    stop_on_fail: emptyParameters
      ? WORKFLOW_TYPES.stop
      : (parameters.stop_on_fail ?? WORKFLOW_TYPES.continue),
    stop_on_success: emptyParameters
      ? WORKFLOW_TYPES.stop
      : (parameters.stop_on_success ?? WORKFLOW_TYPES.continue),
    ticket_system_name: parameters.ticket_system_name ?? '',
  };
};

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
    : { value: MIN_DURATION_VALUE, unit: TIME_UNITS.second },
});

/**
 * Creates change state action parameters form.
 *
 * @param {Object} [parameters={}] Action parameters.
 * @returns {Object}
 */
export const changeStateActionParametersToForm = (parameters = {}) => ({
  ...defaultActionParametersToForm(parameters),
  state: parameters.state ?? ALARM_STATES.minor,
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
  ticket_resources: parameters.ticket_resources ?? false,
  ticket_system_name: parameters.ticket_system_name ?? '',
  mapping: objectToTextPairs(parameters.ticket_data),
});

/**
 * Creates pbehavior action parameters form.
 *
 * @param {Object} [parameters={}]
 * @param {string} [timezone]
 * @returns {Object}
 */
const pbehaviorActionParametersToForm = (parameters = {}, timezone) => {
  const pbehaviorForm = {
    ...defaultActionForwardAuthorToForm(parameters),
    ...pbehaviorToForm(parameters, null, timezone),
  };

  pbehaviorForm.start_on_trigger = !!parameters.start_on_trigger;
  pbehaviorForm.duration = durationToForm(parameters.duration);

  return pbehaviorForm;
};

/**
 * Prepare default parameters for all action types.
 *
 * @returns {Object}
 */
const prepareDefaultActionParameters = () => ({
  [ACTION_TYPES.snooze]: snoozeActionParametersToForm(),
  [ACTION_TYPES.unsnooze]: defaultActionParametersToForm(),
  [ACTION_TYPES.pbehavior]: pbehaviorActionParametersToForm(),
  [ACTION_TYPES.pbehaviorRemove]: defaultActionForwardAuthorToForm(),
  [ACTION_TYPES.changeState]: changeStateActionParametersToForm(),
  [ACTION_TYPES.ack]: defaultActionParametersToForm(),
  [ACTION_TYPES.ackremove]: defaultActionParametersToForm(),
  [ACTION_TYPES.assocticket]: associateTicketActionParametersToForm(),
  [ACTION_TYPES.cancel]: defaultActionParametersToForm(),
  [ACTION_TYPES.webhook]: webhookActionParametersToForm(),
  [ACTION_TYPES.pause]: pauseActionParametersToForm(),
  [ACTION_TYPES.resume]: resumeActionParametersToForm(),
});

/**
 * Converts action parameters to form for all supported scenario actions.
 *
 * @param {Object} [action={}] Action entity.
 * @param {string} [timezone]
 * @returns {Object}
 */
export const actionParametersToForm = (action = {}, timezone) => {
  const parameters = prepareDefaultActionParameters();

  if (!action.type || !action.parameters) {
    return parameters;
  }

  const parametersPreparers = {
    [ACTION_TYPES.snooze]: snoozeActionParametersToForm,
    [ACTION_TYPES.webhook]: webhookActionParametersToForm,
    [ACTION_TYPES.pbehavior]: pbehaviorActionParametersToForm,
    [ACTION_TYPES.assocticket]: associateTicketActionParametersToForm,
    [ACTION_TYPES.changeState]: changeStateActionParametersToForm,
    [ACTION_TYPES.pause]: pauseActionParametersToForm,
    [ACTION_TYPES.resume]: resumeActionParametersToForm,
  };

  const prepareParametersToForm = parametersPreparers[action.type];

  parameters[action.type] = prepareParametersToForm
    ? prepareParametersToForm(action.parameters, timezone)
    : defaultActionParametersToForm(action.parameters);

  return parameters;
};

/**
 * Converts action entity to form.
 *
 * @param {Object} [action={}]
 * @param {string} [timezone]
 * @returns {Object}
 */
export const actionToForm = (action = {}, timezone) => ({
  type: action.type ?? ACTION_TYPES.snooze,
  key: uid(),
  parameters: actionParametersToForm(action, timezone),
  drop_scenario_if_not_matched: !!action.drop_scenario_if_not_matched,
  emit_trigger_success: !!action.emit_trigger_success,
  emit_trigger_fail: !!action.emit_trigger_fail,
  comment: action.comment ?? '',
  patterns: filterPatternsToForm(action, [PATTERNS_FIELDS.alarm, PATTERNS_FIELDS.entity]),
});

/**
 * Converts webhook parameters form to entity.
 *
 * @param {Object} [parameters={}]
 * @returns {Object}
 */
export const formToWebhookActionParameters = (parameters = {}) => ({
  declare_ticket: formToDeclareTicketRuleWebhookDeclareTicket(parameters.declare_ticket),
  request: formToRequest(parameters.request),
  auth_token: formToRequestAuthToken(parameters.auth_token, parameters.request.auth?.type),
  multiple_urls: parameters.multiple_urls ?? false,
  empty_response: parameters.empty_response,
  skip_for_child: parameters.skip_for_child,
  skip_for_instruction: parameters.skip_for_instruction,
  stop_on_fail: parameters.stop_on_fail,
  stop_on_success: parameters.stop_on_success,
  ticket_system_name: parameters.ticket_system_name,
});

/**
 * Converts pbehavior parameters form to entity.
 *
 * @param {Object} [parameters={}]
 * @param {string} [timezone]
 * @returns {Object}
 */
export const formToPbehaviorActionParameters = (parameters = {}, timezone) => {
  const pbehavior = formToPbehavior(omit(parameters, ['start_on_trigger', 'duration']), timezone);

  if (parameters.start_on_trigger) {
    pbehavior.start_on_trigger = parameters.start_on_trigger;
    pbehavior.duration = parameters.duration;
  }

  return pbehaviorToRequest(pbehavior);
};

/**
 * Converts associate ticket parameters form to entity.
 *
 * @param {Object} form
 * @returns {Object}
 */
export const formToAssociateTicketParameters = form => ({
  ...omit(form, ['mapping']),
  ticket_data: textPairsToObject(form.mapping),
});

/**
 * Converts action form to action parameters.
 *
 * @param {Object} form
 * @param {string} [timezone]
 * @returns {Object}
 */
const formToActionParameters = (form, timezone) => {
  const {
    forward_author: forwardAuthor,
    author,
    ...parametersByCurrentType
  } = form.parameters?.[form.type] ?? {};

  const parametersPreparers = {
    [ACTION_TYPES.webhook]: formToWebhookActionParameters,
    [ACTION_TYPES.pbehavior]: formToPbehaviorActionParameters,
    [ACTION_TYPES.assocticket]: formToAssociateTicketParameters,
  };

  const prepareParametersToAction = parametersPreparers[form.type];
  const parameters = prepareParametersToAction
    ? prepareParametersToAction(parametersByCurrentType, timezone)
    : parametersByCurrentType;

  if (!isPbehaviorActionType(form.type)) {
    parameters.forward_author = forwardAuthor;

    if (!forwardAuthor) {
      parameters.author = author;
    }
  }

  return parameters;
};

/**
 * Converts action form to API entity.
 *
 * @param {Object} form Action form.
 * @param {string} [timezone]
 * @returns {Object}
 */
export const formToAction = (form = {}, timezone) => {
  const parameters = formToActionParameters(form, timezone);
  const action = {
    ...omit(form, ['key', 'patterns', 'emit_trigger_fail']),
    ...formFilterToPatterns(form.patterns, [PATTERNS_FIELDS.alarm, PATTERNS_FIELDS.entity]),
    parameters,
  };

  if (isWebhookActionType(form.type) && parameters.declare_ticket) {
    action.emit_trigger_fail = form.emit_trigger_fail;
  }

  return action;
};
