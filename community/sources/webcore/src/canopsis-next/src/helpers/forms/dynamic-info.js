import { unsetSeveralFieldsWithConditions } from '@/helpers/immutable';

import { getConditionsForRemovingEmptyPatterns } from './shared/patterns';

/**
 * Convert a dynamic information's object to a dynamic information's form object
 *
 * @param {Object} dynamicInfo
 * @returns {Object}
 */
export function dynamicInfoToForm(dynamicInfo = {}) {
  return {
    general: {
      _id: dynamicInfo._id || '',
      name: dynamicInfo.name || '',
      description: dynamicInfo.description || '',
      disable_during_periods: dynamicInfo.disable_during_periods || [],
    },
    infos: dynamicInfo.infos ? [...dynamicInfo.infos] : [],
    patterns: {
      alarm_patterns: dynamicInfo.alarm_patterns || [],
      entity_patterns: dynamicInfo.entity_patterns || [],
    },
  };
}

/**
 * Convert a dynamic information's form object to a API compatible dynamic info object
 *
 * @param {Object} form
 * @returns {Object}
 */
export function formToDynamicInfo(form) {
  const idCondition = value => value === '';

  const dynamicInfo = {
    ...form.general,
    ...form.patterns,

    infos: form.infos,
  };

  return unsetSeveralFieldsWithConditions(dynamicInfo, {
    ...getConditionsForRemovingEmptyPatterns(['alarm_patterns', 'entity_patterns']),

    _id: idCondition,
  });
}
