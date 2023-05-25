import { get, isNumber } from 'lodash';

import { COLORS } from '@/config';
import { COLOR_INDICATOR_TYPES, ENTITIES_STATES_STYLES, EVENT_ENTITY_COLORS_BY_TYPE } from '@/constants';

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
export const getEntityStateColor = value => get(ENTITIES_STATES_STYLES, [value, 'color']);

/**
 * Get color for a entity by colorIndicator and isGrey parameters
 *
 * @param {Service | Entity | {}} [entity = {}]
 * @param {string} [colorIndicator = COLOR_INDICATOR_TYPES.state]
 * @returns {string|*}
 */
export const getEntityColor = (entity = {}, colorIndicator = COLOR_INDICATOR_TYPES.state) => {
  if (entity.is_grey) {
    return COLORS.state.pause;
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
