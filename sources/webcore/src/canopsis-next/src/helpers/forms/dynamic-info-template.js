import uuid from '@/helpers/uuid';
import uid from '@/helpers/uid';

/**
 * Generate dynamic info template form name object by dynamic info template name string
 *
 * @param {string} [name = '']
 * @returns {{key: string, value: string}}
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
 * @param {Array<string>} [names = []]
 * @returns {{_id: string, title: string, names: {key: string, value: string}[]}}
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
 * @param {Array<{ key: string, value: string }>} [names = []]
 * @returns {{_id: string, title: string, names: string[]}}
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
