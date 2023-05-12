import { ENTITIES_STATES, ENTITIES_STATUSES } from '@/constants';

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
