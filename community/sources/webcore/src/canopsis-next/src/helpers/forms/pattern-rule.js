import uid from '@/helpers/uid';
import { isSimpleRule } from '@/helpers/treeview';

/**
 * @typedef {Object} PatternRuleFormAdvancedField
 * @property {string} key
 * @property {string} operator
 * @property {string|number|boolean|null} value
 */

/**
 * @typedef {Object} PatternRuleForm
 * @property {string} field
 * @property {string|number|boolean|null} value
 * @property {boolean} advancedMode
 * @property {Array<PatternRuleFormAdvancedField>} advancedFields
 */

/**
 * Convert pattern rule to form
 *
 * @param {string} [field = '']
 * @param {string} [value = '']
 * @returns {PatternRuleForm}
 */
export function patternRuleToForm({ field = '', value = '' } = {}) {
  const isSimple = isSimpleRule(value);
  const form = {
    field,
    value: '',
    advancedMode: !isSimple,
    advancedFields: [],
  };

  if (isSimple) {
    form.value = value;
  } else {
    form.advancedFields = Object.entries(value)
      .map(([fieldKey, fieldValue]) => ({ key: uid(), operator: fieldKey, value: fieldValue }));
  }

  return form;
}

/**
 * Convert form to pattern rule
 *
 * @param {PatternRuleForm} form
 * @returns {Object}
 */
export function formToPatternRule(form) {
  if (!form.advancedMode) {
    return {
      field: form.field,
      value: form.value,
    };
  }

  const value = form.advancedFields.reduce((acc, field) => {
    acc[field.operator] = field.value;

    return acc;
  }, {});

  return {
    value,

    field: form.field,
  };
}
