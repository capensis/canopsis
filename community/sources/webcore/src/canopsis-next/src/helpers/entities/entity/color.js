import { isNumber } from 'lodash';

import { COLORS, CSS_COLORS_VARS } from '@/config';
import { COLOR_INDICATOR_TYPES, ENTITIES_STATES, EVENT_ENTITY_COLORS_BY_TYPE } from '@/constants';

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
  [ENTITIES_STATES.ok]: CSS_COLORS_VARS.state.ok,
  [ENTITIES_STATES.minor]: CSS_COLORS_VARS.state.minor,
  [ENTITIES_STATES.major]: CSS_COLORS_VARS.state.major,
  [ENTITIES_STATES.critical]: CSS_COLORS_VARS.state.critical,
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
    return CSS_COLORS_VARS.state.pause;
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
