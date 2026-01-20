import { mapValues } from 'lodash';

import { PATTERN_OPERATORS } from '@/constants';

/**
 * @typedef { 'string' | 'int' | 'timestamp' | 'duration' | 'reference' | 'object' | 'string_array' } PatternFieldType
 */

/**
 * @typedef {Object} PatternField
 * @property {string} name
 * @property {PatternFieldType} type
 * @property {boolean} enabled
 * @property {boolean} alias
 * @property {boolean} only_absolute_time_cond
 */

/**
 * @typedef {
 *   'alarm_pattern'
 *   | 'entity_pattern'
 *   | 'pbehavior_pattern'
 *   | 'total_entity_pattern'
 *   | 'service_weather_pattern'
 * } PatternFieldsKey
 */

/**
 * @typedef {Object<PatternFieldsKey, PatternField[]>} PatternFields
 */

/**
 * @typedef {Object} PatternFieldForm
 * @property {string} value
 * @property {boolean} alias
 * @property {Object} options
 * @property {boolean} [options.disabled]
 */

/**
 * @typedef {Object<PatternFieldsKey, PatternFieldForm[]>} PatternFieldsForm
 */

/**
 * Convert patterns fields to form
 *
 * @param {PatternFields | {}} [patternsFields = {}]
 * @return {PatternFieldsForm}
 */
export const patternsFieldsToForm = (patternsFields = {}) => mapValues(patternsFields, fields => (
  fields.map((field) => {
    const formField = {
      value: field.name,
      alias: field.alias,
      options: {
        disabled: !field.enabled,
      },
    };

    if (field.only_absolute_time_cond) {
      formField.options.operators = [PATTERN_OPERATORS.inRangeDates];
    }

    return formField;
  })
));
