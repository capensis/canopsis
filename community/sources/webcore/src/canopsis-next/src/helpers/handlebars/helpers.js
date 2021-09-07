import {
  get, isFunction, isNumber, isObject, unescape, isString,
} from 'lodash';
import Handlebars from 'handlebars';
import axios from 'axios';

import { DATETIME_FORMATS } from '@/constants';

import dateFilter from '@/filters/date';
import durationFilter from '@/filters/duration';

import i18n from '@/i18n';

/**
 * Prepare object attributes from `{ key: value, keySecond: valueSecond }` format
 * to `'escape(key)=escape(value) escape(keySecond)=escape(valueSecond)'` format.
 *
 * @param {Object} attributes
 * @returns {string}
 */
function prepareAttributes(attributes) {
  return Object.entries(attributes)
    .map(([key, value]) => `${Handlebars.escapeExpression(key)}="${Handlebars.escapeExpression(value)}"`)
    .join(' ');
}

/**
 * Convert date to long format
 *
 * Example: {{date 1000000}} -> 12/01/1970 20:46:40
 *
 * @param {string|number} date
 * @returns {string}
 */
export function timestampHelper(date) {
  let result = '';

  if (date) {
    result = dateFilter(date, 'long');
  }

  return result;
}

/**
 * Create special internal router link
 *
 * @param {Object} options
 * @returns {Handlebars.SafeString}
 */
export function internalLinkHelper(options) {
  const { href, text, ...attributes } = options.hash;
  const path = href.replace(window.location.origin, '');

  const link = `<router-link to="${path}" ${prepareAttributes(attributes)}>${text}</router-link>`;

  return new Handlebars.SafeString(link);
}

/**
 * Return duration by seconds
 *
 * First example: {{duration 120}} -> 2 mins
 * Second example: {{duration 12000}} -> 3 hrs 20 mins
 *
 * @param second
 * @returns {String}
 */
export function durationHelper(second) {
  return durationFilter(second, undefined, DATETIME_FORMATS.refreshFieldFormat);
}

/**
 * Return icon by alarm state
 *
 * Example {{state 0}} -> draw green element with ok text
 *
 * @param state
 * @returns {Handlebars.SafeString}
 */
export function alarmStateHelper(state) {
  return new Handlebars.SafeString(`<c-alarm-chip value="${state}"></c-alarm-chip>`);
}

/**
 * Pass response of a request to the child block
 *
 * Example:
 * {{#request method="get" url="https://test.com" path="data.users" variable="users"
 * username="test" password="test" headers='{ "test": "test2" }'}}
 *   {{#each users}}
 *     <li>{{login}}</li>
 *   {{/each}}
 * {{/request}}
 *
 * @param options
 * @returns {Promise<string|*>}
 */
export async function requestHelper(options) {
  const {
    method = 'get',
    url,
    headers,
    path,
    variable,
    username,
    password,
  } = options.hash;

  if (!url) {
    throw new Error('helper {{request}}: \'url\' is required');
  }

  try {
    const axiosOptions = {
      method,
      url: unescape(url),
    };

    if (headers) {
      axiosOptions.headers = JSON.parse(headers);
    }

    if (username || password) {
      axiosOptions.auth = { username, password };
    }

    const { data } = await axios(axiosOptions);

    if (isFunction(options.fn)) {
      const value = path ? get(data, path) : data;
      const context = variable ? { [variable]: value } : value;

      return options.fn(context);
    }

    return '';
  } catch (err) {
    console.error(err);

    const { status } = err.response || {};

    switch (status) {
      case 401:
        return i18n.t('handlebars.requestHelper.errors.unauthorized');
      case 408:
        return i18n.t('handlebars.requestHelper.errors.timeout');
      default:
        return i18n.t('handlebars.requestHelper.errors.other');
    }
  }
}

/**
 * Compare two parameters
 *
 * Number example: {{#compare 12 '>' 10}}PRINT SOMETHING{{/compare}}
 * String example: {{#compare 'test' '==' 'test'}}PRINT SOMETHING{{/compare}}
 * String regex example: {{#compare 'TEST' 'regex' 'est$' flags='i'}}PRINT SOMETHING{{/compare}}
 *
 * @param {string|number|null} a
 * @param {string} operator
 * @param {string|number|null} b
 * @param {Object} options
 * @returns {*}
 */
export function compareHelper(a, operator, b, options = {}) {
  if (arguments.length < 4) {
    throw new Error('handlebars Helper {{compare}} expects 4 arguments');
  }

  const flags = get(options, ['hash', 'flags']);
  let result;

  if (['regex', 'regexp'].includes(operator)) {
    result = new RegExp(b, flags).test(a);
  } else {
    let preparedA = a;
    let preparedB = b;

    if (flags && flags.search('i') !== -1) {
      preparedA = String(a).toLowerCase();
      preparedB = String(b).toLowerCase();
    }

    switch (operator) {
      case '==':
        result = preparedA == preparedB; // eslint-disable-line eqeqeq
        break;
      case '===':
        result = preparedA === preparedB;
        break;
      case '!=':
        result = preparedA != preparedB; // eslint-disable-line eqeqeq
        break;
      case '!==':
        result = preparedA !== preparedB;
        break;
      case '<':
      case '&lt;':
        result = preparedA < preparedB;
        break;
      case '>':
      case '&gt;':
        result = preparedA > preparedB;
        break;
      case '<=':
      case '&lte;':
      case '&lt;=':
        result = preparedA <= preparedB;
        break;
      case '>=':
      case '&gte;':
      case '&gt;=':
        result = preparedA >= preparedB;
        break;
      case 'typeof':
        result = typeof preparedA === preparedB; // eslint-disable-line valid-typeof
        break;
      default:
        throw new Error(`helper {{compare}}: invalid operator: '${operator}'`);
    }
  }

  if (isFunction(options.fn) && isFunction(options.inverse)) {
    return result ? options.fn(this) : options.inverse(this);
  }

  return result;
}

/**
 * Concat every primitive arguments
 *
 * Example: {{concat "example" object.field}}
 * Example with request helper: {{#request url=(concat "http://example.com/" object.field)}}something{{/request}}
 *
 * @param {...any} args
 * @returns {string}
 */
export function concatHelper(...args) {
  return args.reduce((acc, arg) => (!isObject(arg) ? acc + arg : acc), '');
}

/**
 * Sum for every number arguments
 *
 * Example: {{sum 1 2 3 4 5}}
 *
 * @param {...numbers} args
 * @returns {string}
 */
export function sumHelper(...args) {
  return args.reduce((acc, arg) => (isNumber(arg) ? acc + arg : acc), 0);
}

/**
 * Subtracting one number from the second
 *
 * Example: {{minus 10 1}}
 *
 * @param {number} a
 * @param {number} b
 * @returns {number}
 */
export function minusHelper(a, b) {
  if (!isNumber(a)) {
    throw new TypeError('expected the first argument to be a number');
  }

  if (!isNumber(b)) {
    throw new TypeError('expected the second argument to be a number');
  }

  return a - b;
}

/**
 * Multiple two numbers
 *
 * Example: {{mul 2 4}}
 *
 * @param {number} a
 * @param {number} b
 * @returns {number}
 */
export function mulHelper(a, b) {
  if (!isNumber(a)) {
    throw new TypeError('expected the first argument to be a number');
  }

  if (!isNumber(b)) {
    throw new TypeError('expected the second argument to be a number');
  }

  return a * b;
}

/**
 * Division of two numbers
 *
 * Example: {{divide 10 2}}
 *
 * @param {number} a
 * @param {number} b
 * @returns {number}
 */
export function divideHelper(a, b) {
  if (!isNumber(a)) {
    throw new TypeError('expected the first argument to be a number');
  }

  if (!isNumber(b)) {
    throw new TypeError('expected the second argument to be a number');
  }

  return a / b;
}

/**
 * Capitalize the first word in a string.
 *
 * @param {string} str
 * @returns {string}
 */
export function capitalizeHelper(str) {
  if (!isString(str)) {
    return '';
  }

  return str.charAt(0).toUpperCase() + str.slice(1);
}

/**
 * Capitalize all words in a string.
 *
 * @param {string} str
 * @returns {string}
 */
export function capitalizeAllHelper(str) {
  if (!isString(str)) {
    return '';
  }

  return str.replace(/\w\S*/g, capitalizeHelper);
}

/**
 * Lowercase all of characters in a string
 *
 * Example: {{lowercase 'test'}}
 *
 * @param {string|Object} str
 * @returns {string}
 */
export function lowercaseHelper(str) {
  if (isObject(str) && str.fn) {
    return str.fn(this).toLowerCase();
  }

  if (!isString(str)) {
    return '';
  }

  return str.toLowerCase();
}

/**
 * Uppercase all of characters in a string
 *
 * Example: {{uppercase 'test'}}
 *
 * @param {string|Object} str
 * @returns {string}
 */
export function uppercaseHelper(str) {
  if (isObject(str) && str.fn) {
    return str.fn(this).toUpperCase();
  }

  if (!isString(str)) {
    return '';
  }

  return str.toUpperCase();
}
