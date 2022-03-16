import { cloneDeep, omit, isEmpty } from 'lodash';

import { ENTITIES_STATES, ACTION_TYPES } from '@/constants';

import { objectToTextPairs, textPairsToObject } from '../text-pairs';
import uid from '../uid';
import { durationToForm } from '../date/duration';
import { getLocaleTimezone } from '../date/date';

import { formToPbehavior, pbehaviorToForm, pbehaviorToRequest } from './planning-pbehavior';

/**
 * @typedef {Duration} RetryDuration
 * @property {number} count
 */

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
 * @property {string} author
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
 * @typedef {Object} ActionWebhookRequestParameter
 * @property {string} method
 * @property {string} url
 * @property {{ username: string, password: string }} auth
 * @property {Object} headers
 * @property {string} payload
 * @property {boolean} skip_verify
 */

/**
 * @typedef {ActionWebhookRequestParameter} ActionWebhookRequestFormParameter
 * @property {TextPairObject[]} headers
 */

/**
 * @typedef {Object} ActionWebhookParameters
 * @property {ActionWebhookRequestParameter} request
 * @property {?Object} [declare_ticket]
 * @property {boolean} declare_ticket.empty_response
 * @property {boolean} declare_ticket.is_regexp
 * @property {number} retry_count
 * @property {Duration} retry_delay
 */

/**
 * @typedef {ActionWebhookParameters} ActionWebhookFormParameters
 * @property {ActionWebhookRequestFormParameter} request
 * @property {RetryDuration} retry
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
 * @typedef {Object} Action
 * @property {ActionType} type
 * @property {boolean} drop_scenario_if_not_matched
 * @property {boolean} emit_trigger
 * @property {Object} patterns
 * @property {string} comment
 * @property {Object[]} alarm_patterns
 * @property {Object[]} entity_patterns
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
 * @property {Object[]} patterns.alarm_patterns
 * @property {Object[]} patterns.entity_patterns
 * @property {Object.<ActionType, ActionFormParameters>} parameters
 */

/**
 * Convert action parameters to form
 *
 * @param {ActionDefaultParameters | {}} [parameters = {}]
 * @returns {ActionDefaultParameters}
 */
const defaultActionParametersToForm = (parameters = {}) => ({
  output: parameters.output ?? '',
  author: parameters.author ?? '',
});

/**
 * Convert webhook request field to form object
 *
 * @param {ActionWebhookRequestParameter} [request]
 * @return {ActionWebhookRequestFormParameter}
 */
const webhookActionRequestParametersToForm = (request = {}) => ({
  method: request.method || '',
  url: request.url || '',
  auth: request.auth,
  headers: request.headers ? objectToTextPairs(request.headers) : [],
  payload: request.payload || '{}',
  skip_verify: !!request.skip_verify,
});

/**
 * Convert action webhook parameters to form
 *
 * @param {ActionWebhookParameters} [parameters = {}]
 * @returns {ActionWebhookFormParameters}
 */
const webhookActionParametersToForm = (parameters = {}) => {
  const { empty_response: emptyResponse, is_regexp: isRegexp, ...variables } = parameters.declare_ticket || {};

  return {
    declare_ticket: objectToTextPairs(variables),
    empty_response: !!emptyResponse,
    is_regexp: !!isRegexp,
    retry: parameters.retry_delay
      ? { count: parameters.retry_count, ...durationToForm(parameters.retry_delay) }
      : { count: '', unit: '', value: '' },
    request: webhookActionRequestParametersToForm(parameters.request),
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
export const actionToForm = (action = {}, timezone = getLocaleTimezone()) => {
  const type = action.type || ACTION_TYPES.snooze;

  return {
    type,
    key: uid(),
    parameters: actionParametersToForm(action, timezone),
    drop_scenario_if_not_matched: !!action.drop_scenario_if_not_matched,
    emit_trigger: !!action.emit_trigger,
    comment: action.comment || '',
    patterns: {
      alarm_patterns: action.alarm_patterns ? cloneDeep(action.alarm_patterns) : [],
      entity_patterns: action.entity_patterns ? cloneDeep(action.entity_patterns) : [],
    },
  };
};

/**
 * Convert pbehavior parameters to action
 *
 * @param {ActionWebhookFormParameters | {}} [parameters = {}]
 * @return {ActionWebhookParameters}
 */
export const formToWebhookActionParameters = (parameters = {}) => {
  const webhook = {
    declare_ticket: null,
    request: {
      ...parameters.request,
      payload: parameters.request.payload,
      headers: textPairsToObject(parameters.request.headers),
    },
  };

  if (parameters.retry.value) {
    webhook.retry_count = parameters.retry.count;
    webhook.retry_delay = parameters.retry;
  }

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
 * Convert form to action
 *
 * @param {ActionForm} form
 * @param {string} [timezone]
 * @returns {Action}
 */
export const formToAction = (form, timezone) => {
  const parametersByCurrentType = form.parameters[form.type];

  const parametersPreparers = {
    [ACTION_TYPES.webhook]: formToWebhookActionParameters,
    [ACTION_TYPES.pbehavior]: formToPbehaviorActionParameters,
  };

  const prepareParametersToAction = parametersPreparers[form.type];
  const parameters = prepareParametersToAction
    ? prepareParametersToAction(parametersByCurrentType, timezone)
    : { ...parametersByCurrentType };

  return {
    ...omit(form, ['key', 'patterns']),
    ...form.patterns,
    parameters,
  };
};
