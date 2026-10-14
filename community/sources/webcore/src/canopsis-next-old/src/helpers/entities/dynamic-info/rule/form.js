import { DYNAMIC_INFO_INFORMATION_TYPES, PATTERNS_FIELDS } from '@/constants';

import { filterPatternsToForm, formFilterToPatterns } from '@/helpers/entities/filter/form';

import { dynamicInfoInformationToForm, formToDynamicInfoInformation } from '../information/form';

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
 * @param {DynamicInfo} dynamicInfo
 * @returns {DynamicInfoForm}
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
    patterns: filterPatternsToForm(dynamicInfo, [PATTERNS_FIELDS.alarm, PATTERNS_FIELDS.entity]),
  };
};

/**
 * Convert a dynamic information's form object to a API compatible dynamic info object
 *
 * @param {DynamicInfoForm} form
 * @returns {DynamicInfo}
 */
export const formToDynamicInfo = (form) => {
  const { patterns, infos = [], ...dynamicInfo } = form;

  return {
    ...dynamicInfo,
    infos: infos.map(formToDynamicInfoInformation),
    ...formFilterToPatterns(patterns, [PATTERNS_FIELDS.alarm, PATTERNS_FIELDS.entity]),
  };
};
