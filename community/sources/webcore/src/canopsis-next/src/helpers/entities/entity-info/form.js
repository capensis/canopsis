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
 * @property {ContextEntityInfoValuePrimitive} value
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
    value,
  };
};
