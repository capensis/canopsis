import { groupBy } from 'lodash';

import {
  ENTITIES_STATES_STYLES_ICONS,
  ENTITIES_STATES_STYLES_TEXT,
  ENTITIES_STATUSES_STYLES_ICONS,
  ENTITIES_STATUSES_STYLES_TEXT,
  UNKNOWN_VALUE_STYLE,
  DATETIME_FORMATS,
} from '@/constants';

import { getEntityEventColor, getEntityStateColor, getEntityStatusColor } from '@/helpers/entities/entity/color';
import { getEntityEventIcon } from '@/helpers/entities/entity/icons';
import { convertDateToString } from '@/helpers/date/date';
import { addKeyInEntities } from '@/helpers/array';

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
 * Return entity status text by status value
 *
 * @param {number} [value]
 * @returns {string}
 */
export const getEntityStatusText = value => ENTITIES_STATUSES_STYLES_TEXT[value];

/**
 * Return entity status icon by status value
 *
 * @param {number} [value]
 * @returns {string}
 */
export const getEntityStatusIcon = value => ENTITIES_STATUSES_STYLES_ICONS[value];

/**
 * Return object that contains the status style
 * @param value The status value
 * @returns {*} Object with the color, icon and text associated
 */
export const formatStatus = value => ({
  icon: getEntityStatusIcon(value) ?? UNKNOWN_VALUE_STYLE.icon,
  text: getEntityStatusText(value) ?? UNKNOWN_VALUE_STYLE.text,
  color: getEntityStatusColor(value) ?? UNKNOWN_VALUE_STYLE.color,
});

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

/**
 * Get grouped steps by date
 *
 * @param {AlarmEvent[]} steps
 * @return {Object.<string, AlarmEvent[]>}
 */
export const groupAlarmSteps = steps => (
  groupBy(addKeyInEntities(steps), step => convertDateToString(step.t, DATETIME_FORMATS.short))
);
