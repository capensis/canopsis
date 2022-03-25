import { QUICK_RANGES, TIME_UNITS } from '@/constants';

import uid from '@/helpers/uid';

/**
 * @typedef { 'alarm' | 'entity' | 'pbehavior' } PatternTypes
 */

/**
 * @typedef {Object} PatternRuleCondition
 * @property {string} type
 * @property {Object | string | number} value
 */

/**
 * @typedef {Object} PatternRule
 * @property {PatternRuleCondition} cond
 * @property {string} field
 * @property {string} field_type
 */

/**
 * @typedef {PatternRule[]} PatternRules
 */

/**
 * @typedef {PatternRules[]} PatternGroups
 */

/**
 * @typedef {Object} Pattern
 * @property {string | Symbol} id
 * @property {PatternTypes} type
 * @property {PatternGroups} alarm_pattern
 * @property {PatternGroups} entity_pattern
 * @property {PatternGroups} pbehavior_pattern
 * @property {PatternGroups} event_pattern
 */

/**
 * @typedef {Object} PatternRuleRangeForm
 * @property {string} type
 * @property {string | number} [from]
 * @property {string | number} [to]
 */

/**
 * @typedef {Object} PatternRuleForm
 * @property {string} key
 * @property {string} attribute
 * @property {string} operator
 * @property {string} field
 * @property {string} dictionary
 * @property {number | string} value
 * @property {PatternRuleRangeForm} range
 * @property {Duration} duration
 */

/**
 * @typedef {Object} PatternGroupForm
 * @property {PatternRuleForm[]} rules
 * @property {string} key
 */

/**
 * @typedef {Pattern} PatternForm
 * @property {PatternGroupForm[]} groups
 */

/**
 * Convert pattern rule to form
 *
 * @param {PatternRule} rule
 * @return {PatternRuleForm}
 */
export const patternRuleToForm = (rule = {}) => ({
  /** TODO: Should be finished with create filter form */
  key: uid(),
  attribute: rule.field || 'v.component',
  operator: 'equal',
  field: '',
  dictionary: '',
  value: rule?.cond?.value ?? '',
  range: {
    type: QUICK_RANGES.last1Hour.value,
    from: 0,
    to: 0,
  },
  duration: {
    value: 1,
    unit: TIME_UNITS.second,
  },
  ...rule,
});

/**
 * Convert pattern rules to form
 *
 * @param {PatternRules} rules
 * @return {PatternRuleForm[]}
 */
export const patternRulesToForm = (rules = [undefined]) => rules.map(patternRuleToForm);

/**
 * Convert pattern rules to group form
 *
 * @param {PatternRules} rules
 * @return {PatternGroupForm}
 */
export const patternRulesToGroup = rules => ({
  key: uid(),
  rules: patternRulesToForm(rules),
});

const patternsToGroups = (patterns = [undefined]) => patterns.map(patternRulesToGroup);

/**
 * Convert pattern to pattern form
 *
 * @param {Pattern} pattern
 * @return {PatternForm}
 */
export const patternToForm = (pattern = {}) => ({
  ...pattern,
  id: pattern.id || '',
  groups: patternsToGroups(
    pattern.alarm_pattern
    || pattern.entity_pattern
    || pattern.pbehavior_pattern
    || pattern.event_pattern,
  ),
});
