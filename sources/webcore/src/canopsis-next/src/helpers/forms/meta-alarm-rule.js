import { cloneDeep } from 'lodash';
import { META_ALARMS_RULE_TYPES } from '@/constants';

/**
 * Convert meta alarm filter rule to form
 *
 * @param {Object} [rule={}]
 * @returns {Object}
 */
export function metaAlarmRuleToForm(rule = {}) {
  return {
    _id: rule._id || '',
    type: rule.type || META_ALARMS_RULE_TYPES.relation,
    description: rule.description || '',
    alarm_patterns: rule.alarm_patterns ? cloneDeep(rule.alarm_patterns) : {},
  };
}
