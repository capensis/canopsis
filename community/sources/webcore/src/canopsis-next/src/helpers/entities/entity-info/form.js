import { isArray } from 'lodash';

import { primitiveArrayToForm, formToPrimitiveArray } from '@/helpers/entities/shared/form';

/**
 * @typedef {string | number | boolean | null} ContextEntityInfoValuePrimitive
 */

/**
 * @typedef {Object} ContextEntityInfo
 * @property {string} name
 * @property {string} description
 * @property {ContextEntityInfoValuePrimitive[] | ContextEntityInfoValuePrimitive} value
 */

/**
 * @typedef {ContextEntityInfo} ContextEntityInfoForm
 * @property {ContextEntityInfoValuePrimitive | Array<{value: string} & ObjectKey>} value
 */

/**
 * Convert entity info object to form
 *
 * @param {ContextEntityInfo} entityInfo
 * @returns {ContextEntityInfoForm}
 */
export const entityInfoToForm = (entityInfo = {}) => {
  const { name = '', description = '', value = '' } = entityInfo;

  const formValue = isArray(value) ? primitiveArrayToForm(value) : value;

  return {
    name,
    description,
    value: formValue,
  };
};

/**
 * Convert entity info form to entity info object
 *
 * @param {ContextEntityInfoForm} form
 * @returns {ContextEntityInfo}
 */
export const formToEntityInfo = (form = {}) => {
  const { name = '', description = '', value = '' } = form;

  const entityValue = isArray(value) ? formToPrimitiveArray(value) : value;

  return {
    name,
    description,
    value: entityValue,
  };
};
