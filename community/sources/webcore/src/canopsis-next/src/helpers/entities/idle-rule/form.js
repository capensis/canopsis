import { omit, pick } from 'lodash';

import {
  ACTION_TYPES,
  IDLE_RULE_ALARM_CONDITIONS,
  IDLE_RULE_TYPES,
  MIN_DURATION_VALUE,
  MIN_PRIORITY,
  PATTERNS_FIELDS,
  TIME_UNITS,
} from '@/constants';

import { durationToForm } from '@/helpers/date/duration';
import {
  actionParametersToForm,
  formToAction,
  isAssociateTicketActionType,
} from '@/helpers/entities/action';
import { filterPatternsToForm, formFilterToPatterns } from '@/helpers/entities/filter/form';

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
 * Check is idle rule entity type.
 *
 * @param {string} type
 * @returns {boolean}
 */
export const isIdleRuleEntityType = type => type === IDLE_RULE_TYPES.entity;

const DEFAULT_IDLE_RULE_SNOOZE_DURATION = { value: MIN_DURATION_VALUE, unit: TIME_UNITS.second };
const DEFAULT_IDLE_RULE_PBEHAVIOR_DURATION = { value: MIN_DURATION_VALUE, unit: TIME_UNITS.minute };

/**
 * Prepare finite numeric form values with a minimum while preserving empty invalid form states.
 *
 * @param {*} value
 * @param {number} min
 * @returns {*}
 */
const prepareMinNumber = (value, min) => {
  if (typeof value !== 'number') {
    return value;
  }

  return Math.max(value, min);
};

/**
 * Prepare duration fields and fill missing values from a default duration.
 *
 * @param {Duration} duration
 * @param {Duration} defaultDuration
 * @returns {Duration}
 */
const prepareDuration = (duration, defaultDuration) => ({
  ...defaultDuration,
  ...duration,
  value: prepareMinNumber(duration?.value ?? defaultDuration.value, MIN_DURATION_VALUE),
  unit: duration?.unit ?? defaultDuration.unit,
});

/**
 * Convert idle rule operation parameters to form parameters with idle-rule specific defaults.
 *
 * @param {Action} operation
 * @returns {Object}
 */
const idleRuleActionParametersToForm = (operation = {}) => {
  const type = operation.type ?? ACTION_TYPES.snooze;
  const parameters = actionParametersToForm({
    type,
    parameters: operation.parameters,
  });
  const hasExistingSnoozeDuration =
    type === ACTION_TYPES.snooze && operation.parameters?.duration;
  const hasExistingPbehaviorDuration =
    type === ACTION_TYPES.pbehavior && operation.parameters?.duration;

  if (!hasExistingSnoozeDuration) {
    parameters[ACTION_TYPES.snooze] = {
      ...parameters[ACTION_TYPES.snooze],
      duration: prepareDuration(null, DEFAULT_IDLE_RULE_SNOOZE_DURATION),
    };
  }

  if (!hasExistingPbehaviorDuration) {
    parameters[ACTION_TYPES.pbehavior] = {
      ...parameters[ACTION_TYPES.pbehavior],
      duration: prepareDuration(null, DEFAULT_IDLE_RULE_PBEHAVIOR_DURATION),
    };
  }

  return parameters;
};

/**
 * Prepare numeric idle-rule fields with lower bounds before saving.
 *
 * @param {Object} form Idle rule form.
 * @returns {Object}
 */
export const prepareIdleRuleNumbers = (form) => {
  const preparedForm = {
    ...form,
    priority: prepareMinNumber(form.priority, MIN_PRIORITY),
    duration: form.duration
      ? {
        ...form.duration,
        value: prepareMinNumber(form.duration.value, MIN_DURATION_VALUE),
      }
      : form.duration,
  };
  const snoozeParameters = form.operation?.parameters?.[ACTION_TYPES.snooze];
  const pbehaviorParameters = form.operation?.parameters?.[ACTION_TYPES.pbehavior];

  if (snoozeParameters?.duration) {
    preparedForm.operation = {
      ...preparedForm.operation,
      parameters: {
        ...preparedForm.operation.parameters,
        [ACTION_TYPES.snooze]: {
          ...snoozeParameters,
          duration: prepareDuration(snoozeParameters.duration, DEFAULT_IDLE_RULE_SNOOZE_DURATION),
        },
      },
    };
  }

  if (pbehaviorParameters && form.operation?.type === ACTION_TYPES.pbehavior) {
    preparedForm.operation = {
      ...preparedForm.operation,
      parameters: {
        ...preparedForm.operation.parameters,
        [ACTION_TYPES.pbehavior]: {
          ...pbehaviorParameters,
          duration: prepareDuration(
            pbehaviorParameters.duration,
            DEFAULT_IDLE_RULE_PBEHAVIOR_DURATION,
          ),
        },
      },
    };
  }

  return preparedForm;
};

/**
 * Convert idle rule object to form compatible object.
 *
 * @param {Object} [idleRule={}] Idle rule entity.
 * @returns {Object} Idle rule form.
 */
export const idleRuleToForm = (idleRule = {}) => {
  const type = idleRule.operation?.type ?? ACTION_TYPES.snooze;
  const parameters = idleRuleActionParametersToForm(idleRule.operation);

  return {
    enabled: idleRule.enabled ?? true,
    name: idleRule.name ?? '',
    description: idleRule.description ?? '',
    comment: idleRule.comment ?? '',
    type: idleRule.type ?? IDLE_RULE_TYPES.alarm,
    duration: idleRule.duration
      ? durationToForm(idleRule.duration)
      : { value: MIN_DURATION_VALUE, unit: TIME_UNITS.minute },
    priority: idleRule.priority,
    disable_during_periods: idleRule.disable_during_periods ?? [],
    alarm_condition: idleRule.alarm_condition ?? IDLE_RULE_ALARM_CONDITIONS.lastEvent,
    operation: {
      type,
      parameters,
    },
    patterns: filterPatternsToForm(idleRule, [PATTERNS_FIELDS.entity, PATTERNS_FIELDS.alarm]),
  };
};

/**
 * Convert form object to idle rule API compatible object.
 *
 * @param {Object} form Idle rule form.
 * @returns {Object} Idle rule entity.
 */
export const formToIdleRule = (form) => {
  const preparedForm = prepareIdleRuleNumbers(form);
  const isEntityType = form.type === IDLE_RULE_TYPES.entity;
  const idleRule = omit(preparedForm, ['alarm_condition', 'operation', 'patterns', 'comment']);

  if (!isEntityType) {
    idleRule.alarm_condition = preparedForm.alarm_condition;
    idleRule.operation = pick(formToAction(preparedForm.operation), ['type', 'parameters']);

    if (isAssociateTicketActionType(idleRule.operation.type)) {
      idleRule.comment = preparedForm.comment;
    }
  }

  return {
    ...idleRule,
    ...formFilterToPatterns(
      preparedForm.patterns,
      isEntityType ? [PATTERNS_FIELDS.entity] : [PATTERNS_FIELDS.entity, PATTERNS_FIELDS.alarm],
    ),
  };
};
