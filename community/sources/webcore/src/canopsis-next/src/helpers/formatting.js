import { ENTITIES_STATES_STYLES, ENTITY_STATUS_STYLES, UNKNOWN_VALUE_STYLE } from '@/constants';

import { getEntityEventColor } from '@/helpers/color';
import { getEntityEventIcon } from '@/helpers/icon';

/**
 * Return object that contains the state style
 * @param value The state value
 * @returns {*} Object with the color, icon and text associated
 */
export function formatState(value) {
  if (!ENTITIES_STATES_STYLES[value]) {
    return UNKNOWN_VALUE_STYLE;
  }

  return ENTITIES_STATES_STYLES[value];
}

/**
 * Return object that contains the status style
 * @param value The status value
 * @returns {*} Object with the color, icon and text associated
 */
export function formatStatus(value) {
  if (!ENTITY_STATUS_STYLES[value]) {
    return UNKNOWN_VALUE_STYLE;
  }

  return ENTITY_STATUS_STYLES[value];
}

/**
 * Return object that contains the event style
 *
 * @param {string} event The event name
 * @returns {*} Object with the color, icon and text associated
 */
export const formatEvent = (event) => {
  const icon = getEntityEventIcon(event);

  if (!icon) {
    return UNKNOWN_VALUE_STYLE;
  }

  return {
    icon,
    color: getEntityEventColor(event),
  };
};

/**
 * Return object that contains the step style
 *
 * @param {AlarmEvent} step
 * @returns {Object}
 */
export const formatStep = (step) => {
  if (step._t.startsWith('status')) {
    return formatStatus(step.val);
  }

  if (step._t.startsWith('state')) {
    return formatState(step.val);
  }

  return formatEvent(step._t);
};
