import { OLD_PATTERNS_FIELDS, PATTERN_CUSTOM_ITEM_VALUE, PATTERNS_FIELDS } from '@/constants';

import { formGroupsToPatternRules, patternToForm } from '@/helpers/forms/pattern';

/**
 * @typedef {Object} FilterPatterns
 * @property {PatternGroups} [alarm_pattern]
 * @property {string} [corporate_alarm_pattern]
 *
 * @property {PatternGroups} [entity_pattern]
 * @property {string} [corporate_entity_pattern]
 *
 * @property {PatternGroups} [pbehavior_pattern]
 * @property {string} [corporate_pbehavior_pattern]
 *
 * @property {PatternGroups} [event_pattern]
 * @property {string} [corporate_event_pattern]
 *
 * @property {PatternGroups} [total_entity_pattern]
 * @property {string} [corporate_total_entity_pattern]
 *
 * @property {PatternGroups} [weather_service_pattern]
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
 * @param {PatternsFields} [fields = [PATTERNS_FIELDS.alarm, PATTERNS_FIELDS.entity, PATTERNS_FIELDS.pbehavior]]
 * @param {string[]} [oldFields = []]
 * @return {FilterPatterns}
 */
export const filterPatternsToForm = (
  filter = {},
  fields = [PATTERNS_FIELDS.alarm, PATTERNS_FIELDS.entity, PATTERNS_FIELDS.pbehavior],
  oldFields = [OLD_PATTERNS_FIELDS.alarm, OLD_PATTERNS_FIELDS.entity],
) => fields.reduce((acc, field, index) => {
  const oldField = oldFields[index];
  const {
    [`corporate_${field}`]: id,
    [field]: pattern,
    [oldField]: oldMongoQuery,
  } = filter;

  acc[field] = patternToForm({
    [field]: pattern,
    old_mongo_query: !pattern?.length ? oldMongoQuery : undefined,
    is_corporate: !!id && id !== PATTERN_CUSTOM_ITEM_VALUE,
    id,
  });

  return acc;
}, {});

/**
 * Convert filter object to filter form
 *
 * @param {Object} [filter = {}]
 * @param {Array} [fields]
 * @returns {FilterForm}
 */
export const filterToForm = (filter = {}, fields) => ({
  title: filter.title ?? '',
  old_mongo_query: filter.old_mongo_query,
  is_private: filter.is_private ?? false,
  ...filterPatternsToForm(filter, fields),
});

/**
 * Convert patterns form to patterns
 *
 * @param {FilterPatternsForm} [form = {}]
 * @param {PatternsFields} [fields = [PATTERNS_FIELDS.alarm, PATTERNS_FIELDS.pbehavior, PATTERNS_FIELDS.entity]]
 * @param {boolean} [corporate = false]
 * @return {{}}
 */
export const formFilterToPatterns = (
  form = {},
  fields = [
    PATTERNS_FIELDS.alarm,
    PATTERNS_FIELDS.entity,
    PATTERNS_FIELDS.pbehavior,
  ],
  corporate = false,
) => fields.reduce((acc, field) => {
  const patterns = form[field];

  if (!patterns) {
    return acc;
  }

  if ((corporate || patterns.is_corporate) && patterns.id !== PATTERN_CUSTOM_ITEM_VALUE) {
    acc[`corporate_${field}`] = patterns.id;
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
 * @param {PatternsFields} [fields]
 * @param {boolean} [corporate]
 * @returns {Filter}
 */
export const formToFilter = (form, fields, corporate) => ({
  title: form.title,
  is_private: form.is_private,
  ...formFilterToPatterns(form, fields, corporate),
});
