import { formGroupsToPatternRules, patternToForm } from '@/helpers/forms/pattern';

import { PATTERN_CUSTOM_ITEM_VALUE, PATTERNS_FIELDS } from '@/constants';

import { addKeyInEntities, removeKeyFromEntities } from '../entities';

/**
 * @typedef {Object} FilterPatterns
 * @property {PatternGroups} [alarm_pattern]
 * @property {string} [corporate_alarm_pattern]
 * @property {PatternGroups} [entity_pattern]
 * @property {string} [corporate_entity_pattern]
 * @property {PatternGroups} [pbehavior_pattern]
 * @property {string} [corporate_pbehavior_pattern]
 * @property {PatternGroups} [event_pattern]
 * @property {string} [corporate_event_pattern]
 * @property {Object} [old_mongo_query]
 */

/**
 * @typedef {FilterPatterns} Filter
 * @property {string} title
 * @property {string} author
 * @property {boolean} is_private
 */

/**
 * @typedef {Object} FilterFormRules
 * @property {string} field
 * @property {string} operator
 * @property {any} input
 */

/**
 * @typedef {Object} FilterPatternsForm
 * @property {PatternGroupsForm} [alarm_pattern]
 * @property {PatternGroupsForm} [entity_pattern]
 * @property {PatternGroupsForm} [pbehavior_pattern]
 */

/**
 * @typedef {Filter & FilterPatternsForm} FilterForm
 */

/**
 * Convert filter patterns to form
 *
 * @param {Filter} filter
 * @param {PatternsFields} [fields = [PATTERNS_FIELDS.alarm, PATTERNS_FIELDS.pbehavior, PATTERNS_FIELDS.entity]]
 * @return {FilterPatterns}
 */
export const filterPatternsToForm = (
  filter = {},
  fields = [PATTERNS_FIELDS.alarm, PATTERNS_FIELDS.pbehavior, PATTERNS_FIELDS.entity],
) => fields.reduce((acc, field) => {
  const { [`corporate_${field}`]: id, [field]: pattern } = filter;

  acc[field] = patternToForm({ [field]: pattern, id });

  return acc;
}, {});

/**
 * Convert filter object to filter form
 *
 * @param {Object} [filter = {}]
 * @returns {FilterForm}
 */
export const filterToForm = (filter = {}) => ({
  title: filter.title ?? '',
  old_mongo_query: filter.old_mongo_query,
  is_private: filter.is_private ?? false,
  ...filterPatternsToForm(filter),
});

/**
 * Convert patterns form to patterns
 *
 * @param {FilterPatternsForm} form
 * @param {PatternsFields} [fields = [PATTERNS_FIELDS.alarm, PATTERNS_FIELDS.pbehavior, PATTERNS_FIELDS.entity]]
 * @return {{}}
 */
export const formFilterToPatterns = (
  form,
  fields = [
    PATTERNS_FIELDS.alarm,
    PATTERNS_FIELDS.entity,
    PATTERNS_FIELDS.pbehavior,
  ],
) => fields.reduce((acc, field) => {
  const patterns = form[field];

  if (!patterns) {
    return acc;
  }

  if (patterns.id !== PATTERN_CUSTOM_ITEM_VALUE) {
    acc[`corporate_${field}`] = patterns.id;

    return acc;
  }

  if (patterns.groups) {
    acc[field] = formGroupsToPatternRules(patterns.groups);
  }

  return acc;
}, {});

/**
 * Convert filter form to filter
 *
 * @param {FilterForm} form
 * @param {PatternsFields} fields
 * @returns {Filter}
 */
export const formToFilter = (form, fields) => ({
  title: form.title,
  is_private: form.is_private,
  ...formFilterToPatterns(form, fields),
});

/**
 * Convert filters to filters form
 *
 * @param {Array} [filters = []]
 * @returns {Array}
 */
export const filtersToForm = (filters = []) => addKeyInEntities(filters);

/**
 * Convert filters form to filters object
 *
 * @param {Array} [filters = []]
 * @returns {Array}
 */
export const formToFilters = (filters = []) => removeKeyFromEntities(filters);
