import { omit, pick, isNumber } from 'lodash';

import {
  DEFAULT_TIME_INTERVAL,
  META_ALARMS_RULE_TYPES,
  META_ALARMS_THRESHOLD_TYPES,
  PATTERNS_FIELDS,
} from '@/constants';

import { durationToForm } from '@/helpers/date/duration';
import { filterPatternsToForm, formFilterToPatterns } from '@/helpers/entities/filter/form';
import { formToPrimitiveArray, primitiveArrayToForm } from '@/helpers/entities/shared/form';
import { addKeyInEntity, removeKeyFromEntities } from '@/helpers/array';

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
 * @typedef {Object} MetaAlarmRuleTemplatesConfig
 * @property {string} component_template
 * @property {string} resource_template
 */

/**
 * @typedef {MetaAlarmRuleTemplatesConfig} MetaAlarmRuleAttributeConfig
 */

/**
 * @typedef {MetaAlarmRuleAttributeConfig} MetaAlarmRuleAttributeConfigForm
 */

/**
 * @typedef {MetaAlarmRuleTemplatesConfig} MetaAlarmRuleTimeBasedConfig
 * @property {Object} time_interval
 */

/**
 * @typedef {MetaAlarmRuleAttributeConfig} MetaAlarmRuleTimeBasedConfigForm
 * @property {Interval} time_interval
 * @property {Interval} [child_inactive_delay]
 */

/**
 * @typedef {MetaAlarmRuleTimeBasedConfig & MetaAlarmRuleAttributeConfig} MetaAlarmRuleComplexConfig
 * @property {number} [threshold_rate]
 * @property {number} [threshold_count]
 */

/**
 * @typedef { MetaAlarmRuleTimeBasedConfigForm & MetaAlarmRuleAttributeConfigForm } MetaAlarmRuleComplexConfigForm
 * @property { 'thresholdRate' | 'thresholdCount' } threshold_type
 */

/**
 * @typedef {MetaAlarmRuleComplexConfig & MetaAlarmRuleTemplatesConfig} MetaAlarmRuleValueGroupConfig
 * @property {string[]} value_paths
 */

/**
 * @typedef {MetaAlarmRuleComplexConfigForm & MetaAlarmRuleTemplatesConfig} MetaAlarmRuleValueGroupConfigForm
 * @property {{ key: string, value: string }[]} value_paths
 */

/**
 * @typedef {Object} MetaAlarmRuleCorelConfig
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
 * @typedef {Object} MetaAlarmRuleTags
 * @property {boolean} copy_from_children
 * @property {string[]} filter_by_label
 */

/**
 * @typedef {MetaAlarmRuleTags} MetaAlarmRuleTagsForm
 */

/**
 * @typedef {Object} MetaAlarmRuleInfosItem
 * @property {string} name
 * @property {string | boolean | number | null} value
 * @property {string} description
 * @property {boolean} copy_from_children
 */

/**
 * @typedef {FilterPatterns} MetaAlarmRule
 * @property {string} _id
 * @property {MetaAlarmRuleType} type
 * @property {string} name
 * @property {boolean} auto_resolve
 * @property {string} output_template
 * @property {MetaAlarmRuleTags} tags
 * @property {MetaAlarmRuleInfosItem[]} infos
 * @property {MetaAlarmRuleConfig} [config]
 */

/**
 * @typedef {MetaAlarmRuleInfosItem} MetaAlarmRuleInfosItemForm
 * @property {string} key
 */

/**
 * @typedef {MetaAlarmRule} MetaAlarmRuleForm
 * @property {FilterPatternsForm} patterns
 * @property {MetaAlarmRuleInfosItemForm[]} infos
 * @property {MetaAlarmRuleConfigForm} [config]
 */

/**
 * Check meta alarm type is attribute
 *
 * @param {MetaAlarmRuleType} type
 * @return {boolean}
 */
export const isAttributeMetaAlarmRuleType = type => type === META_ALARMS_RULE_TYPES.attribute;

/**
 * Check meta alarm type is timebased
 *
 * @param {MetaAlarmRuleType} type
 * @return {boolean}
 */
export const isTimebasedMetaAlarmRuleType = type => type === META_ALARMS_RULE_TYPES.timebased;

/**
 * Check meta alarm type is complex
 *
 * @param {MetaAlarmRuleType} type
 * @return {boolean}
 */
export const isComplexMetaAlarmRuleType = type => type === META_ALARMS_RULE_TYPES.complex;

/**
 * Check meta alarm type is valuegroup
 *
 * @param {MetaAlarmRuleType} type
 * @return {boolean}
 */
export const isValueGroupMetaAlarmRuleType = type => type === META_ALARMS_RULE_TYPES.valuegroup;

/**
 * Check meta alarm type is corel
 *
 * @param {MetaAlarmRuleType} type
 * @return {boolean}
 */
export const isCorelMetaAlarmRuleType = type => type === META_ALARMS_RULE_TYPES.corel;

/**
 * Check meta alarm type is manualgroup
 *
 * @param {MetaAlarmRuleType} type
 * @return {boolean}
 */
export const isManualGroupMetaAlarmRuleType = type => type === META_ALARMS_RULE_TYPES.manualgroup;

/**
 * Check meta alarm type is auto
 *
 * @param {MetaAlarmRuleType} type
 * @return {boolean}
 */
export const isAutoMetaAlarmRuleType = type => type && type !== META_ALARMS_RULE_TYPES.manualgroup;

/**
 * Check meta alarm type has a patterns
 *
 * @param {MetaAlarmRuleType} type
 * @return {boolean}
 */
export const isMetaAlarmRuleTypeHasPatterns = type => isAttributeMetaAlarmRuleType(type)
  || isComplexMetaAlarmRuleType(type)
  || isValueGroupMetaAlarmRuleType(type)
  || isCorelMetaAlarmRuleType(type);

/**
 * Check meta alarm type has a total entity patterns
 *
 * @param {MetaAlarmRuleType} type
 * @return {boolean}
 */
export const isMetaAlarmRuleTypeHasTotalEntityPatterns = type => isComplexMetaAlarmRuleType(type)
  || isValueGroupMetaAlarmRuleType(type);

/**
 * Convert meta alarm rule to patterns
 *
 * @param {MetaAlarmRule} rule
 * @return {FilterPatterns}
 */
export const metaAlarmFilterPatternsToForm = rule => filterPatternsToForm(
  rule,
  [PATTERNS_FIELDS.alarm, PATTERNS_FIELDS.event, PATTERNS_FIELDS.entity, PATTERNS_FIELDS.totalEntity],
);

/**
 * Convert meta alarm rule tags to form
 *
 * @param {MetaAlarmRuleTags} [tags = {}]
 * @returns {MetaAlarmRuleTagsForm}
 */
export const metaAlarmRuleTagsToForm = (tags = {}) => ({
  copy_from_children: !!tags.copy_from_children,
  filter_by_label: tags.filter_by_label ?? [],
});

/**
 * Convert meta alarm rule infos item to form
 *
 * @param {MetaAlarmRuleInfosItem} [item = {}]
 * @returns {MetaAlarmRuleInfosItemForm}
 */
export const metaAlarmRuleInfosItemToForm = (item = {}) => addKeyInEntity({
  copy_from_children: !!item.copy_from_children,
  name: item.name ?? '',
  value: item.value ?? '',
  description: item.description ?? '',
});

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
    output_template: rule.output_template ?? '{{ .LastChild.Alarm.Value.State.Message }}',
    patterns: metaAlarmFilterPatternsToForm(rule),
    tags: metaAlarmRuleTagsToForm(rule.tags || {}),
    infos: (rule.infos || []).map(metaAlarmRuleInfosItemToForm),
    config: {
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
      child_inactive_delay: config.child_inactive_delay ? durationToForm(config.child_inactive_delay) : undefined,
      component_template: config.component_template ?? '',
      resource_template: config.resource_template ?? '',
    },
  };
};

/**
 * Checks if the given alarm rule type has template configuration fields.
 *
 * @param {string} type - The type of the alarm rule to check.
 * @returns {boolean} Returns `true` if the `type` has template configuration fields, otherwise `false`.
 */
export const hasTemplateConfigFields = type => [
  META_ALARMS_RULE_TYPES.timebased,
  META_ALARMS_RULE_TYPES.attribute,
  META_ALARMS_RULE_TYPES.complex,
  META_ALARMS_RULE_TYPES.valuegroup,
].includes(type);

/**
 * Convert form to meta alarm rule tags
 *
 * @param {MetaAlarmRuleTagsForm} form
 * @returns {MetaAlarmRuleTags}
 */
export const formToMetaAlarmRuleTags = (form = {}) => {
  const result = pick(form, ['copy_from_children']);

  if (result.copy_from_children && form.filter_by_label?.length) {
    result.filter_by_label = form.filter_by_label;
  }

  return result;
};

/**
 * Convert form to meta alarm rule
 *
 * @param {MetaAlarmRuleForm} [form={}]
 * @returns {MetaAlarmRule}
 */
export const formToMetaAlarmRule = (form = {}) => {
  const patternsFields = [
    PATTERNS_FIELDS.alarm,
    PATTERNS_FIELDS.entity,
  ];

  if ([META_ALARMS_RULE_TYPES.complex, META_ALARMS_RULE_TYPES.valuegroup].includes(form.type)) {
    patternsFields.push(PATTERNS_FIELDS.totalEntity);
  }

  const metaAlarmRule = {
    ...omit(form, ['config', 'patterns', 'infos', 'tags']),
    ...formFilterToPatterns(form.patterns, patternsFields),

    infos: removeKeyFromEntities(form.infos),
    tags: formToMetaAlarmRuleTags(form.tags),
  };

  switch (form.type) {
    case META_ALARMS_RULE_TYPES.complex:
    case META_ALARMS_RULE_TYPES.valuegroup:
    case META_ALARMS_RULE_TYPES.corel: {
      const isComplex = isComplexMetaAlarmRuleType(form.type);
      const isValueGroup = isValueGroupMetaAlarmRuleType(form.type);
      const isCorel = isCorelMetaAlarmRuleType(form.type);

      const thresholdField = isCorel || form.config.threshold_type === META_ALARMS_THRESHOLD_TYPES.thresholdCount
        ? 'threshold_rate'
        : 'threshold_count';

      const fields = ['threshold_type', 'component_template', 'resource_template', thresholdField];

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

      metaAlarmRule.config = config;
      break;
    }
    case META_ALARMS_RULE_TYPES.timebased:
      metaAlarmRule.config = pick(form.config, ['time_interval']);
      break;
  }

  if (hasTemplateConfigFields(metaAlarmRule.type)) {
    metaAlarmRule.config = {
      ...metaAlarmRule.config,
      ...pick(form.config, ['component_template', 'resource_template']),
    };
  }

  return metaAlarmRule;
};
