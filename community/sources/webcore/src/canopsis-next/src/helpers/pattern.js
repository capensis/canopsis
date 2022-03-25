import { isArray, isBoolean, isEmpty, isNan, isNull, isNumber, isUndefined } from 'lodash';

import {
  PATTERN_INPUT_TYPES,
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
} from '@/constants';

/**
 * @typedef { 'string' | 'number' | 'infos' | 'date' | 'duration' } PatternRuleType
 */

/**
 * @typedef { 'string' | 'number' | 'boolean' | 'null' | 'array' } PatternValueType
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
 * Return value type
 *
 * @param {PatternValue} value
 * @return {string}
 */
export const getValueType = (value) => {
  if (isBoolean(value)) {
    return PATTERN_INPUT_TYPES.boolean;
  }

  if (isNumber(value)) {
    return PATTERN_INPUT_TYPES.number;
  }

  if (isNull(value)) {
    return PATTERN_INPUT_TYPES.null;
  }

  if (isArray(value)) {
    return PATTERN_INPUT_TYPES.array;
  }

  return PATTERN_INPUT_TYPES.string;
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
    case PATTERN_INPUT_TYPES.number:
      return Number(preparedValue) || 0;
    case PATTERN_INPUT_TYPES.boolean:
      return Boolean(preparedValue);
    case PATTERN_INPUT_TYPES.string:
      return (isNan(preparedValue) || isNull(preparedValue))
        ? ''
        : String(preparedValue);
    case PATTERN_INPUT_TYPES.null:
      return null;
    case PATTERN_INPUT_TYPES.array:
      return preparedValue ? [preparedValue] : [];
    default:
      return undefined;
  }
};

/**
 * Get operators by type of value
 *
 * @param {PatternValueType} valueType
 * @return {string[]}
 */
export const getOperatorsByValueType = (valueType) => {
  switch (valueType) {
    case PATTERN_INPUT_TYPES.number:
      return PATTERN_NUMBER_OPERATORS;
    case PATTERN_INPUT_TYPES.array:
      return PATTERN_ARRAY_OPERATORS;
    case PATTERN_INPUT_TYPES.boolean:
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
  if (ruleType === PATTERN_RULE_TYPES.duration) {
    return PATTERN_DURATION_OPERATORS;
  }

  if (ruleType === PATTERN_RULE_TYPES.infos && rule.field === PATTERN_RULE_INFOS_FIELDS.name) {
    return PATTERN_INFOS_NAME_OPERATORS;
  }

  const valueType = getValueType(rule.value);

  return getOperatorsByValueType(valueType);
};

/**
 * Get value type by operator
 *
 * @param {string} operator
 * @return {PatternValueType | undefined}
 */
export const getValueTypeByOperator = (operator) => {
  if (isOperatorForArray(operator)) {
    return PATTERN_INPUT_TYPES.array;
  }

  if (isOperatorForString(operator)) {
    return PATTERN_INPUT_TYPES.string;
  }

  if (isOperatorForNumber(operator)) {
    return PATTERN_INPUT_TYPES.number;
  }

  if (isOperatorForBoolean(operator)) {
    return PATTERN_INPUT_TYPES.boolean;
  }

  if (isOperatorForNull(operator)) {
    return PATTERN_INPUT_TYPES.null;
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
  const valueType = getValueType(value);
  const operatorValueType = getValueTypeByOperator(operator);

  if (valueType === operatorValueType) {
    return value;
  }

  return convertValueByType(value, operatorValueType);
};
