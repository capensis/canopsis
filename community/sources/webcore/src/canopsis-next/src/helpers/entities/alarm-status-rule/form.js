import { omit } from 'lodash';

import { PATTERNS_FIELDS, TIME_UNITS } from '@/constants';

import { durationToForm } from '@/helpers/date/duration';
import { filterPatternsToForm, formFilterToPatterns } from '@/helpers/entities/filter/form';

/**
 * @typedef {FilterPatterns} AlarmStatusRule
 * @property {string} name
 * @property {string} description
 * @property {Duration} duration
 * @property {number} priority
 * @property {number} [freq_limit]
 */

/**
 * @typedef {Object} AlarmStatusRuleForm
 * @property {string} name
 * @property {string} description
 * @property {Duration} duration
 * @property {number} priority
 * @property {FilterPatternsForm} patterns
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
    name: rule.name ?? '',
    duration: rule.duration
      ? durationToForm(rule.duration)
      : { value: 1, unit: TIME_UNITS.minute },
    priority: rule.priority,
    description: rule.description ?? '',
    patterns: filterPatternsToForm(rule, [PATTERNS_FIELDS.alarm, PATTERNS_FIELDS.entity]),
  };

  if (flapping) {
    form.freq_limit = rule.freq_limit ?? 1;
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
  ...formFilterToPatterns(form.patterns),
});
