import { uid } from '@/helpers/uid';

/**
 * @typedef {Object & ObjectKey} KeyValue
 * @property {string} value
 */

/**
 * @typedef {DynamicInfoTemplate} DynamicInfoTemplateForm
 * @property {KeyValue[]} names
 */

/**
 * Create a dynamic info template name form item.
 *
 * @param {string} [name=''] - Template info name.
 * @returns {{ key: string, value: string }}
 */
export const generateTemplateFormName = (name = '') => ({
  key: uid('dynamic-info-template-name'),
  value: name,
});

/**
 * Convert a dynamic info template object to a form object.
 *
 * @param {Object} [template={}] - Dynamic info template entity.
 * @returns {{ _id: string, title: string, names: Array<{ key: string, value: string }> }}
 */
export const templateToForm = (template = {}) => ({
  _id: template._id ?? uid('dynamic-info-template'),
  title: template.title ?? '',
  names: (template.names ?? []).map(generateTemplateFormName),
});

/**
 * Convert a dynamic info template form object to an API compatible object.
 *
 * @param {Object} [form={}] - Dynamic info template form object.
 * @returns {{ _id: string, title: string, names: string[] }}
 */
export const formToTemplate = (form = {}) => ({
  _id: form._id ?? uid('dynamic-info-template'),
  title: form.title ?? '',
  names: (form.names ?? []).map(({ value }) => value),
});

/**
 * Convert a dynamic info template to dynamic info information form items.
 *
 * @param {Object} [template={}] - Dynamic info template entity.
 * @returns {Array<{ name: string, value: string }>}
 */
export const templateToDynamicInfoInfos = (template = {}) =>
  (template.names ?? []).map(name => ({ name, value: '' }));
