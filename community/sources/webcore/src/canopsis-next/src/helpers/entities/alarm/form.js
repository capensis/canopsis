import { ENTITIES_STATES, ENTITIES_STATUSES, TIME_UNITS } from '@/constants';

/**
 * @typedef {Object} AlarmOldLink
 * @property {string} label
 * @property {string} link
 */

/**
 * @typedef {Object} AlarmLink
 * @property {string} icon_name
 * @property {string} label
 * @property {string} url
 * @property {boolean} [single]
 * @property {string} [rule_id]
 * @property {LinkRuleAction} [action]
 */

/**
 * @typedef {Object<string, AlarmLink[] | AlarmOldLink[]>} AlarmLinks
 */

/**
 * @typedef {Object} AlarmEvent
 * @property {number} val
 * @property {string} a
 * @property {number} t
 * @property {string} _t
 * @property {string} m
 */

/**
 * @typedef {Object} Correlation
 * @property {number} total
 * @property {Alarm} data
 */

/**
 * @typedef {Object} AlarmRule
 * @property {string} id
 * @property {string} name
 */

/**
 * @typedef {Object} AlarmAssignedInstructionExecution
 * @property {string} _id
 * @property {number} status
 */

/**
 * @typedef {Object} AlarmAssignedInstruction
 * @property {string} _id
 * @property {string} name
 * @property {?AlarmAssignedInstructionExecution} execution
 */

/**
 * @typedef {Object} AlarmValuePbehaviorInfo
 * @property {string} canonical_type
 * @property {string} icon_name
 * @property {string} id
 * @property {string} name
 * @property {string} reason
 * @property {string} reason_name
 * @property {string} type
 * @property {string} type_name
 * @property {number} timestamp
 */

/**
 * @typedef {Object} AlarmValue
 * @property {string[]} long_output_history
 * @property {number} last_event_date
 * @property {string} connector_name
 * @property {string} initial_long_output
 * @property {string} output
 * @property {string} initial_output
 * @property {string[]} children
 * @property {Object} extra
 * @property {number} last_update_date
 * @property {number} resolved
 * @property {string} resource
 * @property {number} creation_date
 * @property {string} display_name
 * @property {number} total_state_changes
 * @property {string[]} tags
 * @property {string} long_output
 * @property {string} component
 * @property {string} connector
 * @property {Object} infos_rule_version
 * @property {InfosObject} infos
 * @property {string[]} parents
 * @property {AlarmEvent} ack
 * @property {AlarmEvent} state
 * @property {AlarmEvent} status
 * @property {AlarmEvent} ticket
 * @property {AlarmEvent} snooze
 * @property {AlarmEvent} canceled
 * @property {AlarmEvent} last_comment
 * @property {AlarmEvent[]} [steps]
 * @property {AlarmValuePbehaviorInfo} [pbehavior_info]
 */

/**
 * @typedef {Pbehavior} AlarmPbehavior
 * @property {Comment} last_comment
 */

/**
 * @typedef {Object} Alarm
 * @property {string} _id
 * @property {Entity} entity
 * @property {boolean} is_meta_alarm
 * @property {AlarmAssignedInstruction[]} [assigned_instructions]
 * @property {number} [instruction_execution_icon]
 * @property {string[]} running_manual_instructions
 * @property {string[]} running_auto_instructions
 * @property {string[]} failed_manual_instructions
 * @property {string[]} failed_auto_instructions
 * @property {string[]} successful_manual_instructions
 * @property {string[]} successful_auto_instructions
 * @property {boolean} [children_instructions]
 * @property {InfosObject} infos
 * @property {AlarmRule} rule
 * @property {Correlation} consequences
 * @property {Correlation} causes
 * @property {AlarmPbehavior} pbehavior
 * @property {AlarmLinks} links
 * @property {string[]} tags
 * @property {number} t
 * @property {AlarmValue} v
 */

/**
 * Check alarm state is ok
 *
 * @param {Alarm} alarm
 * @returns {boolean}
 */
export const isAlarmStateOk = alarm => alarm.v?.state?.val === ENTITIES_STATES.ok;

/**
 * Check alarm status is closed
 *
 * @param {Alarm} alarm
 * @returns {boolean}
 */
export const isAlarmStatusClosed = alarm => alarm.v?.status?.val === ENTITIES_STATUSES.closed;

/**
 * Check alarm status is cancelled
 *
 * @param {Alarm} alarm
 * @returns {boolean}
 */
export const isAlarmStatusCancelled = alarm => alarm.v?.status?.val === ENTITIES_STATUSES.cancelled;

/**
 * Check alarm status is ongoing
 *
 * @param {Alarm} alarm
 * @returns {boolean}
 */
export const isAlarmStatusOngoing = alarm => alarm.v?.status?.val === ENTITIES_STATUSES.ongoing;

/**
 * Check alarm status is flapping
 *
 * @param {Alarm} alarm
 * @returns {boolean}
 */
export const isAlarmStatusFlapping = alarm => alarm.v?.status?.val === ENTITIES_STATUSES.flapping;

/**
 * @typedef {Object} SnoozeAction
 * @property {number} duration
 * @property {string} comment
 */

/**
 * @typedef {SnoozeAction} SnoozeActionForm
 * @property {Duration} duration
 */

/**
 * Checks if alarm is cancelled
 *
 * @param {Alarm} alarm - alarm entity
 * @returns {boolean}
 */
export const isCancelledAlarmStatus = alarm => alarm.v.status.val === ENTITIES_STATUSES.cancelled;

/**
 * Checks if alarm is closed
 *
 * @param {Alarm} alarm - alarm entity
 * @returns {boolean}
 */
export const isClosedAlarmStatus = alarm => alarm.v.status.val === ENTITIES_STATUSES.closed;

/**
 * Checks if alarm is resolved
 *
 * @param alarm - alarm entity
 * @returns {boolean}
 */
export const isResolvedAlarm = alarm => !!alarm.v.resolved;

/**
 * Checks if action available for alarm
 *
 * @param {Alarm} alarm - alarm entity
 * @param {Widget} widget - alarm entity
 * @returns {boolean}
 */
export const isActionAvailableForAlarm = (alarm, widget) => {
  /**
   * When alarm state is ok, but we enable actions with ok state
   */
  if (isAlarmStateOk(alarm)) {
    return widget?.parameters?.isActionsAllowWithOkState;
  }

  /**
   * When alarm is cancelled, but not resolved we can uncancel by mass actions
   */
  if (isCancelledAlarmStatus(alarm)) {
    return !isResolvedAlarm(alarm);
  }

  return !isResolvedAlarm(alarm)
    && !isClosedAlarmStatus(alarm);
};

/**
 * Checks if alarm have critical state
 *
 * @param alarm - alarm entity
 * @returns {boolean}
 */
export const isAlarmStateNotOk = alarm => ENTITIES_STATES.ok !== alarm.v.state.val;

/**
 * Convert snooze object to form snooze
 *
 * @param {SnoozeAction} snooze
 * @returns {SnoozeActionForm}
 */
export const snoozeToForm = (snooze = {}) => ({
  duration: {
    unit: snooze.duration?.unit ?? TIME_UNITS.minute,
    value: snooze.duration?.seconds ?? 1,
  },
  comment: snooze.comment ?? '',
});
