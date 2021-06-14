import { isArray } from 'lodash';

import uid from '../uid';

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
 * @property {Array<{ key: string, value: ContextEntityInfoValuePrimitive }> | ContextEntityInfoValuePrimitive} value
 */

/**
 * Convert entity info object to form
 *
 * @param {ContextEntityInfo} entityInfo
 * @returns {ContextEntityInfoForm}
 */
export const entityInfoToForm = (entityInfo = {}) => {
  const { name = '', description = '', value = '' } = entityInfo;

  return {
    name,
    description,
    value: isArray(value) ? value.map(v => ({ key: uid(), value: v })) : value,
  };
};

/**
 * Convert form to entity info object
 *
 * @param {ContextEntityInfoForm} form
 * @returns {ContextEntityInfo}
 */
export const formToEntityInfo = (form) => {
  const { name = '', description = '', value = '' } = form;

  return {
    name,
    description,
    value: isArray(value) ? value.map(item => item.value) : value,
  };
};
