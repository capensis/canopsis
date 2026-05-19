import { DYNAMIC_INFO_INFORMATION_TYPES } from '@/constants';

import { addKeyInEntity, removeKeyFromEntity } from '@/helpers/array';

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
 * @param {DynamicInfoInformation} [info={}] - The dynamic info information object
 * @returns {DynamicInfoInformationForm} The form data structure with name and value fields
 */
export const dynamicInfoInformationToForm = (info = {}) => addKeyInEntity({
  name: info.name ?? '',
  type: info.type ?? DYNAMIC_INFO_INFORMATION_TYPES.setToInfo,
  value: info.value ?? '',
});

/**
 * Converts a dynamic info information form object to a dynamic info information object
 *
 * @param {DynamicInfoInformationForm} [form={}] - The dynamic info information form object
 * @returns {DynamicInfoInformation} The dynamic info information object
 */
export const formToDynamicInfoInformation = (form = {}) => removeKeyFromEntity(form);
