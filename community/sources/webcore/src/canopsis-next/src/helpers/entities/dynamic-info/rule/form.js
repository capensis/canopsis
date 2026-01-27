import { isArray } from 'lodash';

import { PATTERNS_FIELDS } from '@/constants';

import { filterPatternsToForm, formFilterToPatterns } from '@/helpers/entities/filter/form';
import { primitiveArrayToForm, formToPrimitiveArray } from '@/helpers/entities/shared/form';

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
  const infos = (dynamicInfo.infos || []).map((info) => {
    const value = isArray(info.value) && (!info.value.length || !info.value[0]?.key)
      ? primitiveArrayToForm(info.value)
      : info.value;

    return {
      ...info,
      value,
    };
  });

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
  const { patterns, infos, ...dynamicInfo } = form;

  const convertedInfos = (infos || []).map((info) => {
    const value = isArray(info.value) && info.value.length > 0 && info.value[0]?.key
      ? formToPrimitiveArray(info.value)
      : info.value;

    return {
      ...info,
      value,
    };
  });

  return {
    ...dynamicInfo,
    infos: convertedInfos,
    ...formFilterToPatterns(patterns, [PATTERNS_FIELDS.alarm, PATTERNS_FIELDS.entity]),
  };
};
