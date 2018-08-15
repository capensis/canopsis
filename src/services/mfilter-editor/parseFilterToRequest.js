import isEmpty from 'lodash/isEmpty';

import { OPERATORS } from '@/constants';

function parseFilterRuleToRequest(rule) {
  const parsedRule = {};

  if (rule.operator === '' || rule.field === '') {
    throw new Error('Invalid Rule');
  }

  /**
   * Determine the rule syntax based on the rule operator
   */
  switch (rule.operator) {
    case OPERATORS.equal: {
      parsedRule[rule.field] = rule.input;
      break;
    }
    case OPERATORS.notEqual: {
      parsedRule[rule.field] = {
        $ne: rule.input,
      };
      break;
    }
    case OPERATORS.in: {
      parsedRule[rule.field] = {
        $in: [rule.input],
      };
      break;
    }
    case OPERATORS.notIn: {
      parsedRule[rule.field] = {
        $nin: [rule.input],
      };
      break;
    }
    case OPERATORS.beginsWith: {
      parsedRule[rule.field] = {
        $regex: `^${rule.input}`,
      };
      break;
    }
    case OPERATORS.doesntBeginWith: {
      parsedRule[rule.field] = {
        $regex: `^(?!'${rule.input}')`,
      };
      break;
    }
    case OPERATORS.contains: {
      parsedRule[rule.field] = {
        $regex: rule.input,
      };
      break;
    }
    case OPERATORS.doesntContains: {
      parsedRule[rule.field] = {
        $regex: `^((?!'${rule.input}').)*$`,
        $options: 's',
      };
      break;
    }
    case OPERATORS.endsWith: {
      parsedRule[rule.field] = {
        $regex: `${rule.input}$`,
      };
      break;
    }
    case OPERATORS.doesntEndWith: {
      parsedRule[rule.field] = {
        $regex: `(?<!'${rule.input}')$`,
      };
      break;
    }
    case OPERATORS.isEmpty: {
      parsedRule[rule.field] = '';
      break;
    }
    case OPERATORS.isNotEmpty: {
      parsedRule[rule.field] = {
        $ne: '',
      };
      break;
    }
    case OPERATORS.isNull: {
      parsedRule[rule.field] = null;
      break;
    }
    case OPERATORS.isNotNull: {
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

  if (isEmpty(group.groups) && isEmpty(group.rules)) {
    throw new Error('Empty Group');
  }

  /**
   * Parse each rule of a group and add it to the parsedGroup array
   */
  group.rules.map((rule) => {
    try {
      parsedGroup[group.condition].push(parseFilterRuleToRequest(rule));
    } catch (e) {
      parsedGroup[group.condition].push({});
    }

    return parsedGroup;
  });

  /**
   * Parse each group of the group ans add it to the parsedGroup array
   */
  group.groups.map((item) => {
    if (isEmpty(group.groups) && isEmpty(group.rules)) {
      throw new Error('Empty Group');
    }

    try {
      return parsedGroup[group.condition].push(parseFilterGroupToRequest(item));
    } catch (e) {
      return parsedGroup[group.condition].push({});
    }
  });

  return parsedGroup;
}

/**
 * @param {Array} filter
 *
 * @description Take the filter data (from the form of the visual editor)
 * , and return a JSON object representing the request
 */
export default function parseFilterToRequest(filter) {
  const request = {};

  request[filter[0].condition] = [];

  filter[0].rules.map((rule) => {
    try {
      return request[filter[0].condition].push(parseFilterRuleToRequest(rule));
    } catch (e) {
      return request[filter[0].condition].push({});
    }
  });

  /**
   * Map on the filter's group array, and call itself recursively
   * to parse Groups, and group's groups, etc
   */

  filter[0].groups.map((group) => {
    try {
      return request[filter[0].condition].push(parseFilterGroupToRequest(group));
    } catch (e) {
      return request[filter[0].condition].push({});
    }
  });

  return request;
}
