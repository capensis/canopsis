import { omit, pick } from 'lodash';

import {
  ACTION_TYPES,
  IDLE_RULE_ALARM_CONDITIONS,
  IDLE_RULE_TYPES,
  OLD_PATTERNS_FIELDS,
  PATTERNS_FIELDS,
  TIME_UNITS,
} from '@/constants';

import { enabledToForm } from '@/helpers/forms/shared/common';
import { durationToForm } from '@/helpers/date/duration';
import { formToAction, actionParametersToForm } from '@/helpers/forms/action';
import { filterPatternsToForm, formFilterToPatterns } from '@/helpers/forms/filter';
import { isAssociateTicketActionType } from '@/helpers/entities/action';

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
 * @property {string} comment
 * @property {string} description
 * @property {IdleRuleType} type
 * @property {Duration} duration
 * @property {number} priority
 * @property {DisableDuringPeriods[]} disable_during_periods
 * @property {IdleRuleAlarmCondition} alarm_condition
 * @property {Action} operation
 */

/**
 * @typedef {IdleRule} IdleRuleForm
 * @property {ActionForm} operation
 * @property {FilterPatternsForm} patterns
 */

/**
 * Check is idle rule entity type
 *
 * @param {IdleRuleType} type
 * @returns {boolean}
 */
export const isIdleRuleEntityType = type => type === IDLE_RULE_TYPES.entity;

/**
 * Convert idle rule object to form compatible object
 *
 * @param {IdleRule} [idleRule = {}]
 * @return {IdleRuleForm}
 */
export const idleRuleToForm = (idleRule = {}) => {
  const type = idleRule.operation?.type ?? ACTION_TYPES.snooze;
  const parameters = actionParametersToForm({
    type,
    parameters: idleRule.operation?.parameters,
  });

  return {
    enabled: enabledToForm(idleRule.enabled),
    name: idleRule.name ?? '',
    description: idleRule.description ?? '',
    comment: idleRule.comment ?? '',
    type: idleRule.type ?? IDLE_RULE_TYPES.alarm,
    duration: idleRule.duration
      ? durationToForm(idleRule.duration)
      : { value: 1, unit: TIME_UNITS.minute },
    priority: idleRule.priority ?? 1,
    disable_during_periods: idleRule.disable_during_periods ?? [],
    alarm_condition: idleRule.alarm_condition ?? IDLE_RULE_ALARM_CONDITIONS.lastEvent,
    operation: {
      type,
      parameters,
    },
    patterns: filterPatternsToForm(
      idleRule,
      [PATTERNS_FIELDS.entity, PATTERNS_FIELDS.alarm],
      [OLD_PATTERNS_FIELDS.entity, OLD_PATTERNS_FIELDS.alarm],
    ),
  };
};

/**
 * Convert form object to idle API compatible object
 *
 * @param {IdleRuleForm} form
 * @return {IdleRule}
 */
export const formToIdleRule = (form) => {
  const isEntityType = form.type === IDLE_RULE_TYPES.entity;
  const idleRule = omit(form, ['alarm_condition', 'operation', 'patterns', 'comment']);

  if (!isEntityType) {
    idleRule.alarm_condition = form.alarm_condition;
    idleRule.operation = pick(formToAction(form.operation), ['type', 'parameters']);

    if (isAssociateTicketActionType(idleRule.operation.type)) {
      idleRule.comment = form.comment;
    }
  }

  return {
    ...idleRule,
    ...formFilterToPatterns(
      form.patterns,
      isEntityType ? [PATTERNS_FIELDS.entity] : [PATTERNS_FIELDS.entity, PATTERNS_FIELDS.alarm],
    ),
  };
};
