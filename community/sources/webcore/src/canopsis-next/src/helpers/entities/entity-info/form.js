/**
 * @typedef {string | number | boolean | null | Date} ContextEntityInfoValue
 */

/**
 * @typedef {Object} ContextEntityInfo
 * @property {string} name
 * @property {string} description
 * @property {ContextEntityInfoValue[] | ContextEntityInfoValue} value
 */

/**
 * @typedef {ContextEntityInfo} ContextEntityInfoForm
 * @property {ContextEntityInfoValue} value
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
