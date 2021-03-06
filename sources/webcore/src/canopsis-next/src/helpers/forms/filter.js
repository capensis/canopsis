import { cloneDeep, isEmpty, isString } from 'lodash';

import { FILTER_DEFAULT_VALUES } from '@/constants';

import { addKeyInEntity, removeKeyFromEntity } from '../entities';
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

/**
 * Convert filter string to object
 *
 * @param {Object|string} [filter = {}]
 * @returns {Object}
 */
export function filterToObject(filter = {}) {
  try {
    return isString(filter) ? JSON.parse(filter) : filter;
  } catch (err) {
    console.warn(err);

    return {};
  }
}

/**
 * Convert filters to filters form
 *
 * @param {Array} [filters = []]
 * @returns {Array}
 */
export const filtersToForm = (filters = []) => cloneDeep(addKeyInEntity(filters));

/**
 * Convert filters form to filters object
 *
 * @param {Array} [filters = []]
 * @returns {Array}
 */
export const formToFilters = (filters = []) => removeKeyFromEntity(filters);
