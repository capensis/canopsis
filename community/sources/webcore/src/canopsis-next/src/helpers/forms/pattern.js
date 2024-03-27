import { isBoolean, omit } from 'lodash';

import {
  ALARM_PATTERN_FIELDS,
  ENTITY_PATTERN_FIELDS,
  EVENT_FILTER_PATTERN_FIELDS,
  PATTERN_CONDITIONS,
  PATTERN_CUSTOM_ITEM_VALUE,
  PATTERN_INFOS_NAME_OPERATORS,
  PATTERN_OPERATORS,
  PATTERN_QUICK_RANGES,
  PATTERN_RULE_INFOS_FIELDS,
  PATTERN_TYPES,
  QUICK_RANGES,
  PATTERNS_FIELDS,
  SERVICE_WEATHER_PATTERN_FIELDS,
  TIME_UNITS, PATTERN_FIELD_TYPES,
} from '@/constants';

import uid from '@/helpers/uid';
import {
  getObjectPatternRuleField,
  isDatePatternRuleField,
  isInfosPatternRuleField,
  isExtraInfosPatternRuleField,
  isObjectPatternRuleField,
} from '@/helpers/pattern';
import { convertDateToTimestamp } from '@/helpers/date/date';
import {
  getDiffBetweenStartAndStopQuickInterval,
  getQuickRangeByDiffBetweenStartAndStop,
} from '@/helpers/date/date-intervals';
import { durationToForm, toSeconds } from '@/helpers/date/duration';

/**
 * @typedef { 'alarm' | 'entity' | 'pbehavior' } PatternTypes
 */

/**
 * @typedef { 'alarm_pattern' | 'entity_pattern' | 'pbehavior_pattern' | 'total_entity_pattern' } PatternsField
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
 * @property {PatternRuleRangeCondition | Duration | string | number | boolean | []} value
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
 * @property {string} fieldType
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
 * @property {Object} old_mongo_query
 */

/**
 * Convert pattern rule to form
 *
 * @param {PatternRule} rule
 * @return {PatternRuleForm}
 */
export const patternRuleToForm = (rule = {}) => {
  const form = {
    key: uid(),
    attribute: rule.field ?? '',
    operator: '',
    field: '',
    fieldType: rule.field_type ?? PATTERN_FIELD_TYPES.string,
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
  const isEntityComponentInfos = rule.field?.startsWith(ENTITY_PATTERN_FIELDS.componentInfos);
  const isDuration = rule.field === ALARM_PATTERN_FIELDS.duration;
  const isInfos = isAlarmInfos || isEntityInfos || isEntityComponentInfos;
  const isExtraInfos = !isInfos && rule.field?.startsWith(EVENT_FILTER_PATTERN_FIELDS.extraInfos);
  const patternObjectField = getObjectPatternRuleField(rule.field);

  switch (rule.cond.type) {
    case PATTERN_CONDITIONS.equal: {
      if (isBoolean(rule.cond.value)) {
        if (rule.field === SERVICE_WEATHER_PATTERN_FIELDS.grey) {
          form.operator = rule.cond.value ? PATTERN_OPERATORS.isGrey : PATTERN_OPERATORS.isNotGrey;
        }
      }

      if (!form.operator) {
        form.operator = PATTERN_OPERATORS.equal;
        form.value = rule.cond.value;
      }
      break;
    }
    case PATTERN_CONDITIONS.notEqual:
      form.operator = PATTERN_OPERATORS.notEqual;
      form.value = rule.cond.value;
      break;

    case PATTERN_CONDITIONS.greater:
      form.operator = isDuration
        ? PATTERN_OPERATORS.longer
        : PATTERN_OPERATORS.higher;
      form.value = isDuration ? '' : rule.cond.value;
      break;
    case PATTERN_CONDITIONS.less:
      form.operator = isDuration
        ? PATTERN_OPERATORS.shorter
        : PATTERN_OPERATORS.lower;
      form.value = isDuration ? '' : rule.cond.value;
      break;

    case PATTERN_CONDITIONS.exist:
      if (rule.cond.value) {
        form.operator = {
          [ALARM_PATTERN_FIELDS.snooze]: PATTERN_OPERATORS.snoozed,
          [ALARM_PATTERN_FIELDS.ack]: PATTERN_OPERATORS.acked,
          [ALARM_PATTERN_FIELDS.canceled]: PATTERN_OPERATORS.canceled,
          [ALARM_PATTERN_FIELDS.ticket]: PATTERN_OPERATORS.ticketAssociated,
          [ALARM_PATTERN_FIELDS.activationDate]: PATTERN_OPERATORS.activated,
        }[rule.field];
      } else {
        form.operator = {
          [ALARM_PATTERN_FIELDS.snooze]: PATTERN_OPERATORS.notSnoozed,
          [ALARM_PATTERN_FIELDS.ack]: PATTERN_OPERATORS.notAcked,
          [ALARM_PATTERN_FIELDS.canceled]: PATTERN_OPERATORS.notCanceled,
          [ALARM_PATTERN_FIELDS.ticket]: PATTERN_OPERATORS.ticketNotAssociated,
          [ALARM_PATTERN_FIELDS.activationDate]: PATTERN_OPERATORS.inactive,
        }[rule.field];
      }

      if (!form.operator) {
        form.operator = rule.cond.value === true
          ? PATTERN_OPERATORS.exist
          : PATTERN_OPERATORS.notExist;
      }

      if (rule.field === ALARM_PATTERN_FIELDS.activationDate) {
        form.attribute = ALARM_PATTERN_FIELDS.activated;
      }
      break;

    case PATTERN_CONDITIONS.hasEvery:
      form.operator = PATTERN_OPERATORS.hasEvery;
      form.value = rule.cond.value;
      break;
    case PATTERN_CONDITIONS.isOneOf:
      form.operator = PATTERN_OPERATORS.isOneOf;
      form.value = rule.cond.value;
      break;
    case PATTERN_CONDITIONS.hasOneOf:
      form.operator = rule.field === ALARM_PATTERN_FIELDS.tags
        ? PATTERN_OPERATORS.with
        : PATTERN_OPERATORS.hasOneOf;
      form.value = rule.cond.value;
      break;
    case PATTERN_CONDITIONS.hasNot:
      form.operator = rule.field === ALARM_PATTERN_FIELDS.tags
        ? PATTERN_OPERATORS.without
        : PATTERN_OPERATORS.hasNot;
      form.value = rule.cond.value;
      break;
    case PATTERN_CONDITIONS.isNotOneOf:
      form.operator = PATTERN_OPERATORS.isNotOneOf;
      form.value = rule.cond.value;
      break;
    case PATTERN_CONDITIONS.isEmpty:
      form.operator = rule.cond.value === true
        ? PATTERN_OPERATORS.isEmpty
        : PATTERN_OPERATORS.isNotEmpty;
      form.value = [];
      break;

    case PATTERN_CONDITIONS.contains:
      form.operator = PATTERN_OPERATORS.contains;
      form.value = rule.cond.value;
      break;
    case PATTERN_CONDITIONS.notContains:
      form.operator = PATTERN_OPERATORS.notContains;
      form.value = rule.cond.value;
      break;
    case PATTERN_CONDITIONS.beginsWith:
      form.operator = PATTERN_OPERATORS.beginsWith;
      form.value = rule.cond.value;
      break;
    case PATTERN_CONDITIONS.notBeginWith:
      form.operator = PATTERN_OPERATORS.notBeginWith;
      form.value = rule.cond.value;
      break;
    case PATTERN_CONDITIONS.endsWith:
      form.operator = PATTERN_OPERATORS.endsWith;
      form.value = rule.cond.value;
      break;
    case PATTERN_CONDITIONS.notEndWith:
      form.operator = PATTERN_OPERATORS.notEndWith;
      form.value = rule.cond.value;
      break;
    case PATTERN_CONDITIONS.regexp:
      form.operator = PATTERN_OPERATORS.regexp;
      form.value = rule.cond.value;
      break;

    case PATTERN_CONDITIONS.relativeTime: {
      const { value, unit } = rule.cond.value;

      const seconds = toSeconds(value, unit);

      form.range.type = getQuickRangeByDiffBetweenStartAndStop(seconds, PATTERN_QUICK_RANGES).value;
      break;
    }
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

  if (isExtraInfos || isInfos) {
    if (isExtraInfos) {
      form.attribute = EVENT_FILTER_PATTERN_FIELDS.extraInfos;
      form.dictionary = rule.field.slice(EVENT_FILTER_PATTERN_FIELDS.extraInfos.length + 1);
    } else if (isInfos) {
      if (isAlarmInfos) {
        form.attribute = ALARM_PATTERN_FIELDS.infos;
      } else if (isEntityInfos) {
        form.attribute = ENTITY_PATTERN_FIELDS.infos;
      } else {
        form.attribute = ENTITY_PATTERN_FIELDS.componentInfos;
      }

      form.dictionary = rule.field.slice(form.attribute.length + 1);
    }

    form.field = PATTERN_INFOS_NAME_OPERATORS.includes(rule.cond.type)
      ? PATTERN_RULE_INFOS_FIELDS.name
      : PATTERN_RULE_INFOS_FIELDS.value;
  }

  if (patternObjectField) {
    form.attribute = patternObjectField;
    form.dictionary = rule.field.replace(`${patternObjectField}.`, '');
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

export const patternsToGroups = (patterns = []) => patterns.map(patternRulesToGroup);

/**
 * Convert pattern to pattern form
 *
 * @param {Pattern} pattern
 * @return {PatternForm}
 */
export const patternToForm = (pattern = {}) => ({
  ...omit(pattern, [PATTERNS_FIELDS.alarm, PATTERNS_FIELDS.entity, PATTERNS_FIELDS.pbehavior, PATTERNS_FIELDS.event]),

  title: pattern.title ?? '',
  id: pattern.id ?? PATTERN_CUSTOM_ITEM_VALUE,
  type: pattern.type ?? PATTERN_TYPES.alarm,
  is_corporate: pattern.is_corporate ?? false,
  old_mongo_query: pattern.old_mongo_query,
  groups: patternsToGroups(
    pattern.alarm_pattern
    || pattern.entity_pattern
    || pattern.pbehavior_pattern
    || pattern.event_pattern
    || pattern.total_entity_pattern
    || pattern.weather_service_pattern,
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
    value: {
      value: getDiffBetweenStartAndStopQuickInterval(range.type),
      unit: TIME_UNITS.second,
    },
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
    field: rule.attribute === ALARM_PATTERN_FIELDS.activated
      ? ALARM_PATTERN_FIELDS.activationDate
      : rule.attribute,
    cond: {
      value: rule.value,
      type: PATTERN_CONDITIONS.equal,
    },
  };

  const isInfos = isInfosPatternRuleField(rule.attribute);
  const isExtraInfos = isExtraInfosPatternRuleField(rule.attribute);
  const isDate = isDatePatternRuleField(rule.attribute);
  const isObject = isObjectPatternRuleField(rule.attribute);

  if (isInfos || isExtraInfos || isObject) {
    pattern.field = [rule.attribute, rule.dictionary].join('.');
  }

  if (isDate) {
    pattern.cond = formDateIntervalConditionToPatternRuleCondition(rule.range);

    return pattern;
  }

  if ((isExtraInfos || isInfos) && rule.field !== PATTERN_RULE_INFOS_FIELDS.name) {
    pattern.field_type = rule.fieldType;
  }

  switch (rule.operator) {
    case PATTERN_OPERATORS.equal:
      pattern.cond.type = PATTERN_CONDITIONS.equal;
      break;
    case PATTERN_OPERATORS.notEqual:
      pattern.cond.type = PATTERN_CONDITIONS.notEqual;
      break;
    case PATTERN_OPERATORS.contains:
      pattern.cond.type = PATTERN_CONDITIONS.contains;
      break;
    case PATTERN_OPERATORS.notContains:
      pattern.cond.type = PATTERN_CONDITIONS.notContains;
      break;

    case PATTERN_OPERATORS.beginsWith:
      pattern.cond.type = PATTERN_CONDITIONS.beginsWith;
      break;
    case PATTERN_OPERATORS.notBeginWith:
      pattern.cond.type = PATTERN_CONDITIONS.notBeginWith;
      break;
    case PATTERN_OPERATORS.endsWith:
      pattern.cond.type = PATTERN_CONDITIONS.endsWith;
      break;
    case PATTERN_OPERATORS.notEndWith:
      pattern.cond.type = PATTERN_CONDITIONS.notEndWith;
      break;
    case PATTERN_OPERATORS.regexp:
      pattern.cond.type = PATTERN_CONDITIONS.regexp;
      break;

    case PATTERN_OPERATORS.hasEvery:
      pattern.cond.type = PATTERN_CONDITIONS.hasEvery;
      break;
    case PATTERN_OPERATORS.isOneOf:
      pattern.cond.type = PATTERN_CONDITIONS.isOneOf;
      break;
    case PATTERN_OPERATORS.hasOneOf:
      pattern.cond.type = PATTERN_CONDITIONS.hasOneOf;
      break;
    case PATTERN_OPERATORS.hasNot:
      pattern.cond.type = PATTERN_CONDITIONS.hasNot;
      break;
    case PATTERN_OPERATORS.isNotOneOf:
      pattern.cond.type = PATTERN_CONDITIONS.isNotOneOf;
      break;
    case PATTERN_OPERATORS.isEmpty:
      pattern.cond.type = PATTERN_CONDITIONS.isEmpty;
      pattern.cond.value = true;
      break;
    case PATTERN_OPERATORS.isNotEmpty:
      pattern.cond.type = PATTERN_CONDITIONS.isEmpty;
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

    case PATTERN_OPERATORS.isGrey:
      pattern.cond.type = PATTERN_CONDITIONS.equal;
      pattern.cond.value = true;
      break;
    case PATTERN_OPERATORS.isNotGrey:
      pattern.cond.type = PATTERN_CONDITIONS.equal;
      pattern.cond.value = false;
      break;

    case PATTERN_OPERATORS.exist:
    case PATTERN_OPERATORS.ticketAssociated:
    case PATTERN_OPERATORS.canceled:
    case PATTERN_OPERATORS.snoozed:
    case PATTERN_OPERATORS.acked:
    case PATTERN_OPERATORS.activated:
      pattern.cond.type = PATTERN_CONDITIONS.exist;
      pattern.cond.value = true;
      break;
    case PATTERN_OPERATORS.notExist:
    case PATTERN_OPERATORS.ticketNotAssociated:
    case PATTERN_OPERATORS.notCanceled:
    case PATTERN_OPERATORS.notSnoozed:
    case PATTERN_OPERATORS.notAcked:
    case PATTERN_OPERATORS.inactive:
      pattern.cond.type = PATTERN_CONDITIONS.exist;
      pattern.cond.value = false;
      break;
    case PATTERN_OPERATORS.with:
      pattern.cond.type = PATTERN_CONDITIONS.hasOneOf;
      break;
    case PATTERN_OPERATORS.without:
      pattern.cond.type = PATTERN_CONDITIONS.hasNot;
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
