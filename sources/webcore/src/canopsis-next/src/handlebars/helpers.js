import { isFunction, isObject } from 'lodash';

function isOptions(val) {
  return isObject(val) && isObject(val.hash);
}

function isBlock(options) {
  return isOptions(options) && isFunction(options.fn) && isFunction(options.inverse);
}

function value(val, context, options) {
  if (isOptions(val)) {
    return value(null, val, options);
  }

  if (isOptions(context)) {
    return value(val, {}, context);
  }

  if (isBlock(options)) {
    return val ? options.fn(context) : options.inverse(context);
  }

  return val;
}

export function compare(a, operator, b, options) {
  if (arguments.length < 4) {
    throw new Error('handlebars Helper {{compare}} expects 4 arguments');
  }

  let result;

  switch (operator) {
    case '==':
      result = a == b; // eslint-disable-line eqeqeq
      break;
    case '===':
      result = a === b;
      break;
    case '!=':
      result = a != b; // eslint-disable-line eqeqeq
      break;
    case '!==':
      result = a !== b;
      break;
    case '<':
    case '&lt;':
      result = a < b;
      break;
    case '>':
    case '&gt;':
      result = a > b;
      break;
    case '<=':
    case '&lte;':
      result = a <= b;
      break;
    case '>=':
    case '&gte;':
      result = a >= b;
      break;
    case 'typeof':
      result = typeof a === b; // eslint-disable-line valid-typeof
      break;
    case 'regex':
    case 'regexp':
      result = new RegExp(b, options.flags).test(b);
      break;
    default: {
      throw new Error(`helper {{compare}}: invalid operator: '${operator}'`);
    }
  }

  return value(result, this, options);
}
