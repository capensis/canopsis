import { MAX_PBEHAVIOR_DEFAULT_TSTOP, PBEHAVIOR_TYPE_TYPES, WEATHER_ENTITY_PBEHAVIOR_DEFAULT_TITLE } from '@/constants';

import { formToPbehavior, pbehaviorToRequest } from '@/helpers/forms/planning-pbehavior';
import uid from '@/helpers/uid';

/**
 * Check if pbehavior is active
 *
 * @param {Pbehavior} pbehavior
 * @return {boolean}
 */
export const isActivePbehavior = pbehavior => pbehavior.type.type === PBEHAVIOR_TYPE_TYPES.active;

/**
 * Check if pbehavior is paused
 *
 * @param {Pbehavior} pbehavior
 * @return {boolean}
 */
export const isPausedPbehavior = pbehavior => pbehavior.type.type === PBEHAVIOR_TYPE_TYPES.pause;

/**
 * Check if pbehaviors have a active type
 *
 * @param {Pbehavior[]} pbehaviors
 * @return {boolean}
 */
export const hasActivePbehavior = pbehaviors => pbehaviors.some(isActivePbehavior);

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
  name: `${WEATHER_ENTITY_PBEHAVIOR_DEFAULT_TITLE}-${entity.name}-${uid()}`,
  tstart: new Date(),
  tstop: new Date(MAX_PBEHAVIOR_DEFAULT_TSTOP * 1000),
  comments: [{
    message: comment,
  }],
  filter: {
    _id: { $in: [entity._id] },
  },
}));
