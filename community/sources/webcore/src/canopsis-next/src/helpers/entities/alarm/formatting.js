import {
  ALARM_STATES_ICONS,
  ALARM_STATES_TEXTS,
  ALARM_STATUSES_ICONS,
  ALARM_STATUSES_TEXTS,
  ALARM_UNKNOWN_VALUE,
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
 * @returns {string}
 */
export const getAlarmStatusIcon = value => ALARM_STATUSES_ICONS[value];

/**
 * Return object that contains the status style
 *
 * @param {number} status
 * @returns {{ icon: string, text: string, color: string }}
 */
export const formatAlarmStatus = status => ({
  icon: getAlarmStatusIcon(status) ?? ALARM_UNKNOWN_VALUE.icon,
  text: getAlarmStatusText(status) ?? ALARM_UNKNOWN_VALUE.text,
  color: getAlarmStatusColor(status) ?? ALARM_UNKNOWN_VALUE.color,
});
