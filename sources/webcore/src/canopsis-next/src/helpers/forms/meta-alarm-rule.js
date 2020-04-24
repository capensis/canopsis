import { cloneDeep, omit, pick } from 'lodash';
import { META_ALARMS_RULE_TYPES } from '@/constants';

/**
 * Convert meta alarm rule to form
 *
 * @param {Object} [rule={}]
 * @returns {Object}
 */
export function metaAlarmRuleToForm(rule = {}) {
  return {
    _id: rule._id || '',
    type: rule.type || META_ALARMS_RULE_TYPES.attribute,
    name: rule.name || '',
    config: {
      alarm_patterns: rule.alarm_patterns ? cloneDeep(rule.alarm_patterns) : [],
      entity_patterns: rule.entity_patterns ? cloneDeep(rule.entity_patterns) : [],
      event_patterns: rule.event_patterns ? cloneDeep(rule.event_patterns) : [],
      threshold_rate: rule.threshold_rate || 1,
      threshold_count: rule.threshold_count || 1,
      time_interval: rule.time_interval || 1,
    },
  };
}

/**
 * Convert form to meta alarm rul
 *
 * @param {Object} [form={}]
 * @returns {Object}
 */
export function formToMetaAlarmRule(form = {}) {
  const metaAlarmRule = omit(form, ['config']);

  metaAlarmRule.config = {
    [META_ALARMS_RULE_TYPES.attribute]: pick(form.config, ['alarm_patterns', 'entity_patterns', 'event_patterns']),
    [META_ALARMS_RULE_TYPES.complex]: pick(form.config, ['threshold_rate', 'threshold_count']),
    [META_ALARMS_RULE_TYPES.relation]: null,
    [META_ALARMS_RULE_TYPES.timebased]: pick(form.config, ['time_interval']),
  }[form.type];

  return metaAlarmRule;
}
