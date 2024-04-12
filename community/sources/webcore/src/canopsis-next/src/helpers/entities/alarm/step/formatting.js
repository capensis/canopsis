import { ALARM_STEPS_COLORS, ALARM_STEPS_ICONS, ALARM_UNKNOWN_VALUE } from '@/constants';

import { formatAlarmState, formatAlarmStatus } from '../formatting';

/**
 * Get color for entity event
 *
 * @param {string} type
 */
export const getAlarmStepColor = type => ALARM_STEPS_COLORS[type];

/**
 * Get color for entity event
 *
 * @param {string} type
 */
export const getAlarmStepIcon = type => ALARM_STEPS_ICONS[type];

/**
 * Return object that contains the step style for horizontal timeline
 *
 * @param {AlarmStep} step
 * @returns {{ icon: string, color: string }}
 */
export const formatNotificationAlarmStep = (step) => {
  if (step._t.startsWith('status')) {
    return formatAlarmStatus(step.val);
  }

  if (step._t.startsWith('state')) {
    return formatAlarmState(step.val);
  }

  const icon = getAlarmStepIcon(step._t);

  if (!icon) {
    return ALARM_UNKNOWN_VALUE;
  }

  return {
    icon,
    color: getAlarmStepColor(step._t),
  };
};
