import { AVAILABILITY_DISPLAY_PARAMETERS, AVAILABILITY_FIELDS, AVAILABILITY_SHOW_TYPE } from '@/constants';

/**
 * Return field name by parameter and show type
 *
 * @param {number} displayParameter
 * @param {number} showType
 * @return {string}
 */
export const getAvailabilityFieldByDisplayParameterAndShowType = (displayParameter, showType) => {
  const isPercentType = showType === AVAILABILITY_SHOW_TYPE.percent;

  if (displayParameter === AVAILABILITY_DISPLAY_PARAMETERS.uptime) {
    return isPercentType ? AVAILABILITY_FIELDS.uptimeShare : AVAILABILITY_FIELDS.uptimeDuration;
  }

  return isPercentType ? AVAILABILITY_FIELDS.downtimeShare : AVAILABILITY_FIELDS.downtimeDuration;
};
