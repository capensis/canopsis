import { isNumber, omit } from 'lodash';

import { DEPRECATED_TRIGGERS } from '@/constants';

import { durationToForm } from '@/helpers/date/duration';
import { actionToForm, formToAction } from '@/helpers/entities/action/form';
import { flattenErrorMap } from '@/helpers/entities/shared/form';

/**
 * @typedef {Object} Trigger
 * @property {string} type
 * @property {number} [threshold]
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
 * Convert scenario entity to form.
 *
 * @param {Object} [scenario={}]
 * @param {string} [timezone]
 * @returns {Object}
 */
export const scenarioToForm = (scenario = {}, timezone) => ({
  name: scenario.name || '',
  priority: scenario.priority,
  enabled: scenario.enabled ?? true,
  delay: scenario.delay ? durationToForm(scenario.delay) : { value: undefined, unit: undefined },
  triggers: scenario.triggers ? [...scenario.triggers] : [],
  disable_during_periods: scenario.disable_during_periods
    ? [...scenario.disable_during_periods]
    : [],
  actions: scenario.actions
    ? scenario.actions.map(action => actionToForm(action, timezone))
    : [actionToForm(undefined, timezone)],
});

/**
 * Convert scenario form to API entity.
 *
 * @param {Object} form
 * @param {string} [timezone]
 * @returns {Object}
 */
export const formToScenario = (form = {}, timezone) => ({
  ...omit(form, ['delay', 'actions']),
  delay: isNumber(form.delay?.value) ? form.delay : undefined,
  actions: form.actions.map(action => formToAction(action, timezone)),
});

/**
 * Convert server errors to scenario form paths.
 *
 * @param {Object} errors
 * @param {Object} form
 * @returns {Object}
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
