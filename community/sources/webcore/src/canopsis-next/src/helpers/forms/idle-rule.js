import { omit, pick } from 'lodash';

import { IDLE_RULE_ALARM_CONDITIONS, IDLE_RULE_TYPES, PATTERNS_FIELDS, TIME_UNITS } from '@/constants';

import { enabledToForm } from '@/helpers/forms/shared/common';
import { durationToForm } from '@/helpers/date/duration';
import { formToAction, actionToForm } from '@/helpers/forms/action';
import { filterPatternsToForm, formFilterToPatterns } from '@/helpers/forms/filter';

/**
 * @typedef { 'alarm' | 'entity' } IdleRuleType
 */

/**
 * @typedef { 'last_event' | 'last_update' } IdleRuleAlarmCondition
 */

/**
 * @typedef {FilterPatterns} IdleRule
 * @property {boolean} enabled
 * @property {string} name
 * @property {string} description
 * @property {IdleRuleType} type
 * @property {Duration} duration
 * @property {number} priority
 * @property {DisableDuringPeriods[]} disable_during_periods
 * @property {Object[]} entity_patterns
 * @property {Object[]} alarm_patterns
 * @property {IdleRuleAlarmCondition} alarm_condition
 * @property {Action} operation
 */

/**
 * @typedef {IdleRule & FilterPatternsForm} IdleRuleForm
 * @property {ActionForm} operation
 */

/**
 * Convert idle rule object to form compatible object
 *
 * @param {IdleRule} [idleRule = {}]
 * @return {IdleRuleForm}
 */
export const idleRuleToForm = (idleRule = {}) => ({
  enabled: enabledToForm(idleRule.enabled),
  name: idleRule.name ?? '',
  description: idleRule.description ?? '',
  type: idleRule.type ?? IDLE_RULE_TYPES.alarm,
  duration: idleRule.duration
    ? durationToForm(idleRule.duration)
    : { value: 1, unit: TIME_UNITS.minute },
  priority: idleRule.priority ?? 1,
  disable_during_periods: idleRule.disable_during_periods ?? [],
  alarm_condition: idleRule.alarm_condition ?? IDLE_RULE_ALARM_CONDITIONS.lastEvent,
  operation: pick(actionToForm(idleRule.operation), ['type', 'parameters']),

  ...filterPatternsToForm(idleRule, [PATTERNS_FIELDS.entity, PATTERNS_FIELDS.alarm]),
});

/**
 * Convert form object to idle API compatible object
 *
 * @param {IdleRuleForm} form
 * @return {IdleRule}
 */
export const formToIdleRule = (form) => {
  const isEntityType = form.type === IDLE_RULE_TYPES.entity;
  const idleRule = omit(form, [
    'alarm_condition',
    'operation',
    PATTERNS_FIELDS.entity,
    PATTERNS_FIELDS.alarm,
  ]);

  if (!isEntityType) {
    idleRule.alarm_condition = form.alarm_condition;
    idleRule.operation = pick(formToAction(form.operation), ['type', 'parameters']);
  }

  return {
    ...idleRule,
    ...formFilterToPatterns(
      form,
      isEntityType ? [PATTERNS_FIELDS.entity] : [PATTERNS_FIELDS.entity, PATTERNS_FIELDS.alarm],
    ),
  };
};
