import { PATTERNS_FIELDS } from '@/constants';
import { COLORS } from '@/config';

import { filterPatternsToForm, formFilterToPatterns } from '@/helpers/entities/filter/form';

/**
 * @typedef {AlarmTag} AlarmTagForm
 * @property {FilterPatternsForm} patterns
 */

/**
 * Convert tag to patterns
 *
 * @param {AlarmTag} tag
 * @return {FilterPatterns}
 */
export const tagFilterPatternsToForm = tag => filterPatternsToForm(
  tag,
  [PATTERNS_FIELDS.alarm, PATTERNS_FIELDS.entity],
);

/**
 * Convert tag to form object
 *
 * @param {AlarmTag} [tag = {}]
 * @returns {AlarmTagForm}
 */
export const tagToForm = (tag = {}) => ({
  name: tag.name ?? '',
  color: tag.color ?? COLORS.secondary,
  patterns: tagFilterPatternsToForm(tag),
});

/**
 * Convert form patterns to tag patterns
 *
 * @param {FilterPatternsForm} patterns
 * @returns {FilterPatterns}
 */
const formPatternsToTagPatterns = patterns => formFilterToPatterns(
  patterns,
  [PATTERNS_FIELDS.alarm, PATTERNS_FIELDS.entity],
);

/**
 * Convert tag form to user object
 *
 * @param {AlarmTagForm} [form = {}]
 * @returns {AlarmTag}
 */
export const formToTag = (form = {}) => ({
  name: form.name,
  color: form.color,
  ...formPatternsToTagPatterns(form.patterns),
});
