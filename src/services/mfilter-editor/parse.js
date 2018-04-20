export function parseRule2Request(rule) {
  const parsedRule = {};

  // Determine the good syntax for the chosen operator
  switch (rule.operator) {
    case ('equal'): {
      parsedRule[rule.field] = rule.input;
      break;
    }
    case ('not equal'): {
      parsedRule[rule.field] = {
        $ne: rule.input,
      };
      break;
    }
    case ('in'): {
      parsedRule[rule.field] = {
        $in: [rule.input],
      };
      break;
    }
    case ('not in'): {
      parsedRule[rule.field] = {
        $nin: [rule.input],
      };
      break;
    }
    case ('begins with'): {
      parsedRule[rule.field] = {
        $regex: `^${rule.input}`,
      };
      break;
    }
    case ('doesn\'t begin with'): {
      parsedRule[rule.field] = {
        $regex: `^(?!'${rule.input}')`,
      };
      break;
    }
    case ('contains'): {
      parsedRule[rule.field] = {
        $regex: rule.input,
      };
      break;
    }
    case ('doesn\'t contain'): {
      parsedRule[rule.field] = {
        $regex: `^((?!'${rule.input}').)*$`,
        $options: 's',
      };
      break;
    }
    case ('ends with'): {
      parsedRule[rule.field] = {
        $regex: `${rule.input}$`,
      };
      break;
    }
    case ('doesn\'t end with'): {
      parsedRule[rule.field] = {
        $regex: `(?<!'${rule.input}')$`,
      };
      break;
    }
    case ('is empty'): {
      parsedRule[rule.field] = '';
      break;
    }
    case ('is not empty'): {
      parsedRule[rule.field] = {
        $ne: '',
      };
      break;
    }
    case ('is null'): {
      parsedRule[rule.field] = null;
      break;
    }
    case ('is not null'): {
      parsedRule[rule.field] = {
        $ne: null,
      };
      break;
    }
    default: {
      break;
    }
  }

  return parsedRule;
}

export function parseGroup2Request(group) {
  const parsedGroup = {};

  parsedGroup[group.condition] = [];

  // Parse each rule of a group and add it to the parsedGroup array
  group.rules.map((rule) => {
    parsedGroup[group.condition].push(parseRule2Request(rule));
    return parsedGroup;
  });

  // Parse each group of the group ans add it to the parsedGroup array
  group.groups.map((item) => {
    parsedGroup[group.condition].push(parseGroup2Request(item));
    return parsedGroup;
  });

  return parsedGroup;
}

export function parseRule2Filter(rule) {
  const parsedRule = {
    field: '',
    operator: '',
    input: '',
  };

  const [field] = Object.keys(rule);
  parsedRule.field = field;

  // Switch to determine if it's a short syntax for '$eq' and '$eq:'''
  switch (typeof Object.values(rule)[0]) {
    case ('string'): {
      if (Object.values(rule)[0] === '') {
        parsedRule.operator = 'is empty';
      } else {
        const [input] = Object.values(rule);
        parsedRule.input = input;
        parsedRule.operator = 'equal';
      }
      break;
    }

    case ('object'): {
      // Handle the particular 'is null' case
      if (Object.values(rule)[0] === null) {
        parsedRule.operator = 'is null';
        break;
      }

      const value = Object.values(rule)[0];
      const operator = Object.keys(value)[0];

      // Switch to dtermine the right operator, and then assign the right input value
      switch (operator) {
        case ('$eq'): {
          const [input] = Object.values(rule);
          parsedRule.input = input;
          parsedRule.operator = 'equal';
          break;
        }
        case ('$ne'): {
          if (Object.values(value)[0] == null) {
            // CAUSE UNE ERREUR
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
          return null;
        }
      }
      break;
    }

    default: {
      return null;
    }
  }

  return parsedRule;
}

export function parseGroup2Filter(group) {
  const parsedGroup = {
    condition: '',
    groups: [],
    rules: [],
  };

  const [condition] = Object.keys(group);
  parsedGroup.condition = condition;

  Object.values(group)[0].map((item) => {
    if (Object.values(item)[0].constructor === Array) {
      parsedGroup.groups.push(parseGroup2Filter(item));
    } else if (parseRule2Filter(item) != null) {
      parsedGroup.rules.push(parseRule2Filter(item));
    }
    return group;
  });
  return parsedGroup;
}
