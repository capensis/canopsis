import { omit, isNumber } from 'lodash';

import { DEPRECATED_TRIGGERS } from '@/constants';

import { durationToForm } from '@/helpers/date/duration';
import { formToAction, actionToForm } from '@/helpers/entities/action';
import { getLocalTimezone } from '@/helpers/date/date';
import { flattenErrorMap } from '@/helpers/entities/shared/form';

/**
 * @typedef {Object} Trigger
 * @property {string} type
 * @property {number} [threshold]
 */

/**
 * @typedef {Object} Scenario
 * @property {string} name
 * @property {string} author
 * @property {number} priority
 * @property {boolean} enabled
 * @property {Duration} delay
 * @property {Trigger[]} triggers
 * @property {DisableDuringPeriods[]} disable_during_periods
 * @property {Action[]} actions
 */

/**
 * @typedef {Scenario} ScenarioForm
 * @property {ActionForm[]} actions
 */

/**
 * Check trigger is deprecated
 *
 * @param {string} type
 * @returns {boolean}
 */
export const isDeprecatedTrigger = type => DEPRECATED_TRIGGERS.includes(type);
/**
 * Convert scenario to form
 *
 * @param {Scenario} [scenario = {}]
 * @param {string} [timezone = getLocalTimezone()]
 * @returns {ScenarioForm}
 */
export const scenarioToForm = (scenario = {}, timezone = getLocalTimezone()) => ({
  name: scenario.name || '',
  priority: scenario.priority,
  enabled: scenario.enabled ?? true,
  delay: scenario.delay
    ? durationToForm(scenario.delay)
    : { value: undefined, unit: undefined },
  triggers: scenario.triggers ? [...scenario.triggers] : [],
  disable_during_periods: scenario.disable_during_periods ? [...scenario.disable_during_periods] : [],
  actions: scenario.actions
    ? scenario.actions.map(action => actionToForm(action, timezone))
    : [actionToForm(undefined, timezone)],
});

/**
 * Convert form to scenario
 *
 * @param {ScenarioForm} form
 * @param {string} [timezone = getLocalTimezone()]
 * @returns {Scenario}
 */
export const formToScenario = (form, timezone = getLocalTimezone()) => ({
  ...omit(form, ['delay', 'actions']),
  delay: isNumber(form.delay?.value)
    ? form.delay
    : undefined,
  actions: form.actions.map(action => formToAction(action, timezone)),
});

/**
 * Convert error structure to form structure
 *
 * @param {FlattenErrors} errors
 * @param {ScenarioForm} form
 * @return {FlattenErrors}
 */
export const scenarioErrorToForm = (errors, form) => {
  const prepareActionsErrors = (errorsObject) => {
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

  return flattenErrorMap(errors, prepareActionsErrors);
};
