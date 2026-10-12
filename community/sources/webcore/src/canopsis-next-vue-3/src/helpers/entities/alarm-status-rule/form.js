import { omit } from 'lodash';

import {
  MIN_DURATION_VALUE,
  MIN_FREQUENCY_LIMIT,
  MIN_PRIORITY,
  PATTERNS_FIELDS,
  TIME_UNITS,
} from '@/constants';

import { durationToForm } from '@/helpers/date/duration';
import { filterPatternsToForm, formFilterToPatterns } from '@/helpers/entities/filter/form';

/**
 * Prepare finite numeric form values with a minimum while preserving empty invalid form states.
 *
 * @param {number|string|null|undefined} value
 * @param {number} min
 * @returns {number|string|null|undefined}
 */
const prepareMinFormNumber = (value, min) => {
  if (value === '' || value === null || value === undefined) {
    return value;
  }

  const number = Number(value);

  return Number.isFinite(number) ? Math.max(number, min) : value;
};

/**
 * @typedef {FilterPatterns} AlarmStatusRule
 * @property {boolean} enabled
 * @property {string} name
 * @property {string} description
 * @property {Duration} duration
 * @property {number} priority
 * @property {number} [freq_limit]
 */

/**
 * @typedef {Object} AlarmStatusRuleForm
 * @property {boolean} enabled
 * @property {string} name
 * @property {string} description
 * @property {Duration} duration
 * @property {number} priority
 * @property {FilterPatternsForm} patterns
 */

/**
 * Convert alarm status rule object to form compatible object.
 *
 * @param {Object} [rule={}]
 * @param {boolean} [flapping=false]
 * @return {Object}
 */
export const alarmStatusRuleToForm = (rule = {}, flapping = false) => {
  const form = {
    name: rule.name ?? '',
    duration: rule.duration
      ? durationToForm(rule.duration)
      : { value: MIN_DURATION_VALUE, unit: TIME_UNITS.minute },
    priority: rule.priority,
    description: rule.description ?? '',
    enabled: rule.enabled ?? true,
    patterns: filterPatternsToForm(rule, [PATTERNS_FIELDS.alarm, PATTERNS_FIELDS.entity]),
  };

  if (flapping) {
    form.freq_limit = rule.freq_limit ?? MIN_FREQUENCY_LIMIT;
  }

  return form;
};

/**
 * Convert form compatible object to alarm status rule object.
 *
 * @param {Object} form
 * @return {Object}
 */
export const formToAlarmStatusRule = form => ({
  ...prepareAlarmStatusRuleNumbers(omit(form, ['patterns'])),
  ...formFilterToPatterns(form.patterns),
});

/**
 * Prepare numeric fields with non-negative/non-zero lower bounds before saving.
 *
 * @param {Object} form
 * @returns {Object}
 */
export const prepareAlarmStatusRuleNumbers = (form) => {
  const preparedForm = { ...form };

  if (Object.hasOwn(preparedForm, 'priority')) {
    preparedForm.priority = prepareMinFormNumber(preparedForm.priority, MIN_PRIORITY);
  }

  if (preparedForm.duration) {
    preparedForm.duration = {
      ...form.duration,
      value: prepareMinFormNumber(form.duration.value, MIN_DURATION_VALUE),
    };
  }

  if (Object.hasOwn(preparedForm, 'freq_limit')) {
    preparedForm.freq_limit = prepareMinFormNumber(
      preparedForm.freq_limit,
      MIN_FREQUENCY_LIMIT,
    );
  }

  return preparedForm;
};
