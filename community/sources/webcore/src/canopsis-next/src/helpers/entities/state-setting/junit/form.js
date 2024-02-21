import { isNumber } from 'lodash';

import { JUNIT_STATE_SETTING_METHODS, JUNIT_STATE_SETTING_THRESHOLDS_TYPES } from '@/constants';

/**
 * @typedef { 'worst' | 'worst_of_share' } JunitStateSettingMethod
 */

/**
 * @typedef {Object} JunitStateThreshold
 * @property {number} critical
 * @property {number} major
 * @property {number} minor
 * @property {number} type
 */

/**
 * @typedef {Object} JunitStateSettingThresholds
 * @property {JunitStateThreshold} errors
 * @property {JunitStateThreshold} failures
 * @property {JunitStateThreshold} skipped
 */

/**
 * @typedef {Object} JunitStateSetting
 * @property {StateSettingThresholds} junit_thresholds
 * @property {StateSettingMethod} method
 * @property {string} type
 */

/**
 * @typedef {JunitStateThreshold} JunitStateThresholdForm
 */

/**
 * @typedef {JunitStateSettingThresholds} JunitStateSettingThresholdsForm
 */

/**
 * @typedef {JunitStateSetting} JunitStateSettingForm
 */

/**
 * Convert state setting threshold criterion to form state setting threshold criterion
 *
 * @param {JunitStateThreshold} thresholdCriterion
 * @return {JunitStateThresholdForm}
 */
const junitStateThresholdToFormStateThreshold = (thresholdCriterion = {}) => ({
  critical: thresholdCriterion.critical || 40,
  major: thresholdCriterion.major || 25,
  minor: thresholdCriterion.minor || 10,
  type: isNumber(thresholdCriterion.type)
    ? thresholdCriterion.type
    : JUNIT_STATE_SETTING_THRESHOLDS_TYPES.percent,
});

/**
 * Convert state setting thresholds to form state setting thresholds
 *
 * @param {JunitStateSettingThresholds} [thresholds={}]
 * @return {JunitStateSettingThresholdsForm}
 */
const junitStateSettingThresholdsToForm = (thresholds = {}) => ({
  errors: junitStateThresholdToFormStateThreshold(thresholds.errors),
  failures: junitStateThresholdToFormStateThreshold(thresholds.failures),
  skipped: junitStateThresholdToFormStateThreshold(thresholds.skipped),
});

/**
 * Convert state setting to form state setting
 *
 * @param {JunitStateSetting} [stateSetting={}]
 * @return {JunitStateSettingForm}
 */
export const junitStateSettingToForm = (stateSetting = {}) => ({
  junit_thresholds: junitStateSettingThresholdsToForm(stateSetting.junit_thresholds),
  method: stateSetting.method || JUNIT_STATE_SETTING_METHODS.worstOfShare,
  type: stateSetting.type || '',
});

/**
 * Convert form state setting to state setting
 *
 * @param {JunitStateSettingForm} form
 * @return {JunitStateSetting}
 */
export const formToJunitStateSetting = form => ({
  ...form,

  junit_thresholds: form.method === JUNIT_STATE_SETTING_METHODS.worstOfShare ? form.junit_thresholds : undefined,
});
