import { cloneDeep } from 'lodash';

import { patternToForm } from '@/helpers/forms/pattern';

import { addKeyInEntities, removeKeyFromEntities } from '../entities';
import parseFilterToRequest from '../filter/editor/parse-filter-to-request';

/**
 * @typedef {Object} FilterFormRules
 * @property {string} field
 * @property {string} operator
 * @property {any} input
 */

/**
 * @typedef {Object} FilterForm
 * @property {string} title
 * @property {Object<FilterForm>} alarm_pattern
 * @property {Object<FilterFormRules>} entity_pattern
 * @property {Object<FilterFormRules>} pbehavior_pattern
 * @property {Object<FilterFormRules>} event_pattern
 */

/**
 * Convert filter object to filter form
 *
 * @param {Object} [filter = {}]
 * @returns {FilterForm}
 */
export const filterToForm = (filter = {}) => ({
  title: filter.title ?? '',
  alarm_pattern: patternToForm(filter.alarm_pattern),
  entity_pattern: patternToForm(filter.entity_pattern),
  pbehavior_pattern: patternToForm(filter.pbehavior_pattern),
  event_pattern: patternToForm(filter.event_pattern),
});

/**
 * Convert filter form to filter
 *
 * @param {FilterForm} form
 * @returns {Object}
 */
export const formToFilter = form => parseFilterToRequest(form);

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
