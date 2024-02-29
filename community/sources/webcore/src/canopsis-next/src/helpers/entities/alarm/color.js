import { CSS_COLORS_VARS } from '@/config';
import { ALARM_STATES, ALARM_STATUSES } from '@/constants';

/**
 * Get color by entity impact state
 *
 * @param {number} state
 * @returns {string}
 */
export const getAlarmStateColor = state => ({
  [ALARM_STATES.ok]: CSS_COLORS_VARS.state.ok,
  [ALARM_STATES.minor]: CSS_COLORS_VARS.state.minor,
  [ALARM_STATES.major]: CSS_COLORS_VARS.state.major,
  [ALARM_STATES.critical]: CSS_COLORS_VARS.state.critical,
}[state]);

/**
 * Get color by alarm status
 *
 * @param {number} value
 * @returns {string}
 */
export const getAlarmStatusColor = value => ({
  [ALARM_STATUSES.closed]: CSS_COLORS_VARS.status.closed,
  [ALARM_STATUSES.ongoing]: CSS_COLORS_VARS.status.ongoing,
  [ALARM_STATUSES.stealthy]: CSS_COLORS_VARS.status.stealthy,
  [ALARM_STATUSES.flapping]: CSS_COLORS_VARS.status.flapping,
  [ALARM_STATUSES.cancelled]: CSS_COLORS_VARS.status.cancelled,
  [ALARM_STATUSES.noEvents]: CSS_COLORS_VARS.status.noEvents,
}[value]);
