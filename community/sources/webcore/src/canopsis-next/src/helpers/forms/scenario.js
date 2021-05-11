import moment from 'moment-timezone';
import { isUndefined, cloneDeep, omit, isNumber } from 'lodash';

import { ENTITIES_STATES, SCENARIO_ACTION_TYPES, TIME_UNITS } from '@/constants';

import { flattenErrorMap } from '@/helpers/forms/flatten-error-map';
import { objectToTextPairs, textPairsToObject } from '@/helpers/text-pairs';

import uid from '../uid';
import { durationToForm, formToDuration } from '../date/duration';
import { formToPbehavior, pbehaviorToForm, pbehaviorToRequest } from './planning-pbehavior';


/**
 * @typedef {DurationForm} RetryDurationForm
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
 * } ScenarioActionType
 */

/**
 * @typedef {Object} ScenarioActionDefaultParameters
 * @property {string} output
 */

/**
 * @typedef {ScenarioActionDefaultParameters} ScenarioActionSnoozeParameters
 * @property {Duration} duration
 */

/**
 * @typedef {ScenarioActionDefaultParameters} ScenarioActionSnoozeFormParameters
 * @property {DurationForm} duration
 */

/**
 * @typedef {ScenarioActionDefaultParameters} ScenarioActionChangeStateParameters
 * @property {number} state
 */

/**
 * @typedef {Object} ScenarioActionWebhookRequestParameter
 * @property {string} method
 * @property {string} url
 * @property {{ username: string, password: string }} auth
 * @property {Object} headers
 * @property {string} payload
 * @property {boolean} skip_verify
 */

/**
 * @typedef {ScenarioActionWebhookRequestParameter} ScenarioActionWebhookRequestFormParameter
 * @property {TextPairObject[]} headers
 */

/**
 * @typedef {Object} ScenarioActionWebhookParameters
 * @property {ScenarioActionWebhookRequestParameter} request
 * @property {Object} declare_ticket
 * @property {boolean} declare_ticket.empty_response
 * @property {boolean} declare_ticket.is_regexp
 * @property {number} retry_count
 * @property {Duration} retry_delay
 */

/**
 * @typedef {ScenarioActionWebhookParameters} ScenarioActionWebhookFormParameters
 * @property {ScenarioActionWebhookRequestFormParameter} request
 * @property {RetryDurationForm} retry
 * @property {boolean} empty_response
 * @property {boolean} is_regexp
 * @property {TextPairObject[]} declare_ticket
 */

/**
 * @typedef {ScenarioActionDefaultParameters} ScenarioActionAssocTicketParameters
 * @property {string} ticket
 */

/**
 * @typedef {
 *   Pbehavior |
 *   ScenarioActionDefaultParameters |
 *   ScenarioActionSnoozeParameters |
 *   ScenarioActionChangeStateParameters |
 *   ScenarioActionWebhookParameters |
 *   ScenarioActionAssocTicketParameters
 * } ScenarioActionParameters
 */

/**
 * @typedef {Object} ScenarioAction
 * @property {ScenarioActionType} type
 * @property {boolean} drop_scenario_if_not_matched
 * @property {boolean} emit_trigger
 * @property {Object} patterns
 * @property {Object[]} alarm_patterns
 * @property {Object[]} entity_patterns
 * @property {ScenarioActionParameters} parameters
 */

/**
 * @typedef {Object} Scenario
 * @property {string} name
 * @property {string} author
 * @property {number} priority
 * @property {boolean} enabled
 * @property {Duration} delay
 * @property {string[]} triggers
 * @property {DisableDuringPeriods[]} disable_during_periods
 * @property {ScenarioAction[]} actions
 */

/**
 * @typedef {
 *   PbehaviorForm |
 *   ScenarioActionDefaultParameters |
 *   ScenarioActionSnoozeFormParameters |
 *   ScenarioActionChangeStateParameters |
 *   ScenarioActionWebhookFormParameters |
 *   ScenarioActionAssocTicketParameters
 * } ScenarioActionFormParameters
 */

/**
 * @typedef {ScenarioAction} ScenarioActionForm
 * @property {Object[]} patterns.alarm_patterns
 * @property {Object[]} patterns.entity_patterns
 * @property {Object.<ScenarioActionType, ScenarioActionFormParameters>} parameters
 */

/**
 * @typedef {Scenario} ScenarioForm
 * @property {DurationForm} delay
 * @property {ScenarioActionForm[]} actions
 */

/**
 * Convert scenario action parameters to form
 *
 * @param {ScenarioActionDefaultParameters | {}} [parameters = {}]
 * @returns {ScenarioActionDefaultParameters}
 */
const scenarioDefaultActionParametersToForm = (parameters = {}) => ({
  output: parameters.output || '',
});

/**
 * Convert webhook request field to form object
 *
 * @param {ScenarioActionWebhookRequestParameter} [request]
 * @return {ScenarioActionWebhookRequestFormParameter}
 */
const scenarioWebhookActionRequestParametersToForm = (request = {}) => ({
  method: request.method || '',
  url: request.url || '',
  auth: request.auth,
  headers: request.headers ? objectToTextPairs(request.headers) : [],
  payload: request.payload || '{}',
  skip_verify: !!request.skip_verify,
});

/**
 * Convert scenario action webhook parameters to form
 *
 * @param {ScenarioActionWebhookParameters} [parameters = {}]
 * @returns {ScenarioActionWebhookFormParameters}
 */
const scenarioWebhookActionParametersToForm = (parameters = {}) => {
  const { empty_response: emptyResponse, is_regexp: isRegexp, ...variables } = parameters.declare_ticket || {};

  return ({
    declare_ticket: objectToTextPairs(variables),
    empty_response: !!emptyResponse,
    is_regexp: !!isRegexp,
    retry: parameters.retry_delay
      ? { count: parameters.retry_count, ...durationToForm(parameters.retry_delay) }
      : { count: '', unit: '', value: '' },
    request: scenarioWebhookActionRequestParametersToForm(parameters.request),
  });
};

/**
 * Convert scenario action snooze parameters to form
 *
 * @param {ScenarioActionSnoozeParameters | {}} [parameters = {}]
 * @returns {ScenarioActionSnoozeFormParameters}
 */
const scenarioSnoozeActionParametersToForm = (parameters = {}) => ({
  ...scenarioDefaultActionParametersToForm(parameters),
  duration: parameters.duration
    ? durationToForm(parameters.duration)
    : { value: 1, unit: TIME_UNITS.second },
});

/**
 * Convert scenario action snooze parameters to form
 *
 * @param {ScenarioActionChangeStateParameters | {}} [parameters = {}]
 * @returns {ScenarioActionChangeStateParameters}
 */
const scenarioChangeStateActionParametersToForm = (parameters = {}) => ({
  ...scenarioDefaultActionParametersToForm(parameters),
  state: parameters.state || ENTITIES_STATES.minor,
});

/**
 * Convert scenario action assoc ticket parameters to form
 *
 * @param {ScenarioActionAssocTicketParameters | {}} [parameters = {}]
 * @returns {ScenarioActionAssocTicketParameters}
 */
const scenarioAssocTicketActionParametersToForm = (parameters = {}) => ({
  ...scenarioDefaultActionParametersToForm(parameters),
  ticket: parameters.ticket || '',
});

/**
 * Convert scenario action pbehavior parameters to form
 *
 * @param {Pbehavior} [parameters = {}]
 * @param {string} [timezone = moment.tz.guess()]
 * @returns {PbehaviorForm}
 */
const scenarioPbehaviorActionParametersToForm = (parameters = {}, timezone = moment.tz.guess()) => {
  const pbehaviorForm = pbehaviorToForm(parameters, null, timezone);

  pbehaviorForm.start_on_trigger = !!parameters.start_on_trigger;
  pbehaviorForm.duration = durationToForm(parameters.duration);

  return pbehaviorForm;
};


/**
 * Prepare parameters for all scenario action types
 *
 * @returns {Object.<ScenarioActionType, ScenarioActionFormParameters>}
 */
const prepareDefaultScenarioActionParameters = () => ({
  [SCENARIO_ACTION_TYPES.snooze]: scenarioSnoozeActionParametersToForm(),
  [SCENARIO_ACTION_TYPES.pbehavior]: scenarioPbehaviorActionParametersToForm(),
  [SCENARIO_ACTION_TYPES.changeState]: scenarioChangeStateActionParametersToForm(),
  [SCENARIO_ACTION_TYPES.ack]: scenarioDefaultActionParametersToForm(),
  [SCENARIO_ACTION_TYPES.ackremove]: scenarioDefaultActionParametersToForm(),
  [SCENARIO_ACTION_TYPES.assocticket]: scenarioAssocTicketActionParametersToForm(),
  [SCENARIO_ACTION_TYPES.cancel]: scenarioDefaultActionParametersToForm(),
  [SCENARIO_ACTION_TYPES.webhook]: scenarioWebhookActionParametersToForm(),
});

/**
 * Convert scenario action parameters to form
 *
 * @param {ScenarioAction} action
 * @param {string} timezone
 * @returns {Object.<ScenarioActionType, ScenarioActionParameters>}
 */
export const scenarioActionParametersToForm = (action, timezone) => {
  const parameters = prepareDefaultScenarioActionParameters();

  if (!action.type || !action.parameters) {
    return parameters;
  }

  const parametersPreparers = {
    [SCENARIO_ACTION_TYPES.snooze]: scenarioSnoozeActionParametersToForm,
    [SCENARIO_ACTION_TYPES.webhook]: scenarioWebhookActionParametersToForm,
    [SCENARIO_ACTION_TYPES.pbehavior]: scenarioPbehaviorActionParametersToForm,
  };

  const prepareParametersToFormFunction = parametersPreparers[action.type];

  parameters[action.type] = prepareParametersToFormFunction
    ? prepareParametersToFormFunction(action.parameters, timezone)
    : { ...action.parameters };

  return parameters;
};

/**
 * Convert scenario action to form
 *
 * @param {ScenarioAction} [scenarioAction = {}]
 * @param {string} [timezone = moment.tz.guess()]
 * @returns {ScenarioActionForm}
 */
export const scenarioActionToForm = (scenarioAction = {}, timezone = moment.tz.guess()) => {
  const type = scenarioAction.type || SCENARIO_ACTION_TYPES.snooze;

  return {
    type,
    key: uid(),
    parameters: scenarioActionParametersToForm(scenarioAction, timezone),
    drop_scenario_if_not_matched: !!scenarioAction.drop_scenario_if_not_matched,
    emit_trigger: !!scenarioAction.emit_trigger,
    patterns: {
      alarm_patterns: scenarioAction.alarm_patterns ? cloneDeep(scenarioAction.alarm_patterns) : [],
      entity_patterns: scenarioAction.entity_patterns ? cloneDeep(scenarioAction.entity_patterns) : [],
    },
  };
};

/**
 * Convert scenario to form
 *
 * @param {Scenario} [scenario = {}]
 * @param {string} [timezone = moment.tz.guess()]
 * @returns {ScenarioForm}
 */
export const scenarioToForm = (scenario = {}, timezone = moment.tz.guess()) => ({
  name: scenario.name || '',
  priority: scenario.priority || 1,
  enabled: !isUndefined(scenario.enabled) ? scenario.enabled : true,
  delay: scenario.delay
    ? durationToForm(scenario.delay)
    : { value: undefined, unit: undefined },
  triggers: scenario.triggers ? [...scenario.triggers] : [],
  disable_during_periods: scenario.disable_during_periods ? [...scenario.disable_during_periods] : [],
  actions: scenario.actions
    ? scenario.actions.map(action => scenarioActionToForm(action, timezone))
    : [scenarioActionToForm(undefined, timezone)],
});

/**
 * Convert pbehavior parameters to scenario action
 *
 * @param {ScenarioActionWebhookFormParameters | {}} [parameters = {}]
 * @return {ScenarioActionWebhookParameters}
 */
export const formToScenarioWebhookActionParameters = (parameters = {}) => {
  const webhook = {
    declare_ticket: {
      empty_response: parameters.empty_response,
      is_regexp: parameters.is_regexp,
      ...textPairsToObject(parameters.declare_ticket),
    },
    request: {
      ...parameters.request,
      payload: parameters.request.payload,
      headers: textPairsToObject(parameters.request.headers),
    },
  };

  if (parameters.retry.value) {
    webhook.retry_count = parameters.retry.count;
    webhook.retry_delay = formToDuration(parameters.retry);
  }

  return webhook;
};

/**
 * Convert snooze parameters to scenario action
 *
 * @param {ScenarioActionSnoozeFormParameters} parameters
 * @return {ScenarioActionSnoozeParameters}
 */
export const formToScenarioSnoozeActionParameters = (parameters = {}) =>
  ({
    ...parameters,
    duration: formToDuration(parameters.duration),
  });

/**
 * Convert pbehavior parameters to scenario action
 *
 * @param {PbehaviorForm} parameters
 * @param [timezone = moment.tz.guess()]
 * @return {PbehaviorRequest}
 */
export const formToScenarioPbehaviorActionParameters = (parameters = {}, timezone = moment.tz.guess()) => {
  const pbehavior = formToPbehavior(omit(parameters, ['start_on_trigger', 'duration']), timezone);

  if (parameters.start_on_trigger) {
    pbehavior.start_on_trigger = parameters.start_on_trigger;
    pbehavior.duration = formToDuration(parameters.duration);
  }

  return pbehaviorToRequest(pbehavior);
};

/**
 * Convert form to scenario action
 *
 * @param {ScenarioActionForm} form
 * @param {string} [timezone]
 * @returns {ScenarioAction}
 */
export const formToScenarioAction = (form, timezone) => {
  const parametersByCurrentType = form.parameters[form.type];

  const parametersPreparers = {
    [SCENARIO_ACTION_TYPES.snooze]: formToScenarioSnoozeActionParameters,
    [SCENARIO_ACTION_TYPES.webhook]: formToScenarioWebhookActionParameters,
    [SCENARIO_ACTION_TYPES.pbehavior]: formToScenarioPbehaviorActionParameters,
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

/**
 * Convert form to scenario
 *
 * @param {ScenarioForm} form
 * @param {string} [timezone = moment.tz.guess()]
 * @returns {Scenario}
 */
export const formToScenario = (form, timezone = moment.tz.guess()) => ({
  ...omit(form, ['delay', 'actions']),
  delay: form.delay && isNumber(form.delay.value)
    ? formToDuration(form.delay)
    : undefined,
  actions: form.actions.map(action => formToScenarioAction(action, timezone)),
});

/**
 * Convert error structure to form structure
 *
 * @param {FlattenErrors} errors
 * @param {ScenarioForm} form
 * @return {FlattenErrors}
 */
export const scenarioErrorToForm = (errors, form) => {
  const prepareScenarioActionsErrors = (errorsObject) => {
    const { actions, ...errorMessages } = errorsObject;

    if (actions) {
      errorMessages.actions = actions.reduce((acc, messages, index) => {
        const action = form.actions[index];
        acc[action.key] = messages;

        return acc;
      }, {});
    }

    return errorMessages;
  };

  return flattenErrorMap(errors, prepareScenarioActionsErrors);
};
