import {
  ALARM_STATES_ICONS,
  ALARM_STATES_TEXTS,
  ALARM_STATES_CLASSES,
  ALARM_STATUSES_ICONS,
  ALARM_STATUSES_RESOLVED_ICONS,
  ALARM_STATUSES_TEXTS,
  ALARM_UNKNOWN_VALUE,
  ALARM_STATES_UNKNOWN_CLASS,
} from '@/constants';

import { getAlarmStateColor, getAlarmStatusColor } from './color';

/**
 * Return entity state icon by state value
 *
 * @param {number} [value]
 * @returns {string}
 */
export const getAlarmStateIcon = value => ALARM_STATES_ICONS[value];

/**
 * Return entity state text by state value
 *
 * @param {number} [value]
 * @returns {string}
 */
export const getAlarmStateText = value => ALARM_STATES_TEXTS[value];

/**
 * Return object that contains the state style
 *
 * @param {number} state
 * @returns {{ icon: string, text: string, color: string }}
 */
export const formatAlarmState = state => ({
  icon: getAlarmStateIcon(state) ?? ALARM_UNKNOWN_VALUE.icon,
  text: getAlarmStateText(state) ?? ALARM_UNKNOWN_VALUE.text,
  color: getAlarmStateColor(state) ?? ALARM_UNKNOWN_VALUE.color,
  class: ALARM_STATES_CLASSES[state] ?? ALARM_STATES_UNKNOWN_CLASS,
});

/**
 * Return entity status text by status value
 *
 * @param {number} [value]
 * @returns {string}
 */
export const getAlarmStatusText = value => ALARM_STATUSES_TEXTS[value];

/**
 * Return entity status icon by status value
 *
 * @param {number} [value]
 * @param {boolean} resolved
 * @returns {string}
 */
export const getAlarmStatusIcon = (value, resolved) => (
  resolved
    ? ALARM_STATUSES_RESOLVED_ICONS[value] ?? ALARM_STATUSES_ICONS[value]
    : ALARM_STATUSES_ICONS[value]
);

/**
 * Return object that contains the status style
 *
 * @param {number} status
 * @param {boolean} resolved
 * @returns {{ icon: string, text: string, color: string }}
 */
export const formatAlarmStatus = (status, resolved) => ({
  icon: getAlarmStatusIcon(status, resolved) ?? ALARM_UNKNOWN_VALUE.icon,
  text: getAlarmStatusText(status) ?? ALARM_UNKNOWN_VALUE.text,
  color: getAlarmStatusColor(status) ?? ALARM_UNKNOWN_VALUE.color,
});
