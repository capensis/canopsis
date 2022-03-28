import { omit, pick, isNumber } from 'lodash';

import { DEFAULT_TIME_INTERVAL, META_ALARMS_RULE_TYPES, META_ALARMS_THRESHOLD_TYPES, PATTERNS_FIELDS } from '@/constants';

import { durationToForm } from '@/helpers/date/duration';
import { filterPatternsToForm, formFilterToPatterns } from '@/helpers/forms/filter';

import { formToPrimitiveArray, primitiveArrayToForm } from './shared/common';

/**
 * @typedef {
 *   'relation' |
 *   'timebased' |
 *   'attribute' |
 *   'complex' |
 *   'valuegroup' |
 *   'corel'
 * } MetaAlarmRuleType
 */

/**
 * @typedef {FilterPatterns} MetaAlarmRuleAttributeConfig
 */

/**
 * @typedef {Object} MetaAlarmRuleAttributeConfigForm
 * @property {FilterPatternsForm} patterns
 */

/**
 * @typedef {Object} MetaAlarmRuleTimeBasedConfig
 * @property {Object} time_interval
 */

/**
 * @typedef {MetaAlarmRuleAttributeConfig} MetaAlarmRuleTimeBasedConfigForm
 * @property {Interval} time_interval
 */

/**
 * @typedef {MetaAlarmRuleTimeBasedConfig & MetaAlarmRuleAttributeConfig} MetaAlarmRuleComplexConfig
 * @property {number} [threshold_rate]
 * @property {number} [threshold_count]
 */

/**
 * @typedef {
 *   MetaAlarmRuleTimeBasedConfigForm &
 *   MetaAlarmRuleAttributeConfigForm
 * } MetaAlarmRuleComplexConfigForm
 * @property { 'thresholdRate' | 'thresholdCount' } threshold_type
 */

/**
 * @typedef {MetaAlarmRuleComplexConfig} MetaAlarmRuleValueGroupConfig
 * @property {string[]} value_paths
 */

/**
 * @typedef {MetaAlarmRuleComplexConfigForm} MetaAlarmRuleValueGroupConfigForm
 * @property {{ key: string, value: string }[]} value_paths
 */

/**
 * @typedef {MetaAlarmRuleAttributeConfig} MetaAlarmRuleCorelConfig
 * @property {string} corel_id
 * @property {string} corel_status
 * @property {string} corel_parent
 * @property {string} corel_child
 * @property {string} threshold_count
 */

/**
 * @typedef {MetaAlarmRuleCorelConfig} MetaAlarmRuleCorelConfigForm
 */

/**
 * @typedef {
 *   MetaAlarmRuleAttributeConfig &
 *   MetaAlarmRuleTimeBasedConfig &
 *   MetaAlarmRuleComplexConfig &
 *   MetaAlarmRuleValueGroupConfig &
 *   MetaAlarmRuleCorelConfig
 * } MetaAlarmRuleConfig
 */

/**
 * @typedef {
 *   MetaAlarmRuleAttributeConfigForm &
 *   MetaAlarmRuleTimeBasedConfigForm &
 *   MetaAlarmRuleComplexConfigForm &
 *   MetaAlarmRuleValueGroupConfigForm &
 *   MetaAlarmRuleCorelConfigForm
 * } MetaAlarmRuleConfigForm
 */

/**
 * @typedef {Object} MetaAlarmRule
 * @property {string} _id
 * @property {MetaAlarmRuleType} type
 * @property {string} name
 * @property {boolean} auto_resolve
 * @property {string} output_template
 * @property {MetaAlarmRuleConfig} [config]
 */

/**
 * @typedef {MetaAlarmRule} MetaAlarmRuleForm
 * @property {MetaAlarmRuleConfigForm} [config]
 */

/**
 * Convert meta alarm rule to form
 *
 * @param {MetaAlarmRule} [rule={}]
 * @returns {MetaAlarmRuleForm}
 */
export const metaAlarmRuleToForm = (rule = {}) => {
  const config = rule.config ?? {};

  return {
    _id: rule._id ?? '',
    type: rule.type ?? META_ALARMS_RULE_TYPES.attribute,
    name: rule.name ?? '',
    auto_resolve: !!rule.auto_resolve,
    output_template: rule.output_template ?? '{{ .Children.Alarm.Value.State.Message }}',
    config: {
      /** TODO: Do we need to use total entity pattern and event pattern ? */
      patterns: filterPatternsToForm(config, [PATTERNS_FIELDS.entity, PATTERNS_FIELDS.alarm, PATTERNS_FIELDS.event]),
      value_paths: config.value_paths ? primitiveArrayToForm(config.value_paths) : [],
      threshold_rate: config.threshold_rate ? config.threshold_rate * 100 : 100,
      threshold_count: config.threshold_count ?? 1,
      corel_id: config.corel_id ?? '',
      corel_status: config.corel_status ?? '',
      corel_parent: config.corel_parent ?? '',
      corel_child: config.corel_child ?? '',
      threshold_type: isNumber(config.threshold_count)
        ? META_ALARMS_THRESHOLD_TYPES.thresholdCount
        : META_ALARMS_THRESHOLD_TYPES.thresholdRate,
      time_interval: durationToForm(config.time_interval ?? DEFAULT_TIME_INTERVAL),
    },
  };
};

/**
 * Convert form to meta alarm rule
 *
 * @param {MetaAlarmRuleForm} [form={}]
 * @returns {MetaAlarmRule}
 */
export const formToMetaAlarmRule = (form = {}) => {
  const metaAlarmRule = omit(form, ['config']);

  switch (form.type) {
    case META_ALARMS_RULE_TYPES.attribute: {
      metaAlarmRule.config = formFilterToPatterns(
        form.config.patterns,
        [PATTERNS_FIELDS.alarm, PATTERNS_FIELDS.event, PATTERNS_FIELDS.entity],
      );
      break;
    }
    case META_ALARMS_RULE_TYPES.corel:
    case META_ALARMS_RULE_TYPES.complex:
    case META_ALARMS_RULE_TYPES.valuegroup: {
      const isComplex = form.type === META_ALARMS_RULE_TYPES.complex;
      const isValueGroup = form.type === META_ALARMS_RULE_TYPES.valuegroup;
      const isCorel = form.type === META_ALARMS_RULE_TYPES.corel;

      const thresholdField = isCorel || form.config.threshold_type === META_ALARMS_THRESHOLD_TYPES.thresholdCount
        ? 'threshold_rate'
        : 'threshold_count';

      const fields = ['threshold_type', thresholdField];

      if (isComplex || isCorel) {
        fields.push('value_paths');
      }

      const config = omit(form.config, fields);

      if (isValueGroup) {
        config.value_paths = formToPrimitiveArray(config.value_paths);
      }

      if (config.threshold_rate) {
        config.threshold_rate /= 100;
      }

      metaAlarmRule.config = {
        ...metaAlarmRule.config,
        ...formFilterToPatterns(config.patterns),
      };
      break;
    }
    case META_ALARMS_RULE_TYPES.timebased:
      metaAlarmRule.config = pick(form.config, ['time_interval']);
      break;
  }

  return metaAlarmRule;
};
