import isEmpty from 'lodash/isEmpty';

/**
 * @param Object
 * @description Determine the operator and the input value of a rule
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
      parsedRule.operator = 'is empty';
    } else {
      const [input] = Object.values(rule);
      parsedRule.input = input;
      parsedRule.operator = 'equal';
    }
  } else if (typeof ruleValue === 'object') {
    /**
     * Switch to determine the right operator, and then assign the right input value
     */
    switch (operator) {
      case ('$eq'): {
        const [input] = Object.values(rule);
        parsedRule.input = input;
        parsedRule.operator = 'equal';
        break;
      }
      case ('$ne'): {
        if (Object.values(ruleValue)[0] == null) {
          parsedRule.operator = 'is not null';
        } else if (Object.values(ruleValue)[0] === '') {
          parsedRule.operator = 'is not empty';
        } else {
          const [inputObject] = Object.values(rule);
          const [input] = Object.values(inputObject);
          parsedRule.input = input;
          parsedRule.operator = 'not equal';
        }
        break;
      }
      case ('$in'): {
        const [inputArray] = Object.values(ruleValue);
        const [input] = inputArray;
        parsedRule.input = input;
        parsedRule.operator = 'in';
        break;
      }
      case ('$nin'): {
        const [inputArray] = Object.values(ruleValue);
        const [input] = inputArray;
        parsedRule.input = input;
        parsedRule.operator = 'not in';
        break;
      }
      default: {
        /**
         * Throw an error if the operator was not found.
         */
        throw new Error('Operator not found');
      }
    }
  }

  return parsedRule;
}

function parseRuleToFilter(rule) {
  const parsedRule = {
    field: '',
    operator: '',
    input: '',
  };

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
  const parsedGroup = {
    condition: '$and',
    groups: [],
    rules: [],
  };

  const groupContent = Object.values(group)[0];

  if (isEmpty(groupContent)) {
    return parsedGroup;
  }

  const [condition] = Object.keys(group);
  parsedGroup.condition = condition;

  /**
  * Map over the items of a group.
  * If the item is an array -> It's a group.
  * Else -> It's a rule.
  */
  groupContent.map((item) => {
    if (Array.isArray(Object.values(item)[0])) {
      parsedGroup.groups.push(parseGroupToFilter(item));
    } else {
      parsedGroup.rules.push(parseRuleToFilter(item));
    }
    return group;
  });

  return parsedGroup;
}
