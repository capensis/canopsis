import { DYNAMIC_INFO_INFORMATION_TYPES } from '@/constants';

import { uid } from '@/helpers/uid';

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
