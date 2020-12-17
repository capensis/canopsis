import { isEmpty, cloneDeep } from 'lodash';

import { FILTER_DEFAULT_VALUES } from '@/constants';

import parseGroupToFilter from '../filter/editor/parse-group-to-filter';
import parseFilterToRequest from '../filter/editor/parse-filter-to-request';

/**
 * @typedef {Object} FilterFormRules
 * @property {string} field
 * @property {string} operator
 * @property {any} input
 */

/**
 * @typedef {Object} FilterForm
 * @property {string} condition
 * @property {Object<FilterForm>} groups
 * @property {Object<FilterFormRules>} rules
 */

/**
 * Convert filter object to filter form
 *
 * @param {Object} [filter = {}]
 * @returns {FilterForm}
 */
export function filterToForm(filter = {}) {
  if (isEmpty(filter)) {
    return cloneDeep(FILTER_DEFAULT_VALUES.group);
  }

  return parseGroupToFilter(filter);
}

/**
 * Convert filter form to filter
 *
 * @param {FilterForm} form
 * @returns {Object}
 */
export function formToFilter(form) {
  return parseFilterToRequest(form);
}
