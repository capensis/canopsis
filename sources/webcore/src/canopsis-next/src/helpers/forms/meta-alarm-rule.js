import { cloneDeep, omit, pick, isNumber } from 'lodash';
import moment from 'moment';

import { DEFAULT_TIME_INTERVAL, META_ALARMS_RULE_TYPES, META_ALARMS_THRESHOLD_TYPES } from '@/constants';

import { convertDurationToIntervalObject } from '@/helpers/date';
import { unsetSeveralFieldsWithConditions } from '@/helpers/immutable';

import { getConditionsForRemovingEmptyPatterns } from './shared/patterns';

/**
 * Convert meta alarm rule to form
 *
 * @param {Object} [rule={}]
 * @returns {Object}
 */
export function metaAlarmRuleToForm(rule = {}) {
  const config = rule.config || {};

  return {
    _id: rule._id || '',
    type: rule.type || META_ALARMS_RULE_TYPES.attribute,
    name: rule.name || '',
    config: {
      value_path: config.value_path || '',
      auto_resolve: !!config.auto_resolve,
      alarm_patterns: config.alarm_patterns ? cloneDeep(config.alarm_patterns) : [],
      entity_patterns: config.entity_patterns ? cloneDeep(config.entity_patterns) : [],
      event_patterns: config.event_patterns ? cloneDeep(config.event_patterns) : [],
      threshold_rate: config.threshold_rate || 1,
      threshold_count: config.threshold_count || 1,
      threshold_type: isNumber(config.threshold_count)
        ? META_ALARMS_THRESHOLD_TYPES.thresholdCount
        : META_ALARMS_THRESHOLD_TYPES.thresholdRate,
      time_interval: config.time_interval
        ? convertDurationToIntervalObject(config.time_interval)
        : DEFAULT_TIME_INTERVAL,
    },
  };
}

/**
 * Convert form to meta alarm rule
 *
 * @param {Object} [form={}]
 * @returns {Object}
 */
export function formToMetaAlarmRule(form = {}) {
  const metaAlarmRule = omit(form, ['config']);

  switch (form.type) {
    case META_ALARMS_RULE_TYPES.attribute: {
      const config = pick(form.config, ['alarm_patterns', 'entity_patterns', 'event_patterns']);

      metaAlarmRule.config = unsetSeveralFieldsWithConditions(
        config,
        getConditionsForRemovingEmptyPatterns(),
      );
      break;
    }
    case META_ALARMS_RULE_TYPES.complex:
    case META_ALARMS_RULE_TYPES.valuegroup: {
      const isComplex = form.type === META_ALARMS_RULE_TYPES.complex;

      const thresholdField = form.config.threshold_type === META_ALARMS_THRESHOLD_TYPES.thresholdCount || !isComplex
        ? 'threshold_rate'
        : 'threshold_count';

      const fields = ['threshold_type', thresholdField];

      if (isComplex) {
        fields.push('value_path');
      }

      const config = omit(form.config, fields);

      metaAlarmRule.config = unsetSeveralFieldsWithConditions(
        config,
        getConditionsForRemovingEmptyPatterns(),
      );
      break;
    }
    case META_ALARMS_RULE_TYPES.timebased:
      metaAlarmRule.config = pick(form.config, ['time_interval', 'auto_resolve']);
      break;
  }

  if (metaAlarmRule.config && metaAlarmRule.config.time_interval) {
    const { unit, interval } = metaAlarmRule.config.time_interval;

    metaAlarmRule.config.time_interval = moment.duration(
      interval,
      unit,
    ).asSeconds();
  }

  return metaAlarmRule;
}
