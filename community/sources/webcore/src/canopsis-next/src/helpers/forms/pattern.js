import { isBoolean } from 'lodash';

import {
  ALARM_PATTERN_FIELDS,
  ENTITY_PATTERN_FIELDS,
  PATTERN_CONDITIONS,
  PATTERN_CUSTOM_ITEM_VALUE,
  PATTERN_INFOS_NAME_OPERATORS,
  PATTERN_INPUT_TYPES,
  PATTERN_OPERATORS,
  PATTERN_QUICK_RANGES,
  PATTERN_RULE_INFOS_FIELDS,
  PATTERN_TYPES,
  QUICK_RANGES,
} from '@/constants';

import uid from '@/helpers/uid';
import { getValueType } from '@/helpers/pattern';
import { convertDateToTimestamp } from '@/helpers/date/date';
import {
  getDiffBetweenStartAndStopQuickInterval,
  getQuickRangeByDiffBetweenStartAndStop,
} from '@/helpers/date/date-intervals';
import { durationToForm } from '@/helpers/date/duration';

/**
 * @typedef { 'alarm' | 'entity' | 'pbehavior' } PatternTypes
 */

/**
 * @typedef { 'alarm_pattern' | 'entity_pattern' | 'pbehavior_pattern' } PatternsField
 */

/**
 * @typedef {PatternsField[]} PatternsFields
 */

/**
 * @typedef {Object} PatternRuleRangeCondition
 * @property {number} from
 * @property {number} to
 */

/**
 * @typedef {Object} PatternRuleCondition
 * @property {string} type
 * @property {PatternRuleRangeCondition | Duration | string | number} value
 */

/**
 * @typedef {Object} PatternRule
 * @property {PatternRuleCondition} cond
 * @property {string} field
 * @property {string} [field_type]
 */

/**
 * @typedef {PatternRule[]} PatternRules
 */

/**
 * @typedef {PatternRules[]} PatternGroups
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
 * @typedef {PatternRuleForm[]} PatternGroupsForm
 */

/**
 * @typedef {Pattern} PatternForm
 * @property {PatternGroupsForm} groups
 */

/**
 * Check pattern is date
 *
 * @param {string} value
 * @return {boolean}
 */
const isDatePatternRule = value => [
  ALARM_PATTERN_FIELDS.creationDate,
  ALARM_PATTERN_FIELDS.lastEventDate,
  ALARM_PATTERN_FIELDS.lastUpdateDate,
  ALARM_PATTERN_FIELDS.ackAt,
  ALARM_PATTERN_FIELDS.resolvedAt,
  ENTITY_PATTERN_FIELDS.lastEventDate,
].includes(value);

/**
 * Check pattern is infos
 *
 * @param {string} value
 * @return {boolean}
 */
const isInfosPatternRuleAttribute = value => [
  ALARM_PATTERN_FIELDS.infos,
  ENTITY_PATTERN_FIELDS.infos,
].includes(value);

/**
 * Convert pattern rule to form
 *
 * @param {PatternRule} rule
 * @return {PatternRuleForm}
 */
export const patternRuleToForm = (rule = {}) => {
  const form = {
    key: uid(),
    attribute: '',
    operator: '',
    field: '',
    dictionary: '',
    value: '',
    range: {
      type: QUICK_RANGES.last1Hour.value,
      from: 0,
      to: 0,
    },
    duration: durationToForm(),
  };

  if (!rule.cond) {
    return form;
  }

  const isAlarmInfos = rule.field?.startsWith(ALARM_PATTERN_FIELDS.infos);
  const isEntityInfos = rule.field?.startsWith(ENTITY_PATTERN_FIELDS.infos);
  const isDuration = rule.field === ALARM_PATTERN_FIELDS.duration;
  const isInfos = isAlarmInfos || isEntityInfos;

  switch (rule.cond.type) {
    case PATTERN_CONDITIONS.equal: {
      if (isBoolean(rule.cond.value)) {
        if (rule.cond.value) {
          form.operator = {
            [ALARM_PATTERN_FIELDS.snooze]: PATTERN_OPERATORS.snoozed,
            [ALARM_PATTERN_FIELDS.ack]: PATTERN_OPERATORS.acked,
            [ALARM_PATTERN_FIELDS.canceled]: PATTERN_OPERATORS.canceled,
            [ALARM_PATTERN_FIELDS.ticket]: PATTERN_OPERATORS.ticketAssociated,
          }[rule.field];
        } else {
          form.operator = {
            [ALARM_PATTERN_FIELDS.snooze]: PATTERN_OPERATORS.notSnoozed,
            [ALARM_PATTERN_FIELDS.ack]: PATTERN_OPERATORS.notAcked,
            [ALARM_PATTERN_FIELDS.canceled]: PATTERN_OPERATORS.notCanceled,
            [ALARM_PATTERN_FIELDS.ticket]: PATTERN_OPERATORS.ticketNotAssociated,
          }[rule.field];
        }
      }

      if (!form.operator) {
        form.operator = PATTERN_OPERATORS.equal;
        form.value = rule.cond.value;
      }
      break;
    }
    case PATTERN_CONDITIONS.notEqual: {
      form.operator = PATTERN_OPERATORS.notEqual;
      form.value = rule.cond.value;
      break;
    }

    case PATTERN_CONDITIONS.greater:
      form.operator = isDuration
        ? PATTERN_OPERATORS.longer
        : PATTERN_OPERATORS.higher;
      form.value = rule.cond.value;
      break;
    case PATTERN_CONDITIONS.less:
      form.operator = isDuration
        ? PATTERN_OPERATORS.shorter
        : PATTERN_OPERATORS.longer;
      form.value = rule.cond.value;
      break;

    case PATTERN_CONDITIONS.exist:
      form.operator = rule.cond.value === true
        ? PATTERN_OPERATORS.exist
        : PATTERN_OPERATORS.notExist;
      break;

    case PATTERN_CONDITIONS.hasEvery:
      form.operator = PATTERN_OPERATORS.hasEvery;
      form.value = rule.cond.value;
      break;
    case PATTERN_CONDITIONS.hasOneOf:
      form.operator = PATTERN_OPERATORS.hasOneOf;
      form.value = rule.cond.value;
      break;
    case PATTERN_CONDITIONS.hasNot:
      form.operator = PATTERN_OPERATORS.hasNot;
      form.value = rule.cond.value;
      break;
    case PATTERN_CONDITIONS.isEmpty:
      form.operator = rule.cond.value === true
        ? PATTERN_OPERATORS.isEmpty
        : PATTERN_OPERATORS.isNotEmpty;
      break;

    case PATTERN_CONDITIONS.regexp: {
      const notContainsMatch = rule.cond.value.match(/^\^\(\(\?!(?<value>.+)\)\.\)\*\$$/);

      if (notContainsMatch?.groups?.value) {
        form.operator = PATTERN_OPERATORS.notContains;
        form.value = notContainsMatch.groups.value;
        break;
      }

      const notEndWithMatch = rule.cond.value.match(/^\(\?<!(?<value>.+)\)\$$/);

      if (notEndWithMatch?.groups?.value) {
        form.operator = PATTERN_OPERATORS.notEndWith;
        form.value = notEndWithMatch.groups.value;
        break;
      }

      const notBeginWithMatch = rule.cond.value.match(/^\^\(\?!(?<value>.+)\)$/);

      if (notBeginWithMatch?.groups?.value) {
        form.operator = PATTERN_OPERATORS.notBeginWith;
        form.value = notBeginWithMatch.groups.value;
        break;
      }

      const endsWithMatch = rule.cond.value.match(/(?<value>.+)\$$/);

      if (endsWithMatch?.groups?.value) {
        form.operator = PATTERN_OPERATORS.endsWith;
        form.value = endsWithMatch.groups.value;
        break;
      }

      const beginsWithMatch = rule.cond.value.match(/^\^(?<value>.+)/);

      if (beginsWithMatch?.groups?.value) {
        form.operator = PATTERN_OPERATORS.beginsWith;
        form.value = beginsWithMatch.groups.value;
        break;
      }

      form.operator = PATTERN_OPERATORS.contains;
      form.value = rule.cond.value;
      break;
    }
    case PATTERN_CONDITIONS.relativeTime:
      form.range.type = getQuickRangeByDiffBetweenStartAndStop(rule.cond.value, PATTERN_QUICK_RANGES).value;
      break;
    case PATTERN_CONDITIONS.absoluteTime:
      form.range = {
        type: QUICK_RANGES.custom.value,
        ...rule.cond.value,
      };
      break;
  }

  if (isDuration) {
    form.duration = durationToForm(rule.cond.value);
  }

  if (isInfos) {
    const infosPrefix = isAlarmInfos ? ALARM_PATTERN_FIELDS.infos : ENTITY_PATTERN_FIELDS.infos;

    form.attribute = rule.field.slice(0, infosPrefix.length);
    form.dictionary = rule.field.slice(infosPrefix.length + 1);

    form.field = PATTERN_INFOS_NAME_OPERATORS.includes(rule.cond.type)
      ? PATTERN_RULE_INFOS_FIELDS.name
      : PATTERN_RULE_INFOS_FIELDS.value;
  } else {
    form.attribute = rule.field;
  }

  return form;
};

/**
 * Convert pattern rules to form
 *
 * @param {PatternRules} [rules = []]
 * @return {PatternRuleForm[]}
 */
export const patternRulesToForm = (rules = []) => rules.map(patternRuleToForm);

/**
 * Convert pattern rules to group form
 *
 * @param {PatternRules} [rules = []]
 * @return {PatternGroupForm}
 */
export const patternRulesToGroup = rules => ({
  key: uid(),
  rules: patternRulesToForm(rules),
});

const patternsToGroups = (patterns = []) => patterns.map(patternRulesToGroup);

/**
 * Convert pattern to pattern form
 *
 * @param {Pattern} pattern
 * @return {PatternForm}
 */
export const patternToForm = (pattern = {}) => ({
  ...pattern,
  title: pattern.title ?? '',
  id: pattern.id ?? PATTERN_CUSTOM_ITEM_VALUE,
  type: pattern.type ?? PATTERN_TYPES.alarm,
  is_corporate: pattern.is_corporate ?? false,
  groups: patternsToGroups(
    pattern.alarm_pattern
    || pattern.entity_pattern
    || pattern.pbehavior_pattern
    || pattern.event_pattern,
  ),
});

/**
 * Convert range to pattern condition
 *
 * @param {PatternRuleRangeForm} range
 * @return {PatternRuleCondition}
 */
export const formDateIntervalConditionToPatternRuleCondition = (range) => {
  if (range.type === QUICK_RANGES.custom.value) {
    return {
      type: PATTERN_CONDITIONS.absoluteTime,
      value: {
        from: convertDateToTimestamp(range.from),
        to: convertDateToTimestamp(range.to),
      },
    };
  }

  return {
    value: getDiffBetweenStartAndStopQuickInterval(range.type),
    type: PATTERN_CONDITIONS.relativeTime,
  };
};

/**
 * Convert pattern form rule to pattern rule
 *
 * @param {PatternRuleForm} rule
 * @return {PatternRule}
 */
export const formRuleToPatternRule = (rule) => {
  const pattern = {
    field: rule.attribute,
    cond: {
      value: rule.value,
      type: PATTERN_CONDITIONS.equal,
    },
  };

  const isInfos = isInfosPatternRuleAttribute(rule.attribute);
  const isDate = isDatePatternRule(rule.attribute);

  if (isInfos) {
    pattern.field = [rule.attribute, rule.dictionary].join('.');
  }

  if (isDate) {
    pattern.cond = formDateIntervalConditionToPatternRuleCondition(rule.range);

    return pattern;
  }

  if (isInfos && rule.field !== PATTERN_RULE_INFOS_FIELDS.name) {
    pattern.field_type = getValueType(rule.value);
  }

  switch (rule.operator) {
    case PATTERN_OPERATORS.equal:
      pattern.cond.type = PATTERN_CONDITIONS.equal;
      break;
    case PATTERN_OPERATORS.notEqual:
      pattern.cond.type = PATTERN_CONDITIONS.notEqual;
      break;
    case PATTERN_OPERATORS.contains:
      pattern.cond.type = PATTERN_CONDITIONS.regexp;
      break;
    case PATTERN_OPERATORS.notContains:
      pattern.cond.type = PATTERN_CONDITIONS.regexp;
      pattern.cond.value = `^((?!${rule.value}).)*$`;
      break;

    case PATTERN_OPERATORS.beginsWith:
      pattern.cond.type = PATTERN_CONDITIONS.regexp;
      pattern.cond.value = `^${rule.value}`;
      break;
    case PATTERN_OPERATORS.notBeginWith:
      pattern.cond.type = PATTERN_CONDITIONS.regexp;
      pattern.cond.value = `^(?!${rule.value})`;
      break;
    case PATTERN_OPERATORS.endsWith:
      pattern.cond.type = PATTERN_CONDITIONS.regexp;
      pattern.cond.value = `${rule.value}$`;
      break;
    case PATTERN_OPERATORS.notEndWith:
      pattern.cond.type = PATTERN_CONDITIONS.regexp;
      pattern.cond.value = `(?<!${rule.value})$`;
      break;

    case PATTERN_OPERATORS.exist:
      pattern.cond.type = PATTERN_CONDITIONS.exist;
      pattern.cond.value = true;
      break;
    case PATTERN_OPERATORS.notExist:
      pattern.cond.type = PATTERN_CONDITIONS.exist;
      pattern.cond.value = false;
      break;

    case PATTERN_OPERATORS.hasEvery:
      pattern.cond.type = PATTERN_CONDITIONS.hasEvery;
      pattern.field_type = PATTERN_INPUT_TYPES.array;
      break;
    case PATTERN_OPERATORS.hasOneOf:
      pattern.cond.type = PATTERN_CONDITIONS.hasOneOf;
      pattern.field_type = PATTERN_INPUT_TYPES.array;
      break;
    case PATTERN_OPERATORS.hasNot:
      pattern.cond.type = PATTERN_CONDITIONS.hasNot;
      pattern.field_type = PATTERN_INPUT_TYPES.array;
      break;
    case PATTERN_OPERATORS.isEmpty:
      pattern.cond.type = PATTERN_CONDITIONS.isEmpty;
      pattern.field_type = PATTERN_INPUT_TYPES.array;
      pattern.cond.value = true;
      break;
    case PATTERN_OPERATORS.isNotEmpty:
      pattern.cond.type = PATTERN_CONDITIONS.isEmpty;
      pattern.field_type = PATTERN_INPUT_TYPES.array;
      pattern.cond.value = false;
      break;

    case PATTERN_OPERATORS.higher:
      pattern.cond.type = PATTERN_CONDITIONS.greater;
      break;

    case PATTERN_OPERATORS.lower:
      pattern.cond.type = PATTERN_CONDITIONS.less;
      break;

    case PATTERN_OPERATORS.longer:
      pattern.cond.type = PATTERN_CONDITIONS.greater;
      pattern.cond.value = rule.duration;
      break;

    case PATTERN_OPERATORS.shorter:
      pattern.cond.type = PATTERN_CONDITIONS.less;
      pattern.cond.value = rule.duration;
      break;

    case PATTERN_OPERATORS.ticketAssociated:
    case PATTERN_OPERATORS.canceled:
    case PATTERN_OPERATORS.snoozed:
    case PATTERN_OPERATORS.acked:
      pattern.cond.type = PATTERN_CONDITIONS.equal;
      pattern.cond.value = true;
      break;
    case PATTERN_OPERATORS.ticketNotAssociated:
    case PATTERN_OPERATORS.notCanceled:
    case PATTERN_OPERATORS.notSnoozed:
    case PATTERN_OPERATORS.notAcked:
      pattern.cond.type = PATTERN_CONDITIONS.equal;
      pattern.cond.value = false;
      break;
  }

  return pattern;
};

/**
 * Convert form group to pattern rules
 *
 * @param {PatternGroupForm} group
 * @return {PatternRules}
 */
export const formGroupToPatternRules = group => group.rules.map(formRuleToPatternRule);

/**
 * Convert form groups to pattern rules
 *
 * @param {PatternGroupsForm} groups
 * @return {PatternGroups}
 */
export const formGroupsToPatternRules = groups => groups.map(formGroupToPatternRules);

/**
 * Convert pattern form to pattern
 *
 * @param {PatternForm} form
 * @return {Pattern}
 */
export const formToPattern = (form) => {
  const { groups, id, ...pattern } = form;

  switch (form.type) {
    case PATTERN_TYPES.alarm:
      pattern.alarm_pattern = formGroupsToPatternRules(groups);
      break;
    case PATTERN_TYPES.entity:
      pattern.entity_pattern = formGroupsToPatternRules(groups);
      break;
    case PATTERN_TYPES.pbehavior:
      pattern.pbehavior_pattern = formGroupsToPatternRules(groups);
      break;
  }

  return pattern;
};
