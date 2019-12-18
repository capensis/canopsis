import { isEmpty, isObject, cloneDeep, isNull } from 'lodash';

import { FILTER_OPERATORS, FILTER_DEFAULT_VALUES, FILTER_MONGO_OPERATORS } from '@/constants';
import uid from '@/helpers/uid';

/**
 * @description Determine the operator and the input value of a rule
 * @param {Object} rule
 */
function ruleOperatorAndInput(rule) {
  const parsedRule = {
    operator: '',
    input: '',
  };

  const ruleValue = Object.values(rule)[0];

  /**
   * Switch to determine if it's a short syntax for '$eq' and '$eq:'''
   */
  if (!isObject(ruleValue)) {
    if (ruleValue === '') {
      parsedRule.operator = FILTER_OPERATORS.isEmpty;
    } else {
      const [input] = Object.values(rule);
      if (isNull(input)) {
        parsedRule.operator = FILTER_OPERATORS.isNull;
      } else {
        parsedRule.operator = FILTER_OPERATORS.equal;
        parsedRule.input = input;
      }
    }
  } else {
    const operator = Object.keys(ruleValue)[0];

    /**
     * Switch to determine the right operator, and then assign the right input value
     */
    switch (operator) {
      case FILTER_MONGO_OPERATORS.equal: {
        const [input] = Object.values(rule);
        parsedRule.input = input;
        parsedRule.operator = FILTER_OPERATORS.equal;
        break;
      }
      case FILTER_MONGO_OPERATORS.notEqual: {
        if (Object.values(ruleValue)[0] === null) {
          parsedRule.operator = FILTER_OPERATORS.isNotNull;
        } else if (Object.values(ruleValue)[0] === '') {
          parsedRule.operator = FILTER_OPERATORS.isNotEmpty;
        } else {
          const [input] = Object.values(ruleValue);
          parsedRule.input = input;
          parsedRule.operator = FILTER_OPERATORS.notEqual;
        }
        break;
      }
      case FILTER_MONGO_OPERATORS.in: {
        const [inputArray] = Object.values(ruleValue);
        parsedRule.input = inputArray;
        parsedRule.operator = FILTER_OPERATORS.in;
        break;
      }
      case FILTER_MONGO_OPERATORS.notIn: {
        const [inputArray] = Object.values(ruleValue);
        parsedRule.input = inputArray;
        parsedRule.operator = FILTER_OPERATORS.notIn;
        break;
      }
      case FILTER_MONGO_OPERATORS.regex: {
        const [input] = Object.values(ruleValue);
        parsedRule.input = input;
        parsedRule.operator = FILTER_OPERATORS.contains;
        break;
      }
      default: {
        /**
         * Throw an error if the operator was not found.
         */
        const [input] = Object.values(ruleValue);

        parsedRule.input = input;
        parsedRule.operator = operator;
      }
    }
  }

  return parsedRule;
}

function parseRuleToFilter(rule) {
  const parsedRule = cloneDeep(FILTER_DEFAULT_VALUES.rule);

  if (isEmpty(rule)) {
    return parsedRule;
  }

  const [field] = Object.keys(rule);
  parsedRule.field = field;
  const { operator, input } = ruleOperatorAndInput(rule);

  parsedRule.operator = operator;
  parsedRule.input = input;

  return parsedRule;
}


export default function parseGroupToFilter(group) {
  const parsedGroup = cloneDeep(FILTER_DEFAULT_VALUES.group);
  let groupContent = Object.values(group)[0];

  if (!isObject(groupContent)) {
    groupContent = [{ ...group }];
  } else {
    const [condition] = Object.keys(group);

    /**
     * It there is group without condition wrapper. Ex.: { _id: '123' } instead of { $and: [{ _id: '123' }] }
     */
    if (![FILTER_MONGO_OPERATORS.and, FILTER_MONGO_OPERATORS.or].includes(condition)) {
      parsedGroup.condition = FILTER_DEFAULT_VALUES.condition;
      groupContent = [{ ...group }];
    } else {
      parsedGroup.condition = condition;
    }
  }

  if (isEmpty(groupContent)) {
    return parsedGroup;
  }

  /**
  * Map over the items of a group.
  * If the item is an array -> It's a group.
  * Else -> It's a rule.
  */
  groupContent.forEach((item) => {
    if (Array.isArray(Object.values(item)[0])) {
      parsedGroup.groups[uid('group')] = parseGroupToFilter(item);
    } else {
      Object.entries(item).forEach(([key, value]) => {
        parsedGroup.rules[uid('rules')] = parseRuleToFilter({ [key]: value });
      }, {});
    }
  });

  return parsedGroup;
}
