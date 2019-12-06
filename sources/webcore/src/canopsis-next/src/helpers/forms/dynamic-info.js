import { unsetSeveralFieldsWithConditions } from '@/helpers/immutable';

export function dynamicInfoToForm(dynamicInfo = {}) {
  return {
    general: {
      _id: dynamicInfo._id || '',
      name: dynamicInfo.name || '',
      description: dynamicInfo.description || '',
    },
    infos: dynamicInfo.infos || [],
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

export function formToDynamicInfo(form) {
  let dynamicInfo = {
    ...form.general,
    ...form.patterns,
    infos: form.infos,
  };

  dynamicInfo = removeEmptyPatternsAndIdFromDynamicInfo(dynamicInfo);

  return dynamicInfo;
}
