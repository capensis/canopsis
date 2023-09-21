import { isNumber } from 'lodash';

import { COLORS, CSS_COLOR_VARS } from '@/config';
import { COLOR_INDICATOR_TYPES, ENTITIES_STATES, ENTITIES_STATUSES, EVENT_ENTITY_COLORS_BY_TYPE } from '@/constants';

/**
 * Get color by entity impact state
 *
 * @param {number} value
 * @returns {string}
 */
export const getImpactStateColor = value => COLORS.impactState[value];

/**
 * Get color by entity impact state
 *
 * @param {number} value
 * @returns {string}
 */
export const getEntityStateColor = value => ({
  [ENTITIES_STATES.ok]: CSS_COLOR_VARS.state.ok,
  [ENTITIES_STATES.minor]: CSS_COLOR_VARS.state.minor,
  [ENTITIES_STATES.major]: CSS_COLOR_VARS.state.major,
  [ENTITIES_STATES.critical]: CSS_COLOR_VARS.state.critical,
}[value]);

/**
 * Get color by entity status
 *
 * @param {number} value
 * @returns {string}
 */
export const getEntityStatusColor = value => ({
  [ENTITIES_STATUSES.closed]: CSS_COLOR_VARS.status.closed,
  [ENTITIES_STATUSES.ongoing]: CSS_COLOR_VARS.status.ongoing,
  [ENTITIES_STATUSES.stealthy]: CSS_COLOR_VARS.status.stealthy,
  [ENTITIES_STATUSES.flapping]: CSS_COLOR_VARS.status.flapping,
  [ENTITIES_STATUSES.cancelled]: CSS_COLOR_VARS.status.cancelled,
  [ENTITIES_STATUSES.noEvents]: CSS_COLOR_VARS.status.noEvents,
}[value]);

/**
 * Get color for a entity by colorIndicator and isGrey parameters
 *
 * @param {Service | Entity | {}} [entity = {}]
 * @param {string} [colorIndicator = COLOR_INDICATOR_TYPES.state]
 * @returns {string|*}
 */
export const getEntityColor = (entity = {}, colorIndicator = COLOR_INDICATOR_TYPES.state) => {
  if (entity.is_grey) {
    return CSS_COLOR_VARS.state.pause;
  }

  if (colorIndicator === COLOR_INDICATOR_TYPES.state) {
    const state = isNumber(entity.state) ? entity.state : entity.state?.val;

    return getEntityStateColor(state);
  }

  return getImpactStateColor(entity.impact_state);
};

/**
 * Get color for entity event
 *
 * @param {string} type
 */
export const getEntityEventColor = type => EVENT_ENTITY_COLORS_BY_TYPE[type];
