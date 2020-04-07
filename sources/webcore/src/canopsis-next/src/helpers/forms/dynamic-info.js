import { unsetSeveralFieldsWithConditions } from '@/helpers/immutable';

/**
 * Convert a dynamic information's object to a dynamic information's form object
 * @param {Object} dynamicInfo
 * @returns {Object}
 */
export function dynamicInfoToForm(dynamicInfo = {}) {
  return {
    general: {
      _id: dynamicInfo._id || '',
      name: dynamicInfo.name || '',
      description: dynamicInfo.description || '',
    },
    infos: dynamicInfo.infos ? [...dynamicInfo.infos] : [],
    patterns: {
      alarm_patterns: dynamicInfo.alarm_patterns || [],
      entity_patterns: dynamicInfo.entity_patterns || [],
    },
  };
}

/**
 * Remove empty "patterns" (alarm_pattern, entity_pattern) fields from dynamic info
 * @param {Object} dynamicInfo
 * @returns {Object}
 */
function removeEmptyPatternsAndIdFromDynamicInfo(dynamicInfo) {
  const idCondition = value => value === '';
  const patternsCondition = value => !value || !value.length;

  return unsetSeveralFieldsWithConditions(dynamicInfo, {
    _id: idCondition,
    alarm_patterns: patternsCondition,
    entity_patterns: patternsCondition,
  });
}

/**
 * Convert a dynamic information's form object to a API compatible dynamic info object
 * @param {Object} form
 * @returns {Object}
 */
export function formToDynamicInfo(form) {
  const dynamicInfo = {
    ...form.general,
    ...form.patterns,
    infos: form.infos,
  };

  return removeEmptyPatternsAndIdFromDynamicInfo(dynamicInfo);
}
