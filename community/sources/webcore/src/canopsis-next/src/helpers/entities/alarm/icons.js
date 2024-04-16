import { ALARM_LIST_ACTIONS_TYPES_ICONS } from '@/constants';

/**
 * Get icon for alarm action
 *
 * @param {string} type
 */
export const getAlarmActionIcon = type => ALARM_LIST_ACTIONS_TYPES_ICONS[type];
