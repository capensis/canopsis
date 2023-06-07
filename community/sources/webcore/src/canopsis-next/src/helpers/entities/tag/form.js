import { PATTERNS_FIELDS } from '@/constants';
import { COLORS } from '@/config';

import { filterPatternsToForm, formFilterToPatterns } from '@/helpers/entities/filter/form';

/**
 * @typedef {AlarmTag} AlarmTagForm
 * @property {FilterPatternsForm} patterns
 */

/**
 * Convert user to form object
 *
 * @param {AlarmTag} [tag = {}]
 * @returns {AlarmTagForm}
 */
export const tagToForm = (tag = {}) => ({
  name: tag.name ?? '',
  color: tag.color ?? COLORS.secondary,
  patterns: filterPatternsToForm(
    tag,
    [PATTERNS_FIELDS.alarm, PATTERNS_FIELDS.entity],
  ),
});

/**
 * Convert user form to user object
 *
 * @param {AlarmTagForm} [form = {}]
 * @returns {AlarmTag}
 */
export const formToTag = (form = {}) => ({
  name: form.name,
  color: form.color,
  ...formFilterToPatterns(form.patterns, [PATTERNS_FIELDS.alarm, PATTERNS_FIELDS.entity]),
});
