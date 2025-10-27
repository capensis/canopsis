import { DYNAMIC_INFO_INFORMATION_TYPES } from '@/constants';

/**
 * @typedef {Object} DynamicInfoInformation
 * @property {string} [name] - The name of the dynamic info information
 * @property {string} [value] - The value of the dynamic info information
 */

/**
 * @typedef {Object} DynamicInfoInformationForm
 * @property {string} name - The name field for the form
 * @property {string} value - The value field for the form
 */

/**
 * Converts a dynamic info information object to form data structure
 *
 * @param {DynamicInfoInformation} [info={}] - The dynamic info information object
 * @returns {DynamicInfoInformationForm} The form data structure with name and value fields
 */
export const dynamicInfoInformationToForm = (info = {}) => ({
  name: info.name ?? '',
  type: info.type ?? DYNAMIC_INFO_INFORMATION_TYPES.setToInfo,
  value: info.value ?? '',
});
