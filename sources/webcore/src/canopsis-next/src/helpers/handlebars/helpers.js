import { get, isFunction } from 'lodash';
import Handlebars from 'handlebars';

import dateFilter from '@/filters/date';
import durationFilter from '@/filters/duration';
import { DATETIME_FORMATS, ENTITY_INFOS_TYPE } from '@/constants';

/**
 * Prepare object attributes from `{ key: value, keySecond: valueSecond }` format
 * to `'escape(key)=escape(value) escape(keySecond)=escape(valueSecond)'` format.
 *
 * @param {Object} attributes
 * @returns {string}
 */
function prepareAttributes(attributes) {
  return Object.entries(attributes)
    .map(([key, value]) =>
      `${Handlebars.escapeExpression(key)}="${Handlebars.escapeExpression(value)}"`)
    .join(' ');
}

/**
 * Convert date to long format
 *
 * @param {string|number} date
 * @returns {string}
 */
export function timestamp(date) {
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
export function internalLink(options) {
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
 * Return date in long format
 *
 * Example: {{date 1000000}} -> 12/01/1970 20:46:40
 *
 * @param time
 * @returns {string}
 */
export function dateHelper(time) {
  return dateFilter(time, 'long');
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
  return new Handlebars.SafeString(`<alarm-chips type="${ENTITY_INFOS_TYPE.state}" value="${state}"></alarm-chips>`);
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
export function compare(a, operator, b, options = {}) {
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
    case '&lt;=':
      result = a <= b;
      break;
    case '>=':
    case '&gte;':
    case '&gt;=':
      result = a >= b;
      break;
    case 'typeof':
      result = typeof a === b; // eslint-disable-line valid-typeof
      break;
    case 'regex':
    case 'regexp':
      result = new RegExp(b, get(options, ['hash', 'flags'])).test(a);
      break;
    default:
      throw new Error(`helper {{compare}}: invalid operator: '${operator}'`);
  }

  if (isFunction(options.fn) && isFunction(options.inverse)) {
    return result ? options.fn(this) : options.inverse(this);
  }

  return result;
}
