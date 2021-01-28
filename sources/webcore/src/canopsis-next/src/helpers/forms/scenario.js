import moment from 'moment-timezone';
import { isUndefined, cloneDeep, omit, isNumber } from 'lodash';

import { ENTITIES_STATES, SCENARIO_ACTION_TYPES } from '@/constants';

import uid from '../uid';
import { durationToForm, formToDuration } from '../date/duration';
import { formToPbehavior, pbehaviorToForm, pbehaviorToRequest } from './planning-pbehavior';

/**
 * @typedef {
 *   "snooze" |
 *   "pbehavior" |
 *   "changestate" |
 *   "webhook" |
 *   "ack" |
 *   "ackremove" |
 *   "assocticket" |
 *   "declareticket" |
 *   "cancel"
 * } ScenarioActionType
 */

/**
 * @typedef {Object} ScenarioActionDefaultParameters
 * @property {string} author
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
 */

/**
 * @typedef {Object} ScenarioActionWebhookParameters
 * @property {ScenarioActionWebhookRequestParameter} request
 * @property {Object} declare_ticket
 * @property {number} retry_count
 * @property {Duration} retry_delay
 */

/**
 * @typedef {Object} ScenarioActionWebhookFormParameters
 * @property {ScenarioActionWebhookRequestParameter} request
 * @property {Object} declare_ticket
 * @property {number} retry_count
 * @property {DurationForm} retry_delay
 */

/**
 * @typedef {ScenarioActionDefaultParameters} ScenarioActionAssocTicketParameters
 * @property {string} ticket
 */

/**
 * @typedef {Object} ScenarioAction
 * @property {ScenarioActionType} type
 * @property {boolean} drop_scenario_if_not_matched
 * @property {boolean} emit_trigger
 * @property {Object[]} alarm_patterns
 * @property {Object[]} entity_patterns
 * @property {
 *   Pbehavior |
 *   ScenarioActionDefaultParameters |
 *   ScenarioActionSnoozeParameters |
 *   ScenarioActionChangeStateParameters |
 *   ScenarioActionWebhookParameters |
 *   ScenarioActionAssocTicketParameters
 * } parameters
 */

/**
 * @typedef {Object} Scenario
 * @property {string} name
 * @property {string} author
 * @property {number} priority
 * @property {boolean} enabled
 * @property {Duration} delay
 * @property {string[]} triggers
 * @property {string[]} disable_during_periods
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
 * } ScenarioActionParameters
 */

/**
 * @typedef {ScenarioAction} ScenarioActionForm
 * @property {Object.<ScenarioActionType, ScenarioActionParameters>} parameters
 */

/**
 * @typedef {Scenario} ScenarioForm
 * @property {DurationForm} delay
 * @property {ScenarioActionForm[]} actions
 */

/**
 * Convert scenario action parameters to form
 *
 * @param {ScenarioActionDefaultParameters} [parameters]
 * @returns {ScenarioActionDefaultParameters}
 */
const scenarioDefaultActionParametersToForm = (parameters = {}) => ({
  output: parameters.output || '',
  author: parameters.author || '',
});

/**
 * Convert scenario action webhook parameters to form
 *
 * @param {ScenarioActionWebhookParameters} [parameters]
 * @returns {ScenarioActionWebhookFormParameters}
 */
const scenarioWebhookActionParametersToForm = (parameters = {}) => ({
  declare_ticket: cloneDeep(parameters.declare_ticket || {}),
  retry_count: parameters.retry_count || '',
  retry_delay: durationToForm(parameters.retry_delay),
  request: parameters.request || {
    method: '',
    url: '',
    auth: {
      username: '',
      password: '',
    },
    headers: {},
    payload: '',
  },
});

/**
 * Convert scenario action snooze parameters to form
 *
 * @param {ScenarioActionSnoozeParameters} [parameters]
 * @returns {ScenarioActionSnoozeFormParameters}
 */
const scenarioSnoozeActionParametersToForm = (parameters = {}) => ({
  ...scenarioDefaultActionParametersToForm(parameters),
  duration: durationToForm(parameters.duration),
});

/**
 * Convert scenario action snooze parameters to form
 *
 * @param {ScenarioActionChangeStateParameters} [parameters={}]
 * @returns {ScenarioActionChangeStateParameters}
 */
const scenarioChangeStateActionParametersToForm = (parameters = {}) => ({
  ...scenarioDefaultActionParametersToForm(parameters),
  state: parameters.state || ENTITIES_STATES.minor,
});

/**
 * Convert scenario action assoc ticket parameters to form
 *
 * @param {ScenarioActionAssocTicketParameters} [parameters={}]
 * @returns {ScenarioActionAssocTicketParameters}
 */
const scenarioAssocTicketActionParametersToForm = (parameters = {}) => ({
  ...scenarioDefaultActionParametersToForm(parameters),
  ticket: parameters.ticket || '',
});

/**
 * Convert scenario action pbehavior parameters to form
 *
 * @param {Pbehavior} [parameters={}]
 * @param {string} [timezone=moment.tz.guess()]
 * @returns {PbehaviorForm}
 */
const scenarioPbehaviorActionParametersToForm = (parameters = {}, timezone = moment.tz.guess()) =>
  pbehaviorToForm(parameters, null, timezone);

/**
 *
 * @returns {Object.<ScenarioActionType, ScenarioActionParameters>}
 */
const prepareDefaultScenarioActionParameters = () => ({
  [SCENARIO_ACTION_TYPES.snooze]: scenarioSnoozeActionParametersToForm(),
  [SCENARIO_ACTION_TYPES.pbehavior]: scenarioPbehaviorActionParametersToForm(),
  [SCENARIO_ACTION_TYPES.changeState]: scenarioChangeStateActionParametersToForm(),
  [SCENARIO_ACTION_TYPES.ack]: scenarioDefaultActionParametersToForm(),
  [SCENARIO_ACTION_TYPES.ackremove]: scenarioDefaultActionParametersToForm(),
  [SCENARIO_ACTION_TYPES.assocticket]: scenarioAssocTicketActionParametersToForm(),
  [SCENARIO_ACTION_TYPES.declareticket]: scenarioDefaultActionParametersToForm(),
  [SCENARIO_ACTION_TYPES.cancel]: scenarioDefaultActionParametersToForm(),
  [SCENARIO_ACTION_TYPES.webhook]: scenarioWebhookActionParametersToForm(),
});

/**
 * Convert scenario action to form
 *
 * @param {ScenarioAction} [scenarioAction = {}]
 * @param {string} [timezone = moment.tz.guess()]
 * @returns {ScenarioActionForm}
 */
export const scenarioActionToForm = (scenarioAction = {}, timezone = moment.tz.guess()) => {
  const type = scenarioAction.type || SCENARIO_ACTION_TYPES.snooze;
  const parameters = prepareDefaultScenarioActionParameters();

  const parametersPreparers = {
    [SCENARIO_ACTION_TYPES.snooze]: scenarioSnoozeActionParametersToForm,
    [SCENARIO_ACTION_TYPES.webhook]: scenarioWebhookActionParametersToForm,
    [SCENARIO_ACTION_TYPES.pbehavior]: scenarioPbehaviorActionParametersToForm,
  };

  const prepareParametersToFormFunction = parametersPreparers[type];

  if (scenarioAction.parameters) {
    parameters[type] = prepareParametersToFormFunction
      ? prepareParametersToFormFunction(scenarioAction.parameters, timezone)
      : { ...scenarioAction.parameters };
  }

  return {
    type,
    parameters,
    key: uid(),
    drop_scenario_if_not_matched: !isUndefined(scenarioAction.drop_scenario_if_not_matched)
      ? scenarioAction.drop_scenario_if_not_matched
      : true,
    emit_trigger: !isUndefined(scenarioAction.emit_trigger) ? scenarioAction.emit_trigger : true,
    alarm_patterns: scenarioAction.alarm_patterns ? cloneDeep(scenarioAction.alarm_patterns) : [],
    entity_patterns: scenarioAction.entity_patterns ? cloneDeep(scenarioAction.entity_patterns) : [],
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
  author: scenario.author || '',
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
 * Convert form to scenario action
 *
 * @param {ScenarioActionForm} form
 * @param {string} [timezone = moment.tz.guess()]
 * @returns {ScenarioAction}
 */
export const formToScenarioAction = (form, timezone = moment.tz.guess()) => {
  const parametersByCurrentType = form.parameters[form.type];

  const parametersPreparers = {
    [SCENARIO_ACTION_TYPES.snooze]: (parameters = {}) =>
      ({ ...parameters, duration: formToDuration(parameters.duration) }),

    [SCENARIO_ACTION_TYPES.webhook]: (parameters = {}) =>
      ({ ...parameters, retry_delay: formToDuration(parameters.retry_delay) }),

    [SCENARIO_ACTION_TYPES.pbehavior]: (parameters = {}) => {
      const pbehavior = formToPbehavior(parameters, timezone);

      if (parameters.start_on_trigger) {
        pbehavior.start_on_trigger = parameters.start_on_trigger;
        pbehavior.duration = formToDuration(parameters.duration);
      }

      return pbehaviorToRequest(pbehavior);
    },
  };

  const prepareParametersToAction = parametersPreparers[form.type];

  return {
    ...omit(form, ['key']),

    parameters: prepareParametersToAction
      ? prepareParametersToAction(parametersByCurrentType)
      : { ...parametersByCurrentType },
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
