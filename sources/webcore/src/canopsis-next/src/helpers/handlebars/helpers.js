import { get, isFunction, unescape } from 'lodash';
import Handlebars from 'handlebars';
import axios from 'axios';

import dateFilter from '@/filters/date';
import durationFilter from '@/filters/duration';
import { DATETIME_FORMATS, ENTITY_INFOS_TYPE } from '@/constants';

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
    .map(([key, value]) =>
      `${Handlebars.escapeExpression(key)}="${Handlebars.escapeExpression(value)}"`)
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
  return new Handlebars.SafeString(`<alarm-chips type="${ENTITY_INFOS_TYPE.state}" value="${state}"></alarm-chips>`);
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
    headers = '{}',
    path,
    variable,
    username,
    password,
  } = options.hash;

  if (!url) {
    throw new Error('helper {{request}}: \'url\' is required');
  }

  try {
    const { data } = await axios({
      method,
      url: unescape(url),
      auth: { username, password },
      headers: JSON.parse(headers),
    });

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
