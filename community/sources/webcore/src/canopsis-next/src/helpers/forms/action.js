import { omit, isEmpty } from 'lodash';

import { ENTITIES_STATES, ACTION_TYPES, PATTERNS_FIELDS, OLD_PATTERNS_FIELDS } from '@/constants';

import { filterPatternsToForm, formFilterToPatterns } from '@/helpers/forms/filter';

import { objectToTextPairs, textPairsToObject } from '../text-pairs';
import uid from '../uid';
import { durationToForm } from '../date/duration';
import { getLocaleTimezone } from '../date/date';

import { formToPbehavior, pbehaviorToForm, pbehaviorToRequest } from './planning-pbehavior';
import { requestToForm, retryToForm, formToRequest, formToRetry } from './shared/request';

/**
 * @typedef {
 *   "snooze" |
 *   "pbehavior" |
 *   "changestate" |
 *   "webhook" |
 *   "ack" |
 *   "ackremove" |
 *   "assocticket" |
 *   "cancel"
 * } ActionType
 */

/**
 * @typedef {Object} ActionDefaultParameters
 * @property {string} output
 * @property {boolean} [forward_author]
 * @property {string} [author]
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
 * @typedef {ActionDefaultParameters} ActionAssocTicketParameters
 * @property {string} ticket
 */

/**
 * @typedef {Object} ActionWebhookParameters
 * @property {Request} request
 * @property {?Object} [declare_ticket]
 * @property {boolean} declare_ticket.empty_response
 * @property {boolean} declare_ticket.is_regexp
 * @property {number} retry_count
 * @property {Duration} retry_delay
 * @property {boolean} [forward_author]
 * @property {string} [author]
 */

/**
 * @typedef {ActionWebhookParameters} ActionWebhookFormParameters
 * @property {RequestForm} request
 * @property {boolean} empty_response
 * @property {boolean} is_regexp
 * @property {TextPairObject[]} declare_ticket
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
 * @typedef {FilterPatterns} Action
 * @property {ActionType} type
 * @property {boolean} drop_scenario_if_not_matched
 * @property {boolean} emit_trigger
 * @property {string} comment
 * @property {ActionParameters} parameters
 */

/**
 * @typedef {
 *   PbehaviorForm |
 *   ActionDefaultParameters |
 *   ActionSnoozeParameters |
 *   ActionChangeStateParameters |
 *   ActionWebhookFormParameters |
 *   ActionAssocTicketParameters
 * } ActionFormParameters
 */

/**
 * @typedef {Action} ActionForm
 * @property {FilterPatternsForm} patterns
 * @property {Object.<ActionType, ActionFormParameters>} parameters
 */

/**
 * Check action type is pbehavior
 *
 * @param {ActionType} type
 * @return {boolean}
 */
export const isPbehaviorActionType = type => type === ACTION_TYPES.pbehavior;

/**
 * Convert action parameters to form
 *
 * @param {ActionDefaultParameters | {}} [parameters = {}]
 * @returns {ActionDefaultParameters}
 */
const defaultActionParametersToForm = (parameters = {}) => ({
  output: parameters.output ?? '',
  forward_author: parameters.forward_author ?? true,
  author: parameters.author ?? '',
});

/**
 * Convert action webhook parameters to form
 *
 * @param {ActionWebhookParameters} [parameters = {}]
 * @returns {ActionWebhookFormParameters}
 */
const webhookActionParametersToForm = (parameters = {}) => {
  const { empty_response: emptyResponse, is_regexp: isRegexp, ...variables } = parameters.declare_ticket ?? {};

  return {
    forward_author: parameters.forward_author ?? true,
    author: parameters.author ?? '',
    declare_ticket: objectToTextPairs(variables),
    empty_response: !!emptyResponse,
    is_regexp: !!isRegexp,
    retry: retryToForm(parameters),
    request: requestToForm(parameters.request),
  };
};

/**
 * Convert action snooze parameters to form
 *
 * @param {ActionSnoozeParameters | {}} [parameters = {}]
 * @returns {ActionSnoozeParameters}
 */
const snoozeActionParametersToForm = (parameters = {}) => ({
  ...defaultActionParametersToForm(parameters),
  duration: durationToForm(parameters.duration),
});

/**
 * Convert action snooze parameters to form
 *
 * @param {ActionChangeStateParameters | {}} [parameters = {}]
 * @returns {ActionChangeStateParameters}
 */
const changeStateActionParametersToForm = (parameters = {}) => ({
  ...defaultActionParametersToForm(parameters),
  state: parameters.state ?? ENTITIES_STATES.minor,
});

/**
 * Convert action assoc ticket parameters to form
 *
 * @param {ActionAssocTicketParameters | {}} [parameters = {}]
 * @returns {ActionAssocTicketParameters}
 */
const assocTicketActionParametersToForm = (parameters = {}) => ({
  ...defaultActionParametersToForm(parameters),
  ticket: parameters.ticket ?? '',
});

/**
 * Convert action pbehavior parameters to form
 *
 * @param {Pbehavior} [parameters = {}]
 * @param {string} [timezone = getLocaleTimezone()]
 * @returns {PbehaviorForm}
 */
const pbehaviorActionParametersToForm = (parameters = {}, timezone = getLocaleTimezone()) => {
  const pbehaviorForm = pbehaviorToForm(parameters, null, timezone);

  pbehaviorForm.start_on_trigger = !!parameters.start_on_trigger;
  pbehaviorForm.duration = durationToForm(parameters.duration);

  return pbehaviorForm;
};

/**
 * Prepare parameters for all action types
 *
 * @returns {Object.<ActionType, ActionFormParameters>}
 */
const prepareDefaultActionParameters = () => ({
  [ACTION_TYPES.snooze]: snoozeActionParametersToForm(),
  [ACTION_TYPES.pbehavior]: pbehaviorActionParametersToForm(),
  [ACTION_TYPES.changeState]: changeStateActionParametersToForm(),
  [ACTION_TYPES.ack]: defaultActionParametersToForm(),
  [ACTION_TYPES.ackremove]: defaultActionParametersToForm(),
  [ACTION_TYPES.assocticket]: assocTicketActionParametersToForm(),
  [ACTION_TYPES.cancel]: defaultActionParametersToForm(),
  [ACTION_TYPES.webhook]: webhookActionParametersToForm(),
});

/**
 * Convert action parameters to form
 *
 * @param {Action} action
 * @param {string} timezone
 * @returns {Object.<ActionType, ActionParameters>}
 */
export const actionParametersToForm = (action, timezone) => {
  const parameters = prepareDefaultActionParameters();

  if (!action.type || !action.parameters) {
    return parameters;
  }

  const parametersPreparers = {
    [ACTION_TYPES.snooze]: snoozeActionParametersToForm,
    [ACTION_TYPES.webhook]: webhookActionParametersToForm,
    [ACTION_TYPES.pbehavior]: pbehaviorActionParametersToForm,
  };

  const prepareParametersToFormFunction = parametersPreparers[action.type];

  parameters[action.type] = prepareParametersToFormFunction
    ? prepareParametersToFormFunction(action.parameters, timezone)
    : { ...action.parameters };

  return parameters;
};

/**
 * Convert action to form
 *
 * @param {Action} [action = {}]
 * @param {string} [timezone = getLocaleTimezone()]
 * @returns {ActionForm}
 */
export const actionToForm = (action = {}, timezone = getLocaleTimezone()) => ({
  type: action.type ?? ACTION_TYPES.snooze,
  key: uid(),
  parameters: actionParametersToForm(action, timezone),
  drop_scenario_if_not_matched: !!action.drop_scenario_if_not_matched,
  emit_trigger: !!action.emit_trigger,
  comment: action.comment ?? '',
  patterns: filterPatternsToForm(
    action,
    [PATTERNS_FIELDS.alarm, PATTERNS_FIELDS.entity],
    [OLD_PATTERNS_FIELDS.alarm, OLD_PATTERNS_FIELDS.entity],
  ),
});

/**
 * Convert pbehavior parameters to action
 *
 * @param {ActionWebhookFormParameters | {}} [parameters = {}]
 * @return {ActionWebhookParameters}
 */
export const formToWebhookActionParameters = (parameters = {}) => {
  const webhook = {
    declare_ticket: null,
    request: formToRequest(parameters.request),

    ...formToRetry(parameters.retry),
  };

  if (parameters.empty_response || parameters.is_regexp || !isEmpty(parameters.declare_ticket)) {
    webhook.declare_ticket = {
      empty_response: parameters.empty_response,
      is_regexp: parameters.is_regexp,

      ...textPairsToObject(parameters.declare_ticket),
    };
  }

  return webhook;
};

/**
 * Convert pbehavior parameters to action
 *
 * @param {PbehaviorForm | {}} [parameters = {}]
 * @param {string} [timezone = getLocaleTimezone()]
 * @return {PbehaviorRequest}
 */
export const formToPbehaviorActionParameters = (parameters = {}, timezone = getLocaleTimezone()) => {
  const pbehavior = formToPbehavior(omit(parameters, ['start_on_trigger', 'duration']), timezone);

  if (parameters.start_on_trigger) {
    pbehavior.start_on_trigger = parameters.start_on_trigger;
    pbehavior.duration = parameters.duration;
  }

  return pbehaviorToRequest(pbehavior);
};

/**
 * Convert form to action parameters
 *
 * @param {ActionForm} form
 * @param {string} [timezone]
 * @returns {ActionParameters}
 */
const formToActionParameters = (form, timezone) => {
  const parametersByCurrentType = form.parameters[form.type];

  const parametersPreparers = {
    [ACTION_TYPES.webhook]: formToWebhookActionParameters,
    [ACTION_TYPES.pbehavior]: formToPbehaviorActionParameters,
  };

  const prepareParametersToAction = parametersPreparers[form.type];
  const parameters = prepareParametersToAction
    ? prepareParametersToAction(parametersByCurrentType, timezone)
    : omit(parametersByCurrentType, ['author', 'forward_author']);

  if (!isPbehaviorActionType(form.type)) {
    parameters.forward_author = parametersByCurrentType.forward_author;

    if (!parameters.forward_author) {
      parameters.author = parametersByCurrentType.author;
    }
  }

  return parameters;
};

/**
 * Convert form to action
 *
 * @param {ActionForm} form
 * @param {string} [timezone]
 * @returns {Action}
 */
export const formToAction = (form, timezone) => ({
  ...omit(form, ['key', 'patterns']),
  ...formFilterToPatterns(form.patterns, [PATTERNS_FIELDS.alarm, PATTERNS_FIELDS.entity]),
  parameters: formToActionParameters(form, timezone),
});
