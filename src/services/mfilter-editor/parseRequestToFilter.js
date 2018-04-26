/**
 * @param Object
 * @description Determine the operator and the input value of a rule
 */
function ruleOperatorAndInput(rule) {
  const parsedRule = {
    operator: '',
    input: '',
  };

  /**
   * Switch to determine if it's a short syntax for '$eq' and '$eq:'''
   */
  if (typeof Object.values(rule)[0] === 'string') {
    if (Object.values(rule)[0] === '') {
      parsedRule.operator = 'is empty';
    } else {
      const [input] = Object.values(rule);
      parsedRule.input = input;
      parsedRule.operator = 'equal';
    }
  } else if (typeof Object.values(rule)[0] === 'object') {
    /**
     * Handle the particular 'is null' case
     */
    if (Object.values(rule)[0] === null) {
      parsedRule.operator = 'is null';
    }

    const value = Object.values(rule)[0];
    const operator = Object.keys(value)[0];

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
        if (Object.values(value)[0] == null) {
          // TODO: CAUSE UNE ERREUR
          parsedRule.operator = 'is not null';
        } else if (Object.values(value)[0] === '') {
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
        const [inputArray] = Object.values(value);
        const [input] = inputArray;
        parsedRule.input = input;
        parsedRule.operator = 'in';
        break;
      }
      case ('$nin'): {
        const [inputArray] = Object.values(value);
        const [input] = inputArray;
        parsedRule.input = input;
        parsedRule.operator = 'not in';
        break;
      }
      case ('$regex'): {
        // TODO : Finish for all regexp cases
        const beginsWith = /^\^/;
        if (Object.values(value)[0].match(beginsWith) != null) {
          parsedRule.operator = 'begins with';
          parsedRule.input = Object.values(value)[0].replace('^', '');
        }
        break;
      }
      default: {
        return parsedRule;
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

  const [field] = Object.keys(rule);
  parsedRule.field = field;

  const { operator, input } = ruleOperatorAndInput(rule);

  parsedRule.operator = operator;
  parsedRule.input = input;

  return parsedRule;
}

export default function parseGroupToFilter(group) {
  const parsedGroup = {
    condition: '',
    groups: [],
    rules: [],
  };

  const [condition] = Object.keys(group);
  parsedGroup.condition = condition;

  /**
  * Map over the items of a group.
  * If the item is an array -> It's a group.
  * Else -> It's a rule.
  */
  Object.values(group)[0].map((item) => {
    if (Array.isArray(Object.values(item)[0])) {
      parsedGroup.groups.push(parseGroupToFilter(item));
    } else if (parseRuleToFilter(item) != null) {
      parsedGroup.rules.push(parseRuleToFilter(item));
    }
    return group;
  });
  return parsedGroup;
}
