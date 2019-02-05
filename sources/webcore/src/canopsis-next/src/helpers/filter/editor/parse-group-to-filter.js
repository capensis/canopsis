import { isEmpty, isObject, cloneDeep } from 'lodash';

import { FILTER_OPERATORS, FILTER_DEFAULT_VALUES } from '@/constants';
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
  const operator = Object.keys(ruleValue)[0];

  /**
   * Switch to determine if it's a short syntax for '$eq' and '$eq:'''
   */
  if (typeof ruleValue === 'string') {
    if (Object.values(rule)[0] === '') {
      parsedRule.operator = FILTER_OPERATORS.isEmpty;
    } else {
      const [input] = Object.values(rule);
      parsedRule.input = input;
      parsedRule.operator = FILTER_OPERATORS.equal;
    }
  } else if (typeof ruleValue === 'object') {
    /**
     * Switch to determine the right operator, and then assign the right input value
     */
    switch (operator) {
      case ('$eq'): {
        const [input] = Object.values(rule);
        parsedRule.input = input;
        parsedRule.operator = FILTER_OPERATORS.equal;
        break;
      }
      case ('$ne'): {
        if (Object.values(ruleValue)[0] === null) {
          parsedRule.operator = FILTER_OPERATORS.isNotNull;
        } else if (Object.values(ruleValue)[0] === '') {
          parsedRule.operator = FILTER_OPERATORS.isNotEmpty;
        } else {
          const [inputObject] = Object.values(rule);
          const [input] = Object.values(inputObject);
          parsedRule.input = input;
          parsedRule.operator = FILTER_OPERATORS.notEqual;
        }
        break;
      }
      case ('$in'): {
        const [inputArray] = Object.values(ruleValue);
        const [input] = inputArray;
        parsedRule.input = input;
        parsedRule.operator = FILTER_OPERATORS.in;
        break;
      }
      case ('$nin'): {
        const [inputArray] = Object.values(ruleValue);
        const [input] = inputArray;
        parsedRule.input = input;
        parsedRule.operator = FILTER_OPERATORS.notIn;
        break;
      }
      default: {
        /**
         * Throw an error if the operator was not found.
         */
        const [inputObject] = Object.values(rule);
        const [input] = Object.values(inputObject);

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
    parsedGroup.condition = condition;
  }

  if (isEmpty(groupContent)) {
    return parsedGroup;
  }

  /**
  * Map over the items of a group.
  * If the item is an array -> It's a group.
  * Else -> It's a rule.
  */
  groupContent.map((item) => {
    if (Array.isArray(Object.values(item)[0])) {
      parsedGroup.groups[uid('group')] = parseGroupToFilter(item);
    } else {
      parsedGroup.rules[uid('rules')] = parseRuleToFilter(item);
    }
    return group;
  });

  return parsedGroup;
}
