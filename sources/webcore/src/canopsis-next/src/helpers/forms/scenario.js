import moment from 'moment';
import { isUndefined, cloneDeep, omit } from 'lodash';

import { SCENARIO_ACTION_TYPES } from '@/constants';

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
 * @typedef {ScenarioAction} ScenarioActionForm
 * @property {
 *   PbehaviorForm |
 *   ScenarioActionDefaultParameters |
 *   ScenarioActionSnoozeFormParameters |
 *   ScenarioActionChangeStateParameters |
 *   ScenarioActionWebhookFormParameters |
 *   ScenarioActionAssocTicketParameters
 * } parameters
 */

/**
 * @typedef {Scenario} ScenarioForm
 * @property {DurationForm} delay
 * @property {ScenarioActionForm[]} actions
 */

/**
 * Convert scenario action to form
 *
 * @param {ScenarioAction} [scenarioAction = {}]
 * @param {string} [timezone = moment.tz.guess()]
 * @returns {ScenarioActionForm}
 */
export function scenarioActionToForm(scenarioAction = {}, timezone = moment.tz.guess()) {
  const parametersPreparers = {
    [SCENARIO_ACTION_TYPES.snooze]: (parameters = {}) =>
      ({ ...parameters, duration: durationToForm(parameters.duration) }),

    [SCENARIO_ACTION_TYPES.webhook]: (parameters = {}) =>
      ({
        request: cloneDeep(parameters.request || {}),
        declare_ticket: cloneDeep(parameters.declare_ticket || {}),
        retry_count: parameters.retry_count || '',
        retry_delay: durationToForm(parameters.retry_delay),
      }),

    [SCENARIO_ACTION_TYPES.pbehavior]: (parameters = {}) =>
      pbehaviorToForm(parameters, null, timezone),
  };

  return {
    key: uid(),
    type: scenarioAction.type,
    drop_scenario_if_not_matched: scenarioAction.drop_scenario_if_not_matched,
    emit_trigger: scenarioAction.emit_trigger,
    alarm_patterns: scenarioAction.alarm_patterns ? cloneDeep(scenarioAction.alarm_patterns) : [],
    entity_patterns: scenarioAction.entity_patterns ? cloneDeep(scenarioAction.entity_patterns) : [],

    parameters: parametersPreparers[scenarioAction.type]
      ? parametersPreparers[scenarioAction.type](scenarioAction.parameters)
      : { ...scenarioAction.parameters },
  };
}

/**
 * Convert scenario to form
 *
 * @param {Scenario} [scenario = {}]
 * @param {string} [timezone = moment.tz.guess()]
 * @returns {ScenarioForm}
 */
export function scenarioToForm(scenario = {}, timezone = moment.tz.guess()) {
  return {
    name: scenario.name || '',
    author: scenario.author || '',
    priority: scenario.priority || 0,
    enabled: isUndefined(scenario.enabled) ? true : scenario.enabled,
    delay: durationToForm(scenario.delay),
    triggers: scenario.triggers ? [...scenario.triggers] : [],
    disable_during_periods: scenario.disable_during_periods ? [...scenario.disable_during_periods] : [],
    actions: scenario.actions ? scenario.actions.map(action => scenarioActionToForm(action, timezone)) : [],
  };
}

/**
 * Convert form to scenario action
 *
 * @param {ScenarioActionForm} form
 * @param {string} [timezone = moment.tz.guess()]
 * @returns {ScenarioAction}
 */
export function formToScenarioAction(form, timezone = moment.tz.guess()) {
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

  return {
    ...omit(form, ['key']),

    parameters: parametersPreparers[form.type]
      ? parametersPreparers[form.type](form.parameters)
      : { ...form.parameters },
  };
}

/**
 * Convert form to scenario
 *
 * @param {ScenarioForm} form
 * @param {string} [timezone = moment.tz.guess()]
 * @returns {Scenario}
 */
export function formToScenario(form, timezone = moment.tz.guess()) {
  return {
    ...omit(form, ['delay', 'actions']),

    delay: formToDuration(form.delay),
    actions: form.actions.map(action => formToScenarioAction(action, timezone)),
  };
}
