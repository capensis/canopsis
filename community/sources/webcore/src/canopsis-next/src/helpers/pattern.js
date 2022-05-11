import { isArray, isBoolean, isEmpty, isNan, isNull, isNumber, isUndefined } from 'lodash';

import {
  PATTERN_FIELD_TYPES,
  PATTERN_ARRAY_OPERATORS,
  PATTERN_BOOLEAN_OPERATORS,
  PATTERN_DURATION_OPERATORS,
  PATTERN_INFOS_NAME_OPERATORS,
  PATTERN_NULL_OPERATORS,
  PATTERN_NUMBER_OPERATORS,
  PATTERN_OPERATORS_WITHOUT_VALUE,
  PATTERN_RULE_INFOS_FIELDS,
  PATTERN_RULE_TYPES,
  PATTERN_STRING_OPERATORS,
  PATTERN_CONDITIONS,
} from '@/constants';

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
 * Check is operator for string
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
 * Check field type is string array
 *
 * @param {PatternFieldType} type
 * @return {boolean}
 */
export const isStringArrayFieldType = type => type === PATTERN_FIELD_TYPES.stringArray;

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
 * Convert any value to type value
 *
 * @param {} type
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
      return (isNan(preparedValue) || isNull(preparedValue))
        ? ''
        : String(preparedValue);
    case PATTERN_FIELD_TYPES.null:
      return null;
    case PATTERN_FIELD_TYPES.stringArray:
      return preparedValue ? [preparedValue] : [];
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
      return PATTERN_ARRAY_OPERATORS;
    case PATTERN_FIELD_TYPES.boolean:
      return PATTERN_BOOLEAN_OPERATORS;
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

  if (isInfosRuleType(ruleType) && rule.field === PATTERN_RULE_INFOS_FIELDS.name) {
    return PATTERN_INFOS_NAME_OPERATORS;
  }

  const fieldType = getFieldType(rule.value);

  return getOperatorsByFieldType(fieldType);
};

/**
 * Get value type by operator
 *
 * @param {string} operator
 * @return {PatternFieldType | undefined}
 */
export const getValueTypeByOperator = (operator) => {
  if (isOperatorForArray(operator)) {
    return PATTERN_FIELD_TYPES.stringArray;
  }

  if (isOperatorForString(operator)) {
    return PATTERN_FIELD_TYPES.string;
  }

  if (isOperatorForNumber(operator)) {
    return PATTERN_FIELD_TYPES.number;
  }

  if (isOperatorForBoolean(operator)) {
    return PATTERN_FIELD_TYPES.boolean;
  }

  if (isOperatorForNull(operator)) {
    return PATTERN_FIELD_TYPES.null;
  }

  return undefined;
};

/**
 * Convert value to operator type
 *
 * @param {PatternValue} value
 * @param {string} operator
 * @return {PatternValue|undefined|*}
 */
export const convertValueByOperator = (value, operator) => {
  const valueType = getFieldType(value);
  const operatorValueType = getValueTypeByOperator(operator);

  if (valueType === operatorValueType) {
    return value;
  }

  return convertValueByType(value, operatorValueType);
};

/**
 * Check condition is valid
 *
 * @param {string} condition
 * @return {boolean}
 */
export const isValidPatternCondition = condition => Object.values(PATTERN_CONDITIONS).includes(condition);
