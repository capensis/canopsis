import { omit, isNumber } from 'lodash';

import { flattenErrorMap } from '@/helpers/forms/flatten-error-map';
import { getLocaleTimezone } from '@/helpers/date/date';

import { durationToForm } from '../date/duration';

import { formToAction, actionToForm } from './action';
import { enabledToForm } from './shared/common';

/**
 * @typedef {Object} Scenario
 * @property {string} name
 * @property {string} author
 * @property {number} priority
 * @property {boolean} enabled
 * @property {Duration} delay
 * @property {string[]} triggers
 * @property {DisableDuringPeriods[]} disable_during_periods
 * @property {Action[]} actions
 */

/**
 * @typedef {Scenario} ScenarioForm
 * @property {ActionForm[]} actions
 */

/**
 * Convert scenario to form
 *
 * @param {Scenario} [scenario = {}]
 * @param {string} [timezone = getLocaleTimezone()]
 * @returns {ScenarioForm}
 */
export const scenarioToForm = (scenario = {}, timezone = getLocaleTimezone()) => ({
  name: scenario.name || '',
  priority: scenario.priority || 1,
  enabled: enabledToForm(scenario.enabled),
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
 * @param {string} [timezone = getLocaleTimezone()]
 * @returns {Scenario}
 */
export const formToScenario = (form, timezone = getLocaleTimezone()) => ({
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
