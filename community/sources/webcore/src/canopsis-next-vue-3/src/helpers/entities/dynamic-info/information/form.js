import { DYNAMIC_INFO_INFORMATION_TYPES } from '@/constants';
import { uid } from '@/helpers/uid';

/**
 * @typedef {Object} DynamicInfoInformation
 * @property {string} type - The type of the dynamic info information
 * @property {string} [name] - The name of the dynamic info information
 * @property {string} [type] - The type of the dynamic info information
 * @property {string} [value] - The value of the dynamic info information
 */

/**
 * @typedef {ObjectKey} DynamicInfoInformationForm
 * @property {string} type - The type field for the form
 * @property {string} name - The name field for the form
 * @property {string} type - The type field for the form
 * @property {string | Array<{value: string} & ObjectKey>} value - The value field for the form
 */

/**
 * Converts a dynamic info information object to form data structure
 *
 * @param {Object} [info={}] - The dynamic info information object
 * @returns {Object} The form data structure with name and value fields
 */
export const dynamicInfoInformationToForm = (info = {}) => ({
  key: info.key ?? uid('dynamic-info-information'),
  name: info.name ?? '',
  type: info.type ?? DYNAMIC_INFO_INFORMATION_TYPES.setToInfo,
  value: info.value ?? '',
});

/**
 * Converts a dynamic info information form object to a dynamic info information object
 *
 * @param {Object} [form={}] - The dynamic info information form object
 * @returns {Object} The dynamic info information object
 */
export const formToDynamicInfoInformation = (form = {}) => {
  const dynamicInfo = { ...form };

  delete dynamicInfo.key;

  return dynamicInfo;
};
