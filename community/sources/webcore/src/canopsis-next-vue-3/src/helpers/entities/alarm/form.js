import { ALARM_STATES, ALARM_STATUSES, TIME_UNITS } from '@/constants';

/**
 * Check alarm state is ok.
 *
 * @param {Object} alarm
 * @returns {boolean}
 */
export const isAlarmStateOk = alarm => alarm.v?.state?.val === ALARM_STATES.ok;

/**
 * Check alarm status is closed.
 *
 * @param {Object} alarm
 * @returns {boolean}
 */
export const isAlarmStatusClosed = alarm => alarm.v?.status?.val === ALARM_STATUSES.closed;

/**
 * Check alarm status is cancelled.
 *
 * @param {Object} alarm
 * @returns {boolean}
 */
export const isAlarmStatusCancelled = alarm => alarm.v?.status?.val === ALARM_STATUSES.cancelled;

/**
 * Check alarm status is ongoing.
 *
 * @param {Object} alarm
 * @returns {boolean}
 */
export const isAlarmStatusOngoing = alarm => alarm.v?.status?.val === ALARM_STATUSES.ongoing;

/**
 * Check alarm status is no events.
 *
 * @param {Object} alarm
 * @returns {boolean}
 */
export const isAlarmStatusNoEvents = alarm => alarm.v?.status?.val === ALARM_STATUSES.noEvents;

/**
 * Check alarm status is flapping.
 *
 * @param {Object} alarm
 * @returns {boolean}
 */
export const isAlarmStatusFlapping = alarm => alarm.v?.status?.val === ALARM_STATUSES.flapping;

/**
 * Check alarm status is unknown.
 *
 * @param {Object} alarm
 * @returns {boolean}
 */
export const isAlarmStatusUnknown = alarm => alarm.v?.status?.val === ALARM_STATUSES.unknown;

/**
 * Checks if alarm is resolved.
 *
 * @param {Object} alarm - Alarm entity.
 * @returns {boolean}
 */
export const isResolvedAlarm = alarm => !!alarm.v?.resolved;

/**
 * Checks if action available for alarm.
 *
 * @param {Object} alarm - Alarm entity.
 * @param {Object} widget - Alarm widget.
 * @returns {boolean}
 */
export const isActionAvailableForAlarm = (alarm, widget) => {
  /**
   * When alarm state is ok, but we enable actions with ok state.
   */
  if (isAlarmStateOk(alarm)) {
    return widget?.parameters?.isActionsAllowWithOkState;
  }

  /**
   * When alarm is cancelled, but not resolved we can uncancel by mass actions.
   */
  if (isAlarmStatusCancelled(alarm)) {
    return !isResolvedAlarm(alarm);
  }

  return !isResolvedAlarm(alarm) && !isAlarmStatusClosed(alarm);
};

/**
 * Checks if alarm has non-OK state.
 *
 * @param {Object} alarm - Alarm entity.
 * @returns {boolean}
 */
export const isAlarmStateNotOk = alarm => ALARM_STATES.ok !== alarm.v?.state?.val;

/**
 * Convert snooze object to form snooze.
 *
 * @param {Object} [snooze={}]
 * @returns {Object}
 */
export const snoozeToForm = (snooze = {}) => ({
  duration: {
    unit: snooze.duration?.unit ?? TIME_UNITS.minute,
    value: snooze.duration?.seconds ?? 1,
  },
  comment: snooze.comment ?? '',
});

/**
 * Convert change state form to event request.
 *
 * @param {Object} form
 * @returns {Object}
 */
export const formToChangeStateEvent = (form = {}) => ({
  state: form.state,
  comment: form.output ?? form.comment ?? '',
});

/**
 * Maps an array of alarm objects to a unique list of their associated entities.
 *
 * @param {Array<Object>} alarms - Alarm objects.
 * @returns {Array<Object>} Unique entity objects extracted from alarms.
 */
export const mapAlarmsEntities = (alarms = []) =>
  Object.values(
    alarms.reduce((acc, alarm) => {
      acc[alarm?.entity?._id] = alarm?.entity;

      return acc;
    }, {}),
  ).filter(Boolean);
