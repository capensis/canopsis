import { isEmpty } from 'lodash';

import { FILTER_OPERATORS, FILTER_MONGO_OPERATORS } from '@/constants';

function parseFilterRuleToRequest(rule) {
  const parsedRule = {};

  if (rule.operator === '' || rule.field === '') {
    throw new Error('Invalid Rule');
  }

  /**
   * Determine the rule syntax based on the rule operator
   */
  switch (rule.operator) {
    case FILTER_OPERATORS.equal: {
      parsedRule[rule.field] = rule.input;
      break;
    }
    case FILTER_OPERATORS.notEqual: {
      parsedRule[rule.field] = {
        [FILTER_MONGO_OPERATORS.notEqual]: rule.input,
      };
      break;
    }
    case FILTER_OPERATORS.hasOneOf: {
      parsedRule[rule.field] = {
        [FILTER_MONGO_OPERATORS.in]: rule.input.map(({ value }) => value),
      };
      break;
    }
    case FILTER_OPERATORS.hasNot: {
      parsedRule[rule.field] = {
        [FILTER_MONGO_OPERATORS.notIn]: rule.input.map(({ value }) => value),
      };
      break;
    }
    case FILTER_OPERATORS.beginsWith: {
      parsedRule[rule.field] = {
        [FILTER_MONGO_OPERATORS.regex]: `^${rule.input}`,
      };
      break;
    }
    case FILTER_OPERATORS.notBeginWith: {
      parsedRule[rule.field] = {
        [FILTER_MONGO_OPERATORS.regex]: `^(?!${rule.input})`,
      };
      break;
    }
    case FILTER_OPERATORS.contains: {
      parsedRule[rule.field] = {
        [FILTER_MONGO_OPERATORS.regex]: rule.input,
      };
      break;
    }
    case FILTER_OPERATORS.notContains: {
      parsedRule[rule.field] = {
        [FILTER_MONGO_OPERATORS.regex]: `^((?!${rule.input}).)*$`,
        $options: 's',
      };
      break;
    }
    case FILTER_OPERATORS.endsWith: {
      parsedRule[rule.field] = {
        [FILTER_MONGO_OPERATORS.regex]: `${rule.input}$`,
      };
      break;
    }
    case FILTER_OPERATORS.notEndWith: {
      parsedRule[rule.field] = {
        [FILTER_MONGO_OPERATORS.regex]: `(?<!${rule.input})$`,
      };
      break;
    }
    case FILTER_OPERATORS.isEmpty: {
      parsedRule[rule.field] = '';
      break;
    }
    case FILTER_OPERATORS.isNotEmpty: {
      parsedRule[rule.field] = {
        [FILTER_MONGO_OPERATORS.notEqual]: '',
      };
      break;
    }
    case FILTER_OPERATORS.isNull: {
      parsedRule[rule.field] = null;
      break;
    }
    case FILTER_OPERATORS.isNotNull: {
      parsedRule[rule.field] = {
        [FILTER_MONGO_OPERATORS.notEqual]: null,
      };
      break;
    }
    case FILTER_OPERATORS.greater: {
      parsedRule[rule.field] = {
        [FILTER_MONGO_OPERATORS.greater]: rule.input,
      };
      break;
    }
    case FILTER_OPERATORS.less: {
      parsedRule[rule.field] = {
        [FILTER_MONGO_OPERATORS.less]: rule.input,
      };
      break;
    }
    default: {
      parsedRule[rule.field] = {};
      parsedRule[rule.field][rule.operator] = rule.input;
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

  const rules = Object.values(group.rules);
  const groups = Object.values(group.groups);

  /**
   * Parse each rule of a group and add it to the parsedGroup array
   */
  rules.map((rule) => {
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
  groups.map((item) => {
    try {
      return parsedGroup[group.condition].push(parseFilterGroupToRequest(item));
    } catch (e) {
      return parsedGroup[group.condition].push({});
    }
  });

  return parsedGroup;
}

/**
 * @param {Object} filter
 *
 * @description Take the filter data (from the form of the visual editor)
 * , and return a JSON object representing the request
 */
export default function parseFilterToRequest(filter) {
  const { condition } = filter;
  const request = {};

  const rules = Object.values(filter.rules);
  const groups = Object.values(filter.groups);

  request[condition] = [];

  rules.map((rule) => {
    try {
      return request[condition].push(parseFilterRuleToRequest(rule));
    } catch (e) {
      return request[condition].push({});
    }
  });

  /**
   * Map on the filter's group array, and call itself recursively
   * to parse Groups, and group's groups, etc
   */

  groups.map((group) => {
    try {
      return request[condition].push(parseFilterGroupToRequest(group));
    } catch (e) {
      return request[condition].push({});
    }
  });

  return request;
}
