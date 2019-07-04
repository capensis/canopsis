export function checkIfRuleIsEmpty(rule) {
  return rule.field === '' || rule.operator === '';
}

export function checkIfGroupIsEmpty(group) {
  const hasValidRule = Object.values(group.rules).some(rule => !checkIfRuleIsEmpty(rule));

  if (hasValidRule) {
    return false;
  }

  const hasValidSubGroup = Object.values(group.groups).some(subGroup => !checkIfGroupIsEmpty(subGroup));

  return !hasValidSubGroup;
}
