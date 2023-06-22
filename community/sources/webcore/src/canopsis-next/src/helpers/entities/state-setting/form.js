import { isNumber } from 'lodash';

import { STATE_SETTING_METHODS, STATE_SETTING_THRESHOLD_TYPES } from '@/constants';

/**
 * @typedef { 'worst' | 'worst_of_share' } StateSettingMethod
 */

/**
 * @typedef {Object} StateThreshold
 * @property {number} critical
 * @property {number} major
 * @property {number} minor
 * @property {number} type
 */

/**
 * @typedef {Object} StateSettingThresholds
 * @property {StateThreshold} errors
 * @property {StateThreshold} failures
 * @property {StateThreshold} skipped
 */

/**
 * @typedef {Object} StateSetting
 * @property {StateSettingThresholds} junit_thresholds
 * @property {StateSettingMethod} method
 * @property {string} type
 */

/**
 * @typedef {StateThreshold} StateThresholdForm
 */

/**
 * @typedef {StateSettingThresholds} StateSettingThresholdsForm
 */

/**
 * @typedef {StateSetting} StateSettingForm
 */

/**
 * Convert state setting threshold criterion to form state setting threshold criterion
 *
 * @param {StateThreshold} thresholdCriterion
 * @return {StateThresholdForm}
 */
const stateThresholdToFormStateThreshold = (thresholdCriterion = {}) => ({
  critical: thresholdCriterion.critical || 40,
  major: thresholdCriterion.major || 25,
  minor: thresholdCriterion.minor || 10,
  type: isNumber(thresholdCriterion.type)
    ? thresholdCriterion.type
    : STATE_SETTING_THRESHOLD_TYPES.percent,
});

/**
 * Convert state setting thresholds to form state setting thresholds
 *
 * @param {StateSettingThresholds} [thresholds={}]
 * @return {StateSettingThresholdsForm}
 */
const thresholdsToFormThresholds = (thresholds = {}) => ({
  errors: stateThresholdToFormStateThreshold(thresholds.errors),
  failures: stateThresholdToFormStateThreshold(thresholds.failures),
  skipped: stateThresholdToFormStateThreshold(thresholds.skipped),
});

/**
 * Convert state setting to form state setting
 *
 * @param {StateSetting} [stateSetting={}]
 * @return {StateSettingForm}
 */
export const stateSettingToForm = (stateSetting = {}) => ({
  junit_thresholds: thresholdsToFormThresholds(stateSetting.junit_thresholds),
  method: stateSetting.method || STATE_SETTING_METHODS.worstOfShare,
  type: stateSetting.type || '',
});

/**
 * Convert form state setting to state setting
 *
 * @param {StateSettingForm} form
 * @return {StateSetting}
 */
export const formToStateSetting = form => ({
  ...form,

  junit_thresholds: form.method === STATE_SETTING_METHODS.worstOfShare ? form.junit_thresholds : undefined,
});
