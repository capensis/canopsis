import { isArray } from 'lodash';

import { DYNAMIC_INFO_INFORMATION_TYPES } from '@/constants';

import { primitiveArrayToForm, formToPrimitiveArray } from '@/helpers/entities/shared/form';

/**
 * @typedef {Object} DynamicInfoInformation
 * @property {string} [name] - The name of the dynamic info information
 * @property {string} [type] - The type of the dynamic info information
 * @property {string} [value] - The value of the dynamic info information
 */

/**
 * @typedef {Object} DynamicInfoInformationForm
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
export const dynamicInfoInformationToForm = (info = {}) => ({
  name: info.name ?? '',
  type: info.type ?? DYNAMIC_INFO_INFORMATION_TYPES.setToInfo,
  value: isArray(info.value) ? primitiveArrayToForm(info.value) : info.value ?? '',
});

/**
 * Converts a dynamic info information form object to a dynamic info information object
 *
 * @param {DynamicInfoInformationForm} form - The dynamic info information form object
 * @returns {DynamicInfoInformation} The dynamic info information object
 */
export const formToDynamicInfoInformation = (form = {}) => ({
  ...form,

  value: isArray(form.value) ? formToPrimitiveArray(form.value) : form.value,
});
