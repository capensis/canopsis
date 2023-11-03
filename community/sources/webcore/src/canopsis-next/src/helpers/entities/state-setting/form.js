import { omit, mapValues, pickBy } from 'lodash';

import {
  ENTITIES_STATES_KEYS,
  PATTERNS_FIELDS,
  STATE_SETTING_THRESHOLDS_METHODS,
  STATE_SETTING_METHODS,
} from '@/constants';

import { filterPatternsToForm, formFilterToPatterns } from '../filter/form';

/**
 * @typedef { 'inherited' | 'dependencies_state' } StateSettingMethod
 */

/**
 * @typedef { 'ok' | 'minor' | 'major' | 'critical' } StateSettingThresholdState
 */

/**
 * @typedef { 'number' | 'share' } StateSettingThresholdMethod
 */

/**
 * @typedef { 'lt' | 'gt' } StateSettingThresholdCondition
 */

/**
 * @typedef {Object} StateSettingThreshold
 * @property {StateSettingThresholdMethod} method
 * @property {StateSettingThresholdState} state
 * @property {StateSettingThresholdCondition} cond
 * @property {number} value
 */

/**
 * @typedef {Object<StateSettingThresholdState, StateSettingThreshold>} StateSettingThresholds
 */

/**
 * @typedef StateSettingInherited
 * @property {PatternGroups} inherited_entity_pattern
 */

/**
 * @typedef StateSettingDependenciesState
 * @property {StateSettingThresholds} state_thresholds
 */

/**
 * @typedef {StateSettingInherited | StateSettingDependenciesState} StateSetting
 * @property {string} title
 * @property {number} priority
 * @property {boolean} enabled
 * @property {StateSettingMethod} method
 */

/**
 * @typedef {StateSettingThreshold} StateSettingThresholdForm
 * @property {boolean} enabled
 */

/**
 * @typedef {StateSettingThresholds} StateSettingThresholdsForm
 */

/**
 * @typedef StateSettingInheritedForm
 * @property {FilterPatternsForm} inherited_entity_pattern
 */

/**
 * @typedef StateSettingDependenciesStateForm
 * @property {StateSettingThresholdsForm} state_thresholds
 */

/**
 * @typedef {StateSettingInheritedForm | StateSettingDependenciesStateForm} StateSettingForm
 * @property {string} title
 * @property {number} priority
 * @property {boolean} enabled
 * @property {StateSettingMethod} method
 * @property {FilterPatternsForm} entity_pattern
 */

/**
 * Convert state setting threshold to form
 *
 * @param {StateSettingThreshold} [threshold = {}]
 * @return {StateSettingThresholdForm}
 */
const stateSettingThresholdToForm = (threshold = {}) => ({
  method: threshold.method ?? STATE_SETTING_THRESHOLDS_METHODS.share,
  state: threshold.state ?? '',
  cond: threshold.cond ?? '',
  value: threshold.value ?? '',
  enabled: !!threshold.method,
});

/**
 * Convert state setting thresholds to form
 *
 * @param {StateSettingThresholds} [thresholds = {}]
 * @return {StateSettingThresholdsForm}
 */
const stateSettingThresholdsToForm = (thresholds = {}) => (
  mapValues(ENTITIES_STATES_KEYS, stateKey => stateSettingThresholdToForm(thresholds[stateKey]))
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
  entity_pattern: filterPatternsToForm(
    { [PATTERNS_FIELDS.entity]: stateSetting.entity_pattern },
    [PATTERNS_FIELDS.entity],
  ),
  inherited_entity_pattern: filterPatternsToForm(
    { [PATTERNS_FIELDS.entity]: stateSetting.inherited_entity_pattern },
    [PATTERNS_FIELDS.entity],
  ),
  state_thresholds: stateSettingThresholdsToForm(stateSetting.state_thresholds),
});

/**
 * Convert form to state setting threshold
 *
 * @param {StateSettingThresholdForm} form
 * @return {StateSettingThreshold}
 */
export const formToStateSettingThreshold = form => omit(form, ['enabled']);

/**
 * Convert form to state setting thresholds
 *
 * @param {StateSettingThresholdsForm} form
 * @return {StateSettingThresholds}
 */
export const formToStateSettingThresholds = form => (
  mapValues(pickBy(form, ({ enabled }) => enabled), formToStateSettingThreshold)
);

/**
 * Convert form to state setting
 *
 * @param {StateSettingForm} form
 * @return {StateSetting}
 */
export const formToStateSetting = (form) => {
  const stateSetting = omit(form, ['entity_pattern', 'inherited_entity_pattern', 'state_thresholds']);

  stateSetting.entity_pattern = formFilterToPatterns(form.entity_pattern)[PATTERNS_FIELDS.entity];

  if (form.method === STATE_SETTING_METHODS.inherited) {
    stateSetting.inherited_entity_pattern = (
      formFilterToPatterns(form.inherited_entity_pattern)[PATTERNS_FIELDS.entity]
    );
  } else {
    stateSetting.state_thresholds = formToStateSettingThresholds(form.state_thresholds);
  }

  return stateSetting;
};
