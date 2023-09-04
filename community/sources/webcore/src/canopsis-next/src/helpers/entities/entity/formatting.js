import {
  ENTITIES_STATES_STYLES_ICONS,
  ENTITIES_STATES_STYLES_TEXT,
  ENTITY_STATUS_STYLES,
  UNKNOWN_VALUE_STYLE,
} from '@/constants';

import { getEntityEventColor, getEntityStateColor } from '@/helpers/entities/entity/color';
import { getEntityEventIcon } from '@/helpers/entities/entity/icons';

/**
 * Return entity state text by state value
 *
 * @param {number} [value]
 * @returns {string}
 */
export const getEntityStateText = value => ENTITIES_STATES_STYLES_TEXT[value];

/**
 * Return entity state icon by state value
 *
 * @param {number} [value]
 * @returns {string}
 */
export const getEntityStateIcon = value => ENTITIES_STATES_STYLES_ICONS[value];

/**
 * Return object that contains the state style
 *
 * @param value The state value
 * @returns {*} Object with the color, icon and text associated
 */
export const formatState = value => ({
  icon: getEntityStateIcon(value) ?? UNKNOWN_VALUE_STYLE.icon,
  text: getEntityStateText(value) ?? UNKNOWN_VALUE_STYLE.text,
  color: getEntityStateColor(value) ?? UNKNOWN_VALUE_STYLE.color,
});

/**
 * Return object that contains the status style
 * @param value The status value
 * @returns {*} Object with the color, icon and text associated
 */
export const formatStatus = (value) => {
  if (!ENTITY_STATUS_STYLES[value]) {
    return UNKNOWN_VALUE_STYLE;
  }

  return ENTITY_STATUS_STYLES[value];
};

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
