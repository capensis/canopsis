/**
 * @typedef { 'maintenance' | 'pause' } DisableDuringPeriods
 */

/**
 * @typedef {Object} DynamicInfo
 * @property {string} _id
 * @property {string} name
 * @property {string} description
 * @property {DisableDuringPeriods[]} disable_during_periods
 * @property {Infos[]} infos
 * @property {Object[]} alarm_patterns
 * @property {Object[]} entity_patterns
 */

/**
 * @typedef {Object} DynamicInfoPatternsForm
 * @property {Object[]} alarm_patterns
 * @property {Object[]} entity_patterns
 */

/**
 * @typedef {DynamicInfo} DynamicInfoForm
 * @property {DynamicInfoPatternsForm} patterns
 */

/**
 * Convert a dynamic information's object to a dynamic information's form object
 *
 * @param {DynamicInfo} dynamicInfo
 * @returns {DynamicInfoForm}
 */
export const dynamicInfoToForm = (dynamicInfo = {}) => ({
  _id: dynamicInfo._id || '',
  name: dynamicInfo.name || '',
  description: dynamicInfo.description || '',
  disable_during_periods: dynamicInfo.disable_during_periods || [],
  infos: dynamicInfo.infos ? [...dynamicInfo.infos] : [],
  patterns: {
    alarm_patterns: dynamicInfo.alarm_patterns || [],
    entity_patterns: dynamicInfo.entity_patterns || [],
  },
});

/**
 * Convert a dynamic information's form object to a API compatible dynamic info object
 *
 * @param {DynamicInfoForm} form
 * @returns {DynamicInfo}
 */
export const formToDynamicInfo = (form) => {
  const { patterns, ...dynamicInfo } = form;

  return {
    ...dynamicInfo,
    ...patterns,
  };
};
