import { isNil, isArray, isBoolean, isEmpty, isNan, isNull, isNumber, isString, isUndefined } from 'lodash';

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
  ALARM_PATTERN_FIELDS,
  ENTITY_PATTERN_FIELDS,
  EVENT_FILTER_PATTERN_FIELDS,
  PATTERN_OPERATORS,
  SERVICE_WEATHER_PATTERN_FIELDS,
  OLD_PATTERNS_FIELDS,
  OLD_PATTERN_FIELDS_TO_NEW_FIELDS,
} from '@/constants';

import { isValidDateInterval } from '@/helpers/date/date';
import { isValidDuration } from '@/helpers/date/duration';

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
 * Check rule is object
 *
 * @param {string} type
 * @return {boolean}
 */
export const isObjectRuleType = type => type === PATTERN_RULE_TYPES.object;

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
  ALARM_PATTERN_FIELDS.resolvedAt,
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
  ALARM_PATTERN_FIELDS.component,
  ALARM_PATTERN_FIELDS.connector,
  ALARM_PATTERN_FIELDS.connectorName,
  ALARM_PATTERN_FIELDS.resource,
  ENTITY_PATTERN_FIELDS.id,
  EVENT_FILTER_PATTERN_FIELDS.component,
  EVENT_FILTER_PATTERN_FIELDS.connector,
  EVENT_FILTER_PATTERN_FIELDS.connectorName,
  EVENT_FILTER_PATTERN_FIELDS.resource,
].includes(value);

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
      return isArray(cond.value);
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

  if (isStringArrayFieldType(fieldType)) {
    if (isArrayCondition(cond.type)) {
      return (isArray(cond.value) && cond.value.every(isString))
        || (isBoolean(cond.value) && cond.type === PATTERN_CONDITIONS.isEmpty);
    }

    return false;
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

  if (
    (isInfosRuleType(ruleType) || isExtraInfosRuleType(ruleType))
    && rule.field === PATTERN_RULE_INFOS_FIELDS.name
  ) {
    return PATTERN_INFOS_NAME_OPERATORS;
  }

  const fieldType = getFieldType(rule.value);

  return getOperatorsByFieldType(fieldType);
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
  const valueType = getFieldType(value);
  const operatorsValueType = getValueTypesByOperator(operator);

  if (operatorsValueType.includes(valueType)) {
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
 * Check if pattern for source is old (was not migrated)
 *
 * @param {Object} source
 * @param {string[]} [oldFields = OLD_PATTERNS_FIELDS.mongoQuery]
 * @returns {boolean}
 */
export const isOldPattern = (source, oldFields = [OLD_PATTERNS_FIELDS.mongoQuery]) => {
  const notEmptyOldFields = oldFields.filter(field => source[field]);

  return !!notEmptyOldFields.length
    && notEmptyOldFields.every(field => (
      OLD_PATTERN_FIELDS_TO_NEW_FIELDS[field].every(newField => !source[newField]?.length)
    ));
};
