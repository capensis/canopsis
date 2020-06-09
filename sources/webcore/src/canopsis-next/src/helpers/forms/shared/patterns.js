export function conditionForRemovingPattern(value) {
  return !value || !value.length;
}


export function getConditionsForRemovingEmptyPatterns(keys = ['alarm_patterns', 'entity_patterns', 'event_patterns']) {
  return keys.reduce((acc, key) => {
    acc[key] = conditionForRemovingPattern;

    return acc;
  }, {});
}
