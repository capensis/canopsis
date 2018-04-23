function parseFilterRuleToRequest(rule) {
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

function parseFilterGroupToRequest(group) {
  const parsedGroup = {};

  parsedGroup[group.condition] = [];

  // Parse each rule of a group and add it to the parsedGroup array
  group.rules.map((rule) => {
    parsedGroup[group.condition].push(parseFilterRuleToRequest(rule));
    return parsedGroup;
  });

  // Parse each group of the group ans add it to the parsedGroup array
  group.groups.map((item) => {
    parsedGroup[group.condition].push(parseFilterGroupToRequest(item));
    return parsedGroup;
  });

  return parsedGroup;
}

export default function parseFilterToRequest(filter) {
  const request = {};

  request[filter[0].condition] = [];

  filter[0].rules.map((rule) => {
    request[filter[0].condition].push(parseFilterRuleToRequest(rule));
    return request[filter[0].condition];
  });
  filter[0].groups.map((group) => {
    request[filter[0].condition].push(parseFilterGroupToRequest(group));
    return request[filter[0].condition];
  });

  return request;
}
