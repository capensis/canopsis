import uuid from '@/helpers/uuid';
import uid from '@/helpers/uid';

/**
 * @typedef {Object} KeyValue
 * @property {string} key
 * @property {string} value
 */

/**
 * @typedef {Object} DynamicInfoTemplate
 * @property {string} _id
 * @property {string} title
 * @property {string[]} names
 */

/**
 * @typedef {DynamicInfoTemplate} DynamicInfoTemplateForm
 * @property {KeyValue[]} names
 */

/**
 * Generate dynamic info template form name object by dynamic info template name string
 *
 * @param {string} [name = '']
 * @returns {KeyValue}
 */
export function generateTemplateFormName(name = '') {
  return {
    key: uid(),
    value: name,
  };
}

/**
 * Convert a dynamic information's template object to a dynamic information's template form object
 *
 * @param {string} [_id = '']
 * @param {string} [title = '']
 * @param {string[]} [names = []]
 * @returns {DynamicInfoTemplateForm}
 */
export function templateToForm({
  _id = uuid(),
  title = '',
  names = [],
} = {}) {
  return {
    _id,
    title,
    names: names.map(generateTemplateFormName),
  };
}

/**
 * Convert a dynamic information's template form object to a dynamic information's template object
 *
 * @param {string} [_id = uuid()]
 * @param {string} [title = '']
 * @param {KeyValue[]} [names = []]
 * @returns {DynamicInfoTemplate}
 */
export function formToTemplate({
  _id = uuid(),
  title = '',
  names = [],
} = {}) {
  return {
    _id,
    title,
    names: names.map(({ value }) => value),
  };
}

/**
 * Convert template to dynamic info's informations array
 *
 * @param {DynamicInfoTemplate} template
 * @returns {{ name: string, value: string }[]}
 */
export function templateToDynamicInfoInfos(template) {
  return template.names.map(name => ({ name, value: '' }));
}
