import {
  isNil,
  isArray,
  isBoolean,
  isEmpty,
  isNan,
  isNull,
  isNumber,
  isString,
  isUndefined,
  omit,
} from 'lodash';

import {
  PATTERN_FIELD_TYPES,
  PATTERN_ARRAY_OPERATORS,
  PATTERN_BOOLEAN_OPERATORS,
  PATTERN_DURATION_OPERATORS,
  PATTERN_EXISTS_OPERATORS,
  PATTERN_NULL_OPERATORS,
  PATTERN_NUMBER_OPERATORS,
  PATTERN_OPERATORS_WITHOUT_VALUE,
  PATTERN_RULE_INFOS_FIELDS,
  PATTERN_RULE_TYPES,
  PATTERN_STRING_OPERATORS,
  PATTERN_CONDITIONS,
  ALARM_PATTERN_FIELDS,
  ENTITY_PATTERN_FIELDS,
  EVENT_FILTER_PATTERN_FIELDS,
  PATTERN_OPERATORS,
  SERVICE_WEATHER_PATTERN_FIELDS,
  QUICK_RANGES,
  PATTERN_QUICK_RANGES,
  PATTERNS_FIELDS,
  PATTERN_CUSTOM_ITEM_VALUE,
  PATTERN_TYPES,
  TIME_UNITS,
} from '@/constants';

import { convertDateToTimestamp, isValidDateInterval } from '@/helpers/date/date';
import { durationToForm, isValidDuration, toSeconds } from '@/helpers/date/duration';
import { uid } from '@/helpers/uid';
import {
  getDiffBetweenStartAndStopQuickInterval,
  getQuickRangeByDiffBetweenStartAndStop,
} from '@/helpers/date/date-intervals';

/**
 * @typedef { 'string' | 'number' | 'infos' | 'date' | 'duration' } PatternRuleType
 */

/**
 * @typedef { 'string' | 'int' | 'bool' | 'null' | 'string_array' } PatternFieldType
 */

/**
 * @typedef {boolean | [] | null | string | number} PatternValue
 */

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
 * @typedef {Object & ObjectKey} PatternRuleForm
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
 * @typedef {Object & ObjectKey} PatternGroupForm
 * @property {PatternRuleForm[]} rules
 */

/**
 * @typedef {PatternRuleForm[]} PatternGroupsForm
 */

/**
 * @typedef {Pattern} PatternForm
 * @property {PatternGroupsForm} groups
 */

/**
 * Check is operator has a value
 *
 * @param {string} operator
 * @return {boolean}
 */
export const isOperatorHasValue = operator => !PATTERN_OPERATORS_WITHOUT_VALUE.includes(operator);

/**
 * Check is operator for array
 *
 * @param {string} operator
 * @return {boolean}
 */
export const isOperatorForArray = operator => PATTERN_ARRAY_OPERATORS.includes(operator);

/**
 * Check is operator for string
 *
 * @param {string} operator
 * @return {boolean}
 */
export const isOperatorForString = operator => PATTERN_STRING_OPERATORS.includes(operator);

/**
 * Check is operator for number
 *
 * @param {string} operator
 * @return {boolean}
 */
export const isOperatorForNumber = operator => PATTERN_NUMBER_OPERATORS.includes(operator);

/**
 * Check is operator for boolean
 *
 * @param {string} operator
 * @return {boolean}
 */
export const isOperatorForBoolean = operator => PATTERN_BOOLEAN_OPERATORS.includes(operator);

/**
 * Check is operator for null
 *
 * @param {string} operator
 * @return {boolean}
 */
export const isOperatorForNull = operator => PATTERN_NULL_OPERATORS.includes(operator);

/**
 * Check rule is infos
 *
 * @param {string} type
 * @return {boolean}
 */
export const isInfosRuleType = type => type === PATTERN_RULE_TYPES.infos;

/**
 * Check rule is extra infos
 *
 * @param {string} type
 * @return {boolean}
 */
export const isExtraInfosRuleType = type => type === PATTERN_RULE_TYPES.extraInfos;

/**
 * Check rule is date
 *
 * @param {string} type
 * @return {boolean}
 */
export const isDateRuleType = type => type === PATTERN_RULE_TYPES.date;

/**
 * Check rule is duration
 *
 * @param {string} type
 * @return {boolean}
 */
export const isDurationRuleType = type => type === PATTERN_RULE_TYPES.duration;

/**
 * Check rule is object
 *
 * @param {string} type
 * @return {boolean}
 */
export const isObjectRuleType = type => type === PATTERN_RULE_TYPES.object;

/**
 * Check field type is string array
 *
 * @param {PatternFieldType} type
 * @return {boolean}
 */
export const isStringArrayFieldType = type => type === PATTERN_FIELD_TYPES.stringArray;

/**
 * Check field type is valid
 *
 * @param {*} value
 * @return {boolean}
 */
export const isValidRuleFieldType = value => Object.values(PATTERN_FIELD_TYPES).includes(value);

/**
 * Check condition is boolean
 *
 * @param {string} condition
 * @return {boolean}
 */
export const isBooleanCondition = condition => [
  PATTERN_CONDITIONS.equal,
  PATTERN_CONDITIONS.notEqual,
  PATTERN_CONDITIONS.exist,
  PATTERN_CONDITIONS.isEmpty,
].includes(condition);

/**
 * Check condition is boolean
 *
 * @param {string} condition
 * @return {boolean}
 */
export const isArrayCondition = condition => [
  PATTERN_CONDITIONS.isEmpty,
  PATTERN_CONDITIONS.hasNot,
  PATTERN_CONDITIONS.hasOneOf,
  PATTERN_CONDITIONS.isOneOf,
  PATTERN_CONDITIONS.isNotOneOf,
  PATTERN_CONDITIONS.hasEvery,
].includes(condition);

/**
 * Check condition is valid
 *
 * @param {string} condition
 * @return {boolean}
 */
export const isValidPatternCondition = condition => Object.values(PATTERN_CONDITIONS).includes(condition);

/**
 * Check pattern field is date
 *
 * @param {string} value
 * @return {boolean}
 */
export const isDatePatternRuleField = value => [
  ALARM_PATTERN_FIELDS.creationDate,
  ALARM_PATTERN_FIELDS.lastEventDate,
  ALARM_PATTERN_FIELDS.lastUpdateDate,
  ALARM_PATTERN_FIELDS.ackAt,
  ALARM_PATTERN_FIELDS.resolved,
  ALARM_PATTERN_FIELDS.activationDate,
  ENTITY_PATTERN_FIELDS.lastEventDate,
].includes(value);

/**
 * Check pattern field is number
 *
 * @param {string} value
 * @return {boolean}
 */
export const isNumberPatternRuleField = value => [
  ALARM_PATTERN_FIELDS.state,
  ALARM_PATTERN_FIELDS.status,
  EVENT_FILTER_PATTERN_FIELDS.state,
  ENTITY_PATTERN_FIELDS.impactLevel,
  SERVICE_WEATHER_PATTERN_FIELDS.state,
].includes(value);

/**
 * Check pattern field is array
 *
 * @param {string} value
 * @return {boolean}
 */
export const isArrayPatternRuleField = value => [
  ALARM_PATTERN_FIELDS.displayName,
  ALARM_PATTERN_FIELDS.output,
  ALARM_PATTERN_FIELDS.longOutput,
  ALARM_PATTERN_FIELDS.initialOutput,
  ALARM_PATTERN_FIELDS.initialLongOutput,
  ALARM_PATTERN_FIELDS.component,
  ALARM_PATTERN_FIELDS.connector,
  ALARM_PATTERN_FIELDS.connectorName,
  ALARM_PATTERN_FIELDS.resource,
  ALARM_PATTERN_FIELDS.tags,
  ALARM_PATTERN_FIELDS.lastComment,
  ALARM_PATTERN_FIELDS.lastCommentInitiator,
  ALARM_PATTERN_FIELDS.ticketMessage,
  ALARM_PATTERN_FIELDS.ticketValue,
  ALARM_PATTERN_FIELDS.ticketInitiator,
  ALARM_PATTERN_FIELDS.ticketData,
  ALARM_PATTERN_FIELDS.ackBy,
  ALARM_PATTERN_FIELDS.ackMessage,
  ALARM_PATTERN_FIELDS.ackInitiator,
  ALARM_PATTERN_FIELDS.canceledInitiator,
  ENTITY_PATTERN_FIELDS.id,
  ENTITY_PATTERN_FIELDS.name,
  ENTITY_PATTERN_FIELDS.category,
  ENTITY_PATTERN_FIELDS.type,
  ENTITY_PATTERN_FIELDS.connector,
  ENTITY_PATTERN_FIELDS.component,
  EVENT_FILTER_PATTERN_FIELDS.component,
  EVENT_FILTER_PATTERN_FIELDS.connector,
  EVENT_FILTER_PATTERN_FIELDS.connectorName,
  EVENT_FILTER_PATTERN_FIELDS.resource,
  EVENT_FILTER_PATTERN_FIELDS.output,
  EVENT_FILTER_PATTERN_FIELDS.longOutput,
  EVENT_FILTER_PATTERN_FIELDS.eventType,
  EVENT_FILTER_PATTERN_FIELDS.sourceType,
  EVENT_FILTER_PATTERN_FIELDS.initiator,
  EVENT_FILTER_PATTERN_FIELDS.author,
].some((field) => {
  /**
   * @TODO: update babel-eslint for resolving problem with templates inside optional chaiging function call
   */
  const start = `${field}.`;

  return value === field || value?.startsWith(start);
});

/**
 * Check pattern field is infos
 *
 * @param {string} value
 * @return {boolean}
 */
export const isInfosPatternRuleField = value => [
  ALARM_PATTERN_FIELDS.infos,
  ENTITY_PATTERN_FIELDS.componentInfos,
  ENTITY_PATTERN_FIELDS.infos,
].some((field) => {
  /**
   * @TODO: update babel-eslint for resolving problem with templates inside optional chaiging function call
   */
  const start = `${field}.`;

  return value === field || value?.startsWith(start);
});

/**
 * Check pattern field is duration
 *
 * @param {string} value
 * @return {boolean}
 */
export const isDurationPatternRuleField = value => value === ALARM_PATTERN_FIELDS.duration;

/**
 * Check pattern is extra infos
 *
 * @param {string} value
 * @return {boolean}
 */
export const isExtraInfosPatternRuleField = (value) => {
  /**
   * @TODO: update babel-eslint for resolving problem with templates inside optional chaiging function call
   */
  const start = `${EVENT_FILTER_PATTERN_FIELDS.extraInfos}.`;

  return value === EVENT_FILTER_PATTERN_FIELDS.extraInfos
    || value?.startsWith(start);
};

/**
 * Get object pattern field
 *
 * @param {string} value
 * @return {string}
 */
export const getObjectPatternRuleField = value => [ALARM_PATTERN_FIELDS.ticketData]
  .find((field) => {
    /**
     * @TODO: update babel-eslint for resolving problem with templates inside optional chaiging function call
     */
    const start = `${field}.`;

    return value === field || value?.startsWith(start);
  });

/**
 * Check pattern field is object
 *
 * @param {string} value
 * @return {boolean}
 */
export const isObjectPatternRuleField = value => !!getObjectPatternRuleField(value);

/**
 * Check rule value is valid without field type
 *
 * @param {PatternRule | *} rule
 * @return {boolean}
 */
export const isValidRuleValueWithoutFieldType = (rule) => {
  const { field, cond } = rule;

  if (isDatePatternRuleField(field)) {
    if (cond.type === PATTERN_CONDITIONS.absoluteTime) {
      return isValidDateInterval(cond?.value);
    }

    if (cond.type === PATTERN_CONDITIONS.relativeTime) {
      return isValidDuration(cond.value);
    }
  }

  if (isDurationPatternRuleField(field)) {
    return isValidDuration(cond?.value);
  }

  if (isNumberPatternRuleField(field)) {
    return isNumber(cond.value);
  }

  if (isArrayPatternRuleField(field)) {
    if (isArrayCondition(cond.type)) {
      return isArray(cond.value) || (isBoolean(cond.value) && cond.type === PATTERN_CONDITIONS.isEmpty);
    }
  }

  if (isBoolean(cond.value)) {
    return isBooleanCondition(cond.type);
  }

  return isString(cond.value);
};

/**
 * Return field type by value
 *
 * @param {PatternValue} value
 * @return {string}
 */
export const getFieldType = (value) => {
  if (isBoolean(value)) {
    return PATTERN_FIELD_TYPES.boolean;
  }

  if (isNumber(value)) {
    return PATTERN_FIELD_TYPES.number;
  }

  if (isNull(value)) {
    return PATTERN_FIELD_TYPES.null;
  }

  if (isArray(value)) {
    return PATTERN_FIELD_TYPES.stringArray;
  }

  return PATTERN_FIELD_TYPES.string;
};

/**
 * Check rule value is valid with field type
 *
 * @param {PatternRule | *} rule
 * @return {boolean}
 */
export const isValidRuleValueWithFieldType = (rule) => {
  const { field, cond, field_type: fieldType } = rule;

  if ([PATTERN_FIELD_TYPES.stringArray, PATTERN_FIELD_TYPES.string].includes(fieldType)) {
    if (isArrayCondition(cond.type)) {
      return (isArray(cond.value) && cond.value.every(isString))
        || (isBoolean(cond.value) && cond.type === PATTERN_CONDITIONS.isEmpty);
    }

    if (fieldType === PATTERN_FIELD_TYPES.stringArray) {
      return false;
    }
  }

  const isInfos = isInfosPatternRuleField(field) || isExtraInfosPatternRuleField(field);

  return isInfos && getFieldType(cond.value) === fieldType;
};

/**
 * Check rule value is valid
 *
 * @param {PatternRule | *} rule
 * @return {boolean}
 */
export const isValidRuleValue = rule => (
  rule.field_type
    ? isValidRuleValueWithFieldType(rule)
    : isValidRuleValueWithoutFieldType(rule)
);

/**
 * Check pattern rule is valid
 *
 * @param {PatternRule | *} rule
 * @return {boolean}
 */
export const isValidPatternRule = rule => !!rule?.field
  && !isNil(rule.cond?.value)
  && !isNil(rule.cond?.type)
  && (!rule.field_type || isValidRuleFieldType(rule.field_type))
  && isValidPatternCondition(rule.cond.type)
  && isValidRuleValue(rule);

/**
 * Convert any value to type value
 *
 * @param {PatternFieldType} type
 * @param {PatternValue} [value]
 * @param [defaultValue]
 * @return {PatternValue | undefined}
 */
export const convertValueByType = (value, type, defaultValue) => {
  if (isEmpty(value) && !isUndefined(defaultValue)) {
    return defaultValue;
  }

  const preparedValue = isArray(value) ? value[0] : value;

  switch (type) {
    case PATTERN_FIELD_TYPES.number:
      return Number(preparedValue) || 0;
    case PATTERN_FIELD_TYPES.boolean:
      return Boolean(preparedValue);
    case PATTERN_FIELD_TYPES.string:
      return (isNan(preparedValue) || isNull(preparedValue) || isUndefined(preparedValue))
        ? ''
        : String(preparedValue);
    case PATTERN_FIELD_TYPES.null:
      return null;
    case PATTERN_FIELD_TYPES.stringArray:
      return preparedValue ? [String(preparedValue)] : [];
    default:
      return undefined;
  }
};

/**
 * Get operators by type of value
 *
 * @param {PatternFieldType} fieldType
 * @return {string[]}
 */
export const getOperatorsByFieldType = (fieldType) => {
  switch (fieldType) {
    case PATTERN_FIELD_TYPES.number:
      return PATTERN_NUMBER_OPERATORS;
    case PATTERN_FIELD_TYPES.stringArray:
      return [
        PATTERN_OPERATORS.hasEvery,
        PATTERN_OPERATORS.hasOneOf,
        PATTERN_OPERATORS.hasNot,
        PATTERN_OPERATORS.isEmpty,
        PATTERN_OPERATORS.isNotEmpty,
      ];
    case PATTERN_FIELD_TYPES.boolean:
      return [PATTERN_OPERATORS.equal];
    default:
      return PATTERN_STRING_OPERATORS;
  }
};

/**
 * Get operators by rule type and rule values
 *
 * @param {PatternRuleForm} rule
 * @param {PatternRuleType} ruleType
 * @return {string[]}
 */
export const getOperatorsByRule = (rule, ruleType) => {
  if (isDurationRuleType(ruleType)) {
    return PATTERN_DURATION_OPERATORS;
  }

  const isAnyInfosType = isInfosRuleType(ruleType) || isExtraInfosRuleType(ruleType);

  if (isAnyInfosType && rule.field === PATTERN_RULE_INFOS_FIELDS.name) {
    return PATTERN_EXISTS_OPERATORS;
  }

  let operators = getOperatorsByFieldType(rule.fieldType);

  if (rule.fieldType === PATTERN_FIELD_TYPES.string || isObjectRuleType(ruleType)) {
    operators = [
      ...operators,
      PATTERN_OPERATORS.isOneOf,
      PATTERN_OPERATORS.isNotOneOf,
    ];
  }

  return operators;
};

/**
 * Get value type by operator
 *
 * @param {string} operator
 * @return {PatternFieldType[]}
 */
export const getValueTypesByOperator = (operator) => {
  const operators = [];

  if (isOperatorForArray(operator)) {
    operators.push(PATTERN_FIELD_TYPES.stringArray);
  }

  if (isOperatorForString(operator)) {
    operators.push(PATTERN_FIELD_TYPES.string);
  }

  if (isOperatorForNumber(operator)) {
    operators.push(PATTERN_FIELD_TYPES.number);
  }

  if (isOperatorForBoolean(operator)) {
    operators.push(PATTERN_FIELD_TYPES.boolean);
  }

  if (isOperatorForNull(operator)) {
    operators.push(PATTERN_FIELD_TYPES.null);
  }

  return operators;
};

/**
 * Convert value to operator type
 *
 * @param {PatternValue} value
 * @param {string} operator
 * @return {PatternValue|undefined|*}
 */
export const convertValueByOperator = (value, operator) => {
  const fieldType = getFieldType(value);

  const operatorsValueType = getValueTypesByOperator(operator);

  if (operatorsValueType.includes(fieldType)) {
    return value;
  }

  return convertValueByType(value, operatorsValueType[0]);
};

/**
 * Create id pattern rule for entity patterns
 *
 * @param {string | string[]} value
 * @returns {PatternGroups}
 */
export const createEntityIdPatternByValue = value => [[{
  field: ENTITY_PATTERN_FIELDS.id,
  cond: {
    type: isArray(value) ? PATTERN_CONDITIONS.isOneOf : PATTERN_CONDITIONS.equal,
    value,
  },
}]];

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

    form.field = PATTERN_EXISTS_OPERATORS.includes(rule.cond.type)
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
 * Convert form groups to pattern rules query
 *
 * @param {PatternGroupsForm} [groups = []]
 * @return {string}
 */
export const formGroupsToPatternRulesQuery = (groups = []) => JSON.stringify(formGroupsToPatternRules(groups));

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
