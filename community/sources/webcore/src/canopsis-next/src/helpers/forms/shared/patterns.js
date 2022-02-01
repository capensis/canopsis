/**
 * Default condition for removing patterns
 *
 * @param {Array|undefined} value
 * @returns {boolean}
 */
export function conditionForRemovingPatterns(value) {
  return !value || !value.length;
}

/**
 * Get conditions for removing empty patterns by unsetSeveralFieldsWithConditions
 *
 * @param {Array} keys
 * @param {Function} condition
 * @returns {Object}
 */
export function getConditionsForRemovingEmptyPatterns(
  keys = ['alarm_patterns', 'entity_patterns', 'event_patterns'],
  condition = conditionForRemovingPatterns,
) {
  return keys.reduce((acc, key) => {
    acc[key] = condition;

    return acc;
  }, {});
}
