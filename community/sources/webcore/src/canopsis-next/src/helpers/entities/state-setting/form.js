import { omit, mapValues, pickBy } from 'lodash';

import {
  ENTITIES_STATES_KEYS,
  PATTERNS_FIELDS,
  STATE_SETTING_CONDITIONS_METHODS,
  STATE_SETTING_METHODS,
} from '@/constants';

import { filterPatternsToForm, formFilterToPatterns } from '../filter/form';

/**
 * @typedef { 'inherited' | 'dependencies_state' } StateSettingMethod
 */

/**
 * @typedef { 'ok' | 'minor' | 'major' | 'critical' } StateSettingConditionState
 */

/**
 * @typedef { 'number' | 'share' } StateSettingConditionMethod
 */

/**
 * @typedef { 'lt' | 'gt' } StateSettingConditionCondition
 */

/**
 * @typedef {Object} StateSettingCondition
 * @property {StateSettingConditionMethod} method
 * @property {StateSettingConditionState} state
 * @property {StateSettingConditionCondition} cond
 * @property {number} value
 */

/**
 * @typedef {Object<StateSettingConditionState, StateSettingCondition>} StateSettingConditions
 */

/**
 * @typedef StateSettingInherited
 * @property {PatternGroups} impacting_patterns
 */

/**
 * @typedef StateSettingDependenciesState
 * @property {StateSettingConditions} conditions
 */

/**
 * @typedef {StateSettingInherited | StateSettingDependenciesState} StateSetting
 * @property {string} title
 * @property {number} priority
 * @property {boolean} enabled
 * @property {StateSettingMethod} method
 */

/**
 * @typedef {StateSettingCondition} StateSettingConditionForm
 * @property {boolean} enabled
 */

/**
 * @typedef {StateSettingConditions} StateSettingConditionsForm
 */

/**
 * @typedef StateSettingInheritedForm
 * @property {FilterPatternsForm} impacting_patterns
 */

/**
 * @typedef StateSettingDependenciesStateForm
 * @property {StateSettingConditionsForm} conditions
 */

/**
 * @typedef {StateSettingInheritedForm | StateSettingDependenciesStateForm} StateSettingForm
 * @property {string} title
 * @property {number} priority
 * @property {boolean} enabled
 * @property {StateSettingMethod} method
 * @property {FilterPatternsForm} rule_patterns
 */

/**
 * Convert state setting condition to form
 *
 * @param {StateSettingCondition} [condition = {}]
 * @return {StateSettingConditionForm}
 */
const stateSettingConditionToForm = (condition = {}) => ({
  method: condition.method ?? STATE_SETTING_CONDITIONS_METHODS.share,
  state: condition.state ?? '',
  cond: condition.cond ?? '',
  value: condition.value ?? '',
  enabled: !!condition.method,
});

/**
 * Convert state setting conditions to form
 *
 * @param {StateSettingConditions} [conditions = {}]
 * @return {StateSettingConditionsForm}
 */
const stateSettingConditionsToForm = (conditions = {}) => (
  mapValues(ENTITIES_STATES_KEYS, stateKey => stateSettingConditionToForm(conditions[stateKey]))
);

/**
 * Convert state setting to form
 *
 * @param {StateSetting} [stateSetting = {}]
 * @return {StateSettingForm}
 */
export const stateSettingToForm = (stateSetting = {}) => ({
  title: stateSetting.title ?? '',
  priority: stateSetting.priority ?? 1,
  enabled: stateSetting.enabled ?? true,
  method: stateSetting.method ?? STATE_SETTING_METHODS.inherited,
  rule_patterns: filterPatternsToForm(
    { [PATTERNS_FIELDS.entity]: stateSetting.rule_patterns },
    [PATTERNS_FIELDS.entity],
  ),
  impacting_patterns: filterPatternsToForm(
    { [PATTERNS_FIELDS.entity]: stateSetting.impacting_patterns },
    [PATTERNS_FIELDS.entity],
  ),
  conditions: stateSettingConditionsToForm(stateSetting.conditions),
});

/**
 * Convert form to state setting condition
 *
 * @param {StateSettingConditionForm} form
 * @return {StateSettingCondition}
 */
export const formToStateSettingCondition = form => omit(form, ['enabled']);

/**
 * Convert form to state setting conditions
 *
 * @param {StateSettingConditionsForm} form
 * @return {StateSettingConditions}
 */
export const formToStateSettingConditions = form => (
  mapValues(pickBy(form, ({ enabled }) => enabled), formToStateSettingCondition)
);

/**
 * Convert form to state setting
 *
 * @param {StateSettingForm} form
 * @return {StateSetting}
 */
export const formToStateSetting = (form) => {
  const stateSetting = omit(form, ['rule_patterns', 'impacting_patterns', 'conditions']);

  stateSetting.rule_patterns = formFilterToPatterns(stateSetting.rule_patterns)[PATTERNS_FIELDS.entity];

  if (form.method === STATE_SETTING_METHODS.inherited) {
    stateSetting.impacting_patterns = formFilterToPatterns(stateSetting.impacting_patterns)[PATTERNS_FIELDS.entity];
  } else {
    stateSetting.conditions = formToStateSettingConditions(form.conditions);
  }

  return stateSetting;
};
