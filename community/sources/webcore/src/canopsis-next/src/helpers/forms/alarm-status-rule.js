import { omit, cloneDeep } from 'lodash';

import { TIME_UNITS } from '@/constants';

import { durationToForm } from '@/helpers/date/duration';

/**
 * @typedef {Object} AlarmStatusRulePatternsForm
 * @property {Object[]} alarm_patterns
 * @property {Object[]} entity_patterns
 */

/**
 * @typedef {Object} AlarmStatusRule
 * @property {string} name
 * @property {string} description
 * @property {Duration} duration
 * @property {number} priority
 * @property {number} [freq_limit]
 * @property {Object[]} entity_patterns
 * @property {Object[]} alarm_patterns
 */

/**
 * @typedef {Object} AlarmStatusRuleForm
 * @property {string} name
 * @property {string} description
 * @property {Duration} duration
 * @property {number} priority
 * @property {AlarmStatusRulePatternsForm} patterns
 */

/**
 * Convert alarm status rule object to form compatible object
 *
 * @param {AlarmStatusRule | {}} [rule = {}]
 * @param {boolean} [flapping = false]
 * @return {AlarmStatusRuleForm}
 */
export const alarmStatusRuleToForm = (rule = {}, flapping = false) => {
  const form = {
    name: rule.name || '',
    duration: rule.duration
      ? durationToForm(rule.duration)
      : { value: 1, unit: TIME_UNITS.minute },
    priority: rule.priority || 1,
    description: rule.description || '',
    patterns: {
      alarm_patterns: rule.alarm_patterns ? cloneDeep(rule.alarm_patterns) : [],
      entity_patterns: rule.entity_patterns ? cloneDeep(rule.entity_patterns) : [],
    },
  };

  if (flapping) {
    form.freq_limit = rule.freq_limit || 1;
  }

  return form;
};

/**
 * Convert form compatible object to alarm status rule object
 *
 * @param {AlarmStatusRuleForm} form
 * @return {AlarmStatusRule}
 */
export const formToAlarmStatusRule = form => ({
  ...omit(form, ['patterns']),
  ...form.patterns,
});
