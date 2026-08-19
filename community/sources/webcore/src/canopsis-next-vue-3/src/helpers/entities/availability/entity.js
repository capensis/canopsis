import { AVAILABILITY_DISPLAY_PARAMETERS, AVAILABILITY_FIELDS, AVAILABILITY_SHOW_TYPE } from '@/constants';

/**
 * @typedef {Object} Availability
 * @property {number} downtime_duration
 * @property {number} downtime_share
 * @property {number} uptime_duration
 * @property {number} uptime_share
 * @property {number} uptime_share_history
 * @property {Entity} entity
 */

/**
 * Return field name by parameter and show type.
 *
 * @param {number} displayParameter
 * @param {number} showType
 * @returns {string}
 */
export const getAvailabilityFieldByDisplayParameterAndShowType = (displayParameter, showType) => {
  const isPercentType = showType === AVAILABILITY_SHOW_TYPE.percent;

  if (displayParameter === AVAILABILITY_DISPLAY_PARAMETERS.uptime) {
    return isPercentType ? AVAILABILITY_FIELDS.uptimeShare : AVAILABILITY_FIELDS.uptimeDuration;
  }

  return isPercentType ? AVAILABILITY_FIELDS.downtimeShare : AVAILABILITY_FIELDS.downtimeDuration;
};

/**
 * Return trend field name by display parameter.
 *
 * @param {number} displayParameter
 * @returns {string}
 */
export const getAvailabilityTrendFieldByDisplayParameter = displayParameter =>
  displayParameter === AVAILABILITY_DISPLAY_PARAMETERS.uptime
    ? AVAILABILITY_FIELDS.uptimeShareHistory
    : AVAILABILITY_FIELDS.downtimeShareHistory;
