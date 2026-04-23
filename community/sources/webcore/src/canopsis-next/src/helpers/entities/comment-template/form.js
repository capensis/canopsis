import { cloneDeep } from 'lodash';

import { uid } from '@/helpers/uid';

/**
 * Convert comment template to form
 *
 * @param {Object} [template={}] - Comment template object
 * @returns {Object} Form object
 */
export const commentTemplateToForm = (template = {}) => ({
  name: template.name ?? '',
  fields: (template.fields || []).map(field => ({
    key: uid(),
    name: field.name ?? '',
    required: field.required ?? false,
  })),
});

/**
 * Convert form to comment template
 *
 * @param {Object} form - Form object
 * @returns {Object} Comment template object
 */
export const formToCommentTemplate = (form) => {
  const template = cloneDeep(form);

  template.fields = template.fields.map(({ name, required }) => ({
    name,
    required,
  }));

  return template;
};
