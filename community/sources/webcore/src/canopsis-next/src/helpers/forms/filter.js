import { cloneDeep, isEmpty, isString } from 'lodash';

import { FILTER_DEFAULT_VALUES, QUICK_RANGES, TIME_UNITS } from '@/constants';

import uid from '@/helpers/uid';

import { addKeyInEntities, removeKeyFromEntities } from '../entities';
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
 * @typedef {Object} FilterRule
 */

/**
 * @typedef {FilterRule} FilterRuleForm
 * @property {string} attribute
 * @property {string} field
 * @property {string} value
 * @property {number|string} from
 * @property {number|string} to
 * @property {string} dictionary
 * @property {string} key
 */

/**
 * @typedef {Object} FilterGroup
 * @property {[]} rules
 */

/**
 * @typedef {FilterGroup} FilterGroupForm
 * @property {FilterRuleForm[]} rules
 * @property {string} key
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
    console.error(err);

    return {};
  }
}

/**
 * Convert filters to filters form
 *
 * @param {Array} [filters = []]
 * @returns {Array}
 */
export const filtersToForm = (filters = []) => cloneDeep(addKeyInEntities(filters));

/**
 * Convert filters form to filters object
 *
 * @param {Array} [filters = []]
 * @returns {Array}
 */
export const formToFilters = (filters = []) => removeKeyFromEntities(filters);

/**
 * Convert filter rule to filter form
 *
 * TODO: Should be finished after backend will be connected
 * @return {FilterRuleForm}
 */
export const filterRuleToForm = (rule = {}) => ({
  key: uid(),
  attribute: 'connector',
  field: '',
  dictionary: '',
  value: '',
  range: {
    type: QUICK_RANGES.last1Hour.value,
    from: 0,
    to: 0,
  },
  duration: {
    value: 1,
    unit: TIME_UNITS.second,
  },
  ...rule,
});

/**
 * Convert filter group to filter form
 *
 * TODO: Should be finished when backend will be connected
 * @return {FilterRuleForm[]}
 */
export const filterGroupRulesToForm = (rules = []) => rules.map(filterRuleToForm);

/**
 * Convert filter group to filter form
 *
 * TODO: Should be finished when backend will be connected
 * @param {FilterGroup} group
 * @return {FilterGroupForm}
 */
export const filterGroupToForm = (group = {}) => ({
  key: uid(),
  rules: filterGroupRulesToForm(group.rules),
});
