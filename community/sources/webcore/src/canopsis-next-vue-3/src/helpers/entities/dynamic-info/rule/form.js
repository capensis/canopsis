import { DYNAMIC_INFO_INFORMATION_TYPES } from '@/constants';

import {
  dynamicInfoInformationToForm,
  formToDynamicInfoInformation,
} from '@/helpers/entities/dynamic-info/information/form';

/**
 * @typedef { 'maintenance' | 'pause' } DisableDuringPeriods
 */

/**
 * @typedef {FilterPatterns} DynamicInfo
 * @property {string} _id
 * @property {boolean} enabled
 * @property {string} name
 * @property {string} description
 * @property {DisableDuringPeriods[]} disable_during_periods
 * @property {Infos[]} infos
 */

/**
 * @typedef {DynamicInfo} DynamicInfoForm
 * @property {FilterPatternsForm} patterns
 */

/**
 * Convert a dynamic information's object to a dynamic information's form object
 *
 * @param {Object} [dynamicInfo={}]
 * @returns {Object}
 */
export const dynamicInfoToForm = (dynamicInfo = {}) => {
  const infos = (
    dynamicInfo.infos?.length
      ? dynamicInfo.infos
      : [{ type: DYNAMIC_INFO_INFORMATION_TYPES.setToInfo }]
  ).map(dynamicInfoInformationToForm);

  return {
    infos,

    _id: dynamicInfo._id ?? '',
    name: dynamicInfo.name ?? '',
    enabled: dynamicInfo.enabled ?? true,
    description: dynamicInfo.description ?? '',
    disable_during_periods: dynamicInfo.disable_during_periods ?? [],
    alarm_pattern: dynamicInfo.alarm_pattern ?? [],
    entity_pattern: dynamicInfo.entity_pattern ?? [],
  };
};

/**
 * Convert a dynamic information's form object to a API compatible dynamic info object
 *
 * @param {Object} form
 * @returns {Object}
 */
export const formToDynamicInfo = (form = {}) => ({
  ...form,
  infos: (form.infos ?? []).map(formToDynamicInfoInformation),
});
