import { isNumber } from 'lodash';

import { durationToForm, formToDuration } from '@/helpers/date/duration';

/**
 * @typedef {Object} ScenarioAction
 * @property {string} type
 * @property {boolean} drop_scenario_if_not_matched
 * @property {boolean} emit_trigger
 * @property {Object} parameters
 * @property {Object} entity_patterns
 * @property {Object} alarm_patterns
 */

/**
 * @typedef {Object} Scenario
 * @property {string} name
 * @property {boolean} enabled
 * @property {number} priority
 * @property {string} [author]
 * @property {Array} triggers
 * @property {Duration|DurationForm} delay
 * @property {Array} disable_during_periods
 * @property {ScenarioAction[]} actions
 */

/**
 * Convert scenario entity to form object
 *
 * @param {Scenario} scenario
 * @returns {Scenario}
 */
export const scenarioToForm = (scenario = {}) => ({
  name: scenario.name || '',
  enabled: !!scenario.enabled,
  priority: scenario.priority || 0,
  triggers: scenario.triggers || [],
  delay: scenario.delay
    ? durationToForm(scenario.delay)
    : { value: undefined, unit: undefined },
  disable_during_periods: scenario.disable_during_periods || [],
  actions: scenario.actions || [],
});

/**
 * Convert form scenario to API object
 *
 * @param {Scenario} form
 * @returns {Scenario}
 */
export const formToScenario = (form = {}) => ({
  ...form,
  delay: form.delay && isNumber(form.delay.value)
    ? formToDuration(form.delay)
    : undefined,
});
