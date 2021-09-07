import { COLORS } from '@/config';
import {
  ENTITIES_STATES_STYLES, ENTITY_STATUS_STYLES, EVENT_ENTITY_STYLE, UNKNOWN_VALUE_STYLE,
} from '@/constants';

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
 * Return object that contains the state style
 *
 * @param {number} value - The state value
 * @returns {*} Object with the color, icon and text associated
 */
export function formatImpactState(value) { // TODO: refactor the function
  if (!COLORS.impactState[value]) {
    return UNKNOWN_VALUE_STYLE;
  }

  return {
    color: COLORS.impactState[value],
    text: `${value} of ${COLORS.impactState.length - 1}`,
  };
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
 * @param event The event name
 * @returns {*} Object with the color, icon and text associated
 */
export function formatEvent(event) {
  if (!EVENT_ENTITY_STYLE[event]) {
    return UNKNOWN_VALUE_STYLE;
  }
  return EVENT_ENTITY_STYLE[event];
}

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
