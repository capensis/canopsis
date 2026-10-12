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

/**
 * Return trend field name by parameter and show type
 *
 * @param {number} displayParameter
 * @return {string}
 */
export const getAvailabilityTrendFieldByDisplayParameter = displayParameter => (
  displayParameter === AVAILABILITY_DISPLAY_PARAMETERS.uptime
    ? AVAILABILITY_FIELDS.uptimeShareHistory
    : AVAILABILITY_FIELDS.downtimeShareHistory
);

/**
 * Convert infos to common infos structure
 *
 * @param {Object.<string, string | string[]>[]} infos
 * @return {InfosObject}
 */
const prepareAvailabilityInfos = infos => (infos
  ? Object.entries(infos).reduce((acc, [name, value]) => {
    acc[name] = {
      value,
      name,
    };

    return acc;
  }, {})
  : {});

/**
 * Convert infos to common infos structure
 *
 * @param {Availability[]} data
 * @return {Availability[]}
 */
export const prepareAvailabilitiesResponse = ({ data }) => data.map(availability => ({
  ...availability,
  entity: {
    ...availability.entity,
    infos: prepareAvailabilityInfos(availability.entity.infos),
    component_infos: prepareAvailabilityInfos(availability.entity.component_infos),
  },
}));
