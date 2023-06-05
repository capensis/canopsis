import { COLORS } from '@/config';

import { PBEHAVIOR_ORIGINS, PBEHAVIOR_TYPE_TYPES, WEATHER_ENTITY_PBEHAVIOR_DEFAULT_TITLE } from '@/constants';

import uid from '@/helpers/uid';
import { formToPbehavior, pbehaviorToRequest } from '@/helpers/forms/planning-pbehavior';
import { getNowTimestamp } from '@/helpers/date/date';

/**
 * Check if pbehavior is paused
 *
 * @param {Pbehavior} pbehavior
 * @return {boolean}
 */
export const isPausedPbehavior = pbehavior => pbehavior.type.type === PBEHAVIOR_TYPE_TYPES.pause;

/**
 * Check is not active pbehavior type
 *
 * @param {string} type
 * @returns {boolean}
 */
export const isNotActivePbehaviorType = type => [
  PBEHAVIOR_TYPE_TYPES.pause,
  PBEHAVIOR_TYPE_TYPES.inactive,
  PBEHAVIOR_TYPE_TYPES.maintenance,
].includes(type);

/**
 * Check if pbehaviors have a paused type
 *
 * @param {Pbehavior[]} pbehaviors
 * @return {boolean}
 */
export const hasPausedPbehavior = pbehaviors => pbehaviors.some(isPausedPbehavior);

/**
 * Create downtime pbehavior, without stop time
 *
 * @param {Entity} entity
 * @param {PbehaviorReason} reason
 * @param {string} comment
 * @param {PbehaviorType} type
 * @return {Pbehavior}
 */
export const createDowntimePbehavior = ({ entity, reason, comment, type }) => pbehaviorToRequest(formToPbehavior({
  reason,
  type,
  origin: PBEHAVIOR_ORIGINS.serviceWeather,
  color: COLORS.secondary,
  name: `${WEATHER_ENTITY_PBEHAVIOR_DEFAULT_TITLE}-${entity.name}-${uid()}`,
  tstart: getNowTimestamp(),
  comment,
  entity: entity._id,
}));

/**
 * Get color for pbehavior
 *
 * @param pbehavior
 * @returns {string}
 */
export const getPbehaviorColor = (pbehavior = {}) => pbehavior.color || pbehavior.type?.color || COLORS.secondary;
